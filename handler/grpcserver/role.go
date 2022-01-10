package grpcserver

import (
	"context"
	"errors"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/protobuf"
	"github.com/n-creativesystem/rbns/protoconv"
	"github.com/n-creativesystem/rbns/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleServer struct {
	*protobuf.UnimplementedRoleServer
	svc service.OrganizationAggregation
}

var _ protobuf.RoleServer = (*roleServer)(nil)

func NewRoleServer(svc service.OrganizationAggregation) protobuf.RoleServer {
	return &roleServer{svc: svc}
}

func (s *roleServer) Create(ctx context.Context, in *protobuf.RoleEntities) (*protobuf.RoleEntities, error) {
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
	roles, err := s.svc.RoleCreate(ctx, in.GetOrganizationId(), names, descriptions)
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
	}
	out := &protobuf.RoleEntities{
		Roles: make([]*protobuf.RoleEntity, len(roles)),
	}
	for idx, role := range roles {
		out.Roles[idx] = protoconv.NewRoleEntityByModel(role)
	}
	return out, nil
}

func (s *roleServer) FindById(ctx context.Context, in *protobuf.RoleKey) (*protobuf.RoleEntity, error) {
	role, err := s.svc.RoleFindById(ctx, in.GetOrganizationId(), in.GetId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return protoconv.NewRoleEntityByModel(*role), nil
}

func (s *roleServer) FindAll(ctx context.Context, in *protobuf.RoleFindAll) (*protobuf.RoleEntities, error) {
	roles, err := s.svc.RoleFindAll(ctx, in.GetOrganizationId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
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

func (s *roleServer) Update(ctx context.Context, in *protobuf.RoleUpdateEntity) (*emptypb.Empty, error) {
	err := s.svc.RoleUpdate(ctx, in.GetOrganizationId(), in.GetId(), in.GetName(), in.GetDescription())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return &emptypb.Empty{}, err
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *roleServer) Delete(ctx context.Context, in *protobuf.RoleKey) (*emptypb.Empty, error) {
	err := s.svc.RoleDelete(ctx, in.GetOrganizationId(), in.GetId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return &emptypb.Empty{}, err
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *roleServer) GetPermissions(ctx context.Context, in *protobuf.RoleKey) (*protobuf.PermissionEntities, error) {
	permissions, err := s.svc.GetRolePermissions(ctx, in.GetOrganizationId(), in.GetId())
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return nil, err
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

func (s *roleServer) AddPermissions(ctx context.Context, in *protobuf.RoleReleationPermissions) (*emptypb.Empty, error) {
	permissionIds := make([]string, len(in.GetPermissions()))
	if len(permissionIds) == 0 {
		return &emptypb.Empty{}, nil
	}
	for idx, permission := range in.GetPermissions() {
		permissionIds[idx] = permission.GetId()
	}
	if err := s.svc.AddRolePermissions(ctx, in.GetOrganizationId(), in.GetId(), permissionIds); err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return &emptypb.Empty{}, err
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *roleServer) DeletePermission(ctx context.Context, in *protobuf.RoleReleationPermission) (*emptypb.Empty, error) {
	return s.DeletePermissions(ctx, &protobuf.RoleReleationPermissions{
		Id:             in.Id,
		OrganizationId: in.GetOrganizationId(),
		Permissions: []*protobuf.PermissionKey{
			{
				Id: in.PermissionId,
			},
		},
	})
}

func (s *roleServer) DeletePermissions(ctx context.Context, in *protobuf.RoleReleationPermissions) (*emptypb.Empty, error) {
	permissionIds := make([]string, len(in.GetPermissions()))
	if len(permissionIds) == 0 {
		return &emptypb.Empty{}, nil
	}
	for idx, permission := range in.GetPermissions() {
		permissionIds[idx] = permission.GetId()
	}
	err := s.svc.DeleteRolePermissions(ctx, in.GetOrganizationId(), in.GetId(), permissionIds)
	if err != nil {
		var statusErr model.ErrorStatus
		if errors.As(err, &statusErr) {
			return &emptypb.Empty{}, err
		}
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
