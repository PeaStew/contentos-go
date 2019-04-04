package app

import (
	"errors"
	"fmt"
	"github.com/coschain/contentos-go/common/constants"
	"github.com/coschain/contentos-go/iservices"
	"github.com/coschain/contentos-go/prototype"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

type TrxCallback func(result *prototype.EstimateTrxResult)

type TrxEntry struct {
	result    *prototype.EstimateTrxResult
	sig       string
	size      int
	signer    string
	signerKey *prototype.PublicKeyType
	callback  TrxCallback
}

func NewTrxMgrEntry(trx *prototype.SignedTransaction, callback TrxCallback) *TrxEntry {
	return &TrxEntry{
		result: &prototype.EstimateTrxResult{
			SigTrx:  trx,
			Receipt: &prototype.TransactionReceiptWithInfo{Status: prototype.StatusSuccess},
		},
		callback: callback,
	}
}

func (e *TrxEntry) SetError(err error) error {
	e.result.Receipt.Status = prototype.StatusError
	e.result.Receipt.ErrorInfo = err.Error()
	return err
}

func (e *TrxEntry) Deliver() {
	if e.callback != nil {
		e.callback(e.result)
	}
}

func (e *TrxEntry) InitCheck() error {
	trx := e.result.SigTrx
	if err := trx.Validate(); err != nil {
		return e.SetError(err)
	}
	e.sig = string(trx.Signature.Sig)
	e.size = proto.Size(trx)
	if e.size > constants.MaxTransactionSize {
		return e.SetError(fmt.Errorf("trx too large, size = %d > %d", e.size, constants.MaxTransactionSize))
	}
	creator := ""
	if creators := trx.GetOpCreatorsMap(); len(creators) != 1 {
		return e.SetError(fmt.Errorf("non-unique trx creators, found %d", len(creators)))
	} else {
		for creator = range creators {
			break
		}
	}
	e.signer = creator
	if signKey, err := trx.ExportPubKeys(prototype.ChainId{Value: 0}); err != nil {
		return e.SetError(fmt.Errorf("cannot export signing key: %s", err.Error()))
	} else {
		e.signerKey = signKey
	}
	return nil
}

func (e *TrxEntry) CheckExpiration(blockTime uint32) error {
	expiration := e.result.SigTrx.GetTrx().GetExpiration().GetUtcSeconds()
	if expiration < blockTime {
		return e.SetError(fmt.Errorf("trx expired, %d < %d", expiration, blockTime))
	}
	if expiration > blockTime + constants.TrxMaxExpirationTime {
		return e.SetError(fmt.Errorf("trx expiration too long, %d > %d + %d", expiration, blockTime, constants.TrxMaxExpirationTime))
	}
	return nil
}

func (e *TrxEntry) CheckTapos(checker *TaposChecker) error {
	if err := checker.Check(e.result.SigTrx.Trx); err != nil {
		return e.SetError(fmt.Errorf("tapos failed: %s", err.Error()))
	}
	return nil
}

func (e *TrxEntry) CheckSignerKey(fetcher *AuthFetcher) error {
	if err := fetcher.CheckPublicKey(e.signer, e.signerKey); err != nil {
		return e.SetError(fmt.Errorf("signature failed: %s", err.Error()))
	}
	return nil
}

func (e *TrxEntry) CheckInBlockTrxs(checker *InBlockTrxChecker) error {
	if checker.Has(e.result.SigTrx) {
		return e.SetError(errors.New("found duplicate in-block trx"))
	}
	return nil
}

func (e *TrxEntry) GetTrxResult() *prototype.EstimateTrxResult {
	return e.result
}
func (e *TrxEntry) GetTrxSize() int {
	return e.size
}
func (e *TrxEntry) GetTrxSigner() string {
	return e.signer
}
func (e *TrxEntry) GetTrxSigningKey() *prototype.PublicKeyType {
	return e.signerKey
}

const (
	sMaxWaitingCount  = constants.TrxMaxExpirationTime * 20000
)

type ITrxMgrPlugin interface {
	BlockApplied(b *prototype.SignedBlock)
	BlockReverted(blockNum uint64)
	BlockCommitted(blockNum uint64)
}

type TrxMgr struct {
	db 				iservices.IDatabaseRW
	log             *logrus.Logger
	headTime		uint32
	waiting 		map[string]*TrxEntry
	waitingLock 	sync.RWMutex
	fetched 		map[string]*TrxEntry
	fetchedLock 	sync.RWMutex
	auth 			*AuthFetcher
	tapos 			*TaposChecker
	history 		*InBlockTrxChecker
	plugins         []ITrxMgrPlugin
}

func NewTrxMgr(db iservices.IDatabaseRW, logger *logrus.Logger, lastBlock, commitBlock uint64) *TrxMgr {
	auth := NewAuthFetcher(db, logger, lastBlock, commitBlock)
	tapos := NewTaposChecker(db, logger, lastBlock)
	history := NewInBlockTrxChecker(db, logger, lastBlock)
	return &TrxMgr{
		db:       db,
		log:      logger,
		headTime: (&DynamicGlobalPropsRW{db:db}).GetProps().GetTime().GetUtcSeconds(),
		waiting:  make(map[string]*TrxEntry),
		fetched:  make(map[string]*TrxEntry),
		auth: auth,
		tapos: tapos,
		history: history,
		plugins: []ITrxMgrPlugin{ auth, tapos, history },
	}
}

func (m *TrxMgr) AddTrx(trx *prototype.SignedTransaction, callback TrxCallback) error {
	entry := NewTrxMgrEntry(trx, callback)
	if trx == nil || trx.Signature == nil {
		err := entry.SetError(errors.New("invalid trx"))
		m.deliverEntry(entry)
		return err
	}
	if m.isProcessingTrx(trx) != nil {
		err := entry.SetError(errors.New("trx already in process"))
		m.deliverEntry(entry)
		return err
	}
	c := make(chan error)
	go func() {
		ok := false
		if entry.InitCheck() != nil || m.checkTrx(entry, atomic.LoadUint32(&m.headTime)) != nil {
			m.deliverEntry(entry)
		} else {
			m.waitingLock.Lock()
			m.fetchedLock.RLock()

			ok = m.addToWaiting(entry) > 0

			m.fetchedLock.RUnlock()
			m.waitingLock.Unlock()
		}
		if !ok {
			c <- errors.New(entry.result.Receipt.ErrorInfo)
		} else {
			c <- nil
		}
	}()
	return <-c
}

func (m *TrxMgr) WaitingCount() int {
	m.waitingLock.RLock()
	defer m.waitingLock.RUnlock()
	return len(m.waiting)
}

func (m *TrxMgr) FetchTrx(blockTime uint32, maxCount, maxSize int) (entries []*TrxEntry) {
	m.log.Debugf("TRXMGR: FetchTrx begin: maxCount=%d, maxSize=%d", maxCount, maxSize)
	t0 := time.Now()
	m.waitingLock.Lock()
	defer m.waitingLock.Unlock()

	m.fetchedLock.Lock()
	defer m.fetchedLock.Unlock()

	t1 := time.Now()
	counter, size := 0, 0
	for s, e := range m.waiting {
		if maxCount > 0 && counter >= maxCount {
			break
		}
		if maxSize > 0 && size >= maxSize {
			break
		}
		if m.checkTrx(e, blockTime) != nil {
			m.deliverEntry(e)
			continue
		}
		entries = append(entries, e)
		m.fetched[s] = e
		delete(m.waiting, s)
		counter++
		size += e.size
	}
	t2 := time.Now()
	m.log.Debugf("TRXMGR: FetchTrx end: maxCount=%d, maxSize=%d, count=%d, size=%d, %v|%v|%v",
		maxCount, maxSize, counter, size,
		t2.Sub(t0), t1.Sub(t0), t2.Sub(t1))
	return
}

func (m *TrxMgr) ReturnTrx(failed bool, entries ...*TrxEntry) {
	m.log.Debug("TRXMGR: ReturnTrx begin")
	dispatch := m.deliverEntry
	t0 := time.Now()
	if !failed {
		m.waitingLock.Lock()
		defer m.waitingLock.Unlock()
		dispatch = func(e *TrxEntry) {
			m.addToWaiting(e)
		}
	}
	m.fetchedLock.Lock()
	defer m.fetchedLock.Unlock()

	t1 := time.Now()
	for _, e := range entries {
		s := string(e.result.SigTrx.Signature.Sig)
		f := m.fetched[s]
		if f != nil {
			dispatch(f)
			delete(m.fetched, s)
		}
	}
	t2 := time.Now()
	m.log.Debug("TRXMGR: ReturnTrx end: #tx=%d, %v|%v|%v", len(entries), t2.Sub(t0), t1.Sub(t0), t2.Sub(t1))
}

func (m *TrxMgr) CheckBlockTrxs(b *prototype.SignedBlock) (entries []*TrxEntry, err error) {
	m.log.Debugf("TRXMGR: CheckBlockTrxs begin %d", b.SignedHeader.Number())
	t0 := time.Now()
	if count := len(b.Transactions); count > 0 {
		blockTime := b.SignedHeader.Header.Timestamp.UtcSeconds
		errs := make([]error, count)
		entries = make([]*TrxEntry, count)
		errIdx := int32(-1)
		var wg sync.WaitGroup
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func(idx int) {
				defer wg.Done()
				var err error
				trx := b.Transactions[idx].SigTrx
				e := NewTrxMgrEntry(trx, nil)
				needInitCheck := true
				if ptrx := m.isProcessingTrx(trx); ptrx != nil {
					needInitCheck = false
					e.sig = ptrx.sig
					e.size = ptrx.size
					e.signer = ptrx.signer
					e.signerKey = ptrx.signerKey
				}
				if needInitCheck {
					err = e.InitCheck()
				}
				if err == nil {
					err = m.checkTrx(e, blockTime)
				}
				if err != nil {
					errs[idx] = err
					atomic.CompareAndSwapInt32(&errIdx, -1, int32(idx))
				} else {
					entries[idx] = e
				}
			}(i)
		}
		wg.Wait()
		if errIdx >= 0 {
			entries = nil
			err = fmt.Errorf("block %d trxs[%d] check failed: %s", b.SignedHeader.Number(), errIdx, errs[errIdx].Error())
		}
	}
	t1 := time.Now()
	m.log.Debugf("TRXMGR: CheckBlockTrxs end %d: #tx=%d, %v", b.SignedHeader.Number(), len(b.Transactions), t1.Sub(t0))
	return
}

