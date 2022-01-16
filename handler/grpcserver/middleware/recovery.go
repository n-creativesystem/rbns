package middleware

import (
	"context"
	"fmt"

	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	invalid codes.Code = 2000
)

var ErrAuthStatus = status.Error(invalid, "invalid api key")

// func apiKeyCheck() grpc_auth.AuthFunc {
// 	return func(ctx context.Context) (context.Context, error) {
// 		token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
// 		if err != nil {
// 			return nil, err
// 		}
// 		if !apiKeyManager.Exists(token) {
// 			return nil, ErrAuthStatus
// 		}
// 		return ctx, nil

// 	}
// }

// func AuthUnaryServerInterceptor() grpc.UnaryServerInterceptor {
// 	return grpc_auth.UnaryServerInterceptor(apiKeyCheck())
// }

func RecoveryFunc() func(ctx context.Context, p interface{}) (err error) {
	return func(ctx context.Context, p interface{}) (err error) {
		logger.ErrorWithContext(ctx, p.(error), fmt.Sprintf("p: %+v\n", p))
		return status.Errorf(codes.Internal, "Unexpected error")
	}
}
