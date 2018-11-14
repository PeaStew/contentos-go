package app

import (
	"bytes"
	"github.com/asaskevich/EventBus"
	"github.com/coschain/contentos-go/app/table"
	"github.com/coschain/contentos-go/common"
	"github.com/coschain/contentos-go/common/constants"
	"github.com/coschain/contentos-go/common/eventloop"
	"github.com/coschain/contentos-go/iservices"
	"github.com/coschain/contentos-go/node"
	"github.com/coschain/contentos-go/prototype"
	"github.com/golang/protobuf/proto"
	"time"
)

type skipFlag uint32

const (
	skip_nothing                skipFlag = 0
	skip_transaction_signatures skipFlag = 1 << 0
	skip_apply_transaction      skipFlag = 1 << 1
)

type Controller struct {
	iservices.IController
	// lock for db write
	// pending_trx_list
	// DB Manager
	ctx    *node.ServiceContext
	evLoop *eventloop.EventLoop

	db      iservices.IDatabaseService
	noticer EventBus.Bus
	skip    skipFlag

	_pending_tx        []*prototype.TransactionWrapper
	_isProducing       bool
	_currentTrxId      *prototype.Sha256
	_current_op_in_trx uint16
}

func (c *Controller) getDb() (iservices.IDatabaseService,error) {
	s, err := c.ctx.Service(iservices.DB_SERVER_NAME)
	if err != nil {
		return nil, err
	}
	db := s.(iservices.IDatabaseService)
	return db, nil
}

// for easy test
func (c *Controller) SetDB(db iservices.IDatabaseService) {
	c.db = db
}

// service constructor
func NewController(ctx *node.ServiceContext) (*Controller, error) {
	return &Controller{ctx: ctx}, nil
}

func (c *Controller) Start(node *node.Node) error {
	db,err := c.getDb()
	if err != nil {
		return err
	}
	c.db = db
	c.evLoop = node.MainLoop
	c.noticer = node.EvBus

	c.Open()
	return nil
}

func (c *Controller) Open() {
	var i int32 = 0
	dgpWrap := table.NewSoDynamicGlobalPropertiesWrap(c.db,&i)
	if !dgpWrap.CheckExist() {
		c.initGenesis()
	}
}

func (c *Controller) Stop() error {
	return nil
}

func (c *Controller) setProducing(b bool) {
	c._isProducing = b
}

func (c *Controller) PushTrx(trx *prototype.SignedTransaction) *prototype.TransactionInvoice {
	// this function may be cross routines ? use channel or lock ?
	oldSkip := c.skip
	defer func() {
		c.setProducing(false)
		c.skip = oldSkip
	}()

	// check maximum_block_size
	var i int32 = 0
	dgpWrap := table.NewSoDynamicGlobalPropertiesWrap(c.db,&i)
	if  proto.Size(trx) > int(dgpWrap.GetMaximumBlockSize() - 256) {
		panic("transaction is too large")
	}

	c.setProducing(true)
	return c._pushTrx(trx)
}

func (c *Controller) _pushTrx(trx *prototype.SignedTransaction) *prototype.TransactionInvoice {
	defer func() {
		// undo sub session
		if err := recover(); err != nil {
			c.db.EndTransaction(false)
			panic(err)
		}
	}()
	// start a new undo session when first transaction come after push block
	if len(c._pending_tx) == 0 {
		c.db.BeginTransaction()
	}

	trxWrp := &prototype.TransactionWrapper{}
	trxWrp.SigTrx = trx

	// start a sub undo session for applyTransaction
	c.db.BeginTransaction()

	c._applyTransaction(trxWrp)
	c._pending_tx = append(c._pending_tx, trxWrp)

	// commit sub session
	c.db.EndTransaction(true)

	c.NotifyTrxPending(trx)
	return trxWrp.Invoice
}

func (c *Controller) PushBlock(blk *prototype.SignedBlock) {

}

func (c *Controller) GenerateBlock(accountName string, timestamp uint32,
					 prev common.BlockID) *prototype.SignedBlock {
	return nil
}

func (c *Controller) NotifyOpPostExecute(on *prototype.OperationNotification) {
	c.noticer.Publish(constants.NOTICE_OP_POST, on)
}

func (c *Controller) NotifyOpPreExecute(on *prototype.OperationNotification) {
	c.noticer.Publish(constants.NOTICE_OP_PRE, on)
}

func (c *Controller) NotifyTrxPreExecute(trx *prototype.SignedTransaction) {
	c.noticer.Publish(constants.NOTICE_TRX_PRE, trx)
}

func (c *Controller) NotifyTrxPostExecute(trx *prototype.SignedTransaction) {
	c.noticer.Publish(constants.NOTICE_TRX_POST, trx)
}

func (c *Controller) NotifyTrxPending(trx *prototype.SignedTransaction) {
	c.noticer.Publish(constants.NOTICE_TRX_PENDING, trx)
}

