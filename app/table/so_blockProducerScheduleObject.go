package table

import (
	"errors"
	fmt "fmt"
	"reflect"
	"strings"

	"github.com/coschain/contentos-go/common/encoding/kope"
	"github.com/coschain/contentos-go/iservices"
	prototype "github.com/coschain/contentos-go/prototype"
	proto "github.com/golang/protobuf/proto"
)

////////////// SECTION Prefix Mark ///////////////
var (
	BlockProducerScheduleObjectIdUniTable                       uint32 = 1798653281
	BlockProducerScheduleObjectCurrentShuffledBlockProducerCell uint32 = 3160089807
	BlockProducerScheduleObjectIdCell                           uint32 = 1857477157
	BlockProducerScheduleObjectPubKeyCell                       uint32 = 562239348
)

////////////// SECTION Wrap Define ///////////////
type SoBlockProducerScheduleObjectWrap struct {
	dba      iservices.IDatabaseRW
	mainKey  *int32
	mKeyFlag int    //the flag of the main key exist state in db, -1:has not judged; 0:not exist; 1:already exist
	mKeyBuf  []byte //the buffer after the main key is encoded with prefix
	mBuf     []byte //the value after the main key is encoded
}

func NewSoBlockProducerScheduleObjectWrap(dba iservices.IDatabaseRW, key *int32) *SoBlockProducerScheduleObjectWrap {
	if dba == nil || key == nil {
		return nil
	}
	result := &SoBlockProducerScheduleObjectWrap{dba, key, -1, nil, nil}
	return result
}

func (s *SoBlockProducerScheduleObjectWrap) CheckExist() bool {
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

func (s *SoBlockProducerScheduleObjectWrap) Create(f func(tInfo *SoBlockProducerScheduleObject)) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	if s.mainKey == nil {
		return errors.New("the main key is nil")
	}
	val := &SoBlockProducerScheduleObject{}
	f(val)
	if s.CheckExist() {
		return errors.New("the main key is already exist")
	}
	keyBuf, err := s.encodeMainKey()
	if err != nil {
		return err

	}
	err = s.saveAllMemKeys(val, true)
	if err != nil {
		s.delAllMemKeys(false, val)
		return err
	}

	// update srt list keys
	if err = s.insertAllSortKeys(val); err != nil {
		s.delAllSortKeys(false, val)
		s.dba.Delete(keyBuf)
		s.delAllMemKeys(false, val)
		return err
	}

	//update unique list
	if sucNames, err := s.insertAllUniKeys(val); err != nil {
		s.delAllSortKeys(false, val)
		s.delUniKeysWithNames(sucNames, val)
		s.dba.Delete(keyBuf)
		s.delAllMemKeys(false, val)
		return err
	}

	return nil
}

