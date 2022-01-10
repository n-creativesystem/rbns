package grpcserver

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/n-creativesystem/rbns/handler/grpcserver/middleware"
	"github.com/n-creativesystem/rbns/logger"
	"github.com/n-creativesystem/rbns/protobuf"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func New(
	pSrv protobuf.PermissionServer,
	rSrv protobuf.RoleServer,
	oSrv protobuf.OrganizationServer,
	uSrv protobuf.UserServer,
	// reSrv protobuf.ResourceServer,
	tenantMiddleware middleware.Tenant,
	// apiKey middleware.ApiKey,
) *grpc.Server {

	log := logger.New("grpc server")
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				logger.UnaryServerInterceptor(log),
				grpc_recovery.UnaryServerInterceptor(
					grpc_recovery.WithRecoveryHandler(middleware.RecoveryFunc()),
				),
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				grpc_validator.UnaryServerInterceptor(),
				// apiKey.UnaryServerInterceptor(),
				tenantMiddleware.UnaryServerInterceptor(),
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				logger.StreamServerInterceptor(log),
				grpc_recovery.StreamServerInterceptor(
					grpc_recovery.WithRecoveryHandler(middleware.RecoveryFunc()),
				),
				otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer()),
				grpc_validator.StreamServerInterceptor(),
				// apiKey.StreamServerInterceptor(),
				tenantMiddleware.StreamServerInterceptor(),
			),
		),
	)
	protobuf.RegisterPermissionServer(server, pSrv)
	protobuf.RegisterRoleServer(server, rSrv)
	protobuf.RegisterOrganizationServer(server, oSrv)
	protobuf.RegisterUserServer(server, uSrv)
	// protobuf.RegisterResourceServer(server, reSrv)
	healthRegister(server)
	reflection.Register(server)
	return server
}
