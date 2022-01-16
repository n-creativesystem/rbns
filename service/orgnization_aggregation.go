package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
)

type OrganizationAggregation interface {
	AddUsers(ctx context.Context, organizationId string, userId []string) error
	DeleteUsers(ctx context.Context, organizationId string, userId []string) error

	RoleCreate(ctx context.Context, organizationID string, names, descriptions []string) ([]model.Role, error)
	RoleFindById(ctx context.Context, organizationID string, roleId string) (*model.Role, error)
	RoleFindAll(ctx context.Context, organizationID string) ([]model.Role, error)
	RoleUpdate(ctx context.Context, organizationID string, roleId string, name, description string) error
	RoleDelete(ctx context.Context, organizationID string, roleId string) error
	GetRolePermissions(ctx context.Context, organizationID string, roleId string) ([]model.Permission, error)
	AddRolePermissions(ctx context.Context, organizationID string, roleId string, permissionIds []string) error
	DeleteRolePermissions(ctx context.Context, organizationID string, roleId string, permissionIds []string) error

	AddUserRoles(ctx context.Context, organizationId, userID string, roleIds []string) error
	DeleteUserRoles(ctx context.Context, organizationId, userID string, roleIds []string) error
}

type OrganizationAggregationImpl struct {
	permissionService   Permission
	organizationService Organization
	userService         User
	log                 logger.Logger
	telemetry           telemetryFunc
}

func NewOrganizationAggregation(permissionService Permission, organizationService Organization, userService User) *OrganizationAggregationImpl {
	return &OrganizationAggregationImpl{
		permissionService:   permissionService,
		organizationService: organizationService,
		userService:         userService,
		log:                 logger.New("organization aggregation service"),
		telemetry:           createSpanWithPrefix("organization aggregation service"),
	}
}

func (svc *OrganizationAggregationImpl) AddUsers(ctx context.Context, organizationID string, userId []string) (e error) {
	svc.telemetry(ctx, "add users", func(ctx context.Context) {
		org, err := svc.organizationService.FindById(ctx, organizationID)
		if err != nil {
			e = err
			return
		}
		users, err := svc.userService.FindByIds(ctx, userId)
		if err != nil {
			e = err
			return
		}

		cmd := model.AddOrganizationUserCommand{
			PrimaryCommand: model.PrimaryCommand{
				ID: org.ID,
			},
			User: append([]model.User{}, users...),
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
	})
	return
}

func (svc *OrganizationAggregationImpl) DeleteUsers(ctx context.Context, organizationID string, userId []string) (e error) {
	svc.telemetry(ctx, "", func(ctx context.Context) {
		org, err := svc.organizationService.FindById(ctx, organizationID)
		if err != nil {
			e = err
			return
		}
		users, err := svc.userService.FindByIds(ctx, userId)
		if err != nil {
			e = err
			return
		}

		cmd := model.DeleteOrganizationUserCommand{
			PrimaryCommand: model.PrimaryCommand{
				ID: org.ID,
			},
			User: append([]model.User{}, users...),
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
	})
	return
}

func (svc *OrganizationAggregationImpl) RoleCreate(ctx context.Context, organizationID string, names, descriptions []string) (out []model.Role, e error) {
	svc.telemetry(ctx, "role create", func(ctx context.Context) {
		org, err := svc.organizationService.FindById(ctx, organizationID)
		if err != nil {
			e = err
			return
		}
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
			e = err
			return
		}
		out = make([]model.Role, 0, len(cmd.Roles))
		for _, role := range cmd.Roles {
			out = append(out, *role.Result)
		}
	})
	return
}

func (svc *OrganizationAggregationImpl) RoleFindById(ctx context.Context, organizationID string, roleId string) (out *model.Role, e error) {
	svc.telemetry(ctx, "role find by id", func(ctx context.Context) {
		org, err := svc.organizationService.FindById(ctx, organizationID)
		if err != nil {
			e = err
			return
		}

		id, err := model.NewID(roleId)
		if err != nil {
			e = err
			return
		}
		query := model.GetRoleByIDQuery{
			Organization: org,
			PrimaryQuery: model.PrimaryQuery{
				ID: id,
			},
		}
		if err := bus.Dispatch(ctx, &query); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
			e = err
			return
		}
		out, e = query.Result, nil
	})
	return
}