func (c *Controller) NotifyBlockApply(block *prototype.SignedBlock) {
	c.noticer.Publish(constants.NOTICE_BLOCK_APPLY, block)
}

// calculate reward for creator and witness
func (c *Controller) processBlock() {
}

func (c *Controller) _applyTransaction(trxWrp *prototype.TransactionWrapper) {
	defer func() {
		if err := recover(); err != nil {
			trxWrp.Invoice.Status = 500
			panic("_applyTransaction failed")
		} else {
			trxWrp.Invoice.Status = 200
			return
		}
	}()

	trx := trxWrp.SigTrx
	var err error
	c._currentTrxId, err = trx.Id()
	if err != nil {
		panic("get trx id failed")
	}

	trx.Validate()

	// trx duplicate check
	transactionObjWrap := table.NewSoTransactionObjectWrap(c.db, c._currentTrxId)
	if transactionObjWrap.CheckExist() {
		panic("Duplicate transaction check failed")
	}

	if c.skip&skip_transaction_signatures == 0 {
		postingGetter := func(name string) *prototype.Authority {
			account := &prototype.AccountName{Value:name}
			authWrap := table.NewSoAccountAuthorityObjectWrap(c.db,account)
			auth := authWrap.GetPosting()
			if auth == nil {
				panic("no posting auth")
			}
			return auth
		}
		activeGetter := func(name string) *prototype.Authority {
			account := &prototype.AccountName{Value:name}
			authWrap := table.NewSoAccountAuthorityObjectWrap(c.db,account)
			auth := authWrap.GetPosting()
			if auth == nil {
				panic("no posting auth")
			}
			return auth
		}
		ownerGetter := func(name string) *prototype.Authority {
			account := &prototype.AccountName{Value:name}
			authWrap := table.NewSoAccountAuthorityObjectWrap(c.db,account)
			auth := authWrap.GetPosting()
			if auth == nil {
				panic("no posting auth")
			}
			return auth
		}

		tmpChainId := prototype.ChainId{Value: 0}
		trx.VerifyAuthority(tmpChainId, 2,postingGetter,activeGetter,ownerGetter)
		// @ check_admin
	}

	// TaPos and expired check
	var i int32 = 0
	dgpWrap := table.NewSoDynamicGlobalPropertiesWrap(c.db,&i)
	blockNum := dgpWrap.GetHeadBlockNumber()
	if blockNum > 0 {
		uniWrap := table.UniBlockSummaryObjectIdWrap{}
		idWrap := uniWrap.UniQueryId(&trx.Trx.RefBlockNum)
		if !idWrap.CheckExist() {
			panic("no refBlockNum founded")
		} else {
			blockId := idWrap.GetBlockId()
			summaryId := uint32(blockId.Hash[1])
			if trx.Trx.RefBlockPrefix != summaryId {
				panic("transaction tapos failed")
			}
		}
		// get head time
		if trx.Trx.Expiration.UtcSeconds > uint32(time.Now().Second()+30) {
			panic("transaction expiration too long")
		}
		if uint32(time.Now().Second()) > trx.Trx.Expiration.UtcSeconds {
			panic("transaction has expired")
		}
	}

	// insert trx into DB unique table
	obj := &table.SoTransactionObject{}
	obj.TrxId = c._currentTrxId
	obj.Expiration = &prototype.TimePointSec{UtcSeconds: 100}
	if !transactionObjWrap.CreateTransactionObject(obj) {
		panic("create transactionObject failed")
	}

	c.NotifyTrxPreExecute(trx)

	// process operation
	c._current_op_in_trx = 0
	for _, op := range trx.Trx.Operations {
		c.applyOperation(op)
		c._current_op_in_trx++
	}

	c._currentTrxId = &prototype.Sha256{}
}

func (c *Controller) applyOperation(op *prototype.Operation) {
	n := &prototype.OperationNotification{Op: op}
	c.NotifyOpPreExecute(n)
	eva := getEvaluator(op)
	eva.Apply(op)
	c.NotifyOpPostExecute(n)
}

func getEvaluator(op *prototype.Operation) BaseEvaluator {
	switch op.Op.(type) {
	case *prototype.Operation_Op1:
		return BaseEvaluator(&AccountCreateEvaluator{})
	case *prototype.Operation_Op2:
		return BaseEvaluator(&TransferEvaluator{})
	default:
		panic("no matchable evaluator")
	}
}

func (c *Controller) applyBlock(blk *prototype.SignedBlock) {
	oldFlag := c.skip
	defer func() {
		c.skip = oldFlag
	}()

	c._applyBlock(blk)

	// @ tps update
}

func (c *Controller) _applyBlock(blk *prototype.SignedBlock) {
	//nextBlockNum := blk.Id().BlockNum()

	root := blk.CalculateMerkleRoot()
	if !bytes.Equal(root.Data[:], blk.SignedHeader.Header.TransactionMerkleRoot.Hash) {
		panic("Merkle check failed")
	}
}

