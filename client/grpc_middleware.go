package client

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/n-creativesystem/rbns/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MethodPermissions map[string][]string

type OrganizationUserFunc func(ctx context.Context) (newCtx context.Context, userKey string, organizationName string, err error)

type ServiceRBACFuncOverride interface {
	OrganizationUserFuncOverride(ctx context.Context, fullMethodName string) (newCtx context.Context, userKey string, organizationName string, err error)
}

func UnaryServerInterceptor(client RBNS, mp MethodPermissions, rbacFunc OrganizationUserFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var newCtx context.Context
		var userKey, organizationName string
		var err error
		var permissions []string
		if v, ok := mp[info.FullMethod]; ok {
			permissions = append(permissions, v...)
		}
		if overrideSrv, ok := info.Server.(ServiceRBACFuncOverride); ok {
			newCtx, userKey, organizationName, err = overrideSrv.OrganizationUserFuncOverride(ctx, info.FullMethod)
		} else {
			newCtx, userKey, organizationName, err = rbacFunc(ctx)
		}
		if err != nil {
			return nil, err
		}
		result, err := client.Permissions(newCtx).Check(&protobuf.PermissionCheckRequest{
			UserKey:          userKey,
			OrganizationName: organizationName,
			PermissionNames:  permissions,
		})
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if result.Result {
			return handler(newCtx, req)
		}
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}
}

func StreamServerInterceptor(client RBNS, mp MethodPermissions, rbacFunc OrganizationUserFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var newCtx context.Context
		var userKey, organizationName string
		var err error
		var permissions []string
		if v, ok := mp[info.FullMethod]; ok {
			permissions = append(permissions, v...)
		}
		if overrideSrv, ok := srv.(ServiceRBACFuncOverride); ok {
			newCtx, userKey, organizationName, err = overrideSrv.OrganizationUserFuncOverride(stream.Context(), info.FullMethod)
		} else {
			newCtx, userKey, organizationName, err = rbacFunc(stream.Context())
		}
		if err != nil {
			return err
		}
		result, err := client.Permissions(newCtx).Check(&protobuf.PermissionCheckRequest{
			UserKey:          userKey,
			OrganizationName: organizationName,
			PermissionNames:  permissions,
		})
		if err != nil {
			return err
		}
		if result.Result {
			wrapped := grpc_middleware.WrapServerStream(stream)
			wrapped.WrappedContext = newCtx
			return handler(srv, wrapped)
		}
		return status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}
}