func (s *SoBlockProducerScheduleObjectWrap) getMainKeyBuf() ([]byte, error) {
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

////////////// SECTION LKeys delete/insert ///////////////

func (s *SoBlockProducerScheduleObjectWrap) delAllSortKeys(br bool, val *SoBlockProducerScheduleObject) bool {
	if s.dba == nil {
		return false
	}
	res := true

	return res
}

func (s *SoBlockProducerScheduleObjectWrap) insertAllSortKeys(val *SoBlockProducerScheduleObject) error {
	if s.dba == nil {
		return errors.New("insert sort Field fail,the db is nil ")
	}
	if val == nil {
		return errors.New("insert sort Field fail,get the SoBlockProducerScheduleObject fail ")
	}

	return nil
}

////////////// SECTION LKeys delete/insert //////////////

func (s *SoBlockProducerScheduleObjectWrap) RemoveBlockProducerScheduleObject() bool {
	if s.dba == nil {
		return false
	}
	val := &SoBlockProducerScheduleObject{}
	//delete sort list key
	if res := s.delAllSortKeys(true, nil); !res {
		return false
	}

	//delete unique list
	if res := s.delAllUniKeys(true, nil); !res {
		return false
	}

	err := s.delAllMemKeys(true, val)
	if err == nil {
		s.mKeyBuf = nil
		s.mKeyFlag = -1
		return true
	} else {
		return false
	}
}

////////////// SECTION Members Get/Modify ///////////////
func (s *SoBlockProducerScheduleObjectWrap) getMemKeyPrefix(fName string) uint32 {
	if fName == "CurrentShuffledBlockProducer" {
		return BlockProducerScheduleObjectCurrentShuffledBlockProducerCell
	}
	if fName == "Id" {
		return BlockProducerScheduleObjectIdCell
	}
	if fName == "PubKey" {
		return BlockProducerScheduleObjectPubKeyCell
	}

	return 0
}

func (s *SoBlockProducerScheduleObjectWrap) encodeMemKey(fName string) ([]byte, error) {
	if len(fName) < 1 || s.mainKey == nil {
		return nil, errors.New("field name or main key is empty")
	}
	pre := s.getMemKeyPrefix(fName)
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
	return kope.PackList(list), nil
}

func (s *SoBlockProducerScheduleObjectWrap) saveAllMemKeys(tInfo *SoBlockProducerScheduleObject, br bool) error {
	if s.dba == nil {
		return errors.New("save member Field fail , the db is nil")
	}

	if tInfo == nil {
		return errors.New("save member Field fail , the data is nil ")
	}
	var err error = nil
	errDes := ""
	if err = s.saveMemKeyCurrentShuffledBlockProducer(tInfo); err != nil {
		if br {
			return err
		} else {
			errDes += fmt.Sprintf("save the Field %s fail,error is %s;\n", "CurrentShuffledBlockProducer", err)
		}
	}
	if err = s.saveMemKeyId(tInfo); err != nil {
		if br {
			return err
		} else {
			errDes += fmt.Sprintf("save the Field %s fail,error is %s;\n", "Id", err)
		}
	}
	if err = s.saveMemKeyPubKey(tInfo); err != nil {
		if br {
			return err
		} else {
			errDes += fmt.Sprintf("save the Field %s fail,error is %s;\n", "PubKey", err)
		}
	}

	if len(errDes) > 0 {
		return errors.New(errDes)
	}
	return err
}

func (s *SoBlockProducerScheduleObjectWrap) delAllMemKeys(br bool, tInfo *SoBlockProducerScheduleObject) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	t := reflect.TypeOf(*tInfo)
	errDesc := ""
	for k := 0; k < t.NumField(); k++ {
		name := t.Field(k).Name
		if len(name) > 0 && !strings.HasPrefix(name, "XXX_") {
			err := s.delMemKey(name)
			if err != nil {
				if br {
					return err
				}
				errDesc += fmt.Sprintf("delete the Field %s fail,error is %s;\n", name, err)
			}
		}
	}
	if len(errDesc) > 0 {
		return errors.New(errDesc)
	}
	return nil
}

func (s *SoBlockProducerScheduleObjectWrap) delMemKey(fName string) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	if len(fName) <= 0 {
		return errors.New("the field name is empty ")
	}
	key, err := s.encodeMemKey(fName)
	if err != nil {
		return err
	}
	err = s.dba.Delete(key)
	return err
}

func (s *SoBlockProducerScheduleObjectWrap) saveMemKeyCurrentShuffledBlockProducer(tInfo *SoBlockProducerScheduleObject) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	if tInfo == nil {
		return errors.New("the data is nil")
	}
	val := SoMemBlockProducerScheduleObjectByCurrentShuffledBlockProducer{}
	val.CurrentShuffledBlockProducer = tInfo.CurrentShuffledBlockProducer
	key, err := s.encodeMemKey("CurrentShuffledBlockProducer")
	if err != nil {
		return err
	}
	buf, err := proto.Marshal(&val)
	if err != nil {
		return err
	}
	err = s.dba.Put(key, buf)
	return err
}

func (s *SoBlockProducerScheduleObjectWrap) GetCurrentShuffledBlockProducer() []string {
	res := true
	msg := &SoMemBlockProducerScheduleObjectByCurrentShuffledBlockProducer{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMemKey("CurrentShuffledBlockProducer")
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
				return msg.CurrentShuffledBlockProducer
			}
		}
	}
	if !res {
		var tmpValue []string
		return tmpValue
	}
	return msg.CurrentShuffledBlockProducer
}

