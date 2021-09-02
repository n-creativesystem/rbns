package client

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MethodPermissions map[string][]string

type OrganizationUserFunc func(ctx context.Context) (newCtx context.Context, userKey string, organizationName string, err error)

type ServiceRBACFuncOverride interface {
	OrganizationUserFuncOverride(ctx context.Context, fullMethodName string) []string
}

func UnaryServerInterceptor(client RBNS, mp MethodPermissions, rbacFunc OrganizationUserFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var permissions []string
		if overrideSrv, ok := info.Server.(ServiceRBACFuncOverride); ok {
			p := overrideSrv.OrganizationUserFuncOverride(ctx, info.FullMethod)
			permissions = append(permissions, p...)
		} else {
			if v, ok := mp[info.FullMethod]; ok {
				permissions = append(permissions, v...)
			}
		}
		if len(permissions) > 0 {
			if newCtx, result, err := CheckPermissions(ctx, client, rbacFunc, permissions...); err != nil {
				return nil, err
			} else if result {
				return handler(newCtx, req)
			} else {
				return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
			}
		} else {
			return handler(newCtx, req)
		}
	}
}

func StreamServerInterceptor(client RBNS, mp MethodPermissions, rbacFunc OrganizationUserFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var permissions []string
		if overrideSrv, ok := srv.(ServiceRBACFuncOverride); ok {
			p := overrideSrv.OrganizationUserFuncOverride(stream.Context(), info.FullMethod)
			permissions = append(permissions, p...)
		} else {
			if v, ok := mp[info.FullMethod]; ok {
				permissions = append(permissions, v...)
			}
		}
		if len(permissions) > 0 {
			if newCtx, result, err := CheckPermissions(stream.Context(), client, rbacFunc, permissions...); err != nil {
				return err
			} else if result {
				wrapped := grpc_middleware.WrapServerStream(stream)
				wrapped.WrappedContext = newCtx
				return handler(srv, wrapped)
			} else {
				return status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
			}
		} else {
			wrapped := grpc_middleware.WrapServerStream(stream)
			wrapped.WrappedContext = newCtx
			return handler(srv, wrapped)
		}
	}
}
