package plugins

import (
	"github.com/coschain/contentos-go/iservices"
	"github.com/coschain/contentos-go/node"
)

type PluginMgt struct {
	list []string
}

func NewPluginMgt(list []string) *PluginMgt {
	return &PluginMgt{list: list}
}

func (p *PluginMgt) RegisterTrxPoolDependents(app *node.Node, cfg *node.Config) {
	_ = app.Register(FollowServiceName, func(ctx *node.ServiceContext) (node.Service, error) {
		return NewFollowService(ctx, app.Log)
	})

	_ = app.Register(PostServiceName, func(ctx *node.ServiceContext) (node.Service, error) {
		return NewPostService(ctx)
	})

	_ = app.Register(TrxServiceName, func(ctx *node.ServiceContext) (node.Service, error) {
		return NewTrxSerVice(ctx, app.Log)
	})
}

func (p *PluginMgt) RegisterSQLServices(app *node.Node, cfg *node.Config) {
	for _, l := range p.list  {
		switch l {
		case TrxMysqlServiceName:
			_ = app.Register(TrxMysqlServiceName, func(ctx *node.ServiceContext) (service node.Service, e error) {
				return NewTrxMysqlSerVice(ctx, cfg.Database, app.Log)
			})
		case iservices.DailyStatisticServiceName:
			_ = app.Register(iservices.DailyStatisticServiceName, func(ctx *node.ServiceContext) (node.Service, error) {
				return NewDailyStatisticService(ctx, cfg.Database, app.Log)
			})
		case iservices.BlockLogServiceName:
			_ = app.Register(iservices.BlockLogServiceName, func(ctx *node.ServiceContext) (service node.Service, e error) {
				return NewBlockLogService(ctx, cfg.Database, app.Log)
			})
		case iservices.BlockLogProcessServiceName:
			_ = app.Register(iservices.BlockLogProcessServiceName, func(ctx *node.ServiceContext) (service node.Service, e error) {
				return NewBlockLogProcessService(ctx, cfg.Database, app.Log)
			})
		}
	}
}

func (p *PluginMgt) RegisterIpRestrictService(app *node.Node, cfg *node.Config) {
	for _, l := range p.list {
		if l == iservices.IpRestrictServiceName {
			_ = app.Register(iservices.IpRestrictServiceName, func(ctx *node.ServiceContext) (node.Service, error) {
				return NewIpRestrictService(ctx, app.Log)
			})
			return
		}
	}
}
