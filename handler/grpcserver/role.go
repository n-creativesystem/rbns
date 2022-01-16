package grpcserver

import (
	"context"
	"errors"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/protobuf"
	"github.com/n-creativesystem/rbns/protoconv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *organizationServer) CreateRole(ctx context.Context, in *protobuf.RoleEntities) (*protobuf.RoleEntities, error) {
	inRoles := make([]*protobuf.RoleEntity, len(in.GetRoles()))
	copy(inRoles, in.GetRoles())
	if len(inRoles) == 0 {
		return &protobuf.RoleEntities{
			Roles: make([]*protobuf.RoleEntity, 0),
		}, nil
	}
	names := make([]string, len(inRoles))
	descriptions := make([]string, len(inRoles))
	for idx, role := range inRoles {
		names[idx] = role.Name
		descriptions[idx] = role.Description
	}
	roles, err := s.orgAggregation.RoleCreate(ctx, in.GetOrganizationId(), names, descriptions)
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := &protobuf.RoleEntities{
		Roles: make([]*protobuf.RoleEntity, len(roles)),
	}
	for idx, role := range roles {
		out.Roles[idx] = protoconv.NewRoleEntityByModel(role)
	}
	return out, nil
}

func (s *organizationServer) FindByIdRole(ctx context.Context, in *protobuf.RoleKey) (*protobuf.RoleEntity, error) {
	role, err := s.orgAggregation.RoleFindById(ctx, in.GetOrganizationId(), in.GetId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return protoconv.NewRoleEntityByModel(*role), nil
}

func (s *organizationServer) FindAllRole(ctx context.Context, in *protobuf.RoleFindAll) (*protobuf.RoleEntities, error) {
	roles, err := s.orgAggregation.RoleFindAll(ctx, in.GetOrganizationId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	entities := &protobuf.RoleEntities{
		Roles: make([]*protobuf.RoleEntity, len(roles)),
	}
	for i, role := range roles {
		entities.Roles[i] = protoconv.NewRoleEntityByModel(role)
	}
	return entities, nil
}

func (s *organizationServer) UpdateRole(ctx context.Context, in *protobuf.RoleUpdateEntity) (*emptypb.Empty, error) {
	err := s.orgAggregation.RoleUpdate(ctx, in.GetOrganizationId(), in.GetId(), in.GetName(), in.GetDescription())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *organizationServer) DeleteRole(ctx context.Context, in *protobuf.RoleKey) (*emptypb.Empty, error) {
	err := s.orgAggregation.RoleDelete(ctx, in.GetOrganizationId(), in.GetId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *organizationServer) GetRolePermissions(ctx context.Context, in *protobuf.RoleKey) (*protobuf.PermissionEntities, error) {
	permissions, err := s.orgAggregation.GetRolePermissions(ctx, in.GetOrganizationId(), in.GetId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := protobuf.PermissionEntities{
		Permissions: make([]*protobuf.PermissionEntity, len(permissions)),
	}
	for idx, permission := range permissions {
		res.Permissions[idx] = protoconv.NewPermissionEntityByModel(permission)
	}
	return &res, nil
}

func (s *organizationServer) AddRolePermissions(ctx context.Context, in *protobuf.RoleReleationPermissions) (*emptypb.Empty, error) {
	permissionIds := make([]string, len(in.GetPermissions()))
	if len(permissionIds) == 0 {
		return &emptypb.Empty{}, nil
	}
	for idx, permission := range in.GetPermissions() {
		permissionIds[idx] = permission.GetId()
	}
	if err := s.orgAggregation.AddRolePermissions(ctx, in.GetOrganizationId(), in.GetId(), permissionIds); err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *organizationServer) DeleteRolePermission(ctx context.Context, in *protobuf.RoleReleationPermission) (*emptypb.Empty, error) {
	return s.DeleteRolePermissions(ctx, &protobuf.RoleReleationPermissions{
		Id:             in.Id,
		OrganizationId: in.GetOrganizationId(),
		Permissions: []*protobuf.PermissionKey{
			{
				Id: in.PermissionId,
			},
		},
	})
}

func (s *organizationServer) DeleteRolePermissions(ctx context.Context, in *protobuf.RoleReleationPermissions) (*emptypb.Empty, error) {
	permissionIds := make([]string, len(in.GetPermissions()))
	if len(permissionIds) == 0 {
		return &emptypb.Empty{}, nil
	}
	for idx, permission := range in.GetPermissions() {
		permissionIds[idx] = permission.GetId()
	}
	err := s.orgAggregation.DeleteRolePermissions(ctx, in.GetOrganizationId(), in.GetId(), permissionIds)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return &emptypb.Empty{}, err
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
