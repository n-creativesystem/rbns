package grpcserver

import (
	"context"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/proto"
	"github.com/n-creativesystem/rbns/protoconv"
	"github.com/n-creativesystem/rbns/service"
)

func init() {
	di.MustRegister(newUserServer)
}

type userServer struct {
	*proto.UnimplementedUserServer
	svc service.UserService
}

var _ proto.UserServer = (*userServer)(nil)

func newUserServer(svc service.UserService) proto.UserServer {
	return &userServer{svc: svc}
}

// User
func (s *userServer) Create(ctx context.Context, in *proto.UserEntity) (*proto.Empty, error) {
	roles := make([]string, len(in.GetRoles()))
	for idx, role := range in.GetRoles() {
		roles[idx] = role.GetId()
	}
	err := s.svc.Create(ctx, in.GetKey(), in.GetOrganizationId(), roles...)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (s *userServer) Delete(ctx context.Context, in *proto.UserKey) (*proto.Empty, error) {
	err := s.svc.Delete(ctx, in.GetKey(), in.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, err
}

func (s *userServer) FindByKey(ctx context.Context, in *proto.UserKey) (*proto.UserEntity, error) {
	u, err := s.svc.FindByKey(ctx, in.GetKey(), in.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	out := protoconv.NewUserEntityByModel(*u)
	out.OrganizationId = in.GetOrganizationId()
	return out, nil
}

func (s *userServer) AddRole(ctx context.Context, in *proto.UserRole) (*proto.Empty, error) {
	roles := make([]string, len(in.GetRoles()))
	if len(roles) == 0 {
		return &proto.Empty{}, nil
	}
	for idx, role := range in.GetRoles() {
		roles[idx] = role.GetId()
	}
	err := s.svc.AddRole(ctx, in.GetUser().GetKey(), in.GetUser().GetOrganizationId(), roles)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (s *userServer) DeleteRole(ctx context.Context, in *proto.UserRole) (*proto.Empty, error) {
	roles := make([]string, len(in.GetRoles()))
	if len(roles) == 0 {
		return &proto.Empty{}, nil
	}
	for idx, role := range in.GetRoles() {
		roles[idx] = role.GetId()
	}
	err := s.svc.DeleteRole(ctx, in.GetUser().GetKey(), in.GetUser().GetOrganizationId(), roles)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}