func (c *Controller) initGenesis() {

	// create initminer
	pubKey , _ := prototype.PublicKeyFromWIF(constants.INITMINER_PUBKEY)
	name := &prototype.AccountName{Value:constants.INIT_MINER_NAME}
	newAccountWrap := table.NewSoAccountWrap(c.db,name)
	newAccount := &table.SoAccount{}
	newAccount.Name = name
	newAccount.PubKey = pubKey
	newAccount.CreatedTime = &prototype.TimePointSec{UtcSeconds:0}
	cos := &prototype.Coin{Amount:&prototype.Safe64{Value:constants.INIT_SUPPLY}}
	vest := &prototype.Vest{Amount:&prototype.Safe64{Value:0}}
	newAccount.Balance = cos
	newAccount.VestingShares = vest
	if !newAccountWrap.CreateAccount(newAccount) {
		panic("CreateAccount error")
	}

	// create account authority
	authorityWrap := table.NewSoAccountAuthorityObjectWrap(c.db,name)
	authority := &table.SoAccountAuthorityObject{}
	authority.Account = name

	ownerAuth := &prototype.Authority{
		WeightThreshold: 1,
		KeyAuths: []*prototype.KvKeyAuth{
			&prototype.KvKeyAuth{
				Key: pubKey,
				Weight: 1,
			},
		},
	}
	authority.Posting = ownerAuth
	authority.Active = ownerAuth
	authority.Owner = ownerAuth
	if !authorityWrap.CreateAccountAuthorityObject(authority) {
		panic("CreateAccountAuthorityObject error ")
	}
	// @ create witness_object

	// create dynamic global properties
	var i int32 = 0
	dgpWrap := table.NewSoDynamicGlobalPropertiesWrap(c.db,&i)
	dgp := &table.SoDynamicGlobalProperties{}
	dgp.CurrentWitness = name
	dgp.Time = &prototype.TimePointSec{UtcSeconds:constants.GENESIS_TIME}
	// @ recent_slots_filled
	// @ participation_count
	dgp.CurrentSupply = cos
	dgp.TotalCos = cos
	dgp.MaximumBlockSize = constants.MAX_BLOCK_SIZE
	dgp.TotalVestingShares = &prototype.Vest{Amount:&prototype.Safe64{Value:0}}
	if !dgpWrap.CreateDynamicGlobalProperties(dgp) {
		panic("CreateDynamicGlobalProperties error")
	}

	// create block summary
	for i := uint32(0); i < 0x10000; i++ {
		wrap := table.NewSoBlockSummaryObjectWrap(c.db, &i)
		obj := &table.SoBlockSummaryObject{}
		obj.Id = i
		if !wrap.CreateBlockSummaryObject(obj) {
			panic("CreateBlockSummaryObject error")
		}
	}
}

func (c *Controller) CreateVesting(accountName *prototype.AccountName, cos *prototype.Coin) *prototype.Vest {

	newVesting := prototype.CosToVesting(cos)
	creatorWrap := table.NewSoAccountWrap(c.db,accountName)
	oldVesting := creatorWrap.GetVestingShares()
	oldVesting.Amount.Value += newVesting.Amount.Value
	creatorWrap.MdVestingShares(*oldVesting)

	var i int32 = 0
	dgpWrap := table.NewSoDynamicGlobalPropertiesWrap(c.db,&i)
	originTotal := dgpWrap.GetTotalVestingShares()
	originTotal.Amount.Value += newVesting.Amount.Value
	dgpWrap.MdTotalVestingShares(*originTotal)
	return newVesting
}

func (c *Controller) SubBalance(accountName *prototype.AccountName, cos *prototype.Coin) {
	accountWrap := table.NewSoAccountWrap(c.db,accountName)
	originBalance := accountWrap.GetBalance()
	originBalance.Amount.Value -= cos.Amount.Value
	accountWrap.MdBalance(*originBalance)

	// dynamic glaobal properties
	var i int32 = 0
	dgpWrap := table.NewSoDynamicGlobalPropertiesWrap(c.db,&i)
	originTotal := dgpWrap.GetTotalCos()
	originTotal.Amount.Value -= cos.Amount.Value
	dgpWrap.MdTotalCos(*originTotal)
}

func (c *Controller) AddBalance(accountName *prototype.AccountName, cos *prototype.Coin) {
	accountWrap := table.NewSoAccountWrap(c.db,accountName)
	originBalance := accountWrap.GetBalance()
	originBalance.Amount.Value += cos.Amount.Value
	accountWrap.MdBalance(*originBalance)

	// dynamic glaobal properties
	var i int32 = 0
	dgpWrap := table.NewSoDynamicGlobalPropertiesWrap(c.db,&i)
	originTotal := dgpWrap.GetTotalCos()
	originTotal.Amount.Value += cos.Amount.Value
	dgpWrap.MdTotalCos(*originTotal)
}
