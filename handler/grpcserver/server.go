package grpcserver

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/infra/dao"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/proto"
	"google.golang.org/grpc"
)

type Option func(*config)

type config struct {
	secure bool
}

func WithSecure(conf *config) {
	conf.secure = true
}

func New(db dao.DataBase, opts ...Option) *grpc.Server {
	conf := &config{}
	for _, opt := range opts {
		opt(conf)
	}
	interceptors := []grpc.UnaryServerInterceptor{
		logger.GrpcLogger(),
		Recovery,
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(interceptors...),
		),
	)
	var (
		pSrv    proto.PermissionServer
		rSrv    proto.RoleServer
		oSrv    proto.OrganizationServer
		uSrv    proto.UserServer
		reSrv   proto.ResourceServer
		authSrv authorizationServer
	)
	di.MustInvoke(func(s proto.PermissionServer) {
		pSrv = s
	})
	di.MustInvoke(func(s proto.RoleServer) {
		rSrv = s
	})
	di.MustInvoke(func(s proto.OrganizationServer) {
		oSrv = s
	})
	di.MustInvoke(func(s proto.UserServer) {
		uSrv = s
	})
	di.MustInvoke(func(s proto.ResourceServer) {
		reSrv = s
	})
	di.MustInvoke(func(s authorizationServer) {
		authSrv = s
	})
	proto.RegisterPermissionServer(server, pSrv)
	proto.RegisterRoleServer(server, rSrv)
	proto.RegisterOrganizationServer(server, oSrv)
	proto.RegisterUserServer(server, uSrv)
	proto.RegisterResourceServer(server, reSrv)
	envoyAuthzRegister(server, authSrv)
	healthRegister(server)
	return server
}
