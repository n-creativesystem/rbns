package grpcserver

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	invalid codes.Code = 2000
)

var ErrAuthStatus = status.Error(invalid, "invalid api key")

func apiKeyCheck(apiKey string) grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
		if err != nil {
			return nil, err
		}
		if apiKey != token {
			return nil, ErrAuthStatus
		}
		return ctx, nil

	}
}

func AuthUnaryServerInterceptor(apiKey string) grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(apiKeyCheck(apiKey))
}

func Recovery(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	panicked := true
	defer func() {
		if rErr := recover(); rErr != nil || panicked {
			err = status.Errorf(codes.Internal, "%v", rErr)
		}
	}()
	resp, err := handler(ctx, req)
	panicked = false
	return resp, err
}
