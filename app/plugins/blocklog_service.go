package plugins

import (
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/coschain/contentos-go/app/blocklog"
	"github.com/coschain/contentos-go/common"
	"github.com/coschain/contentos-go/common/constants"
	"github.com/coschain/contentos-go/iservices"
	"github.com/coschain/contentos-go/iservices/service-configs"
	"github.com/coschain/contentos-go/node"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"time"
)

type BlockLogService struct {
	ctx *node.ServiceContext
	config *service_configs.DatabaseConfig
	logger *logrus.Logger
	bus EventBus.Bus
	db *gorm.DB
}

func NewBlockLogService(ctx *node.ServiceContext, config *service_configs.DatabaseConfig, logger *logrus.Logger) (*BlockLogService, error) {
	return &BlockLogService{ctx:ctx, config:config, logger:logger}, nil
}

func (s *BlockLogService) Start(node *node.Node) error  {
	if err := s.initDatabase(); err != nil {
		return err
	}
	s.bus = node.EvBus
	_ = s.bus.Subscribe(constants.NoticeBlockLog, s.onBlockLog)
	_ = s.bus.Subscribe(constants.NoticeLibChange, s.onLibChange)
	return nil
}

func (s *BlockLogService) Stop() error {
	_ = s.bus.Unsubscribe(constants.NoticeBlockLog, s.onBlockLog)
	_ = s.bus.Unsubscribe(constants.NoticeLibChange, s.onLibChange)
	_ = s.db.Close()
	return nil
}

func (s *BlockLogService) Reload(config *node.Config) error {
	return nil
}

func (s *BlockLogService) initDatabase() error {
	connStr := fmt.Sprintf("%s:%s@/%s?charset=utf8mb4&parseTime=True&loc=Local", s.config.User, s.config.Password, s.config.Db)
	if db, err := gorm.Open(s.config.Driver, connStr); err != nil {
		return err
	} else {
		s.db = db
	}
	return nil
}

func (s *BlockLogService) onBlockLog(blockLog *blocklog.BlockLog, blockProducer string) {
	isGenesis := blockLog.BlockNum == 0
	rec := &iservices.BlockLogRecord{
		BlockId:     blockLog.BlockId,
		BlockHeight: blockLog.BlockNum,
		BlockTime:   time.Unix(int64(blockLog.BlockTime), 0),
		BlockProducer: blockProducer,
		Final:       isGenesis,
		JsonLog:     blockLog.ToJsonString(),
	}
	if !s.db.HasTable(rec) {
		s.db.CreateTable(rec)
	}
	s.db.Create(rec)
}

func (s *BlockLogService) onLibChange(blocks []common.ISignedBlock) {
	if count := len(blocks); count > 0 {
		updates := make(map[string][]string)
		for _ , block := range blocks {
			blockId := block.Id()
			tableName := iservices.BlockLogTableNameForBlockHeight(blockId.BlockNum())
			updates[tableName] = append(updates[tableName], fmt.Sprintf("%x", block.Id().Data))
		}
		tx := s.db.Begin()
		for tableName, blockIds := range updates {
			tx.Table(tableName).Where("block_id IN (?)", blockIds).Update("final", true)
		}
		tx.Commit()
	}
}


func init() {
	RegisterSQLTableNamePattern(fmt.Sprintf("%s\\w*", iservices.BlockLogTable))
}
