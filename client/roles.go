package client

import (
	"context"

	"github.com/n-creativesystem/rbns/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Roles interface {
	// DeletePermission is delete permission to the role
	DeletePermission(in *protobuf.RoleReleationPermission, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// GetPermissions is get permission to the role
	GetPermissions(in *protobuf.RoleKey, opts ...grpc.CallOption) (*protobuf.PermissionEntities, error)
	// AddPermissions is add permission to the role
	AddPermissions(in *protobuf.RoleReleationPermissions, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeletePermissions is delete permission to the role
	DeletePermissions(in *protobuf.RoleReleationPermissions, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// FindById is find by id
	FindById(in *protobuf.RoleKey, opts ...grpc.CallOption) (*protobuf.RoleEntity, error)
	// Update is role entity update
	Update(in *protobuf.RoleUpdateEntity, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete is role entity delete
	Delete(in *protobuf.RoleKey, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// RoleCreate is create role
	Create(in *protobuf.RoleEntities, opts ...grpc.CallOption) (*protobuf.RoleEntities, error)
	// FindAll is find roles
	FindAll(in *emptypb.Empty, opts ...grpc.CallOption) (*protobuf.RoleEntities, error)
}

type roleClient struct {
	ctx    context.Context
	client protobuf.RoleClient
}

func (c *roleClient) DeletePermission(in *protobuf.RoleReleationPermission, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.DeletePermission(c.ctx, in, opts...)
}

func (c *roleClient) GetPermissions(in *protobuf.RoleKey, opts ...grpc.CallOption) (*protobuf.PermissionEntities, error) {
	return c.client.GetPermissions(c.ctx, in, opts...)
}

func (c *roleClient) AddPermissions(in *protobuf.RoleReleationPermissions, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.AddPermissions(c.ctx, in, opts...)
}

func (c *roleClient) DeletePermissions(in *protobuf.RoleReleationPermissions, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.DeletePermissions(c.ctx, in, opts...)
}

func (c *roleClient) FindById(in *protobuf.RoleKey, opts ...grpc.CallOption) (*protobuf.RoleEntity, error) {
	return c.client.FindById(c.ctx, in, opts...)
}

func (c *roleClient) Update(in *protobuf.RoleUpdateEntity, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.Update(c.ctx, in, opts...)
}

func (c *roleClient) Delete(in *protobuf.RoleKey, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.Delete(c.ctx, in, opts...)
}

func (c *roleClient) Create(in *protobuf.RoleEntities, opts ...grpc.CallOption) (*protobuf.RoleEntities, error) {
	return c.client.Create(c.ctx, in, opts...)
}

func (c *roleClient) FindAll(in *emptypb.Empty, opts ...grpc.CallOption) (*protobuf.RoleEntities, error) {
	if in == nil {
		in = &emptypb.Empty{}
	}
	return c.client.FindAll(c.ctx, in, opts...)
}
