package grpcserver

import (
	"context"

	"github.com/n-creativesystem/rbns/protobuf"
	"github.com/n-creativesystem/rbns/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type resourceServer struct {
	*protobuf.UnimplementedResourceServer
	svc service.Resource
}

var _ protobuf.ResourceServer = (*resourceServer)(nil)

func newResourceServer(svc service.Resource) protobuf.ResourceServer {
	return &resourceServer{
		svc: svc,
	}
}

func (r *resourceServer) Save(ctx context.Context, req *protobuf.SaveRequest) (*emptypb.Empty, error) {
	err := r.svc.Save(ctx, req.GetMethod(), req.GetUri(), req.GetPermissions()...)
	return &emptypb.Empty{}, err
}

func (r *resourceServer) Authz(ctx context.Context, req *protobuf.AuthzRequest) (*emptypb.Empty, error) {
	return nil, nil
}
