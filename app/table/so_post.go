package table

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/coschain/contentos-go/common/encoding/kope"
	"github.com/coschain/contentos-go/iservices"
	prototype "github.com/coschain/contentos-go/prototype"
	proto "github.com/golang/protobuf/proto"
)

////////////// SECTION Prefix Mark ///////////////
var (
	PostCreatedTable         uint32 = 3346451556
	PostCashoutBlockNumTable uint32 = 1826021466
	PostRewardsTable         uint32 = 2325142906
	PostPostIdUniTable       uint32 = 157486700

	PostPostIdRow uint32 = 3809844522
)

////////////// SECTION Wrap Define ///////////////
type SoPostWrap struct {
	dba         iservices.IDatabaseRW
	mainKey     *uint64
	watcherFlag *PostWatcherFlag
	mKeyFlag    int    //the flag of the main key exist state in db, -1:has not judged; 0:not exist; 1:already exist
	mKeyBuf     []byte //the buffer after the main key is encoded with prefix
	mBuf        []byte //the value after the main key is encoded
	mdFuncMap   map[string]interface{}
}

func NewSoPostWrap(dba iservices.IDatabaseRW, key *uint64) *SoPostWrap {
	if dba == nil || key == nil {
		return nil
	}
	result := &SoPostWrap{dba, key, nil, -1, nil, nil, nil}
	result.initWatcherFlag()
	return result
}

func (s *SoPostWrap) CheckExist() bool {
	if s.dba == nil {
		return false
	}
	if s.mKeyFlag != -1 {
		//if you have already obtained the existence status of the primary key, use it directly
		if s.mKeyFlag == 0 {
			return false
		}
		return true
	}
	keyBuf, err := s.encodeMainKey()
	if err != nil {
		return false
	}

	res, err := s.dba.Has(keyBuf)
	if err != nil {
		return false
	}
	if res == false {
		s.mKeyFlag = 0
	} else {
		s.mKeyFlag = 1
	}
	return res
}

func (s *SoPostWrap) MustExist(errMsgs ...interface{}) *SoPostWrap {
	if !s.CheckExist() {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.MustExist: %v not found", s.mainKey), errMsgs...))
	}
	return s
}

func (s *SoPostWrap) MustNotExist(errMsgs ...interface{}) *SoPostWrap {
	if s.CheckExist() {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.MustNotExist: %v already exists", s.mainKey), errMsgs...))
	}
	return s
}

func (s *SoPostWrap) initWatcherFlag() {
	if s.watcherFlag == nil {
		s.watcherFlag = new(PostWatcherFlag)
		*(s.watcherFlag) = PostWatcherFlagOfDb(s.dba.ServiceId())
	}
}

func (s *SoPostWrap) create(f func(tInfo *SoPost)) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	if s.mainKey == nil {
		return errors.New("the main key is nil")
	}
	val := &SoPost{}
	f(val)
	if s.CheckExist() {
		return errors.New("the main key is already exist")
	}
	keyBuf, err := s.encodeMainKey()
	if err != nil {
		return err

	}

	buf, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	err = s.dba.Put(keyBuf, buf)
	if err != nil {
		return err
	}

	// update srt list keys
	if err = s.insertAllSortKeys(val); err != nil {
		s.delAllSortKeys(false, val)
		s.dba.Delete(keyBuf)
		return err
	}

	//update unique list
	if sucNames, err := s.insertAllUniKeys(val); err != nil {
		s.delAllSortKeys(false, val)
		s.delUniKeysWithNames(sucNames, val)
		s.dba.Delete(keyBuf)
		return err
	}

	s.mKeyFlag = 1

	// call watchers
	s.initWatcherFlag()
	if s.watcherFlag.AnyWatcher {
		ReportTableRecordInsert(s.dba.ServiceId(), s.dba.BranchId(), s.mainKey, val)
	}

	return nil
}

