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
	"github.com/n-creativesystem/rbns/handler/restserver/middleware/auth"
	"github.com/n-creativesystem/rbns/handler/restserver/social"
	"github.com/n-creativesystem/rbns/infra"
	"github.com/n-creativesystem/rbns/service"
	"github.com/n-creativesystem/rbns/storage"
)

var ServiceSet = wire.NewSet(
	wire.Bind(new(service.Tenant), new(*service.TenantImpl)),
	service.NewTenantImpl,
	wire.Bind(new(service.Permission), new(*service.PermissionImpl)),
	service.NewPermissionService,
	wire.Bind(new(service.Organization), new(*service.OrganizationImpl)),
	service.NewOrganizationService,
	wire.Bind(new(service.User), new(*service.UserImpl)),
	service.NewUserService,
	wire.Bind(new(service.OrganizationAggregation), new(*service.OrganizationAggregationImpl)),
	service.NewOrganizationAggregation,
	wire.Bind(new(service.AuthService), new(*service.AuthCache)),
	service.NewAuthCache,
)

var BaseSet = wire.NewSet(
	bus.GetBus,
	config.NewFlags2Config,
	infra.NewFactory,
	storage.Initialize,
	storage.NewKeyPairs,
	storage.NewSessionStore,
	cache.New,
)

var ServerSet = wire.NewSet(
	gateway.New,
	grpcserver.New,
	grpcserver.NewOrganizationService,
	grpcserver.NewPermissionServer,
	grpcserver.NewUserServer,
	auth.NewAuthMiddleware,
	restserver.New,
	middleware.NewTenantMiddleware,
	wire.Bind(new(social.Service), new(*social.SocialService)),
	social.ProvideService,
)

var OSSSet = wire.NewSet(
	BaseSet,
	ServiceSet,
	ServerSet,
)
