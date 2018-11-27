package economist

import (
	"github.com/coschain/contentos-go/app/table"
	"github.com/coschain/contentos-go/common"
	"github.com/coschain/contentos-go/common/constants"
	"github.com/coschain/contentos-go/iservices"
	"github.com/coschain/contentos-go/node"
	"github.com/coschain/contentos-go/prototype"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

var (
	// fixme: the single id should be share with service
	SINGLE_ID int32 = 1
)

type Economist struct {
	ctx               *node.ServiceContext
	db                iservices.IDatabaseService
	rewardAccumulator uint64 // reward accumulator
	vpAccumulator     uint64 // vote power accumulator
	globalProps       *prototype.DynamicProperties
	rewardsKeeper     *prototype.InternalRewardsKeeper
}

func (e *Economist) getDb() (iservices.IDatabaseService, error) {
	s, err := e.ctx.Service("db")
	if err != nil {
		return nil, err
	}
	db := s.(iservices.IDatabaseService)
	return db, nil
}

func New(ctx *node.ServiceContext) (*Economist, error) {

	return &Economist{ctx: ctx}, nil
}

func (e *Economist) Start(node *node.Node) error {
	db, err := e.getDb()
	if err != nil {
		return errors.New("Economist fetch db service error")
	}
	e.db = db
	dgpWrap := table.NewSoGlobalWrap(e.db, &SINGLE_ID)
	if !dgpWrap.CheckExist() {
		return errors.New("the mainkey is already exist")
	}
	e.globalProps = dgpWrap.GetProps()

	keeperWrap := table.NewSoRewardsKeeperWrap(e.db, &SINGLE_ID)
	if !keeperWrap.CheckExist() {
		return errors.New("Economist access rewards keeper error")
	}
	e.rewardsKeeper = keeperWrap.GetKeeper()
	return nil
}

func (e *Economist) Stop() error {
	return nil
}

func (e *Economist) updateRewardsKeeper() error {
	keeper := table.NewSoRewardsKeeperWrap(e.db, &SINGLE_ID)
	success := keeper.MdKeeper(e.rewardsKeeper)
	if !success {
		return errors.New("flush rewards keeper into db error")
	}
	return nil
}

func (e *Economist) Do() error {
	e.decayGlobalVotePower()
	timestamp := e.globalProps.Time.UtcSeconds - uint32(constants.GenesisTime)
	keyPrefix := "cashout:" + string(common.GetBucket(timestamp)) + "_"
	postCashoutList := []string{}
	replyCashoutList := []string{}
	r := regexp.MustCompile(`cashout:(?P<bucket>\d+)_(?P<idx>\d+)`)
	for iter := e.db.NewIterator([]byte(keyPrefix), nil); iter.Valid(); iter.Next() {
		key, err := iter.Key()
		if err != nil {
			return err
		}
		value, err := iter.Value()
		if err != nil {
			return err
		}
		match := r.FindStringSubmatch(string(key))
		if len(match) > 0 {
			idx := match[2]
			switch string(value) {
			case "post":
				postCashoutList = append(postCashoutList, idx)
			case "reply":
				replyCashoutList = append(replyCashoutList, idx)
			}
		}
	}
	if len(postCashoutList) > 0 {
		e.postCashout(postCashoutList)
	}

	if len(postCashoutList) > 0 {
		e.replyCashout(replyCashoutList)
	}

	err := e.updateRewardsKeeper()
	return err
}

func (e *Economist) decayGlobalVotePower() {
	e.globalProps.WeightedVps -= e.globalProps.WeightedVps * constants.BLOCK_INTERVAL / constants.VP_DECAY_TIME
}

func (e *Economist) postCashout(pids []string) {
	posts := []*table.SoPostWrap{}
	var vpAccumulator uint64 = 0
	for _, pidStr := range pids {
		pid, _ := strconv.ParseUint(pidStr, 10, 64)
		post := table.NewSoPostWrap(e.db, &pid)
		vpAccumulator += post.GetWeightedVp()
		posts = append(posts, post)
	}
	blockReward := vpAccumulator * e.globalProps.PostRewards.Value / e.globalProps.WeightedVps
	for _, post := range posts {
		author := post.GetAuthor().Value
		reward := post.GetWeightedVp() * blockReward / vpAccumulator
		if vest, ok := e.rewardsKeeper.Rewards[author]; !ok {
			e.rewardsKeeper.Rewards[author] = &prototype.Vest{Value: reward}
		} else {
			vest.Value += reward
		}
	}
}

// use same algorithm to simplify
func (e *Economist) replyCashout(rids []string) {
	replies := []*table.SoPostWrap{}
	var vpAccumulator uint64 = 0
	for _, pidStr := range rids {
		pid, _ := strconv.ParseUint(pidStr, 10, 64)
		reply := table.NewSoPostWrap(e.db, &pid)
		vpAccumulator += reply.GetWeightedVp()
		replies = append(replies, reply)
	}
	blockReward := vpAccumulator * e.globalProps.ReplyRewards.Value / e.globalProps.WeightedVps
	for _, reply := range replies {
		author := reply.GetAuthor().Value
		reward := reply.GetWeightedVp() * blockReward / vpAccumulator
		if vest, ok := e.rewardsKeeper.Rewards[author]; !ok {
			e.rewardsKeeper.Rewards[author] = &prototype.Vest{Value: reward}
		} else {
			vest.Value += reward
		}
	}
}