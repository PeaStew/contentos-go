package consensus

import (
	"bytes"
	"sync"
	"time"

	"github.com/coschain/contentos-go/common"
	"github.com/coschain/contentos-go/common/prototype"
	"github.com/coschain/contentos-go/config"
	"github.com/coschain/contentos-go/db/forkdb"
)

type Producer struct {
	Name   string
	PubKey *prototype.PublicKeyType
	Weight uint32
}

func (p *Producer) Produce() (common.ISignedBlock, error) {
	return nil, nil
}

type DPoS struct {
	ForkDB     *forkdb.DB
	Producers  []*Producer
	activeNum  uint32
	currentIdx uint32

	stopCh chan struct{}
	wg     sync.WaitGroup
	sync.RWMutex
}

func NewDPoS() *DPoS {
	return &DPoS{
		ForkDB:    forkdb.NewDB(),
		Producers: make([]*Producer, config.ProducerNum),
		stopCh:    make(chan struct{}),
	}
}

func (d *DPoS) CurrentProducer() *Producer {
	d.RLock()
	defer d.RUnlock()
	return d.Producers[0]
}

// Called when a produce round complete, it adds new producers,
// remove unqualified producers and shuffle the block-producing order
func (d *DPoS) shuffle() {}

func (d *DPoS) ActiveProducers() []*Producer {
	d.RLock()
	defer d.RUnlock()
	return d.Producers
}

func (d *DPoS) Start() {
	go d.start()
}

func (d *DPoS) start() {
	d.wg.Add(1)
	defer d.wg.Done()
	for {
		select {
		case <-d.stopCh:
			break
		default:
			if !d.checkGenesis() || !d.checkOurTurn(){
				continue
			}
			b, err := d.GenerateBlock()
			if err != nil {
				d.AddBlock(b)
				// TODO: broadcast block
			}
		}
	}
}

func (d *DPoS) Stop() {
	close(d.stopCh)
	d.wg.Wait()
}

func (d *DPoS) GenerateBlock() (common.ISignedBlock, error) {

	return nil, nil
}

func (d *DPoS) checkGenesis() bool {
	now := time.Now()
	genesisTime := time.Unix(config.GenesisTime, 0)
	if now.After(genesisTime) || now.Equal(genesisTime) {
		return true
	}

	ceil := now.Round(time.Second)
	if ceil.Before(now) {
		ceil = ceil.Add(time.Second)
	}

	if ceil.Before(genesisTime) {
		time.Sleep(ceil.Sub(now))
		return false
	}

	time.Sleep(ceil.Sub(now))
	return true
}

func (d *DPoS) checkOurTurn() bool {
	return true
}

func (d *DPoS) PushBlock(b common.ISignedBlock) error {
	head := d.ForkDB.Head()
	newHead :=d.ForkDB.PushBlock(b)

	if newHead == head {
		// this implies that b is a:
		// 1. detached block or
		// 2. out of range block or
		// 3. head of a non-main branch or
		// 4. illegal block
		return nil
	} else if newHead.Previous() != head.Id() {
		// TODO: swith fork
		return nil
	}

	if err := d.applyBlock(b); err != nil {
		// the block is illegal
		d.ForkDB.MarkAsIllegal(b.Id())
		d.ForkDB.Pop()
		return err
	}

	if bytes.Equal(b.GetSignee().Data, d.Producers[d.activeNum-1].PubKey.Data) {
		d.shuffle()
	}
	return nil
}

func (d *DPoS) RemoveBlock(id common.BlockID) {
	d.ForkDB.Remove(id)
}

func (d *DPoS) ForkRoot(fork1, fork2 common.BlockID) common.BlockID {
	return common.BlockID{}
}

func (d *DPoS) applyBlock(b common.ISignedBlock) error {
	// TODO: state db apply
	return nil
}

func (d *DPoS) popBlock() error {
	// TODO: state db revert
	return nil
}