func (s *SoBlockProducerScheduleObjectWrap) MdCurrentShuffledBlockProducer(p []string) bool {
	if s.dba == nil {
		return false
	}
	key, err := s.encodeMemKey("CurrentShuffledBlockProducer")
	if err != nil {
		return false
	}
	buf, err := s.dba.Get(key)
	if err != nil {
		return false
	}
	ori := &SoMemBlockProducerScheduleObjectByCurrentShuffledBlockProducer{}
	err = proto.Unmarshal(buf, ori)
	sa := &SoBlockProducerScheduleObject{}
	sa.Id = *s.mainKey
	sa.CurrentShuffledBlockProducer = ori.CurrentShuffledBlockProducer

	ori.CurrentShuffledBlockProducer = p
	val, err := proto.Marshal(ori)
	if err != nil {
		return false
	}
	err = s.dba.Put(key, val)
	if err != nil {
		return false
	}
	sa.CurrentShuffledBlockProducer = p

	return true
}

func (s *SoBlockProducerScheduleObjectWrap) saveMemKeyId(tInfo *SoBlockProducerScheduleObject) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	if tInfo == nil {
		return errors.New("the data is nil")
	}
	val := SoMemBlockProducerScheduleObjectById{}
	val.Id = tInfo.Id
	key, err := s.encodeMemKey("Id")
	if err != nil {
		return err
	}
	buf, err := proto.Marshal(&val)
	if err != nil {
		return err
	}
	err = s.dba.Put(key, buf)
	return err
}

func (s *SoBlockProducerScheduleObjectWrap) GetId() int32 {
	res := true
	msg := &SoMemBlockProducerScheduleObjectById{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMemKey("Id")
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

func (s *SoBlockProducerScheduleObjectWrap) saveMemKeyPubKey(tInfo *SoBlockProducerScheduleObject) error {
	if s.dba == nil {
		return errors.New("the db is nil")
	}
	if tInfo == nil {
		return errors.New("the data is nil")
	}
	val := SoMemBlockProducerScheduleObjectByPubKey{}
	val.PubKey = tInfo.PubKey
	key, err := s.encodeMemKey("PubKey")
	if err != nil {
		return err
	}
	buf, err := proto.Marshal(&val)
	if err != nil {
		return err
	}
	err = s.dba.Put(key, buf)
	return err
}

func (s *SoBlockProducerScheduleObjectWrap) GetPubKey() []*prototype.PublicKeyType {
	res := true
	msg := &SoMemBlockProducerScheduleObjectByPubKey{}
	if s.dba == nil {
		res = false
	} else {
		key, err := s.encodeMemKey("PubKey")
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
				return msg.PubKey
			}
		}
	}
	if !res {
		var tmpValue []*prototype.PublicKeyType
		return tmpValue
	}
	return msg.PubKey
}

func (s *SoBlockProducerScheduleObjectWrap) MdPubKey(p []*prototype.PublicKeyType) bool {
	if s.dba == nil {
		return false
	}
	key, err := s.encodeMemKey("PubKey")
	if err != nil {
		return false
	}
	buf, err := s.dba.Get(key)
	if err != nil {
		return false
	}
	ori := &SoMemBlockProducerScheduleObjectByPubKey{}
	err = proto.Unmarshal(buf, ori)
	sa := &SoBlockProducerScheduleObject{}
	sa.Id = *s.mainKey
	sa.PubKey = ori.PubKey

	ori.PubKey = p
	val, err := proto.Marshal(ori)
	if err != nil {
		return false
	}
	err = s.dba.Put(key, val)
	if err != nil {
		return false
	}
	sa.PubKey = p

	return true
}

/////////////// SECTION Private function ////////////////

