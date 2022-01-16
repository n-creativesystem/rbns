package grpcserver

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/n-creativesystem/rbns/handler/grpcserver/middleware"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"github.com/n-creativesystem/rbns/ncsfw/tracer"
	"github.com/n-creativesystem/rbns/protobuf"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func New(
	pSrv protobuf.PermissionServer,
	oSrv protobuf.OrganizationServer,
	uSrv protobuf.UserServer,
	tenantMiddleware middleware.Tenant,
	apiKeyMiddleware middleware.ApiKey,
) *grpc.Server {

	log := logger.New("grpc server")
	otelgrpcOpts := []otelgrpc.Option{otelgrpc.WithTracerProvider(tracer.GetTracerProvider()), otelgrpc.WithPropagators(tracer.GetPropagation())}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				logger.UnaryServerInterceptor(log),
				grpc_recovery.UnaryServerInterceptor(
					grpc_recovery.WithRecoveryHandlerContext(middleware.RecoveryFunc()),
				),
				otelgrpc.UnaryServerInterceptor(otelgrpcOpts...),
				grpc_validator.UnaryServerInterceptor(),
				tenantMiddleware.UnaryServerInterceptor(),
				apiKeyMiddleware.UnaryServerInterceptor(),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				logger.StreamServerInterceptor(log),
				grpc_recovery.StreamServerInterceptor(
					grpc_recovery.WithRecoveryHandlerContext(middleware.RecoveryFunc()),
				),
				otelgrpc.StreamServerInterceptor(otelgrpcOpts...),
				grpc_validator.StreamServerInterceptor(),
				tenantMiddleware.StreamServerInterceptor(),
				apiKeyMiddleware.StreamServerInterceptor(),
			),
		),
	)
	protobuf.RegisterPermissionServer(server, pSrv)
	protobuf.RegisterOrganizationServer(server, oSrv)
	protobuf.RegisterUserServer(server, uSrv)
	healthRegister(server)
	reflection.Register(server)
	return server
}
