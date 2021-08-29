package grpcserver

import (
	"strings"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/protobuf"
	"github.com/n-creativesystem/rbns/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type ApiKeyManagements struct {
	mu     sync.Mutex
	values map[string]struct{}
}

func (apiKey *ApiKeyManagements) Exists(value string) bool {
	apiKey.mu.Lock()
	defer apiKey.mu.Unlock()
	_, ok := apiKey.values[value]
	return ok
}

func (apiKey *ApiKeyManagements) set(key string) {
	apiKey.mu.Lock()
	defer apiKey.mu.Unlock()
	apiKey.values[key] = struct{}{}
}

var (
	apiKeyManager ApiKeyManagements
)

func init() {
	apiKey := viper.GetString("apiKey")
	keys := strings.Split(apiKey, ";")
	apiKeyManager.values = make(map[string]struct{}, len(keys))
	for _, key := range keys {
		apiKeyManager.set(key)
	}
}

type Option func(*config)

type config struct {
	secure bool
}

func WithSecure(conf *config) {
	conf.secure = true
}

func New(reader repository.Reader, writer repository.Writer, logger *logrus.Logger, opts ...Option) *grpc.Server {
	conf := &config{}
	for _, opt := range opts {
		opt(conf)
	}
	interceptors := []grpc.UnaryServerInterceptor{
		grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)),
		Recovery,
		// AuthUnaryServerInterceptor(),
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(interceptors...),
		),
	)

	var (
		pSvc  = service.NewPermissionService(reader, writer)
		rSvc  = service.NewRoleService(reader, writer)
		oSvc  = service.NewOrganizationService(reader, writer)
		uSvc  = service.NewUserService(reader, writer)
		reSvc = service.NewResource(reader, writer)
	)
	var (
		pSrv    protobuf.PermissionServer   = newPermissionServer(pSvc)
		rSrv    protobuf.RoleServer         = newRoleServer(rSvc)
		oSrv    protobuf.OrganizationServer = newOrganizationService(oSvc)
		uSrv    protobuf.UserServer         = newUserServer(uSvc, oSvc)
		reSrv   protobuf.ResourceServer     = newResourceServer(reSvc)
		authSrv authorizationServer         = newAuthz(reSvc)
	)
	protobuf.RegisterPermissionServer(server, pSrv)
	protobuf.RegisterRoleServer(server, rSrv)
	protobuf.RegisterOrganizationServer(server, oSrv)
	protobuf.RegisterUserServer(server, uSrv)
	protobuf.RegisterResourceServer(server, reSrv)
	envoyAuthzRegister(server, authSrv)
	healthRegister(server)
	return server
}
