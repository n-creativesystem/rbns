package grpcserver

import (
	"context"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/proto"
	"github.com/n-creativesystem/rbns/service"
)

func init() {
	di.MustRegister(newResourceServer)
}

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

func (r *resourceServer) Save(ctx context.Context, req *proto.SaveRequest) (*proto.Empty, error) {
	err := r.svc.Save(ctx, req.GetMethod(), req.GetUri(), req.GetPermissions()...)
	return &proto.Empty{}, err
}

func (r *resourceServer) Authz(ctx context.Context, req *proto.AuthzRequest) (*proto.Empty, error) {
	return nil, nil
}