func (s *SoBlockProducerScheduleObjectWrap) update(sa *SoBlockProducerScheduleObject) bool {
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

func (s *SoBlockProducerScheduleObjectWrap) getBlockProducerScheduleObject() *SoBlockProducerScheduleObject {
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

	res := &SoBlockProducerScheduleObject{}
	if proto.Unmarshal(resBuf, res) != nil {
		return nil
	}
	return res
}

func (s *SoBlockProducerScheduleObjectWrap) encodeMainKey() ([]byte, error) {
	if s.mKeyBuf != nil {
		return s.mKeyBuf, nil
	}
	pre := s.getMemKeyPrefix("Id")
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

func (s *SoBlockProducerScheduleObjectWrap) delAllUniKeys(br bool, val *SoBlockProducerScheduleObject) bool {
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

func (s *SoBlockProducerScheduleObjectWrap) delUniKeysWithNames(names map[string]string, val *SoBlockProducerScheduleObject) bool {
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

func (s *SoBlockProducerScheduleObjectWrap) insertAllUniKeys(val *SoBlockProducerScheduleObject) (map[string]string, error) {
	if s.dba == nil {
		return nil, errors.New("insert uniuqe Field fail,the db is nil ")
	}
	if val == nil {
		return nil, errors.New("insert uniuqe Field fail,get the SoBlockProducerScheduleObject fail ")
	}
	sucFields := map[string]string{}
	if !s.insertUniKeyId(val) {
		return sucFields, errors.New("insert unique Field Id fail while insert table ")
	}
	sucFields["Id"] = "Id"

	return sucFields, nil
}

func (s *SoBlockProducerScheduleObjectWrap) delUniKeyId(sa *SoBlockProducerScheduleObject) bool {
	if s.dba == nil {
		return false
	}
	pre := BlockProducerScheduleObjectIdUniTable
	kList := []interface{}{pre}
	if sa != nil {

		sub := sa.Id
		kList = append(kList, sub)
	} else {
		key, err := s.encodeMemKey("Id")
		if err != nil {
			return false
		}
		buf, err := s.dba.Get(key)
		if err != nil {
			return false
		}
		ori := &SoMemBlockProducerScheduleObjectById{}
		err = proto.Unmarshal(buf, ori)
		if err != nil {
			return false
		}
		sub := ori.Id
		kList = append(kList, sub)

	}
	kBuf, err := kope.EncodeSlice(kList)
	if err != nil {
		return false
	}
	return s.dba.Delete(kBuf) == nil
}

func (s *SoBlockProducerScheduleObjectWrap) insertUniKeyId(sa *SoBlockProducerScheduleObject) bool {
	if s.dba == nil || sa == nil {
		return false
	}
	pre := BlockProducerScheduleObjectIdUniTable
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
	val := SoUniqueBlockProducerScheduleObjectById{}
	val.Id = sa.Id

	buf, err := proto.Marshal(&val)

	if err != nil {
		return false
	}

	return s.dba.Put(kBuf, buf) == nil

}

type UniBlockProducerScheduleObjectIdWrap struct {
	Dba iservices.IDatabaseRW
}

func NewUniBlockProducerScheduleObjectIdWrap(db iservices.IDatabaseRW) *UniBlockProducerScheduleObjectIdWrap {
	if db == nil {
		return nil
	}
	wrap := UniBlockProducerScheduleObjectIdWrap{Dba: db}
	return &wrap
}

func (s *UniBlockProducerScheduleObjectIdWrap) UniQueryId(start *int32) *SoBlockProducerScheduleObjectWrap {
	if start == nil || s.Dba == nil {
		return nil
	}
	pre := BlockProducerScheduleObjectIdUniTable
	kList := []interface{}{pre, start}
	bufStartkey, err := kope.EncodeSlice(kList)
	val, err := s.Dba.Get(bufStartkey)
	if err == nil {
		res := &SoUniqueBlockProducerScheduleObjectById{}
		rErr := proto.Unmarshal(val, res)
		if rErr == nil {
			wrap := NewSoBlockProducerScheduleObjectWrap(s.Dba, &res.Id)
			return wrap
		}
	}
	return nil
}
