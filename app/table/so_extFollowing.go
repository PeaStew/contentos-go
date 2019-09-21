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
	ExtFollowingFollowingCreatedOrderTable uint32 = 818811825
	ExtFollowingFollowingInfoUniTable      uint32 = 2172325454

	ExtFollowingFollowingInfoRow uint32 = 2779756882
)

////////////// SECTION Wrap Define ///////////////
type SoExtFollowingWrap struct {
	dba         iservices.IDatabaseRW
	mainKey     *prototype.FollowingRelation
	watcherFlag *ExtFollowingWatcherFlag
	mKeyFlag    int    //the flag of the main key exist state in db, -1:has not judged; 0:not exist; 1:already exist
	mKeyBuf     []byte //the buffer after the main key is encoded with prefix
	mBuf        []byte //the value after the main key is encoded
	mdFuncMap   map[string]interface{}
}

func NewSoExtFollowingWrap(dba iservices.IDatabaseRW, key *prototype.FollowingRelation) *SoExtFollowingWrap {
	if dba == nil || key == nil {
		return nil
	}
	result := &SoExtFollowingWrap{dba, key, nil, -1, nil, nil, nil}
	result.initWatcherFlag()
	return result
}

func (s *SoExtFollowingWrap) CheckExist() bool {
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

func (s *SoExtFollowingWrap) MustExist(errMsgs ...interface{}) *SoExtFollowingWrap {
	if !s.CheckExist() {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowingWrap.MustExist: %v not found", s.mainKey), errMsgs...))
	}
	return s
}

func (s *SoExtFollowingWrap) MustNotExist(errMsgs ...interface{}) *SoExtFollowingWrap {
	if s.CheckExist() {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowingWrap.MustNotExist: %v already exists", s.mainKey), errMsgs...))
	}
	return s
}

func (s *SoExtFollowingWrap) initWatcherFlag() {
	if s.watcherFlag == nil {
		s.watcherFlag = new(ExtFollowingWatcherFlag)
		*(s.watcherFlag) = ExtFollowingWatcherFlagOfDb(s.dba.ServiceId())
	}
}

func (s *SoExtFollowingWrap) create(f func(tInfo *SoExtFollowing)) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	if s.mainKey == nil {
		return errors.New("the main key is nil")
	}
	val := &SoExtFollowing{}
	f(val)
	if val.FollowingInfo == nil {
		val.FollowingInfo = s.mainKey
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

func (s *SoExtFollowingWrap) Create(f func(tInfo *SoExtFollowing), errArgs ...interface{}) *SoExtFollowingWrap {
	err := s.create(f)
	if err != nil {
		panic(bindErrorInfo(fmt.Errorf("SoExtFollowingWrap.Create failed: %s", err.Error()), errArgs...))
	}
	return s
}

func (s *SoExtFollowingWrap) getMainKeyBuf() ([]byte, error) {
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

func (s *SoExtFollowingWrap) modify(f func(tInfo *SoExtFollowing)) error {
	if !s.CheckExist() {
		return errors.New("the SoExtFollowing table does not exist. Please create a table first")
	}
	oriTable := s.getExtFollowing()
	if oriTable == nil {
		return errors.New("fail to get origin table SoExtFollowing")
	}

	curTable := s.getExtFollowing()
	if curTable == nil {
		return errors.New("fail to create current table SoExtFollowing")
	}
	f(curTable)

	//the main key is not support modify
	if !reflect.DeepEqual(curTable.FollowingInfo, oriTable.FollowingInfo) {
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
	err = s.updateExtFollowing(curTable)
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

func (s *SoExtFollowingWrap) Modify(f func(tInfo *SoExtFollowing), errArgs ...interface{}) *SoExtFollowingWrap {
	err := s.modify(f)
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowingWrap.Modify failed: %s", err.Error()), errArgs...))
	}
	return s
}

func (s *SoExtFollowingWrap) SetFollowingCreatedOrder(p *prototype.FollowingCreatedOrder, errArgs ...interface{}) *SoExtFollowingWrap {
	err := s.modify(func(r *SoExtFollowing) {
		r.FollowingCreatedOrder = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowingWrap.SetFollowingCreatedOrder( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoExtFollowingWrap) checkSortAndUniFieldValidity(curTable *SoExtFollowing, fields map[string]bool) error {
	if curTable != nil && fields != nil && len(fields) > 0 {

		if fields["FollowingCreatedOrder"] && curTable.FollowingCreatedOrder == nil {
			return errors.New("sort field FollowingCreatedOrder can't be modified to nil")
		}

	}
	return nil
}

//Get all the modified fields in the table
func (s *SoExtFollowingWrap) getModifiedFields(oriTable *SoExtFollowing, curTable *SoExtFollowing) (map[string]bool, bool, error) {
	if oriTable == nil {
		return nil, false, errors.New("table info is nil, can't get modified fields")
	}
	hasWatcher := false
	fields := make(map[string]bool)

	if !reflect.DeepEqual(oriTable.FollowingCreatedOrder, curTable.FollowingCreatedOrder) {
		fields["FollowingCreatedOrder"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasFollowingCreatedOrderWatcher
	}

	hasWatcher = hasWatcher || s.watcherFlag.WholeWatcher
	return fields, hasWatcher, nil
}

func (s *SoExtFollowingWrap) handleFieldMd(t FieldMdHandleType, so *SoExtFollowing, fields map[string]bool) error {
	if so == nil {
		return errors.New("fail to modify empty table")
	}

	//there is no field need to modify
	if fields == nil || len(fields) < 1 {
		return nil
	}

	errStr := ""

	if fields["FollowingCreatedOrder"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldFollowingCreatedOrder(so.FollowingCreatedOrder, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "FollowingCreatedOrder")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldFollowingCreatedOrder(so.FollowingCreatedOrder, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "FollowingCreatedOrder")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldFollowingCreatedOrder(so.FollowingCreatedOrder, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "FollowingCreatedOrder")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	return nil
}

////////////// SECTION LKeys delete/insert ///////////////

func (s *SoExtFollowingWrap) delSortKeyFollowingCreatedOrder(sa *SoExtFollowing) bool {
	if s.dba == nil || s.mainKey == nil {
		return false
	}
	val := SoListExtFollowingByFollowingCreatedOrder{}
	if sa == nil {
		val.FollowingCreatedOrder = s.GetFollowingCreatedOrder()
		val.FollowingInfo = s.mainKey

	} else {
		val.FollowingCreatedOrder = sa.FollowingCreatedOrder
		val.FollowingInfo = sa.FollowingInfo
	}
	subBuf, err := val.OpeEncode()
	if err != nil {
		return false
	}
	ordErr := s.dba.Delete(subBuf)
	return ordErr == nil
}

func (s *SoExtFollowingWrap) insertSortKeyFollowingCreatedOrder(sa *SoExtFollowing) bool {
	if s.dba == nil || sa == nil {
		return false
	}
	val := SoListExtFollowingByFollowingCreatedOrder{}
	val.FollowingInfo = sa.FollowingInfo
	val.FollowingCreatedOrder = sa.FollowingCreatedOrder
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

func (s *SoExtFollowingWrap) delAllSortKeys(br bool, val *SoExtFollowing) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if !s.delSortKeyFollowingCreatedOrder(val) {
		if br {
			return false
		} else {
			res = false
		}
	}

	return res
}

func (s *SoExtFollowingWrap) insertAllSortKeys(val *SoExtFollowing) error {
	if s.dba == nil {
		return errors.New("insert sort Field fail,the db is nil ")
	}
	if val == nil {
		return errors.New("insert sort Field fail,get the SoExtFollowing fail ")
	}
	if !s.insertSortKeyFollowingCreatedOrder(val) {
		return errors.New("insert sort Field FollowingCreatedOrder fail while insert table ")
	}

	return nil
}

////////////// SECTION LKeys delete/insert //////////////

func (s *SoExtFollowingWrap) removeExtFollowing() error {
	if s.dba == nil {
		return errors.New("database is nil")
	}

	s.initWatcherFlag()

	var oldVal *SoExtFollowing
	if s.watcherFlag.AnyWatcher {
		oldVal = s.getExtFollowing()
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

func (s *SoExtFollowingWrap) RemoveExtFollowing(errMsgs ...interface{}) *SoExtFollowingWrap {
	err := s.removeExtFollowing()
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoExtFollowingWrap.RemoveExtFollowing failed: %s", err.Error()), errMsgs...))
	}
	return s
}

////////////// SECTION Members Get/Modify ///////////////

func (s *SoExtFollowingWrap) GetFollowingCreatedOrder() *prototype.FollowingCreatedOrder {
	res := true
	msg := &SoExtFollowing{}
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
				return msg.FollowingCreatedOrder
			}
		}
	}
	if !res {
		return nil

	}
	return msg.FollowingCreatedOrder
}

func (s *SoExtFollowingWrap) mdFieldFollowingCreatedOrder(p *prototype.FollowingCreatedOrder, isCheck bool, isDel bool, isInsert bool,
	so *SoExtFollowing) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkFollowingCreatedOrderIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldFollowingCreatedOrder(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldFollowingCreatedOrder(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoExtFollowingWrap) delFieldFollowingCreatedOrder(so *SoExtFollowing) bool {
	if s.dba == nil {
		return false
	}

	if !s.delSortKeyFollowingCreatedOrder(so) {
		return false
	}

	return true
}

func (s *SoExtFollowingWrap) insertFieldFollowingCreatedOrder(so *SoExtFollowing) bool {
	if s.dba == nil {
		return false
	}

	if !s.insertSortKeyFollowingCreatedOrder(so) {
		return false
	}

	return true
}

func (s *SoExtFollowingWrap) checkFollowingCreatedOrderIsMetMdCondition(p *prototype.FollowingCreatedOrder) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoExtFollowingWrap) GetFollowingInfo() *prototype.FollowingRelation {
	res := true
	msg := &SoExtFollowing{}
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
				return msg.FollowingInfo
			}
		}
	}
	if !res {
		return nil

	}
	return msg.FollowingInfo
}

////////////// SECTION List Keys ///////////////
type SExtFollowingFollowingCreatedOrderWrap struct {
	Dba iservices.IDatabaseRW
}

func NewExtFollowingFollowingCreatedOrderWrap(db iservices.IDatabaseRW) *SExtFollowingFollowingCreatedOrderWrap {
	if db == nil {
		return nil
	}
	wrap := SExtFollowingFollowingCreatedOrderWrap{Dba: db}
	return &wrap
}

func (s *SExtFollowingFollowingCreatedOrderWrap) GetMainVal(val []byte) *prototype.FollowingRelation {
	res := &SoListExtFollowingByFollowingCreatedOrder{}
	err := proto.Unmarshal(val, res)

	if err != nil {
		return nil
	}
	return res.FollowingInfo

}

func (s *SExtFollowingFollowingCreatedOrderWrap) GetSubVal(val []byte) *prototype.FollowingCreatedOrder {
	res := &SoListExtFollowingByFollowingCreatedOrder{}
	err := proto.Unmarshal(val, res)
	if err != nil {
		return nil
	}
	return res.FollowingCreatedOrder

}

func (m *SoListExtFollowingByFollowingCreatedOrder) OpeEncode() ([]byte, error) {
	pre := ExtFollowingFollowingCreatedOrderTable
	sub := m.FollowingCreatedOrder
	if sub == nil {
		return nil, errors.New("the pro FollowingCreatedOrder is nil")
	}
	sub1 := m.FollowingInfo
	if sub1 == nil {
		return nil, errors.New("the mainkey FollowingInfo is nil")
	}
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
func (s *SExtFollowingFollowingCreatedOrderWrap) ForEachByOrder(start *prototype.FollowingCreatedOrder, end *prototype.FollowingCreatedOrder, lastMainKey *prototype.FollowingRelation,
	lastSubVal *prototype.FollowingCreatedOrder, f func(mVal *prototype.FollowingRelation, sVal *prototype.FollowingCreatedOrder, idx uint32) bool) error {
	if s.Dba == nil {
		return errors.New("the db is nil")
	}
	if (lastSubVal != nil && lastMainKey == nil) || (lastSubVal == nil && lastMainKey != nil) {
		return errors.New("last query param error")
	}
	if f == nil {
		return nil
	}
	pre := ExtFollowingFollowingCreatedOrderTable
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

/////////////// SECTION Private function ////////////////

func (s *SoExtFollowingWrap) update(sa *SoExtFollowing) bool {
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

func (s *SoExtFollowingWrap) getExtFollowing() *SoExtFollowing {
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

	res := &SoExtFollowing{}
	if proto.Unmarshal(resBuf, res) != nil {
		return nil
	}
	return res
}

func (s *SoExtFollowingWrap) updateExtFollowing(so *SoExtFollowing) error {
	if s.dba == nil {
		return errors.New("update fail:the db is nil")
	}

	if so == nil {
		return errors.New("update fail: the SoExtFollowing is nil")
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

func (s *SoExtFollowingWrap) encodeMainKey() ([]byte, error) {
	if s.mKeyBuf != nil {
		return s.mKeyBuf, nil
	}
	pre := ExtFollowingFollowingInfoRow
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

func (s *SoExtFollowingWrap) delAllUniKeys(br bool, val *SoExtFollowing) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if !s.delUniKeyFollowingInfo(val) {
		if br {
			return false
		} else {
			res = false
		}
	}

	return res
}

func (s *SoExtFollowingWrap) delUniKeysWithNames(names map[string]string, val *SoExtFollowing) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if len(names["FollowingInfo"]) > 0 {
		if !s.delUniKeyFollowingInfo(val) {
			res = false
		}
	}

	return res
}

func (s *SoExtFollowingWrap) insertAllUniKeys(val *SoExtFollowing) (map[string]string, error) {
	if s.dba == nil {
		return nil, errors.New("insert uniuqe Field fail,the db is nil ")
	}
	if val == nil {
		return nil, errors.New("insert uniuqe Field fail,get the SoExtFollowing fail ")
	}
	sucFields := map[string]string{}
	if !s.insertUniKeyFollowingInfo(val) {
		return sucFields, errors.New("insert unique Field FollowingInfo fail while insert table ")
	}
	sucFields["FollowingInfo"] = "FollowingInfo"

	return sucFields, nil
}

func (s *SoExtFollowingWrap) delUniKeyFollowingInfo(sa *SoExtFollowing) bool {
	if s.dba == nil {
		return false
	}
	pre := ExtFollowingFollowingInfoUniTable
	kList := []interface{}{pre}
	if sa != nil {
		if sa.FollowingInfo == nil {
			return false
		}

		sub := sa.FollowingInfo
		kList = append(kList, sub)
	} else {
		sub := s.GetFollowingInfo()
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

func (s *SoExtFollowingWrap) insertUniKeyFollowingInfo(sa *SoExtFollowing) bool {
	if s.dba == nil || sa == nil {
		return false
	}

	pre := ExtFollowingFollowingInfoUniTable
	sub := sa.FollowingInfo
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
	val := SoUniqueExtFollowingByFollowingInfo{}
	val.FollowingInfo = sa.FollowingInfo

	buf, err := proto.Marshal(&val)

	if err != nil {
		return false
	}

	return s.dba.Put(kBuf, buf) == nil

}

type UniExtFollowingFollowingInfoWrap struct {
	Dba iservices.IDatabaseRW
}

func NewUniExtFollowingFollowingInfoWrap(db iservices.IDatabaseRW) *UniExtFollowingFollowingInfoWrap {
	if db == nil {
		return nil
	}
	wrap := UniExtFollowingFollowingInfoWrap{Dba: db}
	return &wrap
}

func (s *UniExtFollowingFollowingInfoWrap) UniQueryFollowingInfo(start *prototype.FollowingRelation) *SoExtFollowingWrap {
	if start == nil || s.Dba == nil {
		return nil
	}
	pre := ExtFollowingFollowingInfoUniTable
	kList := []interface{}{pre, start}
	bufStartkey, err := kope.EncodeSlice(kList)
	val, err := s.Dba.Get(bufStartkey)
	if err == nil {
		res := &SoUniqueExtFollowingByFollowingInfo{}
		rErr := proto.Unmarshal(val, res)
		if rErr == nil {
			wrap := NewSoExtFollowingWrap(s.Dba, res.FollowingInfo)

			return wrap
		}
	}
	return nil
}

////////////// SECTION Watchers ///////////////

type ExtFollowingWatcherFlag struct {
	HasFollowingCreatedOrderWatcher bool

	WholeWatcher bool
	AnyWatcher   bool
}

var (
	ExtFollowingTable = &TableInfo{
		Name:    "ExtFollowing",
		Primary: "FollowingInfo",
		Record:  reflect.TypeOf((*SoExtFollowing)(nil)).Elem(),
	}
	ExtFollowingWatcherFlags     = make(map[uint32]ExtFollowingWatcherFlag)
	ExtFollowingWatcherFlagsLock sync.RWMutex
)

func ExtFollowingWatcherFlagOfDb(dbSvcId uint32) ExtFollowingWatcherFlag {
	ExtFollowingWatcherFlagsLock.RLock()
	defer ExtFollowingWatcherFlagsLock.RUnlock()
	return ExtFollowingWatcherFlags[dbSvcId]
}

func ExtFollowingRecordWatcherChanged(dbSvcId uint32) {
	var flag ExtFollowingWatcherFlag
	flag.WholeWatcher = HasTableRecordWatcher(dbSvcId, ExtFollowingTable.Record, "")
	flag.AnyWatcher = flag.WholeWatcher

	flag.HasFollowingCreatedOrderWatcher = HasTableRecordWatcher(dbSvcId, ExtFollowingTable.Record, "FollowingCreatedOrder")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasFollowingCreatedOrderWatcher

	ExtFollowingWatcherFlagsLock.Lock()
	ExtFollowingWatcherFlags[dbSvcId] = flag
	ExtFollowingWatcherFlagsLock.Unlock()
}

////////////// SECTION Json query ///////////////

func ExtFollowingQuery(db iservices.IDatabaseRW, keyJson string) (valueJson string, err error) {
	k := new(prototype.FollowingRelation)
	d := json.NewDecoder(bytes.NewReader([]byte(keyJson)))
	d.UseNumber()
	if err = d.Decode(k); err != nil {
		return
	}
	if v := NewSoExtFollowingWrap(db, k).getExtFollowing(); v == nil {
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
	RegisterTableWatcherChangedCallback(ExtFollowingTable.Record, ExtFollowingRecordWatcherChanged)
	RegisterTableJsonQuery("ExtFollowing", ExtFollowingQuery)
}
