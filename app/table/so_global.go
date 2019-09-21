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
	GlobalIdUniTable uint32 = 155819495

	GlobalIdRow uint32 = 952861743
)

////////////// SECTION Wrap Define ///////////////
type SoGlobalWrap struct {
	dba         iservices.IDatabaseRW
	mainKey     *int32
	watcherFlag *GlobalWatcherFlag
	mKeyFlag    int    //the flag of the main key exist state in db, -1:has not judged; 0:not exist; 1:already exist
	mKeyBuf     []byte //the buffer after the main key is encoded with prefix
	mBuf        []byte //the value after the main key is encoded
	mdFuncMap   map[string]interface{}
}

func NewSoGlobalWrap(dba iservices.IDatabaseRW, key *int32) *SoGlobalWrap {
	if dba == nil || key == nil {
		return nil
	}
	result := &SoGlobalWrap{dba, key, nil, -1, nil, nil, nil}
	result.initWatcherFlag()
	return result
}

func (s *SoGlobalWrap) CheckExist() bool {
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

func (s *SoGlobalWrap) MustExist(errMsgs ...interface{}) *SoGlobalWrap {
	if !s.CheckExist() {
		panic(bindErrorInfo(fmt.Sprintf("SoGlobalWrap.MustExist: %v not found", s.mainKey), errMsgs...))
	}
	return s
}

func (s *SoGlobalWrap) MustNotExist(errMsgs ...interface{}) *SoGlobalWrap {
	if s.CheckExist() {
		panic(bindErrorInfo(fmt.Sprintf("SoGlobalWrap.MustNotExist: %v already exists", s.mainKey), errMsgs...))
	}
	return s
}

func (s *SoGlobalWrap) initWatcherFlag() {
	if s.watcherFlag == nil {
		s.watcherFlag = new(GlobalWatcherFlag)
		*(s.watcherFlag) = GlobalWatcherFlagOfDb(s.dba.ServiceId())
	}
}

func (s *SoGlobalWrap) create(f func(tInfo *SoGlobal)) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	if s.mainKey == nil {
		return errors.New("the main key is nil")
	}
	val := &SoGlobal{}
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

func (s *SoGlobalWrap) Create(f func(tInfo *SoGlobal), errArgs ...interface{}) *SoGlobalWrap {
	err := s.create(f)
	if err != nil {
		panic(bindErrorInfo(fmt.Errorf("SoGlobalWrap.Create failed: %s", err.Error()), errArgs...))
	}
	return s
}

func (s *SoGlobalWrap) getMainKeyBuf() ([]byte, error) {
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

func (s *SoGlobalWrap) modify(f func(tInfo *SoGlobal)) error {
	if !s.CheckExist() {
		return errors.New("the SoGlobal table does not exist. Please create a table first")
	}
	oriTable := s.getGlobal()
	if oriTable == nil {
		return errors.New("fail to get origin table SoGlobal")
	}

	curTable := s.getGlobal()
	if curTable == nil {
		return errors.New("fail to create current table SoGlobal")
	}
	f(curTable)

	//the main key is not support modify
	if !reflect.DeepEqual(curTable.Id, oriTable.Id) {
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
	err = s.updateGlobal(curTable)
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

func (s *SoGlobalWrap) Modify(f func(tInfo *SoGlobal), errArgs ...interface{}) *SoGlobalWrap {
	err := s.modify(f)
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoGlobalWrap.Modify failed: %s", err.Error()), errArgs...))
	}
	return s
}

func (s *SoGlobalWrap) SetProps(p *prototype.DynamicProperties, errArgs ...interface{}) *SoGlobalWrap {
	err := s.modify(func(r *SoGlobal) {
		r.Props = p
	})
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoGlobalWrap.SetProps( %v ) failed: %s", p, err.Error()), errArgs...))
	}
	return s
}

func (s *SoGlobalWrap) checkSortAndUniFieldValidity(curTable *SoGlobal, fields map[string]bool) error {
	if curTable != nil && fields != nil && len(fields) > 0 {

	}
	return nil
}

//Get all the modified fields in the table
func (s *SoGlobalWrap) getModifiedFields(oriTable *SoGlobal, curTable *SoGlobal) (map[string]bool, bool, error) {
	if oriTable == nil {
		return nil, false, errors.New("table info is nil, can't get modified fields")
	}
	hasWatcher := false
	fields := make(map[string]bool)

	if !reflect.DeepEqual(oriTable.Props, curTable.Props) {
		fields["Props"] = true
		hasWatcher = hasWatcher || s.watcherFlag.HasPropsWatcher
	}

	hasWatcher = hasWatcher || s.watcherFlag.WholeWatcher
	return fields, hasWatcher, nil
}

func (s *SoGlobalWrap) handleFieldMd(t FieldMdHandleType, so *SoGlobal, fields map[string]bool) error {
	if so == nil {
		return errors.New("fail to modify empty table")
	}

	//there is no field need to modify
	if fields == nil || len(fields) < 1 {
		return nil
	}

	errStr := ""

	if fields["Props"] {
		res := true
		if t == FieldMdHandleTypeCheck {
			res = s.mdFieldProps(so.Props, true, false, false, so)
			errStr = fmt.Sprintf("fail to modify exist value of %v", "Props")
		} else if t == FieldMdHandleTypeDel {
			res = s.mdFieldProps(so.Props, false, true, false, so)
			errStr = fmt.Sprintf("fail to delete  sort or unique field  %v", "Props")
		} else if t == FieldMdHandleTypeInsert {
			res = s.mdFieldProps(so.Props, false, false, true, so)
			errStr = fmt.Sprintf("fail to insert  sort or unique field  %v", "Props")
		}
		if !res {
			return errors.New(errStr)
		}
	}

	return nil
}

////////////// SECTION LKeys delete/insert ///////////////

func (s *SoGlobalWrap) delAllSortKeys(br bool, val *SoGlobal) bool {
	if s.dba == nil {
		return false
	}
	res := true

	return res
}

func (s *SoGlobalWrap) insertAllSortKeys(val *SoGlobal) error {
	if s.dba == nil {
		return errors.New("insert sort Field fail,the db is nil ")
	}
	if val == nil {
		return errors.New("insert sort Field fail,get the SoGlobal fail ")
	}

	return nil
}

////////////// SECTION LKeys delete/insert //////////////

func (s *SoGlobalWrap) removeGlobal() error {
	if s.dba == nil {
		return errors.New("database is nil")
	}

	s.initWatcherFlag()

	var oldVal *SoGlobal
	if s.watcherFlag.AnyWatcher {
		oldVal = s.getGlobal()
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

func (s *SoGlobalWrap) RemoveGlobal(errMsgs ...interface{}) *SoGlobalWrap {
	err := s.removeGlobal()
	if err != nil {
		panic(bindErrorInfo(fmt.Sprintf("SoGlobalWrap.RemoveGlobal failed: %s", err.Error()), errMsgs...))
	}
	return s
}

////////////// SECTION Members Get/Modify ///////////////

func (s *SoGlobalWrap) GetId() int32 {
	res := true
	msg := &SoGlobal{}
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
				return msg.Id
			}
		}
	}
	if !res {
		var tmpValue int32
		return tmpValue
	}
	return msg.Id
}

func (s *SoGlobalWrap) GetProps() *prototype.DynamicProperties {
	res := true
	msg := &SoGlobal{}
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
				return msg.Props
			}
		}
	}
	if !res {
		return nil

	}
	return msg.Props
}

