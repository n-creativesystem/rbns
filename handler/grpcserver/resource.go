package grpcserver

import (
	"context"

	"github.com/n-creativesystem/rbns/proto"
	"github.com/n-creativesystem/rbns/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type resourceServer struct {
	*proto.UnimplementedResourceServer
	svc service.Resource
}

var _ proto.ResourceServer = (*resourceServer)(nil)

func newResourceServer(svc service.Resource) proto.ResourceServer {
	return &resourceServer{
		svc: svc,
	}
}

func (r *resourceServer) Save(ctx context.Context, req *proto.SaveRequest) (*emptypb.Empty, error) {
	err := r.svc.Save(ctx, req.GetMethod(), req.GetUri(), req.GetPermissions()...)
	return &emptypb.Empty{}, err
}

func (r *resourceServer) Authz(ctx context.Context, req *proto.AuthzRequest) (*emptypb.Empty, error) {
	return nil, nil
}