func (s *SoPostWrap) Create(f func(tInfo *SoPost), errArgs ...interface{}) *SoPostWrap {
	err := s.create(f)
	if err != nil {
		panic(bindErrorInfo(fmt.Errorf("SoPostWrap.Create failed: %s", err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) getMainKeyBuf() ([]byte, error) {
	if s.mainKey == nil {
		return nil, errors.New("the main key is nil")
	}
	if s.mBuf == nil {
		var err error = nil
		s.mBuf, err = kope.Encode(s.mainKey)
		if err != nil {
			return nil, err
		}
	}
	return s.mBuf, nil
}

func (s *SoPostWrap) modify(f func(tInfo *SoPost)) error {
	if !s.CheckExist() {
		return errors.New("the SoPost table does not exist. Please create a table first")
	}
	oriTable := s.getPost()
	if oriTable == nil {
		return errors.New("fail to get origin table SoPost")
	}

	curTable := s.getPost()
	if curTable == nil {
		return errors.New("fail to create current table SoPost")
	}
	f(curTable)

	//the main key is not support modify
	if !reflect.DeepEqual(curTable.PostId, oriTable.PostId) {
		return errors.New("primary key does not support modification")
	}

	s.initWatcherFlag()
	modifiedFields, hasWatcher, err := s.getModifiedFields(oriTable, curTable)
	if err != nil {
		return err
	}

	if modifiedFields == nil || len(modifiedFields) < 1 {
		return nil
	}

	//check whether modify sort and unique field to nil
	err = s.checkSortAndUniFieldValidity(curTable, modifiedFields)
	if err != nil {
		return err
	}

	//check unique
	err = s.handleFieldMd(FieldMdHandleTypeCheck, curTable, modifiedFields)
	if err != nil {
		return err
	}

	//delete sort and unique key
	err = s.handleFieldMd(FieldMdHandleTypeDel, oriTable, modifiedFields)
	if err != nil {
		return err
	}

	//update table
	err = s.updatePost(curTable)
	if err != nil {
		return err
	}

	//insert sort and unique key
	err = s.handleFieldMd(FieldMdHandleTypeInsert, curTable, modifiedFields)
	if err != nil {
		return err
	}

	// call watchers
	if hasWatcher {
		ReportTableRecordUpdate(s.dba.ServiceId(), s.dba.BranchId(), s.mainKey, oriTable, curTable, modifiedFields)
	}

	return nil

}

func (s *SoPostWrap) Modify(f func(tInfo *SoPost), errArgs ...interface{}) *SoPostWrap {
	err := s.modify(f)
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.Modify failed: %s", err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetAuthor(p *prototype.AccountName, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Author = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetAuthor( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetBeneficiaries(p []*prototype.BeneficiaryRouteType, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Beneficiaries = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetBeneficiaries( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetBody(p string, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Body = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetBody( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetCashoutBlockNum(p uint64, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.CashoutBlockNum = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetCashoutBlockNum( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetCategory(p string, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Category = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetCategory( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetChildren(p uint32, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Children = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetChildren( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetCopyright(p uint32, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Copyright = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetCopyright( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetCopyrightMemo(p string, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.CopyrightMemo = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetCopyrightMemo( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetCreated(p *prototype.TimePointSec, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Created = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetCreated( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetDappRewards(p *prototype.Vest, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.DappRewards = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetDappRewards( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetDepth(p uint32, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Depth = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetDepth( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetLastPayout(p *prototype.TimePointSec, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.LastPayout = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetLastPayout( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetParentId(p uint64, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.ParentId = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetParentId( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetRewards(p *prototype.Vest, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Rewards = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetRewards( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetRootId(p uint64, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.RootId = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetRootId( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetTags(p []string, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Tags = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetTags( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetTicket(p uint32, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Ticket = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetTicket( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetTitle(p string, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.Title = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetTitle( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetVoteCnt(p uint64, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.VoteCnt = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetVoteCnt( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) SetWeightedVp(p string, errArgs ...interface{}) *SoPostWrap {
	err := s.modify(func(r *SoPost) {
		r.WeightedVp = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.SetWeightedVp( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoPostWrap) checkSortAndUniFieldValidity(curTable *SoPost, fields map[string]bool) error {
	if curTable != nil && fields != nil && len(fields) > 0 {

		if fields["Created"] && curTable.Created == nil {
			return errors.New("sort field Created can't be modified to nil")
		}

		if fields["Rewards"] && curTable.Rewards == nil {
			return errors.New("sort field Rewards can't be modified to nil")
		}

	}
	return nil
}

//Get all the modified fields in the table
func (s *SoPostWrap) getModifiedFields(oriTable *SoPost, curTable *SoPost) (map[string]bool, bool, error) {
	if oriTable == nil {
		return nil, false, errors.New("table info is nil, can't get modified fields")
	}
	hasWatcher := false
	fields := make(map[string]bool)

	if !reflect.DeepEqual(oriTable.Author, curTable.Author) {
		fields["Author"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasAuthorWatcher
	}

	if !reflect.DeepEqual(oriTable.Beneficiaries, curTable.Beneficiaries) {
		fields["Beneficiaries"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasBeneficiariesWatcher
	}

	if !reflect.DeepEqual(oriTable.Body, curTable.Body) {
		fields["Body"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasBodyWatcher
	}

	if !reflect.DeepEqual(oriTable.CashoutBlockNum, curTable.CashoutBlockNum) {
		fields["CashoutBlockNum"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasCashoutBlockNumWatcher
	}

	if !reflect.DeepEqual(oriTable.Category, curTable.Category) {
		fields["Category"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasCategoryWatcher
	}

	if !reflect.DeepEqual(oriTable.Children, curTable.Children) {
		fields["Children"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasChildrenWatcher
	}

	if !reflect.DeepEqual(oriTable.Copyright, curTable.Copyright) {
		fields["Copyright"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasCopyrightWatcher
	}

	if !reflect.DeepEqual(oriTable.CopyrightMemo, curTable.CopyrightMemo) {
		fields["CopyrightMemo"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasCopyrightMemoWatcher
	}

	if !reflect.DeepEqual(oriTable.Created, curTable.Created) {
		fields["Created"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasCreatedWatcher
	}

	if !reflect.DeepEqual(oriTable.DappRewards, curTable.DappRewards) {
		fields["DappRewards"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasDappRewardsWatcher
	}

	if !reflect.DeepEqual(oriTable.Depth, curTable.Depth) {
		fields["Depth"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasDepthWatcher
	}

	if !reflect.DeepEqual(oriTable.LastPayout, curTable.LastPayout) {
		fields["LastPayout"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasLastPayoutWatcher
	}

	if !reflect.DeepEqual(oriTable.ParentId, curTable.ParentId) {
		fields["ParentId"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasParentIdWatcher
	}

	if !reflect.DeepEqual(oriTable.Rewards, curTable.Rewards) {
		fields["Rewards"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasRewardsWatcher
	}

	if !reflect.DeepEqual(oriTable.RootId, curTable.RootId) {
		fields["RootId"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasRootIdWatcher
	}

	if !reflect.DeepEqual(oriTable.Tags, curTable.Tags) {
		fields["Tags"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasTagsWatcher
	}

	if !reflect.DeepEqual(oriTable.Ticket, curTable.Ticket) {
		fields["Ticket"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasTicketWatcher
	}

	if !reflect.DeepEqual(oriTable.Title, curTable.Title) {
		fields["Title"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasTitleWatcher
	}

	if !reflect.DeepEqual(oriTable.VoteCnt, curTable.VoteCnt) {
		fields["VoteCnt"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasVoteCntWatcher
	}

	if !reflect.DeepEqual(oriTable.WeightedVp, curTable.WeightedVp) {
		fields["WeightedVp"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasWeightedVpWatcher
	}

	hasWatcher = hasWatcher || s.watcherFlag.WholeWatcher
	return fields, hasWatcher, nil
}

func (s *SoPostWrap) handleFieldMd(t FieldMdHandleType, so *SoPost, fields map[string]bool) error {
	if so == nil {
		return errors.New("fail to modify empty table")
	}

	//there is no field need to modify
	if fields == nil || len(fields) < 1 {
		return nil
	}

	errStr := ""

	if fields["Author"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldAuthor(so.Author, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Author")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldAuthor(so.Author, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Author")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldAuthor(so.Author, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Author")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Beneficiaries"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldBeneficiaries(so.Beneficiaries, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Beneficiaries")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldBeneficiaries(so.Beneficiaries, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Beneficiaries")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldBeneficiaries(so.Beneficiaries, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Beneficiaries")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Body"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldBody(so.Body, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Body")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldBody(so.Body, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Body")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldBody(so.Body, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Body")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["CashoutBlockNum"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldCashoutBlockNum(so.CashoutBlockNum, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "CashoutBlockNum")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldCashoutBlockNum(so.CashoutBlockNum, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "CashoutBlockNum")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldCashoutBlockNum(so.CashoutBlockNum, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "CashoutBlockNum")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Category"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldCategory(so.Category, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Category")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldCategory(so.Category, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Category")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldCategory(so.Category, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Category")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Children"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldChildren(so.Children, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Children")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldChildren(so.Children, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Children")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldChildren(so.Children, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Children")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Copyright"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldCopyright(so.Copyright, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Copyright")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldCopyright(so.Copyright, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Copyright")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldCopyright(so.Copyright, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Copyright")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["CopyrightMemo"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldCopyrightMemo(so.CopyrightMemo, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "CopyrightMemo")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldCopyrightMemo(so.CopyrightMemo, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "CopyrightMemo")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldCopyrightMemo(so.CopyrightMemo, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "CopyrightMemo")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Created"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldCreated(so.Created, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Created")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldCreated(so.Created, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Created")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldCreated(so.Created, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Created")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["DappRewards"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldDappRewards(so.DappRewards, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "DappRewards")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldDappRewards(so.DappRewards, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "DappRewards")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldDappRewards(so.DappRewards, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "DappRewards")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Depth"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldDepth(so.Depth, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Depth")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldDepth(so.Depth, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Depth")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldDepth(so.Depth, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Depth")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["LastPayout"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldLastPayout(so.LastPayout, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "LastPayout")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldLastPayout(so.LastPayout, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "LastPayout")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldLastPayout(so.LastPayout, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "LastPayout")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["ParentId"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldParentId(so.ParentId, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "ParentId")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldParentId(so.ParentId, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "ParentId")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldParentId(so.ParentId, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "ParentId")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Rewards"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldRewards(so.Rewards, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Rewards")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldRewards(so.Rewards, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Rewards")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldRewards(so.Rewards, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Rewards")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["RootId"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldRootId(so.RootId, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "RootId")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldRootId(so.RootId, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "RootId")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldRootId(so.RootId, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "RootId")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Tags"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldTags(so.Tags, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Tags")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldTags(so.Tags, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Tags")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldTags(so.Tags, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Tags")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Ticket"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldTicket(so.Ticket, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Ticket")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldTicket(so.Ticket, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Ticket")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldTicket(so.Ticket, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Ticket")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["Title"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldTitle(so.Title, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Title")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldTitle(so.Title, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Title")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldTitle(so.Title, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Title")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["VoteCnt"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldVoteCnt(so.VoteCnt, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "VoteCnt")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldVoteCnt(so.VoteCnt, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "VoteCnt")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldVoteCnt(so.VoteCnt, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "VoteCnt")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["WeightedVp"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldWeightedVp(so.WeightedVp, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "WeightedVp")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldWeightedVp(so.WeightedVp, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "WeightedVp")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldWeightedVp(so.WeightedVp, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "WeightedVp")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	return nil
}

////////////// SECTION LKeys delete/insert ///////////////

func (s *SoPostWrap) delSortKeyCreated(sa *SoPost) bool {
	if s.dba == nil || s.mainKey == nil {
		return false
	}
	val := SoListPostByCreated{}
	if sa == nil {
		val.Created = s.GetCreated()
		val.PostId = *s.mainKey
	} else {
		val.Created = sa.Created
		val.PostId = sa.PostId
	}
	subBuf, err := val.OpeEncode()
	if err != nil {
		return false
	}
	ordErr := s.dba.Delete(subBuf)
	return ordErr == nil
}

func (s *SoPostWrap) insertSortKeyCreated(sa *SoPost) bool {
	if s.dba == nil || sa == nil {
		return false
	}
	val := SoListPostByCreated{}
	val.PostId = sa.PostId
	val.Created = sa.Created
	buf, err := proto.Marshal(&val)
	if err != nil {
		return false
	}
	subBuf, err := val.OpeEncode()
	if err != nil {
		return false
	}
	ordErr := s.dba.Put(subBuf, buf)
	return ordErr == nil
}

func (s *SoPostWrap) delSortKeyCashoutBlockNum(sa *SoPost) bool {
	if s.dba == nil || s.mainKey == nil {
		return false
	}
	val := SoListPostByCashoutBlockNum{}
	if sa == nil {
		val.CashoutBlockNum = s.GetCashoutBlockNum()
		val.PostId = *s.mainKey
	} else {
		val.CashoutBlockNum = sa.CashoutBlockNum
		val.PostId = sa.PostId
	}
	subBuf, err := val.OpeEncode()
	if err != nil {
		return false
	}
	ordErr := s.dba.Delete(subBuf)
	return ordErr == nil
}

func (s *SoPostWrap) insertSortKeyCashoutBlockNum(sa *SoPost) bool {
	if s.dba == nil || sa == nil {
		return false
	}
	val := SoListPostByCashoutBlockNum{}
	val.PostId = sa.PostId
	val.CashoutBlockNum = sa.CashoutBlockNum
	buf, err := proto.Marshal(&val)
	if err != nil {
		return false
	}
	subBuf, err := val.OpeEncode()
	if err != nil {
		return false
	}
	ordErr := s.dba.Put(subBuf, buf)
	return ordErr == nil
}

func (s *SoPostWrap) delSortKeyRewards(sa *SoPost) bool {
	if s.dba == nil || s.mainKey == nil {
		return false
	}
	val := SoListPostByRewards{}
	if sa == nil {
		val.Rewards = s.GetRewards()
		val.PostId = *s.mainKey
	} else {
		val.Rewards = sa.Rewards
		val.PostId = sa.PostId
	}
	subBuf, err := val.OpeEncode()
	if err != nil {
		return false
	}
	ordErr := s.dba.Delete(subBuf)
	return ordErr == nil
}

func (s *SoPostWrap) insertSortKeyRewards(sa *SoPost) bool {
	if s.dba == nil || sa == nil {
		return false
	}
	val := SoListPostByRewards{}
	val.PostId = sa.PostId
	val.Rewards = sa.Rewards
	buf, err := proto.Marshal(&val)
	if err != nil {
		return false
	}
	subBuf, err := val.OpeEncode()
	if err != nil {
		return false
	}
	ordErr := s.dba.Put(subBuf, buf)
	return ordErr == nil
}

func (s *SoPostWrap) delAllSortKeys(br bool, val *SoPost) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if !s.delSortKeyCreated(val) {
		if br {
			return false
		} else {
			res = false
		}
	}
	if !s.delSortKeyCashoutBlockNum(val) {
		if br {
			return false
		} else {
			res = false
		}
	}
	if !s.delSortKeyRewards(val) {
		if br {
			return false
		} else {
			res = false
		}
	}

	return res
}

func (s *SoPostWrap) insertAllSortKeys(val *SoPost) error {
	if s.dba == nil {
		return errors.New("insert sort Field fail,the db is nil ")
	}
	if val == nil {
		return errors.New("insert sort Field fail,get the SoPost fail ")
	}
	if !s.insertSortKeyCreated(val) {
		return errors.New("insert sort Field Created fail while insert table ")
	}
	if !s.insertSortKeyCashoutBlockNum(val) {
		return errors.New("insert sort Field CashoutBlockNum fail while insert table ")
	}
	if !s.insertSortKeyRewards(val) {
		return errors.New("insert sort Field Rewards fail while insert table ")
	}

	return nil
}

////////////// SECTION LKeys delete/insert //////////////

func (s *SoPostWrap) removePost() error {
	if s.dba == nil {
		return errors.New("database is nil")
	}

	s.initWatcherFlag()

	var oldVal *SoPost
	if s.watcherFlag.AnyWatcher {
		oldVal = s.getPost()
	}

	//delete sort list key
	if res := s.delAllSortKeys(true, nil); !res {
		return errors.New("delAllSortKeys failed")
	}

	//delete unique list
	if res := s.delAllUniKeys(true, nil); !res {
		return errors.New("delAllUniKeys failed")
	}

	//delete table
	key, err := s.encodeMainKey()
	if err != nil {
		return fmt.Errorf("encodeMainKey failed: %s", err.Error())
	}
	err = s.dba.Delete(key)
	if err == nil {
		s.mKeyBuf = nil
		s.mKeyFlag = -1

		// call watchers
		if s.watcherFlag.AnyWatcher && oldVal != nil {
			ReportTableRecordDelete(s.dba.ServiceId(), s.dba.BranchId(), s.mainKey, oldVal)
		}
		return nil
	} else {
		return fmt.Errorf("database.Delete failed: %s", err.Error())
	}
}

func (s *SoPostWrap) RemovePost(errMsgs ...interface{}) *SoPostWrap {
	err := s.removePost()
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoPostWrap.RemovePost failed: %s", err.Error()), errMsgs...))
	}
	return s
}

////////////// SECTION Members Get/Modify ///////////////

func (s *SoPostWrap) GetAuthor() *prototype.AccountName {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Author
			}
		}
	}
	if !res {
		return nil

	}
	return msg.Author
}

func (s *SoPostWrap) mdFieldAuthor(p *prototype.AccountName, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkAuthorIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldAuthor(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldAuthor(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldAuthor(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldAuthor(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkAuthorIsMetMdCondition(p *prototype.AccountName) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetBeneficiaries() []*prototype.BeneficiaryRouteType {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Beneficiaries
			}
		}
	}
	if !res {
		var tmpValue []*prototype.BeneficiaryRouteType
		return tmpValue
	}
	return msg.Beneficiaries
}

func (s *SoPostWrap) mdFieldBeneficiaries(p []*prototype.BeneficiaryRouteType, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkBeneficiariesIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldBeneficiaries(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldBeneficiaries(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldBeneficiaries(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldBeneficiaries(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkBeneficiariesIsMetMdCondition(p []*prototype.BeneficiaryRouteType) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetBody() string {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Body
			}
		}
	}
	if !res {
		var tmpValue string
		return tmpValue
	}
	return msg.Body
}

func (s *SoPostWrap) mdFieldBody(p string, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkBodyIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldBody(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldBody(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldBody(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldBody(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkBodyIsMetMdCondition(p string) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetCashoutBlockNum() uint64 {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.CashoutBlockNum
			}
		}
	}
	if !res {
		var tmpValue uint64
		return tmpValue
	}
	return msg.CashoutBlockNum
}

func (s *SoPostWrap) mdFieldCashoutBlockNum(p uint64, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkCashoutBlockNumIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldCashoutBlockNum(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldCashoutBlockNum(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldCashoutBlockNum(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if !s.delSortKeyCashoutBlockNum(so) {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldCashoutBlockNum(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if !s.insertSortKeyCashoutBlockNum(so) {
		return false
	}

	return true
}

func (s *SoPostWrap) checkCashoutBlockNumIsMetMdCondition(p uint64) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetCategory() string {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Category
			}
		}
	}
	if !res {
		var tmpValue string
		return tmpValue
	}
	return msg.Category
}

func (s *SoPostWrap) mdFieldCategory(p string, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkCategoryIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldCategory(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldCategory(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldCategory(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldCategory(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkCategoryIsMetMdCondition(p string) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetChildren() uint32 {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Children
			}
		}
	}
	if !res {
		var tmpValue uint32
		return tmpValue
	}
	return msg.Children
}

func (s *SoPostWrap) mdFieldChildren(p uint32, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkChildrenIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldChildren(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldChildren(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldChildren(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldChildren(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkChildrenIsMetMdCondition(p uint32) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetCopyright() uint32 {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Copyright
			}
		}
	}
	if !res {
		var tmpValue uint32
		return tmpValue
	}
	return msg.Copyright
}

func (s *SoPostWrap) mdFieldCopyright(p uint32, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkCopyrightIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldCopyright(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldCopyright(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldCopyright(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldCopyright(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkCopyrightIsMetMdCondition(p uint32) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetCopyrightMemo() string {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.CopyrightMemo
			}
		}
	}
	if !res {
		var tmpValue string
		return tmpValue
	}
	return msg.CopyrightMemo
}

func (s *SoPostWrap) mdFieldCopyrightMemo(p string, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkCopyrightMemoIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldCopyrightMemo(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldCopyrightMemo(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldCopyrightMemo(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldCopyrightMemo(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkCopyrightMemoIsMetMdCondition(p string) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetCreated() *prototype.TimePointSec {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Created
			}
		}
	}
	if !res {
		return nil

	}
	return msg.Created
}

func (s *SoPostWrap) mdFieldCreated(p *prototype.TimePointSec, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkCreatedIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldCreated(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldCreated(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldCreated(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if !s.delSortKeyCreated(so) {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldCreated(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if !s.insertSortKeyCreated(so) {
		return false
	}

	return true
}

func (s *SoPostWrap) checkCreatedIsMetMdCondition(p *prototype.TimePointSec) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetDappRewards() *prototype.Vest {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.DappRewards
			}
		}
	}
	if !res {
		return nil

	}
	return msg.DappRewards
}

func (s *SoPostWrap) mdFieldDappRewards(p *prototype.Vest, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkDappRewardsIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldDappRewards(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldDappRewards(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldDappRewards(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldDappRewards(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkDappRewardsIsMetMdCondition(p *prototype.Vest) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetDepth() uint32 {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Depth
			}
		}
	}
	if !res {
		var tmpValue uint32
		return tmpValue
	}
	return msg.Depth
}

func (s *SoPostWrap) mdFieldDepth(p uint32, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkDepthIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldDepth(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldDepth(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldDepth(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldDepth(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkDepthIsMetMdCondition(p uint32) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetLastPayout() *prototype.TimePointSec {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.LastPayout
			}
		}
	}
	if !res {
		return nil

	}
	return msg.LastPayout
}

func (s *SoPostWrap) mdFieldLastPayout(p *prototype.TimePointSec, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkLastPayoutIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldLastPayout(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldLastPayout(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldLastPayout(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldLastPayout(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkLastPayoutIsMetMdCondition(p *prototype.TimePointSec) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetParentId() uint64 {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.ParentId
			}
		}
	}
	if !res {
		var tmpValue uint64
		return tmpValue
	}
	return msg.ParentId
}

func (s *SoPostWrap) mdFieldParentId(p uint64, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkParentIdIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldParentId(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldParentId(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldParentId(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldParentId(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkParentIdIsMetMdCondition(p uint64) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetPostId() uint64 {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.PostId
			}
		}
	}
	if !res {
		var tmpValue uint64
		return tmpValue
	}
	return msg.PostId
}

func (s *SoPostWrap) GetRewards() *prototype.Vest {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Rewards
			}
		}
	}
	if !res {
		return nil

	}
	return msg.Rewards
}

func (s *SoPostWrap) mdFieldRewards(p *prototype.Vest, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkRewardsIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldRewards(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldRewards(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldRewards(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if !s.delSortKeyRewards(so) {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldRewards(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if !s.insertSortKeyRewards(so) {
		return false
	}

	return true
}

func (s *SoPostWrap) checkRewardsIsMetMdCondition(p *prototype.Vest) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetRootId() uint64 {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.RootId
			}
		}
	}
	if !res {
		var tmpValue uint64
		return tmpValue
	}
	return msg.RootId
}

func (s *SoPostWrap) mdFieldRootId(p uint64, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkRootIdIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldRootId(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldRootId(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldRootId(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldRootId(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkRootIdIsMetMdCondition(p uint64) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetTags() []string {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Tags
			}
		}
	}
	if !res {
		var tmpValue []string
		return tmpValue
	}
	return msg.Tags
}

func (s *SoPostWrap) mdFieldTags(p []string, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkTagsIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldTags(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldTags(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldTags(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldTags(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkTagsIsMetMdCondition(p []string) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetTicket() uint32 {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Ticket
			}
		}
	}
	if !res {
		var tmpValue uint32
		return tmpValue
	}
	return msg.Ticket
}

func (s *SoPostWrap) mdFieldTicket(p uint32, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkTicketIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldTicket(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldTicket(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldTicket(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldTicket(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkTicketIsMetMdCondition(p uint32) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetTitle() string {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.Title
			}
		}
	}
	if !res {
		var tmpValue string
		return tmpValue
	}
	return msg.Title
}

func (s *SoPostWrap) mdFieldTitle(p string, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkTitleIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldTitle(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldTitle(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldTitle(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldTitle(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkTitleIsMetMdCondition(p string) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetVoteCnt() uint64 {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.VoteCnt
			}
		}
	}
	if !res {
		var tmpValue uint64
		return tmpValue
	}
	return msg.VoteCnt
}

func (s *SoPostWrap) mdFieldVoteCnt(p uint64, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkVoteCntIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldVoteCnt(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldVoteCnt(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldVoteCnt(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldVoteCnt(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkVoteCntIsMetMdCondition(p uint64) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) GetWeightedVp() string {
	res := true
	msg := &SoPost{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMainKey()
		if err != nil {
			res = false
		} else {
			buf, err := s.dba.Get(key)
			if err != nil {
				res = false
			}
			err = proto.Unmarshal(buf, msg)
			if err != nil {
				res = false
			} else {
				return msg.WeightedVp
			}
		}
	}
	if !res {
		var tmpValue string
		return tmpValue
	}
	return msg.WeightedVp
}

func (s *SoPostWrap) mdFieldWeightedVp(p string, isCheck bool, isDel bool, isInsert bool,
	so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkWeightedVpIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldWeightedVp(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldWeightedVp(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoPostWrap) delFieldWeightedVp(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) insertFieldWeightedVp(so *SoPost) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoPostWrap) checkWeightedVpIsMetMdCondition(p string) bool {
	if s.dba == nil {
		return false
	}

	return true
}

////////////// SECTION List Keys ///////////////
type SPostCreatedWrap struct {
	Dba iservices.IDatabaseRW
}

func NewPostCreatedWrap(db iservices.IDatabaseRW) *SPostCreatedWrap {
	if db == nil {
		return nil
	}
	wrap := SPostCreatedWrap{Dba: db}
	return &wrap
}

func (s *SPostCreatedWrap) GetMainVal(val []byte) *uint64 {
	res := &SoListPostByCreated{}
	err := proto.Unmarshal(val, res)

	if err != nil {
		return nil
	}

	return &res.PostId

}

func (s *SPostCreatedWrap) GetSubVal(val []byte) *prototype.TimePointSec {
	res := &SoListPostByCreated{}
	err := proto.Unmarshal(val, res)
	if err != nil {
		return nil
	}
	return res.Created

}

func (m *SoListPostByCreated) OpeEncode() ([]byte, error) {
	pre := PostCreatedTable
	sub := m.Created
	if sub == nil {
		return nil, errors.New("the pro Created is nil")
	}
	sub1 := m.PostId

	kList := []interface{}{pre, sub, sub1}
	kBuf, cErr := kope.EncodeSlice(kList)
	return kBuf, cErr
}

//Query srt by order
//
//start = nil  end = nil (query the db from start to end)
//start = nil (query from start the db)
//end = nil (query to the end of db)
//
//f: callback for each traversal , primary 、sub key、idx(the number of times it has been iterated)
//as arguments to the callback function
//if the return value of f is true,continue iterating until the end iteration;
//otherwise stop iteration immediately
//
//lastMainKey: the main key of the last one of last page
//lastSubVal: the value  of the last one of last page
//
func (s *SPostCreatedWrap) ForEachByOrder(start *prototype.TimePointSec, end *prototype.TimePointSec, lastMainKey *uint64,
	lastSubVal *prototype.TimePointSec, f func(mVal *uint64, sVal *prototype.TimePointSec, idx uint32) bool) error {
	if s.Dba == nil {
		return errors.New("the db is nil")
	}
	if (lastSubVal != nil && lastMainKey == nil) || (lastSubVal == nil && lastMainKey != nil) {
		return errors.New("last query param error")
	}
	if f == nil {
		return nil
	}
	pre := PostCreatedTable
	skeyList := []interface{}{pre}
	if start != nil {
		skeyList = append(skeyList, start)
		if lastMainKey != nil {
			skeyList = append(skeyList, lastMainKey, kope.MinimalKey)
		}
	} else {
		if lastMainKey != nil && lastSubVal != nil {
			skeyList = append(skeyList, lastSubVal, lastMainKey, kope.MinimalKey)
		}
		skeyList = append(skeyList, kope.MinimalKey)
	}
	sBuf, cErr := kope.EncodeSlice(skeyList)
	if cErr != nil {
		return cErr
	}
	eKeyList := []interface{}{pre}
	if end != nil {
		eKeyList = append(eKeyList, end)
	} else {
		eKeyList = append(eKeyList, kope.MaximumKey)
	}
	eBuf, cErr := kope.EncodeSlice(eKeyList)
	if cErr != nil {
		return cErr
	}
	var idx uint32 = 0
	s.Dba.Iterate(sBuf, eBuf, false, func(key, value []byte) bool {
		idx++
		return f(s.GetMainVal(value), s.GetSubVal(value), idx)
	})
	return nil
}

//Query srt by reverse order
//
//f: callback for each traversal , primary 、sub key、idx(the number of times it has been iterated)
//as arguments to the callback function
//if the return value of f is true,continue iterating until the end iteration;
//otherwise stop iteration immediately
//
//lastMainKey: the main key of the last one of last page
//lastSubVal: the value  of the last one of last page
//
func (s *SPostCreatedWrap) ForEachByRevOrder(start *prototype.TimePointSec, end *prototype.TimePointSec, lastMainKey *uint64,
	lastSubVal *prototype.TimePointSec, f func(mVal *uint64, sVal *prototype.TimePointSec, idx uint32) bool) error {
	if s.Dba == nil {
		return errors.New("the db is nil")
	}
	if (lastSubVal != nil && lastMainKey == nil) || (lastSubVal == nil && lastMainKey != nil) {
		return errors.New("last query param error")
	}
	if f == nil {
		return nil
	}
	pre := PostCreatedTable
	skeyList := []interface{}{pre}
	if start != nil {
		skeyList = append(skeyList, start)
		if lastMainKey != nil {
			skeyList = append(skeyList, lastMainKey)
		}
	} else {
		if lastMainKey != nil && lastSubVal != nil {
			skeyList = append(skeyList, lastSubVal, lastMainKey)
		}
		skeyList = append(skeyList, kope.MaximumKey)
	}
	sBuf, cErr := kope.EncodeSlice(skeyList)
	if cErr != nil {
		return cErr
	}
	eKeyList := []interface{}{pre}
	if end != nil {
		eKeyList = append(eKeyList, end)
	}
	eBuf, cErr := kope.EncodeSlice(eKeyList)
	if cErr != nil {
		return cErr
	}
	var idx uint32 = 0
	s.Dba.Iterate(eBuf, sBuf, true, func(key, value []byte) bool {
		idx++
		return f(s.GetMainVal(value), s.GetSubVal(value), idx)
	})
	return nil
}

////////////// SECTION List Keys ///////////////
type SPostCashoutBlockNumWrap struct {
	Dba iservices.IDatabaseRW
}

func NewPostCashoutBlockNumWrap(db iservices.IDatabaseRW) *SPostCashoutBlockNumWrap {
	if db == nil {
		return nil
	}
	wrap := SPostCashoutBlockNumWrap{Dba: db}
	return &wrap
}

func (s *SPostCashoutBlockNumWrap) GetMainVal(val []byte) *uint64 {
	res := &SoListPostByCashoutBlockNum{}
	err := proto.Unmarshal(val, res)

	if err != nil {
		return nil
	}

	return &res.PostId

}

func (s *SPostCashoutBlockNumWrap) GetSubVal(val []byte) *uint64 {
	res := &SoListPostByCashoutBlockNum{}
	err := proto.Unmarshal(val, res)
	if err != nil {
		return nil
	}
	return &res.CashoutBlockNum

}

func (m *SoListPostByCashoutBlockNum) OpeEncode() ([]byte, error) {
	pre := PostCashoutBlockNumTable
	sub := m.CashoutBlockNum

	sub1 := m.PostId

	kList := []interface{}{pre, sub, sub1}
	kBuf, cErr := kope.EncodeSlice(kList)
	return kBuf, cErr
}

//Query srt by order
//
//start = nil  end = nil (query the db from start to end)
//start = nil (query from start the db)
//end = nil (query to the end of db)
//
//f: callback for each traversal , primary 、sub key、idx(the number of times it has been iterated)
//as arguments to the callback function
//if the return value of f is true,continue iterating until the end iteration;
//otherwise stop iteration immediately
//
//lastMainKey: the main key of the last one of last page
//lastSubVal: the value  of the last one of last page
//
func (s *SPostCashoutBlockNumWrap) ForEachByOrder(start *uint64, end *uint64, lastMainKey *uint64,
	lastSubVal *uint64, f func(mVal *uint64, sVal *uint64, idx uint32) bool) error {
	if s.Dba == nil {
		return errors.New("the db is nil")
	}
	if (lastSubVal != nil && lastMainKey == nil) || (lastSubVal == nil && lastMainKey != nil) {
		return errors.New("last query param error")
	}
	if f == nil {
		return nil
	}
	pre := PostCashoutBlockNumTable
	skeyList := []interface{}{pre}
	if start != nil {
		skeyList = append(skeyList, start)
		if lastMainKey != nil {
			skeyList = append(skeyList, lastMainKey, kope.MinimalKey)
		}
	} else {
		if lastMainKey != nil && lastSubVal != nil {
			skeyList = append(skeyList, lastSubVal, lastMainKey, kope.MinimalKey)
		}
		skeyList = append(skeyList, kope.MinimalKey)
	}
	sBuf, cErr := kope.EncodeSlice(skeyList)
	if cErr != nil {
		return cErr
	}
	eKeyList := []interface{}{pre}
	if end != nil {
		eKeyList = append(eKeyList, end)
	} else {
		eKeyList = append(eKeyList, kope.MaximumKey)
	}
	eBuf, cErr := kope.EncodeSlice(eKeyList)
	if cErr != nil {
		return cErr
	}
	var idx uint32 = 0
	s.Dba.Iterate(sBuf, eBuf, false, func(key, value []byte) bool {
		idx++
		return f(s.GetMainVal(value), s.GetSubVal(value), idx)
	})
	return nil
}

////////////// SECTION List Keys ///////////////
type SPostRewardsWrap struct {
	Dba iservices.IDatabaseRW
}

func NewPostRewardsWrap(db iservices.IDatabaseRW) *SPostRewardsWrap {
	if db == nil {
		return nil
	}
	wrap := SPostRewardsWrap{Dba: db}
	return &wrap
}

func (s *SPostRewardsWrap) GetMainVal(val []byte) *uint64 {
	res := &SoListPostByRewards{}
	err := proto.Unmarshal(val, res)

	if err != nil {
		return nil
	}

	return &res.PostId

}

func (s *SPostRewardsWrap) GetSubVal(val []byte) *prototype.Vest {
	res := &SoListPostByRewards{}
	err := proto.Unmarshal(val, res)
	if err != nil {
		return nil
	}
	return res.Rewards

}

func (m *SoListPostByRewards) OpeEncode() ([]byte, error) {
	pre := PostRewardsTable
	sub := m.Rewards
	if sub == nil {
		return nil, errors.New("the pro Rewards is nil")
	}
	sub1 := m.PostId

	kList := []interface{}{pre, sub, sub1}
	kBuf, cErr := kope.EncodeSlice(kList)
	return kBuf, cErr
}

//Query srt by reverse order
//
//f: callback for each traversal , primary 、sub key、idx(the number of times it has been iterated)
//as arguments to the callback function
//if the return value of f is true,continue iterating until the end iteration;
//otherwise stop iteration immediately
//
//lastMainKey: the main key of the last one of last page
//lastSubVal: the value  of the last one of last page
//
func (s *SPostRewardsWrap) ForEachByRevOrder(start *prototype.Vest, end *prototype.Vest, lastMainKey *uint64,
	lastSubVal *prototype.Vest, f func(mVal *uint64, sVal *prototype.Vest, idx uint32) bool) error {
	if s.Dba == nil {
		return errors.New("the db is nil")
	}
	if (lastSubVal != nil && lastMainKey == nil) || (lastSubVal == nil && lastMainKey != nil) {
		return errors.New("last query param error")
	}
	if f == nil {
		return nil
	}
	pre := PostRewardsTable
	skeyList := []interface{}{pre}
	if start != nil {
		skeyList = append(skeyList, start)
		if lastMainKey != nil {
			skeyList = append(skeyList, lastMainKey)
		}
	} else {
		if lastMainKey != nil && lastSubVal != nil {
			skeyList = append(skeyList, lastSubVal, lastMainKey)
		}
		skeyList = append(skeyList, kope.MaximumKey)
	}
	sBuf, cErr := kope.EncodeSlice(skeyList)
	if cErr != nil {
		return cErr
	}
	eKeyList := []interface{}{pre}
	if end != nil {
		eKeyList = append(eKeyList, end)
	}
	eBuf, cErr := kope.EncodeSlice(eKeyList)
	if cErr != nil {
		return cErr
	}
	var idx uint32 = 0
	s.Dba.Iterate(eBuf, sBuf, true, func(key, value []byte) bool {
		idx++
		return f(s.GetMainVal(value), s.GetSubVal(value), idx)
	})
	return nil
}

/////////////// SECTION Private function ////////////////

func (s *SoPostWrap) update(sa *SoPost) bool {
	if s.dba == nil || sa == nil {
		return false
	}
	buf, err := proto.Marshal(sa)
	if err != nil {
		return false
	}

	keyBuf, err := s.encodeMainKey()
	if err != nil {
		return false
	}

	return s.dba.Put(keyBuf, buf) == nil
}

func (s *SoPostWrap) getPost() *SoPost {
	if s.dba == nil {
		return nil
	}
	keyBuf, err := s.encodeMainKey()
	if err != nil {
		return nil
	}
	resBuf, err := s.dba.Get(keyBuf)

	if err != nil {
		return nil
	}

	res := &SoPost{}
	if proto.Unmarshal(resBuf, res) != nil {
		return nil
	}
	return res
}

func (s *SoPostWrap) updatePost(so *SoPost) error {
	if s.dba == nil {
		return errors.New("update fail:the db is nil")
	}

	if so == nil {
		return errors.New("update fail: the SoPost is nil")
	}

	key, err := s.encodeMainKey()
	if err != nil {
		return nil
	}

	buf, err := proto.Marshal(so)
	if err != nil {
		return err
	}

	err = s.dba.Put(key, buf)
	if err != nil {
		return err
	}

	return nil
}

func (s *SoPostWrap) encodeMainKey() ([]byte, error) {
	if s.mKeyBuf != nil {
		return s.mKeyBuf, nil
	}
	pre := PostPostIdRow
	sub := s.mainKey
	if sub == nil {
		return nil, errors.New("the mainKey is nil")
	}
	preBuf, err := kope.Encode(pre)
	if err != nil {
		return nil, err
	}
	mBuf, err := s.getMainKeyBuf()
	if err != nil {
		return nil, err
	}
	list := make([][]byte, 2)
	list[0] = preBuf
	list[1] = mBuf
	s.mKeyBuf = kope.PackList(list)
	return s.mKeyBuf, nil
}

////////////// Unique Query delete/insert/query ///////////////

func (s *SoPostWrap) delAllUniKeys(br bool, val *SoPost) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if !s.delUniKeyPostId(val) {
		if br {
			return false
		} else {
			res = false
		}
	}

	return res
}

func (s *SoPostWrap) delUniKeysWithNames(names map[string]string, val *SoPost) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if len(names["PostId"]) > 0 {
		if !s.delUniKeyPostId(val) {
			res = false
		}
	}

	return res
}

func (s *SoPostWrap) insertAllUniKeys(val *SoPost) (map[string]string, error) {
	if s.dba == nil {
		return nil, errors.New("insert uniuqe Field fail,the db is nil ")
	}
	if val == nil {
		return nil, errors.New("insert uniuqe Field fail,get the SoPost fail ")
	}
	sucFields := map[string]string{}
	if !s.insertUniKeyPostId(val) {
		return sucFields, errors.New("insert unique Field PostId fail while insert table ")
	}
	sucFields["PostId"] = "PostId"

	return sucFields, nil
}

func (s *SoPostWrap) delUniKeyPostId(sa *SoPost) bool {
	if s.dba == nil {
		return false
	}
	pre := PostPostIdUniTable
	kList := []interface{}{pre}
	if sa != nil {

		sub := sa.PostId
		kList = append(kList, sub)
	} else {
		sub := s.GetPostId()

		kList = append(kList, sub)

	}
	kBuf, err := kope.EncodeSlice(kList)
	if err != nil {
		return false
	}
	return s.dba.Delete(kBuf) == nil
}

func (s *SoPostWrap) insertUniKeyPostId(sa *SoPost) bool {
	if s.dba == nil || sa == nil {
		return false
	}

	pre := PostPostIdUniTable
	sub := sa.PostId
	kList := []interface{}{pre, sub}
	kBuf, err := kope.EncodeSlice(kList)
	if err != nil {
		return false
	}
	res, err := s.dba.Has(kBuf)
	if err == nil && res == true {
		//the unique key is already exist
		return false
	}
	val := SoUniquePostByPostId{}
	val.PostId = sa.PostId

	buf, err := proto.Marshal(&val)

	if err != nil {
		return false
	}

	return s.dba.Put(kBuf, buf) == nil

}

type UniPostPostIdWrap struct {
	Dba iservices.IDatabaseRW
}

func NewUniPostPostIdWrap(db iservices.IDatabaseRW) *UniPostPostIdWrap {
	if db == nil {
		return nil
	}
	wrap := UniPostPostIdWrap{Dba: db}
	return &wrap
}

func (s *UniPostPostIdWrap) UniQueryPostId(start *uint64) *SoPostWrap {
	if start == nil || s.Dba == nil {
		return nil
	}
	pre := PostPostIdUniTable
	kList := []interface{}{pre, start}
	bufStartkey, err := kope.EncodeSlice(kList)
	val, err := s.Dba.Get(bufStartkey)
	if err == nil {
		res := &SoUniquePostByPostId{}
		rErr := proto.Unmarshal(val, res)
		if rErr == nil {
			wrap := NewSoPostWrap(s.Dba, &res.PostId)
			return wrap
		}
	}
	return nil
}

////////////// SECTION Watchers ///////////////

type PostWatcherFlag struct {
	HasAuthorWatcher bool

	HasBeneficiariesWatcher bool

	HasBodyWatcher bool

	HasCashoutBlockNumWatcher bool

	HasCategoryWatcher bool

	HasChildrenWatcher bool

	HasCopyrightWatcher bool

	HasCopyrightMemoWatcher bool

	HasCreatedWatcher bool

	HasDappRewardsWatcher bool

	HasDepthWatcher bool

	HasLastPayoutWatcher bool

	HasParentIdWatcher bool

	HasRewardsWatcher bool

	HasRootIdWatcher bool

	HasTagsWatcher bool

	HasTicketWatcher bool

	HasTitleWatcher bool

	HasVoteCntWatcher bool

	HasWeightedVpWatcher bool

	WholeWatcher bool
	AnyWatcher   bool
}

var (
	PostTable = &TableInfo{
		Name:    "Post",
		Primary: "PostId",
		Record:  reflect.TypeOf((*SoPost)(nil)).Elem(),
	}
	PostWatcherFlags     = make(map[uint32]PostWatcherFlag)
	PostWatcherFlagsLock sync.RWMutex
)

func PostWatcherFlagOfDb(dbSvcId uint32) PostWatcherFlag {
	PostWatcherFlagsLock.RLock()
	defer PostWatcherFlagsLock.RUnlock()
	return PostWatcherFlags[dbSvcId]
}

func PostRecordWatcherChanged(dbSvcId uint32) {
	var flag PostWatcherFlag
	flag.WholeWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "")
	flag.AnyWatcher = flag.WholeWatcher

	flag.HasAuthorWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Author")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasAuthorWatcher

	flag.HasBeneficiariesWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Beneficiaries")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasBeneficiariesWatcher

	flag.HasBodyWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Body")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasBodyWatcher

	flag.HasCashoutBlockNumWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "CashoutBlockNum")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasCashoutBlockNumWatcher

	flag.HasCategoryWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Category")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasCategoryWatcher

	flag.HasChildrenWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Children")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasChildrenWatcher

	flag.HasCopyrightWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Copyright")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasCopyrightWatcher

	flag.HasCopyrightMemoWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "CopyrightMemo")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasCopyrightMemoWatcher

	flag.HasCreatedWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Created")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasCreatedWatcher

	flag.HasDappRewardsWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "DappRewards")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasDappRewardsWatcher

	flag.HasDepthWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Depth")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasDepthWatcher

	flag.HasLastPayoutWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "LastPayout")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasLastPayoutWatcher

	flag.HasParentIdWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "ParentId")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasParentIdWatcher

	flag.HasRewardsWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Rewards")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasRewardsWatcher

	flag.HasRootIdWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "RootId")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasRootIdWatcher

	flag.HasTagsWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Tags")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasTagsWatcher

	flag.HasTicketWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Ticket")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasTicketWatcher

	flag.HasTitleWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "Title")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasTitleWatcher

	flag.HasVoteCntWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "VoteCnt")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasVoteCntWatcher

	flag.HasWeightedVpWatcher = HasTableRecordWatcher(dbSvcId, PostTable.Record, "WeightedVp")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasWeightedVpWatcher

	PostWatcherFlagsLock.Lock()
	PostWatcherFlags[dbSvcId] = flag
	PostWatcherFlagsLock.Unlock()
}

////////////// SECTION Json query ///////////////

func PostQuery(db iservices.IDatabaseRW, keyJson string) (valueJson string, err error) {
	k := new(uint64)
	d := json.NewDecoder(bytes.NewReader([]byte(keyJson)))
	d.UseNumber()
	if err = d.Decode(k); err != nil {
		return
	}
	if v := NewSoPostWrap(db, k).getPost(); v == nil {
		err = errors.New("not found")
	} else {
		var jbytes []byte
		if jbytes, err = json.Marshal(v); err == nil {
			valueJson = string(jbytes)
		}
	}
	return
}

func init() {
	RegisterTableWatcherChangedCallback(PostTable.Record, PostRecordWatcherChanged)
	RegisterTableJsonQuery("Post", PostQuery)
}
