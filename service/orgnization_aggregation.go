package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/logger"
)

type OrganizationAggregation interface {
	AddUsers(ctx context.Context, strId string, userId []string) error
	DeleteUsers(ctx context.Context, strId string, userId []string) error

	RoleCreate(ctx context.Context, organizationID string, names, descriptions []string) ([]model.Role, error)
	RoleFindById(ctx context.Context, organizationID string, strId string) (*model.Role, error)
	RoleFindAll(ctx context.Context, organizationID string) ([]model.Role, error)
	RoleUpdate(ctx context.Context, organizationID string, strId string, name, description string) error
	RoleDelete(ctx context.Context, organizationID string, strId string) error
	GetRolePermissions(ctx context.Context, organizationID string, strId string) ([]model.Permission, error)
	AddRolePermissions(ctx context.Context, organizationID string, strId string, permissionIds []string) error
	DeleteRolePermissions(ctx context.Context, organizationID string, strId string, permissionIds []string) error

	AddUserRoles(ctx context.Context, organizationId, userID string, roleIds []string) error
	DeleteUserRoles(ctx context.Context, organizationId, userID string, roleIds []string) error
}

type organizationAggregation struct {
	permissionService   PermissionService
	organizationService OrganizationService
	userService         UserService
	log                 logger.Logger
}

func NewOrganizationAggregation(permissionService PermissionService, organizationService OrganizationService, userService UserService) OrganizationAggregation {
	return &organizationAggregation{
		permissionService:   permissionService,
		organizationService: organizationService,
		userService:         userService,
		log:                 logger.New("organization aggregation service"),
	}
}

