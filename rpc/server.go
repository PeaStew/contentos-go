package rpc

import (
	"github.com/coschain/contentos-go/iservices"
	"github.com/coschain/contentos-go/iservices/service-configs"
	"github.com/coschain/contentos-go/node"
	"github.com/coschain/contentos-go/rpc/pb"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/sirupsen/logrus"
	"net"

	"google.golang.org/grpc"
)

const (
	GRPCMaxRecvMsgSize = 4096
	GRPCServerType     = "tcp"
)

type GRPCServer struct {
	rpcServer *grpc.Server
	ctx       *node.ServiceContext
	api       *APIService
	config    *service_configs.GRPCConfig
	log       *logrus.Logger
}

func NewGRPCServer(ctx *node.ServiceContext, config service_configs.GRPCConfig) (*GRPCServer, error) {
	logService, err := ctx.Service(iservices.LogServerName)
	if err != nil {
		panic(err)
	}

	gi := NewGRPCIntercepter(logService.(iservices.ILog))

	rpc := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(gi.streamRecoveryLoggingInterceptor)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(gi.unaryRecoveryLoggingInterceptor)),
		grpc.MaxRecvMsgSize(GRPCMaxRecvMsgSize))
	api := &APIService{}
	grpcpb.RegisterApiServiceServer(rpc, api)
	srv := &GRPCServer{rpcServer: rpc, ctx: ctx, api: api, config: &config}

	srv.log = logService.(iservices.ILog).GetLog()
	srv.api.log = srv.log
	return srv, nil
}

func (gs *GRPCServer) Start(node *node.Node) error {

	consensus, err := gs.ctx.Service(iservices.ConsensusServerName)
	if err != nil {
		// TODO Mock Test
		//return err
	} else {
		gs.api.consensus = consensus.(iservices.IConsensus)
	}

	db, err := gs.ctx.Service(iservices.DbServerName)
	if err != nil {
		// TODO Mock Test
		//return err
	} else {
		gs.api.db = db.(iservices.IDatabaseService)
	}

	gs.api.mainLoop = node.MainLoop

	err = gs.startGRPC()
	if err != nil {
		return err
	} else {
		gs.log.Infof("GPRC Server Start [ %s ]", gs.config.RPCListen)
	}

	err = gs.startGateway()
	if err != nil {
		return err
	} else {
		gs.log.Infof("Gateway Server Start [ %s ]", gs.config.HTTPListen)
	}

	return nil
}

func (gs *GRPCServer) startGRPC() error {
	gs.log.Infof("RPCListen %v", gs.config.RPCListen)
	listener, err := net.Listen(GRPCServerType, gs.config.RPCListen)
	if err != nil {
		gs.log.Errorf("grpc listener addr: [%s] failure", gs.config.RPCListen)
	}

	go func() {
		grpc.NewServer()
		if err := gs.rpcServer.Serve(listener); err != nil {
			gs.log.Errorf("rpc server start failure, %v", err)
		} else {
			gs.log.Info("rpc server start failure")
		}
	}()

	return nil
}

func (gs *GRPCServer) Stop() error {
	gs.rpcServer.Stop()
	return nil
}

func (gs *GRPCServer) startGateway() error {
	go func() {
		if err := Run(gs.config); err != nil {
			gs.log.Error("rpc gateway start failure")
		} else {
			gs.log.Info("rpc gateway start failure")
		}
	}()
	return nil
}