func (m *TrxMgr) BlockApplied(b *prototype.SignedBlock) {
	m.log.Debugf("TRXMGR: BlockApplied begin %d", b.SignedHeader.Number())
	t0 := time.Now()
	atomic.StoreUint32(&m.headTime, b.SignedHeader.Header.Timestamp.UtcSeconds)

	m.waitingLock.Lock()
	m.fetchedLock.Lock()
	t1 := time.Now()
	for _, txw := range b.Transactions {
		s := string(txw.SigTrx.Signature.Sig)
		if e := m.fetched[s]; e != nil {
			m.deliverEntry(e)
			delete(m.fetched, s)
		}
		if e := m.waiting[s]; e != nil {
			m.deliverEntry(e)
			delete(m.waiting, s)
		}
	}
	t2 := time.Now()
	m.fetchedLock.Unlock()
	m.waitingLock.Unlock()

	m.callPlugins(func(plugin ITrxMgrPlugin) {
		plugin.BlockApplied(b)
	})
	t3 := time.Now()
	m.log.Debugf("TRXMGR: BlockApplied end %d: #tx=%d, %v|%v|%v|%v", b.SignedHeader.Number(), len(b.Transactions), t3.Sub(t0), t1.Sub(t0), t2.Sub(t1), t3.Sub(t2))
	m.log.Debugf("TRXMGR: auth-hit=%v", m.auth.HitRate())
}

