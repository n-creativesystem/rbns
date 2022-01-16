package middleware

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/n-creativesystem/rbns/handler/metadata"
	"github.com/n-creativesystem/rbns/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApiKey interface {
	UnaryServerInterceptor() grpc.UnaryServerInterceptor
	StreamServerInterceptor() grpc.StreamServerInterceptor
}

type ApiKeyImpl struct {
	apiKeyService service.APIKey
}

func (a *ApiKeyImpl) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		apikey := metautils.ExtractIncoming(ctx).Get(metadata.XApiKey)
		loginUser, authErr := a.apiKeyService.Decode(ctx, apikey)
		if authErr != nil {
			return nil, status.Error(codes.Unauthenticated, authErr.Message)
		}
		service.SetCurrentUser(ctx, loginUser)
		return handler(ctx, req)
	}
}

func (a *ApiKeyImpl) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		apikey := metautils.ExtractIncoming(ctx).Get(metadata.XApiKey)
		loginUser, authErr := a.apiKeyService.Decode(ctx, apikey)
		if authErr != nil {
			return status.Error(codes.Unauthenticated, authErr.Message)
		}
		ctx = service.SetCurrentUser(ctx, loginUser)
		ss = &baseSeverStream{
			ServerStream: ss,
			ctx:          ctx,
		}
		return handler(srv, ss)
	}
}
