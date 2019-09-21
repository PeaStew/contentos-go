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
	ExtFollowCountAccountUniTable uint32 = 1673144335

	ExtFollowCountAccountRow uint32 = 1347382161
)

////////////// SECTION Wrap Define ///////////////
type SoExtFollowCountWrap struct {
	dba         iservices.IDatabaseRW
	mainKey     *prototype.AccountName
	watcherFlag *ExtFollowCountWatcherFlag
	mKeyFlag    int    //the flag of the main key exist state in db, -1:has not judged; 0:not exist; 1:already exist
	mKeyBuf     []byte //the buffer after the main key is encoded with prefix
	mBuf        []byte //the value after the main key is encoded
	mdFuncMap   map[string]interface{}
}

func NewSoExtFollowCountWrap(dba iservices.IDatabaseRW, key *prototype.AccountName) *SoExtFollowCountWrap {
	if dba == nil || key == nil {
		return nil
	}
	result := &SoExtFollowCountWrap{dba, key, nil, -1, nil, nil, nil}
	result.initWatcherFlag()
	return result
}

func (s *SoExtFollowCountWrap) CheckExist() bool {
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

func (s *SoExtFollowCountWrap) MustExist(errMsgs ...interface{}) *SoExtFollowCountWrap {
	if !s.CheckExist() {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowCountWrap.MustExist: %v not found", s.mainKey), errMsgs...))
	}
	return s
}

func (s *SoExtFollowCountWrap) MustNotExist(errMsgs ...interface{}) *SoExtFollowCountWrap {
	if s.CheckExist() {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowCountWrap.MustNotExist: %v already exists", s.mainKey), errMsgs...))
	}
	return s
}

func (s *SoExtFollowCountWrap) initWatcherFlag() {
	if s.watcherFlag == nil {
		s.watcherFlag = new(ExtFollowCountWatcherFlag)
		*(s.watcherFlag) = ExtFollowCountWatcherFlagOfDb(s.dba.ServiceId())
	}
}

func (s *SoExtFollowCountWrap) create(f func(tInfo *SoExtFollowCount)) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	if s.mainKey == nil {
		return errors.New("the main key is nil")
	}
	val := &SoExtFollowCount{}
	f(val)
	if val.Account == nil {
		val.Account = s.mainKey
	}
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

func (s *SoExtFollowCountWrap) Create(f func(tInfo *SoExtFollowCount), errArgs ...interface{}) *SoExtFollowCountWrap {
	err := s.create(f)
	if err != nil {
		panic(bindErrorInfo(fmt.Errorf("SoExtFollowCountWrap.Create failed: %s", err.Error()), errArgs...))
	}
	return s
}

func (s *SoExtFollowCountWrap) getMainKeyBuf() ([]byte, error) {
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

func (s *SoExtFollowCountWrap) modify(f func(tInfo *SoExtFollowCount)) error {
	if !s.CheckExist() {
		return errors.New("the SoExtFollowCount table does not exist. Please create a table first")
	}
	oriTable := s.getExtFollowCount()
	if oriTable == nil {
		return errors.New("fail to get origin table SoExtFollowCount")
	}

	curTable := s.getExtFollowCount()
	if curTable == nil {
		return errors.New("fail to create current table SoExtFollowCount")
	}
	f(curTable)

	//the main key is not support modify
	if !reflect.DeepEqual(curTable.Account, oriTable.Account) {
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
	err = s.updateExtFollowCount(curTable)
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

func (s *SoExtFollowCountWrap) Modify(f func(tInfo *SoExtFollowCount), errArgs ...interface{}) *SoExtFollowCountWrap {
	err := s.modify(f)
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowCountWrap.Modify failed: %s", err.Error()), errArgs...))
	}
	return s
}

func (s *SoExtFollowCountWrap) SetFollowerCnt(p uint32, errArgs ...interface{}) *SoExtFollowCountWrap {
	err := s.modify(func(r *SoExtFollowCount) {
		r.FollowerCnt = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowCountWrap.SetFollowerCnt( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoExtFollowCountWrap) SetFollowingCnt(p uint32, errArgs ...interface{}) *SoExtFollowCountWrap {
	err := s.modify(func(r *SoExtFollowCount) {
		r.FollowingCnt = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowCountWrap.SetFollowingCnt( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoExtFollowCountWrap) SetUpdateTime(p *prototype.TimePointSec, errArgs ...interface{}) *SoExtFollowCountWrap {
	err := s.modify(func(r *SoExtFollowCount) {
		r.UpdateTime = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowCountWrap.SetUpdateTime( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoExtFollowCountWrap) checkSortAndUniFieldValidity(curTable *SoExtFollowCount, fields map[string]bool) error {
	if curTable != nil && fields != nil && len(fields) > 0 {

	}
	return nil
}

//Get all the modified fields in the table
func (s *SoExtFollowCountWrap) getModifiedFields(oriTable *SoExtFollowCount, curTable *SoExtFollowCount) (map[string]bool, bool, error) {
	if oriTable == nil {
		return nil, false, errors.New("table info is nil, can't get modified fields")
	}
	hasWatcher := false
	fields := make(map[string]bool)

	if !reflect.DeepEqual(oriTable.FollowerCnt, curTable.FollowerCnt) {
		fields["FollowerCnt"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasFollowerCntWatcher
	}

	if !reflect.DeepEqual(oriTable.FollowingCnt, curTable.FollowingCnt) {
		fields["FollowingCnt"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasFollowingCntWatcher
	}

	if !reflect.DeepEqual(oriTable.UpdateTime, curTable.UpdateTime) {
		fields["UpdateTime"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasUpdateTimeWatcher
	}

	hasWatcher = hasWatcher || s.watcherFlag.WholeWatcher
	return fields, hasWatcher, nil
}

func (s *SoExtFollowCountWrap) handleFieldMd(t FieldMdHandleType, so *SoExtFollowCount, fields map[string]bool) error {
	if so == nil {
		return errors.New("fail to modify empty table")
	}

	//there is no field need to modify
	if fields == nil || len(fields) < 1 {
		return nil
	}

	errStr := ""

	if fields["FollowerCnt"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldFollowerCnt(so.FollowerCnt, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "FollowerCnt")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldFollowerCnt(so.FollowerCnt, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "FollowerCnt")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldFollowerCnt(so.FollowerCnt, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "FollowerCnt")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["FollowingCnt"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldFollowingCnt(so.FollowingCnt, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "FollowingCnt")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldFollowingCnt(so.FollowingCnt, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "FollowingCnt")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldFollowingCnt(so.FollowingCnt, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "FollowingCnt")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	if fields["UpdateTime"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldUpdateTime(so.UpdateTime, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "UpdateTime")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldUpdateTime(so.UpdateTime, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "UpdateTime")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldUpdateTime(so.UpdateTime, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "UpdateTime")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	return nil
}

////////////// SECTION LKeys delete/insert ///////////////

func (s *SoExtFollowCountWrap) delAllSortKeys(br bool, val *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}
	res := true

	return res
}

func (s *SoExtFollowCountWrap) insertAllSortKeys(val *SoExtFollowCount) error {
	if s.dba == nil {
		return errors.New("insert sort Field fail,the db is nil ")
	}
	if val == nil {
		return errors.New("insert sort Field fail,get the SoExtFollowCount fail ")
	}

	return nil
}

////////////// SECTION LKeys delete/insert //////////////

func (s *SoExtFollowCountWrap) removeExtFollowCount() error {
	if s.dba == nil {
		return errors.New("database is nil")
	}

	s.initWatcherFlag()

	var oldVal *SoExtFollowCount
	if s.watcherFlag.AnyWatcher {
		oldVal = s.getExtFollowCount()
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

func (s *SoExtFollowCountWrap) RemoveExtFollowCount(errMsgs ...interface{}) *SoExtFollowCountWrap {
	err := s.removeExtFollowCount()
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowCountWrap.RemoveExtFollowCount failed: %s", err.Error()), errMsgs...))
	}
	return s
}

////////////// SECTION Members Get/Modify ///////////////

func (s *SoExtFollowCountWrap) GetAccount() *prototype.AccountName {
	res := true
	msg := &SoExtFollowCount{}
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
				return msg.Account
			}
		}
	}
	if !res {
		return nil

	}
	return msg.Account
}

func (s *SoExtFollowCountWrap) GetFollowerCnt() uint32 {
	res := true
	msg := &SoExtFollowCount{}
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
				return msg.FollowerCnt
			}
		}
	}
	if !res {
		var tmpValue uint32
		return tmpValue
	}
	return msg.FollowerCnt
}

func (s *SoExtFollowCountWrap) mdFieldFollowerCnt(p uint32, isCheck bool, isDel bool, isInsert bool,
	so *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkFollowerCntIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldFollowerCnt(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldFollowerCnt(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoExtFollowCountWrap) delFieldFollowerCnt(so *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoExtFollowCountWrap) insertFieldFollowerCnt(so *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoExtFollowCountWrap) checkFollowerCntIsMetMdCondition(p uint32) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoExtFollowCountWrap) GetFollowingCnt() uint32 {
	res := true
	msg := &SoExtFollowCount{}
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
				return msg.FollowingCnt
			}
		}
	}
	if !res {
		var tmpValue uint32
		return tmpValue
	}
	return msg.FollowingCnt
}

func (s *SoExtFollowCountWrap) mdFieldFollowingCnt(p uint32, isCheck bool, isDel bool, isInsert bool,
	so *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkFollowingCntIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldFollowingCnt(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldFollowingCnt(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoExtFollowCountWrap) delFieldFollowingCnt(so *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoExtFollowCountWrap) insertFieldFollowingCnt(so *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoExtFollowCountWrap) checkFollowingCntIsMetMdCondition(p uint32) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoExtFollowCountWrap) GetUpdateTime() *prototype.TimePointSec {
	res := true
	msg := &SoExtFollowCount{}
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
				return msg.UpdateTime
			}
		}
	}
	if !res {
		return nil

	}
	return msg.UpdateTime
}

func (s *SoExtFollowCountWrap) mdFieldUpdateTime(p *prototype.TimePointSec, isCheck bool, isDel bool, isInsert bool,
	so *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkUpdateTimeIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldUpdateTime(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldUpdateTime(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoExtFollowCountWrap) delFieldUpdateTime(so *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoExtFollowCountWrap) insertFieldUpdateTime(so *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoExtFollowCountWrap) checkUpdateTimeIsMetMdCondition(p *prototype.TimePointSec) bool {
	if s.dba == nil {
		return false
	}

	return true
}

/////////////// SECTION Private function ////////////////

func (s *SoExtFollowCountWrap) update(sa *SoExtFollowCount) bool {
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

func (s *SoExtFollowCountWrap) getExtFollowCount() *SoExtFollowCount {
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

	res := &SoExtFollowCount{}
	if proto.Unmarshal(resBuf, res) != nil {
		return nil
	}
	return res
}

func (s *SoExtFollowCountWrap) updateExtFollowCount(so *SoExtFollowCount) error {
	if s.dba == nil {
		return errors.New("update fail:the db is nil")
	}

	if so == nil {
		return errors.New("update fail: the SoExtFollowCount is nil")
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

func (s *SoExtFollowCountWrap) encodeMainKey() ([]byte, error) {
	if s.mKeyBuf != nil {
		return s.mKeyBuf, nil
	}
	pre := ExtFollowCountAccountRow
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

func (s *SoExtFollowCountWrap) delAllUniKeys(br bool, val *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if !s.delUniKeyAccount(val) {
		if br {
			return false
		} else {
			res = false
		}
	}

	return res
}

func (s *SoExtFollowCountWrap) delUniKeysWithNames(names map[string]string, val *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if len(names["Account"]) > 0 {
		if !s.delUniKeyAccount(val) {
			res = false
		}
	}

	return res
}

func (s *SoExtFollowCountWrap) insertAllUniKeys(val *SoExtFollowCount) (map[string]string, error) {
	if s.dba == nil {
		return nil, errors.New("insert uniuqe Field fail,the db is nil ")
	}
	if val == nil {
		return nil, errors.New("insert uniuqe Field fail,get the SoExtFollowCount fail ")
	}
	sucFields := map[string]string{}
	if !s.insertUniKeyAccount(val) {
		return sucFields, errors.New("insert unique Field Account fail while insert table ")
	}
	sucFields["Account"] = "Account"

	return sucFields, nil
}

func (s *SoExtFollowCountWrap) delUniKeyAccount(sa *SoExtFollowCount) bool {
	if s.dba == nil {
		return false
	}
	pre := ExtFollowCountAccountUniTable
	kList := []interface{}{pre}
	if sa != nil {
		if sa.Account == nil {
			return false
		}

		sub := sa.Account
		kList = append(kList, sub)
	} else {
		sub := s.GetAccount()
		if sub == nil {
			return true
		}

		kList = append(kList, sub)

	}
	kBuf, err := kope.EncodeSlice(kList)
	if err != nil {
		return false
	}
	return s.dba.Delete(kBuf) == nil
}

func (s *SoExtFollowCountWrap) insertUniKeyAccount(sa *SoExtFollowCount) bool {
	if s.dba == nil || sa == nil {
		return false
	}

	pre := ExtFollowCountAccountUniTable
	sub := sa.Account
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
	val := SoUniqueExtFollowCountByAccount{}
	val.Account = sa.Account

	buf, err := proto.Marshal(&val)

	if err != nil {
		return false
	}

	return s.dba.Put(kBuf, buf) == nil

}

type UniExtFollowCountAccountWrap struct {
	Dba iservices.IDatabaseRW
}

func NewUniExtFollowCountAccountWrap(db iservices.IDatabaseRW) *UniExtFollowCountAccountWrap {
	if db == nil {
		return nil
	}
	wrap := UniExtFollowCountAccountWrap{Dba: db}
	return &wrap
}

func (s *UniExtFollowCountAccountWrap) UniQueryAccount(start *prototype.AccountName) *SoExtFollowCountWrap {
	if start == nil || s.Dba == nil {
		return nil
	}
	pre := ExtFollowCountAccountUniTable
	kList := []interface{}{pre, start}
	bufStartkey, err := kope.EncodeSlice(kList)
	val, err := s.Dba.Get(bufStartkey)
	if err == nil {
		res := &SoUniqueExtFollowCountByAccount{}
		rErr := proto.Unmarshal(val, res)
		if rErr == nil {
			wrap := NewSoExtFollowCountWrap(s.Dba, res.Account)

			return wrap
		}
	}
	return nil
}

////////////// SECTION Watchers ///////////////

type ExtFollowCountWatcherFlag struct {
	HasFollowerCntWatcher bool

	HasFollowingCntWatcher bool

	HasUpdateTimeWatcher bool

	WholeWatcher bool
	AnyWatcher   bool
}

var (
	ExtFollowCountTable = &TableInfo{
		Name:    "ExtFollowCount",
		Primary: "Account",
		Record:  reflect.TypeOf((*SoExtFollowCount)(nil)).Elem(),
	}
	ExtFollowCountWatcherFlags     = make(map[uint32]ExtFollowCountWatcherFlag)
	ExtFollowCountWatcherFlagsLock sync.RWMutex
)

func ExtFollowCountWatcherFlagOfDb(dbSvcId uint32) ExtFollowCountWatcherFlag {
	ExtFollowCountWatcherFlagsLock.RLock()
	defer ExtFollowCountWatcherFlagsLock.RUnlock()
	return ExtFollowCountWatcherFlags[dbSvcId]
}

func ExtFollowCountRecordWatcherChanged(dbSvcId uint32) {
	var flag ExtFollowCountWatcherFlag
	flag.WholeWatcher = HasTableRecordWatcher(dbSvcId, ExtFollowCountTable.Record, "")
	flag.AnyWatcher = flag.WholeWatcher

	flag.HasFollowerCntWatcher = HasTableRecordWatcher(dbSvcId, ExtFollowCountTable.Record, "FollowerCnt")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasFollowerCntWatcher

	flag.HasFollowingCntWatcher = HasTableRecordWatcher(dbSvcId, ExtFollowCountTable.Record, "FollowingCnt")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasFollowingCntWatcher

	flag.HasUpdateTimeWatcher = HasTableRecordWatcher(dbSvcId, ExtFollowCountTable.Record, "UpdateTime")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasUpdateTimeWatcher

	ExtFollowCountWatcherFlagsLock.Lock()
	ExtFollowCountWatcherFlags[dbSvcId] = flag
	ExtFollowCountWatcherFlagsLock.Unlock()
}

////////////// SECTION Json query ///////////////

func ExtFollowCountQuery(db iservices.IDatabaseRW, keyJson string) (valueJson string, err error) {
	k := new(prototype.AccountName)
	d := json.NewDecoder(bytes.NewReader([]byte(keyJson)))
	d.UseNumber()
	if err = d.Decode(k); err != nil {
		return
	}
	if v := NewSoExtFollowCountWrap(db, k).getExtFollowCount(); v == nil {
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
	RegisterTableWatcherChangedCallback(ExtFollowCountTable.Record, ExtFollowCountRecordWatcherChanged)
	RegisterTableJsonQuery("ExtFollowCount", ExtFollowCountQuery)
}
