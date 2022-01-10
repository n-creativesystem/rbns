//go:build wireinject
// +build wireinject

package registry

import (
	"github.com/google/wire"
	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/cache"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/handler/gateway"
	"github.com/n-creativesystem/rbns/handler/grpcserver"
	"github.com/n-creativesystem/rbns/handler/grpcserver/middleware"
	"github.com/n-creativesystem/rbns/handler/restserver"
	rmiddleware "github.com/n-creativesystem/rbns/handler/restserver/middleware"
	"github.com/n-creativesystem/rbns/handler/restserver/social"
	"github.com/n-creativesystem/rbns/infra"
	"github.com/n-creativesystem/rbns/service"
	"github.com/n-creativesystem/rbns/storage"
)

var SuperSet = wire.NewSet(
	bus.GetBus,
	config.NewFlags2Config,
	infra.NewFactory,
	service.NewOrganizationService,
	service.NewPermissionService,
	service.NewUserService,
	service.NewOrganizationAggregation,
	storage.Initialize,
	storage.NewKeyPairs,
	storage.NewSessionStore,
	gateway.New,
	grpcserver.New,
	grpcserver.NewOrganizationService,
	grpcserver.NewPermissionServer,
	grpcserver.NewRoleServer,
	grpcserver.NewUserServer,
	rmiddleware.NewAuthMiddleware,
	restserver.New,
	middleware.NewTenantMiddleware,
	social.ProvideService,
	wire.Bind(new(social.Service), new(*social.SocialService)),
	cache.New,
	service.NewAuthCache,
	wire.Bind(new(service.AuthService), new(*service.AuthCache)),
)
