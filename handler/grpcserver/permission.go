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

type permissionServer struct {
	*protobuf.UnimplementedPermissionServer
	svc service.PermissionService
}

var _ protobuf.PermissionServer = (*permissionServer)(nil)

func newPermissionServer(srv service.PermissionService) protobuf.PermissionServer {
	return &permissionServer{svc: srv}
}

// Permission
func (s *permissionServer) Create(ctx context.Context, in *protobuf.PermissionEntities) (*protobuf.PermissionEntities, error) {
	inPermissions := make([]*protobuf.PermissionEntity, len(in.GetPermissions()))
	copy(inPermissions, in.GetPermissions())
	if len(inPermissions) == 0 {
		return &protobuf.PermissionEntities{
			Permissions: make([]*protobuf.PermissionEntity, 0),
		}, nil
	}
	names := make([]string, len(inPermissions))
	descriptions := make([]string, len(inPermissions))
	for idx, permission := range inPermissions {
		names[idx] = permission.Name
		descriptions[idx] = permission.Description
	}
	permissions, err := s.svc.Create(ctx, names, descriptions)
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := &protobuf.PermissionEntities{
		Permissions: make([]*protobuf.PermissionEntity, len(permissions)),
	}
	for idx, permission := range permissions {
		out.Permissions[idx] = protoconv.NewPermissionEntityByModel(permission)
	}
	return out, err
}

func (s *permissionServer) FindById(ctx context.Context, in *protobuf.PermissionKey) (*protobuf.PermissionEntity, error) {
	permission, err := s.svc.FindById(ctx, in.GetId())
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return protoconv.NewPermissionEntityByModel(*permission), nil
}

func (s *permissionServer) FindAll(ctx context.Context, in *emptypb.Empty) (*protobuf.PermissionEntities, error) {
	permissions, err := s.svc.FindAll(ctx)
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := &protobuf.PermissionEntities{
		Permissions: make([]*protobuf.PermissionEntity, len(permissions)),
	}
	for idx, permission := range permissions {
		out.Permissions[idx] = protoconv.NewPermissionEntityByModel(permission)
	}
	return out, nil
}

func (s *permissionServer) Update(ctx context.Context, in *protobuf.PermissionEntity) (*emptypb.Empty, error) {
	err := s.svc.Update(ctx, in.GetId(), in.GetName(), in.GetDescription())
	if err != nil {
		if err == model.ErrNoData {
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *permissionServer) Delete(ctx context.Context, in *protobuf.PermissionKey) (*emptypb.Empty, error) {
	err := s.svc.Delete(ctx, in.GetId())
	if err != nil {
		if err == model.ErrNoData {
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *permissionServer) Check(ctx context.Context, in *protobuf.PermissionCheckRequest) (*protobuf.PermissionCheckResult, error) {
	result := &protobuf.PermissionCheckResult{
		Result:  false,
		Message: "",
	}
	res, err := s.svc.Check(ctx, in.GetUserKey(), in.GetOrganizationName(), in.GetPermissionNames()...)
	if err != nil {
		if err == model.ErrNoData {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		result.Result = false
		result.Message = err.Error()
		return result, status.Error(codes.Internal, err.Error())
	}
	result.Message = res.GetMsg()
	result.Result = res.IsOk()
	return result, nil
}
