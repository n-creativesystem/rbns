package client

import (
	"context"

	"github.com/n-creativesystem/rbns/protobuf"
)

func CheckPermissions(ctx context.Context, client RBNS, rbacFunc OrganizationUserFunc, permissions ...string) (context.Context, bool, error) {
	var newCtx context.Context
	var userKey, organizationName string
	var err error
	newCtx, userKey, organizationName, err = rbacFunc(ctx)
	if err != nil {
		return ctx, false, err
	}
	result, err := client.Permissions(newCtx).Check(&protobuf.PermissionCheckRequest{
		UserKey:          userKey,
		OrganizationName: organizationName,
		PermissionNames:  permissions,
	})
	if err != nil {
		return newCtx, false, err
	}
	return newCtx, result.Result, nil
}
