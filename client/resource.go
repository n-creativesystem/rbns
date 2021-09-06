package client

import (
	"context"

	"github.com/n-creativesystem/rbns/protobuf"
	"google.golang.org/grpc"
)

type Resource interface {
	Find(in *protobuf.ResourceFindRequest, opts ...grpc.CallOption) (*protobuf.ResourceResponse, error)
	Exists(in *protobuf.ResourceFindRequest, opts ...grpc.CallOption) (*protobuf.ResourceExistsResponse, error)
	Save(in *protobuf.ResourceSaveRequest, opts ...grpc.CallOption) error
	Migration(in *protobuf.ResourceSaveRequest, opts ...grpc.CallOption) error
}

type resource struct {
	ctx    context.Context
	client protobuf.ResourceClient
}

func (r *resource) Find(in *protobuf.ResourceFindRequest, opts ...grpc.CallOption) (*protobuf.ResourceResponse, error) {
	return r.client.Find(r.ctx, in)
}

func (r *resource) Exists(in *protobuf.ResourceFindRequest, opts ...grpc.CallOption) (*protobuf.ResourceExistsResponse, error) {
	return r.client.Exists(r.ctx, in)
}

func (r *resource) Save(in *protobuf.ResourceSaveRequest, opts ...grpc.CallOption) error {
	if _, err := r.client.Save(r.ctx, in); err != nil {
		return err
	}
	return nil
}

func (r *resource) Migration(in *protobuf.ResourceSaveRequest, opts ...grpc.CallOption) error {
	if _, err := r.client.Migration(r.ctx, in); err != nil {
		return err
	}
	return nil
}