func (svc *OrganizationAggregationImpl) RoleFindAll(ctx context.Context, organizationID string) (out []model.Role, e error) {
	svc.telemetry(ctx, "role find all", func(ctx context.Context) {
		org, err := svc.organizationService.FindById(ctx, organizationID)
		if err != nil {
			e = err
			return
		}
		query := model.GetRoleQuery{
			Organization: org,
		}
		if err := bus.Dispatch(ctx, &query); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "query", fmt.Sprintf("%+v", query))
			e = err
			return
		}
		out, err = query.Result, nil
	})
	return
}

func (svc *OrganizationAggregationImpl) RoleUpdate(ctx context.Context, organizationID string, roleId string, name, description string) (e error) {
	svc.telemetry(ctx, "role update", func(ctx context.Context) {
		org, err := svc.organizationService.FindById(ctx, organizationID)
		if err != nil {
			e = err
			return
		}
		n, err := model.NewName(name)
		if err != nil {
			e = err
			return
		}
		id, err := model.NewID(roleId)
		if err != nil {
			e = err
			return
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
			e = err
			return
		}
	})
	return
}

func (svc *OrganizationAggregationImpl) RoleDelete(ctx context.Context, organizationID string, roleId string) (e error) {
	svc.telemetry(ctx, "role delete", func(ctx context.Context) {
		id, err := model.NewID(roleId)
		if err != nil {
			e = err
			return
		}
		cmd := model.DeleteRoleCommand{
			PrimaryCommand: model.PrimaryCommand{
				ID: id,
			},
		}
		if err := bus.Dispatch(ctx, &cmd); err != nil {
			svc.log.ErrorWithContext(ctx, err, "dispatch error", "command", fmt.Sprintf("%+v", cmd))
			e = err
			return
		}
	})
	return
}

func (svc *OrganizationAggregationImpl) GetRolePermissions(ctx context.Context, organizationID string, roleId string) (out []model.Permission, e error) {
	svc.telemetry(ctx, "get role permissions", func(ctx context.Context) {
		role, err := svc.RoleFindById(ctx, organizationID, roleId)
		if err != nil {
			e = err
			return
		}
		out, e = role.Permissions, nil
	})
	return
}

func (svc *OrganizationAggregationImpl) AddRolePermissions(ctx context.Context, organizationID string, roleId string, permissionIds []string) (e error) {
	svc.telemetry(ctx, "add role permissions", func(ctx context.Context) {
		role, err := svc.RoleFindById(ctx, organizationID, roleId)
		if err != nil {
			e = err
			return
		}
		permissions, err := svc.permissionService.FindByIds(ctx, permissionIds)
		if err != nil {
			e = err
			return
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
			e = err
			return
		}
	})
	return
}

func (svc *OrganizationAggregationImpl) DeleteRolePermissions(ctx context.Context, organizationID string, roleId string, permissionIds []string) (e error) {
	svc.telemetry(ctx, "delete role permissions", func(ctx context.Context) {
		role, err := svc.RoleFindById(ctx, organizationID, roleId)
		if err != nil {
			e = err
			return
		}
		permissions, err := svc.permissionService.FindByIds(ctx, permissionIds)
		if err != nil {
			e = err
			return
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
			e = err
			return
		}
	})
	return
}

func (svc *OrganizationAggregationImpl) AddUserRoles(ctx context.Context, organizationId, userID string, roleIds []string) (e error) {
	svc.telemetry(ctx, "add user roles", func(ctx context.Context) {
		org, err := svc.organizationService.FindById(ctx, organizationId)
		if err != nil {
			e = err
			return
		}
		roles, err := svc.RoleFindAll(ctx, org.ID.String())
		id, err := model.NewID(userID)
		if err != nil {
			e = err
			return
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
			e = err
			return
		}
	})
	return
}

func (svc *OrganizationAggregationImpl) DeleteUserRoles(ctx context.Context, organizationId, userID string, roleIds []string) (e error) {
	svc.telemetry(ctx, "delete user roles", func(ctx context.Context) {
		org, err := svc.organizationService.FindById(ctx, organizationId)
		if err != nil {
			e = err
			return
		}
		roles, err := svc.RoleFindAll(ctx, org.ID.String())
		id, err := model.NewID(userID)
		if err != nil {
			e = err
			return
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
			e = err
			return
		}
	})
	return
}
