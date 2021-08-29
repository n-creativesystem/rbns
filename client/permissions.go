package client

import (
	"context"

	"github.com/n-creativesystem/rbns/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Permissions interface {
	// Check is resource check
	Check(in *protobuf.PermissionCheckRequest, opts ...grpc.CallOption) (*protobuf.PermissionCheckResult, error)
	// FindById is find by id
	FindById(in *protobuf.PermissionKey, opts ...grpc.CallOption) (*protobuf.PermissionEntity, error)
	// Update is permission entity update
	Update(in *protobuf.PermissionEntity, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete is permission entity delete
	Delete(in *protobuf.PermissionKey, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Create is create permission
	Create(in *protobuf.PermissionEntities, opts ...grpc.CallOption) (*protobuf.PermissionEntities, error)
	// FindAll is find by application id return permissions
	FindAll(in *emptypb.Empty, opts ...grpc.CallOption) (*protobuf.PermissionEntities, error)
}

type permissionClient struct {
	ctx    context.Context
	client protobuf.PermissionClient
}

func (c *permissionClient) Check(in *protobuf.PermissionCheckRequest, opts ...grpc.CallOption) (*protobuf.PermissionCheckResult, error) {
	return c.client.Check(c.ctx, in, opts...)
}

func (c *permissionClient) FindById(in *protobuf.PermissionKey, opts ...grpc.CallOption) (*protobuf.PermissionEntity, error) {
	return c.client.FindById(c.ctx, in, opts...)
}

func (c *permissionClient) Update(in *protobuf.PermissionEntity, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.Update(c.ctx, in, opts...)
}

func (c *permissionClient) Delete(in *protobuf.PermissionKey, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.Delete(c.ctx, in, opts...)
}

func (c *permissionClient) Create(in *protobuf.PermissionEntities, opts ...grpc.CallOption) (*protobuf.PermissionEntities, error) {
	return c.client.Create(c.ctx, in, opts...)
}

func (c *permissionClient) FindAll(in *emptypb.Empty, opts ...grpc.CallOption) (*protobuf.PermissionEntities, error) {
	return c.client.FindAll(c.ctx, in, opts...)
}
