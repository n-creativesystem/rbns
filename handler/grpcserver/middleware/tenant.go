package middleware

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/ncsfw/tenants"
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
	return &multiTenant{}
}

type multiTenant struct {
}

func (t *multiTenant) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		tenant := metautils.ExtractIncoming(ctx).Get(xTenantID)
		ctx, err = tenants.SetTenantWithContext(ctx, tenant)
		if err != nil {
			return nil, err
		}
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func (t *multiTenant) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		tenant := metautils.ExtractIncoming(ctx).Get(xTenantID)
		ctx, err := tenants.SetTenantWithContext(ctx, tenant)
		if err != nil {
			return err
		}
		ss = &baseSeverStream{
			ServerStream: ss,
			ctx:          ctx,
		}
		return handler(srv, ss)
	}
}