func (s *SoGlobalWrap) mdFieldProps(p *prototype.DynamicProperties, isCheck bool, isDel bool, isInsert bool,
	so *SoGlobal) bool {
	if s.dba == nil {
		return false
	}

	if isCheck {
		res := s.checkPropsIsMetMdCondition(p)
		if !res {
			return false
		}
	}

	if isDel {
		res := s.delFieldProps(so)
		if !res {
			return false
		}
	}

	if isInsert {
		res := s.insertFieldProps(so)
		if !res {
			return false
		}
	}
	return true
}

func (s *SoGlobalWrap) delFieldProps(so *SoGlobal) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoGlobalWrap) insertFieldProps(so *SoGlobal) bool {
	if s.dba == nil {
		return false
	}

	return true
}

func (s *SoGlobalWrap) checkPropsIsMetMdCondition(p *prototype.DynamicProperties) bool {
	if s.dba == nil {
		return false
	}

	return true
}

/////////////// SECTION Private function ////////////////

func (s *SoGlobalWrap) update(sa *SoGlobal) bool {
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

func (s *SoGlobalWrap) getGlobal() *SoGlobal {
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

	res := &SoGlobal{}
	if proto.Unmarshal(resBuf, res) != nil {
		return nil
	}
	return res
}

func (s *SoGlobalWrap) updateGlobal(so *SoGlobal) error {
	if s.dba == nil {
		return errors.New("update fail:the db is nil")
	}

	if so == nil {
		return errors.New("update fail: the SoGlobal is nil")
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

func (s *SoGlobalWrap) encodeMainKey() ([]byte, error) {
	if s.mKeyBuf != nil {
		return s.mKeyBuf, nil
	}
	pre := GlobalIdRow
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

func (s *SoGlobalWrap) delAllUniKeys(br bool, val *SoGlobal) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if !s.delUniKeyId(val) {
		if br {
			return false
		} else {
			res = false
		}
	}

	return res
}

func (s *SoGlobalWrap) delUniKeysWithNames(names map[string]string, val *SoGlobal) bool {
	if s.dba == nil {
		return false
	}
	res := true
	if len(names["Id"]) > 0 {
		if !s.delUniKeyId(val) {
			res = false
		}
	}

	return res
}

func (s *SoGlobalWrap) insertAllUniKeys(val *SoGlobal) (map[string]string, error) {
	if s.dba == nil {
		return nil, errors.New("insert uniuqe Field fail,the db is nil ")
	}
	if val == nil {
		return nil, errors.New("insert uniuqe Field fail,get the SoGlobal fail ")
	}
	sucFields := map[string]string{}
	if !s.insertUniKeyId(val) {
		return sucFields, errors.New("insert unique Field Id fail while insert table ")
	}
	sucFields["Id"] = "Id"

	return sucFields, nil
}

func (s *SoGlobalWrap) delUniKeyId(sa *SoGlobal) bool {
	if s.dba == nil {
		return false
	}
	pre := GlobalIdUniTable
	kList := []interface{}{pre}
	if sa != nil {

		sub := sa.Id
		kList = append(kList, sub)
	} else {
		sub := s.GetId()

		kList = append(kList, sub)

	}
	kBuf, err := kope.EncodeSlice(kList)
	if err != nil {
		return false
	}
	return s.dba.Delete(kBuf) == nil
}

func (s *SoGlobalWrap) insertUniKeyId(sa *SoGlobal) bool {
	if s.dba == nil || sa == nil {
		return false
	}

	pre := GlobalIdUniTable
	sub := sa.Id
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
	val := SoUniqueGlobalById{}
	val.Id = sa.Id

	buf, err := proto.Marshal(&val)

	if err != nil {
		return false
	}

	return s.dba.Put(kBuf, buf) == nil

}

type UniGlobalIdWrap struct {
	Dba iservices.IDatabaseRW
}

func NewUniGlobalIdWrap(db iservices.IDatabaseRW) *UniGlobalIdWrap {
	if db == nil {
		return nil
	}
	wrap := UniGlobalIdWrap{Dba: db}
	return &wrap
}

func (s *UniGlobalIdWrap) UniQueryId(start *int32) *SoGlobalWrap {
	if start == nil || s.Dba == nil {
		return nil
	}
	pre := GlobalIdUniTable
	kList := []interface{}{pre, start}
	bufStartkey, err := kope.EncodeSlice(kList)
	val, err := s.Dba.Get(bufStartkey)
	if err == nil {
		res := &SoUniqueGlobalById{}
		rErr := proto.Unmarshal(val, res)
		if rErr == nil {
			wrap := NewSoGlobalWrap(s.Dba, &res.Id)
			return wrap
		}
	}
	return nil
}

////////////// SECTION Watchers ///////////////

type GlobalWatcherFlag struct {
	HasPropsWatcher bool

	WholeWatcher bool
	AnyWatcher   bool
}

var (
	GlobalTable = &TableInfo{
		Name:    "Global",
		Primary: "Id",
		Record:  reflect.TypeOf((*SoGlobal)(nil)).Elem(),
	}
	GlobalWatcherFlags     = make(map[uint32]GlobalWatcherFlag)
	GlobalWatcherFlagsLock sync.RWMutex
)

func GlobalWatcherFlagOfDb(dbSvcId uint32) GlobalWatcherFlag {
	GlobalWatcherFlagsLock.RLock()
	defer GlobalWatcherFlagsLock.RUnlock()
	return GlobalWatcherFlags[dbSvcId]
}

func GlobalRecordWatcherChanged(dbSvcId uint32) {
	var flag GlobalWatcherFlag
	flag.WholeWatcher = HasTableRecordWatcher(dbSvcId, GlobalTable.Record, "")
	flag.AnyWatcher = flag.WholeWatcher

	flag.HasPropsWatcher = HasTableRecordWatcher(dbSvcId, GlobalTable.Record, "Props")
	flag.AnyWatcher = flag.AnyWatcher || flag.HasPropsWatcher

	GlobalWatcherFlagsLock.Lock()
	GlobalWatcherFlags[dbSvcId] = flag
	GlobalWatcherFlagsLock.Unlock()
}

////////////// SECTION Json query ///////////////

func GlobalQuery(db iservices.IDatabaseRW, keyJson string) (valueJson string, err error) {
	k := new(int32)
	d := json.NewDecoder(bytes.NewReader([]byte(keyJson)))
	d.UseNumber()
	if err = d.Decode(k); err != nil {
		return
	}
	if v := NewSoGlobalWrap(db, k).getGlobal(); v == nil {
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
	RegisterTableWatcherChangedCallback(GlobalTable.Record, GlobalRecordWatcherChanged)
	RegisterTableJsonQuery("Global", GlobalQuery)
}