func (m *TrxMgr) BlockCommitted(blockNum uint64) {
	m.log.Debugf("TRXMGR: BlockCommitted begin %d", blockNum)
	t0 := time.Now()
	m.callPlugins(func(plugin ITrxMgrPlugin) {
		plugin.BlockCommitted(blockNum)
	})
	t1 := time.Now()
	m.log.Debugf("TRXMGR: BlockCommitted end %d: %v", blockNum, t1.Sub(t0))
}

func (m *TrxMgr) BlockReverted(blockNum uint64) {
	m.log.Debugf("TRXMGR: BlockReverted begin %d", blockNum)
	t0 := time.Now()
	m.callPlugins(func(plugin ITrxMgrPlugin) {
		plugin.BlockReverted(blockNum)
	})
	t1 := time.Now()
	m.log.Debugf("TRXMGR: BlockReverted end %d: %v", blockNum, t1.Sub(t0))
}

func (m *TrxMgr) addToWaiting(entries...*TrxEntry) (count int) {
	for _, e := range entries {
		if len(m.waiting) > sMaxWaitingCount {
			_ = e.SetError(errors.New("too many waiting trxs"))
			m.deliverEntry(e)
			continue
		}
		if m.isProcessingNoLock(e.result.SigTrx) != nil {
			_ = e.SetError(errors.New("trx already in process"))
			m.deliverEntry(e)
			continue
		}
		m.waiting[e.sig] = e
		count++
	}
	return
}

func (m *TrxMgr) isProcessingTrx(trx *prototype.SignedTransaction) *TrxEntry {
	m.waitingLock.RLock()
	defer m.waitingLock.RUnlock()
	m.fetchedLock.RLock()
	defer m.fetchedLock.RUnlock()
	return m.isProcessingNoLock(trx)
}

func (m *TrxMgr) isProcessingNoLock(trx *prototype.SignedTransaction) *TrxEntry {
	s := string(trx.Signature.Sig)
	if e := m.waiting[s]; e != nil {
		return e
	}
	return m.fetched[s]
}

func (m *TrxMgr) checkTrx(e *TrxEntry, blockTime uint32) (err error) {
	if err = e.CheckExpiration(blockTime); err != nil {
		return err
	} else if err = e.CheckTapos(m.tapos); err != nil {
		return err
	} else if err = e.CheckSignerKey(m.auth); err != nil {
		return err
	} else if err = e.CheckInBlockTrxs(m.history); err != nil {
		return err
	}
	return
}

func (m *TrxMgr) deliverEntry(e *TrxEntry) {
	go func() {
		e.Deliver()
	}()
}

func (m *TrxMgr) callPlugins(f func(plugin ITrxMgrPlugin)) {
	var wg sync.WaitGroup
	wg.Add(len(m.plugins))
	for i := range m.plugins {
		go func(idx int) {
			defer wg.Done()
			f(m.plugins[idx])
		}(i)
	}
	wg.Wait()
}