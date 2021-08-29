package grpcserver

import (
	"context"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/protobuf"
	"github.com/n-creativesystem/rbns/protoconv"
	"github.com/n-creativesystem/rbns/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userServer struct {
	*protobuf.UnimplementedUserServer
	svc    service.UserService
	orgSvc service.OrganizationService
}

var _ protobuf.UserServer = (*userServer)(nil)

func newUserServer(svc service.UserService, orgSvc service.OrganizationService) protobuf.UserServer {
	return &userServer{svc: svc, orgSvc: orgSvc}
}

// User
func (s *userServer) Create(ctx context.Context, in *protobuf.UserEntity) (*emptypb.Empty, error) {
	roles := make([]string, len(in.GetRoles()))
	for idx, role := range in.GetRoles() {
		roles[idx] = role.GetId()
	}
	err := s.svc.Create(ctx, in.GetKey(), in.GetOrganizationId(), roles...)
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *userServer) Delete(ctx context.Context, in *protobuf.UserKey) (*emptypb.Empty, error) {
	err := s.svc.Delete(ctx, in.GetKey(), in.GetOrganizationId())
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, err
}

func (s *userServer) FindByKey(ctx context.Context, in *protobuf.UserKey) (*protobuf.UserEntity, error) {
	u, err := s.svc.FindByKey(ctx, in.GetKey(), in.GetOrganizationId())
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := protoconv.NewUserEntityByModel(*u)
	out.OrganizationId = in.GetOrganizationId()
	return out, nil
}

func (s *userServer) FindByOrganizationNameAndUserKey(ctx context.Context, in *protobuf.UserKeyByName) (*protobuf.UserEntity, error) {
	org, err := s.orgSvc.FindByName(ctx, in.OrganizationName)
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return s.FindByKey(ctx, &protobuf.UserKey{
		Key:            in.Key,
		OrganizationId: *org.GetID(),
	})
}

func (s *userServer) AddRoles(ctx context.Context, in *protobuf.UserRole) (*emptypb.Empty, error) {
	roles := make([]string, len(in.GetRoles()))
	if len(roles) == 0 {
		return &emptypb.Empty{}, nil
	}
	for idx, role := range in.GetRoles() {
		roles[idx] = role.GetId()
	}
	err := s.svc.AddRole(ctx, in.GetKey(), in.GetOrganizationId(), roles)
	if err != nil {
		if err == model.ErrNoData {
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *userServer) DeleteRoles(ctx context.Context, in *protobuf.UserRole) (*emptypb.Empty, error) {
	roles := make([]string, len(in.GetRoles()))
	if len(roles) == 0 {
		return &emptypb.Empty{}, nil
	}
	for idx, role := range in.GetRoles() {
		roles[idx] = role.GetId()
	}
	err := s.svc.DeleteRole(ctx, in.GetKey(), in.GetOrganizationId(), roles)
	if err != nil {
		if err == model.ErrNoData {
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *userServer) DeleteRole(ctx context.Context, in *protobuf.UserRoleDelete) (*emptypb.Empty, error) {
	return s.DeleteRoles(ctx, &protobuf.UserRole{
		Key:            in.Key,
		OrganizationId: in.OrganizationId,
		Roles: []*protobuf.RoleKey{
			{
				Id: in.RoleId,
			},
		},
	})
}