func (svc *organizationAggregation) AddUsers(ctx context.Context, strId string, userId []string) error {
	org, err := svc.organizationService.FindById(ctx, strId)
	if err != nil {
		return err
	}
	users, err := svc.userService.FindByIds(ctx, userId)
	if err != nil {
		return err
	}

	cmd := model.AddOrganizationUserCommand{
		PrimaryCommand: model.PrimaryCommand{
			ID: org.ID,
		},
		User: append([]model.User{}, users...),
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}

func (svc *organizationAggregation) DeleteUsers(ctx context.Context, strId string, userId []string) error {
	org, err := svc.organizationService.FindById(ctx, strId)
	if err != nil {
		return err
	}
	users, err := svc.userService.FindByIds(ctx, userId)
	if err != nil {
		return err
	}

	cmd := model.DeleteOrganizationUserCommand{
		PrimaryCommand: model.PrimaryCommand{
			ID: org.ID,
		},
		User: append([]model.User{}, users...),
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}

func (svc *organizationAggregation) RoleCreate(ctx context.Context, organizationID string, names, descriptions []string) ([]model.Role, error) {
	org, err := svc.organizationService.FindById(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	var out []model.Role
	cmd := model.AddRoleCommands{
		Organization: org,
		Roles:        make([]*model.AddRoleCommand, 0, len(names)),
	}
	for idx, name := range names {
		n, err := model.NewName(name)
		if err != nil {
			continue
		}
		cmd.Roles = append(cmd.Roles, &model.AddRoleCommand{
			Name:        n,
			Description: descriptions[idx],
		})
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return nil, err
	}
	out = make([]model.Role, 0, len(cmd.Roles))
	for _, role := range cmd.Roles {
		out = append(out, *role.Result)
	}
	return out, nil
}

func (svc *organizationAggregation) RoleFindById(ctx context.Context, organizationID string, strId string) (*model.Role, error) {
	org, err := svc.organizationService.FindById(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	id, err := model.NewID(strId)
	if err != nil {
		return nil, err
	}
	query := model.GetRoleByIDQuery{
		Organization: org,
		PrimaryQuery: model.PrimaryQuery{
			ID: id,
		},
	}
	if err := bus.Dispatch(ctx, &query); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
		return nil, err
	}
	return query.Result, nil
}

func (svc *organizationAggregation) RoleFindAll(ctx context.Context, organizationID string) ([]model.Role, error) {
	org, err := svc.organizationService.FindById(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	query := model.GetRoleQuery{
		Organization: org,
	}
	if err := bus.Dispatch(ctx, &query); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
		return nil, err
	}
	return query.Result, nil
}

func (svc *organizationAggregation) RoleUpdate(ctx context.Context, organizationID string, strId string, name, description string) error {
	org, err := svc.organizationService.FindById(ctx, organizationID)
	if err != nil {
		return err
	}
	n, err := model.NewName(name)
	if err != nil {
		return err
	}
	id, err := model.NewID(strId)
	if err != nil {
		return err
	}
	cmd := model.UpdateRoleCommand{
		Organization: org,
		PrimaryCommand: model.PrimaryCommand{
			ID: id,
		},
		Name:        n,
		Description: description,
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}

func (svc *organizationAggregation) RoleDelete(ctx context.Context, organizationID string, strId string) error {
	id, err := model.NewID(strId)
	if err != nil {
		return err
	}
	cmd := model.DeleteRoleCommand{
		PrimaryCommand: model.PrimaryCommand{
			ID: id,
		},
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}

func (svc *organizationAggregation) GetRolePermissions(ctx context.Context, organizationID string, strId string) ([]model.Permission, error) {
	role, err := svc.RoleFindById(ctx, organizationID, strId)
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

func (svc *organizationAggregation) AddRolePermissions(ctx context.Context, organizationID string, strId string, permissionIds []string) error {
	role, err := svc.RoleFindById(ctx, organizationID, strId)
	if err != nil {
		return err
	}
	permissions, err := svc.permissionService.FindByIds(ctx, permissionIds)
	if err != nil {
		return err
	}

	cmd := model.AddRolePermissionCommand{
		Role:        role,
		Permissions: make([]model.Permission, 0, len(permissions)),
	}
	for _, p := range permissions {
		cmd.Permissions = append(cmd.Permissions, p)
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}

func (svc *organizationAggregation) DeleteRolePermissions(ctx context.Context, organizationID string, strId string, permissionIds []string) error {
	role, err := svc.RoleFindById(ctx, organizationID, strId)
	if err != nil {
		return err
	}
	permissions, err := svc.permissionService.FindByIds(ctx, permissionIds)
	if err != nil {
		return err
	}

	cmd := model.DeleteRolePermissionCommand{
		Role:        role,
		Permissions: make([]model.Permission, 0, len(permissions)),
	}
	for _, p := range permissions {
		cmd.Permissions = append(cmd.Permissions, p)
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}

func (svc *organizationAggregation) AddUserRoles(ctx context.Context, organizationId, userID string, roleIds []string) error {
	org, err := svc.organizationService.FindById(ctx, organizationId)
	if err != nil {
		return err
	}
	roles, err := svc.RoleFindAll(ctx, org.ID.String())
	id, err := model.NewID(userID)
	if err != nil {
		return err
	}
	cmd := model.AddUserRoleCommand{
		Organization: org,
		PrimaryCommand: model.PrimaryCommand{
			ID: id,
		},
		Roles: make([]model.Role, 0, len(roleIds)),
	}
	sort.Strings(roleIds)
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].ID.String() < roles[j].ID.String()
	})

	for _, roleId := range roleIds {
		for _, role := range roles {
			if role.ID.String() == roleId {
				cmd.Roles = append(cmd.Roles, role)
				break
			}
		}
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}

func (svc *organizationAggregation) DeleteUserRoles(ctx context.Context, organizationId, userID string, roleIds []string) error {
	org, err := svc.organizationService.FindById(ctx, organizationId)
	if err != nil {
		return err
	}
	roles, err := svc.RoleFindAll(ctx, org.ID.String())
	id, err := model.NewID(userID)
	if err != nil {
		return err
	}
	cmd := model.DeleteUserRoleCommand{
		Organization: org,
		PrimaryCommand: model.PrimaryCommand{
			ID: id,
		},
		Roles: make([]model.Role, 0, len(roleIds)),
	}
	sort.Strings(roleIds)
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].ID.String() < roles[j].ID.String()
	})

	for _, roleId := range roleIds {
		for _, role := range roles {
			if role.ID.String() == roleId {
				cmd.Roles = append(cmd.Roles, role)
				break
			}
		}
	}
	if err := bus.Dispatch(ctx, &cmd); err != nil {
		svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
		return err
	}
	return nil
}
