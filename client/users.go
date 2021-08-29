package client

import (
	"context"

	"github.com/n-creativesystem/rbns/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Users interface {
	// FindByKey is find organization id and user key
	FindByKey(in *protobuf.UserKey, opts ...grpc.CallOption) (*protobuf.UserEntity, error)
	// FindByOrganizationNameAndUserKey is find organization id and user key
	FindByOrganizationNameAndUserKey(in *protobuf.UserKeyByName, opts ...grpc.CallOption) (*protobuf.UserEntity, error)
	// Delete is delete user
	Delete(in *protobuf.UserKey, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteRole is delete role to user
	DeleteRole(in *protobuf.UserRoleDelete, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// AddRoles is add role to user
	AddRoles(in *protobuf.UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteRoles is delete role to user
	DeleteRoles(in *protobuf.UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Create is create user
	Create(in *protobuf.UserEntity, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type userClient struct {
	ctx    context.Context
	client protobuf.UserClient
}

func (c *userClient) FindByKey(in *protobuf.UserKey, opts ...grpc.CallOption) (*protobuf.UserEntity, error) {
	return c.client.FindByKey(c.ctx, in, opts...)
}

func (c *userClient) FindByOrganizationNameAndUserKey(in *protobuf.UserKeyByName, opts ...grpc.CallOption) (*protobuf.UserEntity, error) {
	return c.client.FindByOrganizationNameAndUserKey(c.ctx, in, opts...)
}

func (c *userClient) Delete(in *protobuf.UserKey, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.Delete(c.ctx, in, opts...)
}

func (c *userClient) DeleteRole(in *protobuf.UserRoleDelete, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.DeleteRole(c.ctx, in, opts...)
}

func (c *userClient) AddRoles(in *protobuf.UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.AddRoles(c.ctx, in, opts...)
}

func (c *userClient) DeleteRoles(in *protobuf.UserRole, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.DeleteRoles(c.ctx, in, opts...)
}

func (c *userClient) Create(in *protobuf.UserEntity, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.Create(c.ctx, in, opts...)
}
