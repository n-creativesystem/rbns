package client

import (
	"context"

	"github.com/n-creativesystem/rbns/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Resource interface {
	Find(id string, opts ...grpc.CallOption) (*protobuf.ResourceResponse, error)
	FindAll(opts ...grpc.CallOption) (*protobuf.ResourceResponses, error)
	Exists(id string, opts ...grpc.CallOption) (*protobuf.ResourceExistsResponse, error)
	Save(id, description string, permissionNames []string, opts ...grpc.CallOption) error
	Migration(id, description string, permissionNames []string, opts ...grpc.CallOption) error
	Delete(id string, opts ...grpc.CallOption) error
}

type resource struct {
	ctx    context.Context
	client protobuf.ResourceClient
}

func (r *resource) Find(id string, opts ...grpc.CallOption) (*protobuf.ResourceResponse, error) {
	return r.client.Find(r.ctx, &protobuf.ResourceKey{Id: id}, opts...)
}

func (r *resource) FindAll(opts ...grpc.CallOption) (*protobuf.ResourceResponses, error) {
	return r.client.FindAll(r.ctx, &emptypb.Empty{}, opts...)
}

func (r *resource) Exists(id string, opts ...grpc.CallOption) (*protobuf.ResourceExistsResponse, error) {
	return r.client.Exists(r.ctx, &protobuf.ResourceKey{Id: id}, opts...)
}

func (r *resource) Save(id, description string, permissionNames []string, opts ...grpc.CallOption) error {
	res, err := r.Exists(id, opts...)
	if err != nil {
		return err
	}
	in := &protobuf.ResourceSaveRequest{
		Id:              id,
		Description:     description,
		PermissionNames: make([]string, 0, len(permissionNames)),
	}
	for _, permission := range permissionNames {
		in.PermissionNames = append(in.PermissionNames, permission)
	}
	if res.IsExists {
		if _, err := r.client.Update(r.ctx, in, opts...); err != nil {
			return err
		}
	} else {
		if _, err := r.client.Create(r.ctx, in, opts...); err != nil {
			return err
		}
	}
	return nil
}

func (r *resource) Migration(id, description string, permissionNames []string, opts ...grpc.CallOption) error {
	in := &protobuf.ResourceSaveRequest{
		Id:              id,
		Description:     description,
		PermissionNames: make([]string, 0, len(permissionNames)),
	}
	for _, permission := range permissionNames {
		in.PermissionNames = append(in.PermissionNames, permission)
	}
	if _, err := r.client.Migration(r.ctx, in, opts...); err != nil {
		return err
	}
	return nil
}

func (r *resource) Delete(id string, opts ...grpc.CallOption) error {
	_, err := r.client.Delete(r.ctx, &protobuf.ResourceKey{Id: id}, opts...)
	return err
}
