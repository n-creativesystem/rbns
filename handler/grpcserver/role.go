package grpcserver

import (
	"context"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/proto"
	"github.com/n-creativesystem/rbns/protoconv"
	"github.com/n-creativesystem/rbns/service"
)

func init() {
	di.MustRegister(newRoleServer)
}

type roleServer struct {
	*proto.UnimplementedRoleServer
	svc service.RoleService
}

var _ proto.RoleServer = (*roleServer)(nil)

func newRoleServer(svc service.RoleService) proto.RoleServer {
	return &roleServer{svc: svc}
}

func (s *roleServer) Create(ctx context.Context, in *proto.RoleEntities) (*proto.RoleEntities, error) {
	inRoles := make([]*proto.RoleEntity, len(in.GetRoles()))
	copy(inRoles, in.GetRoles())
	if len(inRoles) == 0 {
		return &proto.RoleEntities{
			Roles: make([]*proto.RoleEntity, 0),
		}, nil
	}
	names := make([]string, len(inRoles))
	descriptions := make([]string, len(inRoles))
	for idx, role := range inRoles {
		names[idx] = role.Name
		descriptions[idx] = role.Description
	}
	roles, err := s.svc.Create(ctx, names, descriptions)
	if err != nil {
		return nil, err
	}
	out := &proto.RoleEntities{
		Roles: make([]*proto.RoleEntity, len(roles)),
	}
	for idx, role := range roles {
		out.Roles[idx] = protoconv.NewRoleEntityByModel(role)
	}
	return out, nil
}

func (s *roleServer) FindById(ctx context.Context, in *proto.RoleKey) (*proto.RoleEntity, error) {
	role, err := s.svc.FindById(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return protoconv.NewRoleEntityByModel(*role), nil
}

func (s *roleServer) FindAll(ctx context.Context, in *proto.Empty) (*proto.RoleEntities, error) {
	roles, err := s.svc.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	entities := &proto.RoleEntities{
		Roles: make([]*proto.RoleEntity, len(roles)),
	}
	for i, role := range roles {
		entities.Roles[i] = protoconv.NewRoleEntityByModel(role)
	}
	return entities, nil
}

func (s *roleServer) Update(ctx context.Context, in *proto.RoleUpdateEntity) (*proto.Empty, error) {
	err := s.svc.Update(ctx, in.GetId(), in.GetName(), in.GetDescription())
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (s *roleServer) Delete(ctx context.Context, in *proto.RoleKey) (*proto.Empty, error) {
	err := s.svc.Delete(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (s *roleServer) GetPermissions(ctx context.Context, in *proto.RoleKey) (*proto.PermissionEntities, error) {
	permissions, err := s.svc.GetPermissions(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	res := proto.PermissionEntities{
		Permissions: make([]*proto.PermissionEntity, len(permissions)),
	}
	for idx, permission := range permissions {
		res.Permissions[idx] = protoconv.NewPermissionEntityByModel(permission)
	}
	return &res, nil
}

func (s *roleServer) AddPermissions(ctx context.Context, in *proto.RoleReleationPermissions) (*proto.Empty, error) {
	permissionIds := make([]string, len(in.GetPermissions()))
	if len(permissionIds) == 0 {
		return &proto.Empty{}, nil
	}
	for idx, permission := range in.GetPermissions() {
		permissionIds[idx] = permission.GetId()
	}
	if err := s.svc.AddPermissions(ctx, in.GetId(), permissionIds); err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (s *roleServer) DeletePermission(ctx context.Context, in *proto.RoleReleationPermissions) (*proto.Empty, error) {
	permissionIds := make([]string, len(in.GetPermissions()))
	if len(permissionIds) == 0 {
		return &proto.Empty{}, nil
	}
	for idx, permission := range in.GetPermissions() {
		permissionIds[idx] = permission.GetId()
	}
	err := s.svc.DeletePermissions(ctx, in.GetId(), permissionIds)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}
