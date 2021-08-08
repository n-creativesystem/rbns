package grpcserver

import (
	"context"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/proto"
	"github.com/n-creativesystem/rbns/protoconv"
	"github.com/n-creativesystem/rbns/service"
)

func init() {
	di.MustRegister(newPermissionServer)
}

type permissionServer struct {
	*proto.UnimplementedPermissionServer
	svc service.PermissionService
}

var _ proto.PermissionServer = (*permissionServer)(nil)

func newPermissionServer(srv service.PermissionService) proto.PermissionServer {
	return &permissionServer{svc: srv}
}

// Permission
func (s *permissionServer) Create(ctx context.Context, in *proto.PermissionEntities) (*proto.PermissionEntities, error) {
	inPermissions := make([]*proto.PermissionEntity, len(in.GetPermissions()))
	copy(inPermissions, in.GetPermissions())
	if len(inPermissions) == 0 {
		return &proto.PermissionEntities{
			Permissions: make([]*proto.PermissionEntity, 0),
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
		return nil, err
	}
	out := &proto.PermissionEntities{
		Permissions: make([]*proto.PermissionEntity, len(permissions)),
	}
	for idx, permission := range permissions {
		out.Permissions[idx] = protoconv.NewPermissionEntityByModel(permission)
	}
	return out, err
}

func (s *permissionServer) FindById(ctx context.Context, in *proto.PermissionKey) (*proto.PermissionEntity, error) {
	permission, err := s.svc.FindById(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return protoconv.NewPermissionEntityByModel(*permission), nil
}

func (s *permissionServer) FindAll(ctx context.Context, in *proto.Empty) (*proto.PermissionEntities, error) {
	permissions, err := s.svc.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := &proto.PermissionEntities{
		Permissions: make([]*proto.PermissionEntity, len(permissions)),
	}
	for idx, permission := range permissions {
		out.Permissions[idx] = protoconv.NewPermissionEntityByModel(permission)
	}
	return out, nil
}

func (s *permissionServer) Update(ctx context.Context, in *proto.PermissionEntity) (*proto.Empty, error) {
	err := s.svc.Update(ctx, in.GetId(), in.GetName(), in.GetDescription())
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (s *permissionServer) Delete(ctx context.Context, in *proto.PermissionKey) (*proto.Empty, error) {
	err := s.svc.Delete(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (s *permissionServer) Check(ctx context.Context, in *proto.PermissionCheckRequest) (*proto.PermissionCheckResult, error) {
	result := &proto.PermissionCheckResult{
		Result:  false,
		Message: "",
	}
	res, err := s.svc.Check(ctx, in.GetUserKey(), in.GetOrganizationName(), in.GetPermissionNames()...)
	if err != nil {
		result.Result = false
		result.Message = err.Error()
		return result, err
	}
	result.Message = res.GetMsg()
	result.Result = res.IsOk()
	return result, nil
}
