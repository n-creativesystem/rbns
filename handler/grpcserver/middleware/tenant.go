package middleware

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/internal/contexts"
	"google.golang.org/grpc"
)

const (
	xTenantID = "x-tenant-id"
)

type Tenant interface {
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
	StreamServerInterceptor() grpc.StreamServerInterceptor
}

func NewTenantMiddleware(conf *config.Config) Tenant {
	if conf.MultiTenant {
		return &multiTenant{}
	} else {
		return &tenant{}
	}
}

type tenant struct {
}

func (t *tenant) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

func (t *tenant) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, ss)
	}
}

type multiTenant struct {
}

type multiTenantSeverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (ss *multiTenantSeverStream) Context() context.Context {
	return ss.ctx
}

func (t *multiTenant) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		tenant := metautils.ExtractIncoming(ctx).Get(xTenantID)
		ctx = contexts.ToTenantContext(ctx, tenant)
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func (t *multiTenant) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		tenant := metautils.ExtractIncoming(ctx).Get(xTenantID)
		ctx = contexts.ToTenantContext(ctx, tenant)
		ss = &multiTenantSeverStream{
			ServerStream: ss,
			ctx:          ctx,
		}
		return handler(srv, ss)
	}
}
