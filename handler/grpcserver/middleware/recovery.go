package middleware

import (
	"fmt"

	"github.com/n-creativesystem/rbns/logger"
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

func RecoveryFunc() func(p interface{}) error {
	return func(p interface{}) error {
		logger.Error(p.(error), fmt.Sprintf("p: %+v\n", p))
		return status.Errorf(codes.Internal, "Unexpected error")
	}
}
