package client

import (
	"context"

	"github.com/n-creativesystem/rbns/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Organizations interface {
	// FindById is application id and organization id
	FindById(in *protobuf.OrganizationKey, opts ...grpc.CallOption) (*protobuf.OrganizationEntity, error)
	// FindAll is application is return organizations
	FindAll(in *emptypb.Empty, opts ...grpc.CallOption) (*protobuf.OrganizationEntities, error)
	// Update is organization entity update
	Update(in *protobuf.OrganizationUpdateEntity, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete is organization entity delete
	Delete(in *protobuf.OrganizationKey, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Create is create orgnization
	Create(in *protobuf.OrganizationEntity, opts ...grpc.CallOption) (*protobuf.OrganizationEntity, error)
}

type organizationClient struct {
	ctx    context.Context
	client protobuf.OrganizationClient
}

func (c *organizationClient) FindById(in *protobuf.OrganizationKey, opts ...grpc.CallOption) (*protobuf.OrganizationEntity, error) {
	return c.client.FindById(c.ctx, in, opts...)
}

func (c *organizationClient) FindAll(in *emptypb.Empty, opts ...grpc.CallOption) (*protobuf.OrganizationEntities, error) {
	return c.client.FindAll(c.ctx, in, opts...)
}

func (c *organizationClient) Update(in *protobuf.OrganizationUpdateEntity, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.Update(c.ctx, in, opts...)
}

func (c *organizationClient) Delete(in *protobuf.OrganizationKey, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.client.Delete(c.ctx, in, opts...)
}

func (c *organizationClient) Create(in *protobuf.OrganizationEntity, opts ...grpc.CallOption) (*protobuf.OrganizationEntity, error) {
	return c.client.Create(c.ctx, in, opts...)
}
