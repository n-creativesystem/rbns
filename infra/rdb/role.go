package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/entity"
	"github.com/n-creativesystem/rbns/infra/entity/plugins"
	"github.com/n-creativesystem/rbns/infra/rdb/driver"
)

func (f *SQLStore) addRoleBus() {
	bus.AddHandler("sql", f.GetRoleQuery)
	bus.AddHandler("sql", f.GetRoleByIDQuery)
	bus.AddHandler("sql", f.CountRoleByNameQuery)

	bus.AddHandler("sql", f.AddRoleCommand)
	bus.AddHandler("sql", f.AddRoleCommands)
	bus.AddHandler("sql", f.UpdateRoleCommand)
	bus.AddHandler("sql", f.DeleteRoleCommand)
	bus.AddHandler("sql", f.AddRolePermissionCommand)
	bus.AddHandler("sql", f.DeleteRolePermissionCommand)
}

func (f *SQLStore) GetRoleQuery(ctx context.Context, query *model.GetRoleQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var roles []entity.Role
		org, err := getOrganizationByID(sess, plugins.ID(query.Organization.ID.String()), tenant)
		if err != nil {
			return err
		}

		err = sess.Model(org).Preload("Permissions").Association("Roles").Find(&roles)
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		if len(roles) == 0 {
			return model.ErrNoData
		}
		query.Result = make([]model.Role, 0, len(roles))
		for _, role := range roles {
			query.Result = append(query.Result, *role.ConvertModel())
		}
		return nil
	})
}

func (f *SQLStore) GetRoleByIDQuery(ctx context.Context, query *model.GetRoleByIDQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var role entity.Role
		where := entity.Role{
			Model: entity.Model{
				ID:     plugins.ID(query.ID.String()),
				Tenant: tenant,
			},
		}
		org, err := getOrganizationByID(sess, plugins.ID(query.Organization.ID.String()), tenant)
		if err != nil {
			return err
		}
		err = sess.Model(org).Where(&where).Preload("Permissions").Association("Roles").Find(&role)
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		if role.ID == "" {
			return model.ErrNoData
		}
		query.Result = role.ConvertModel()
		return nil

	})
}

func (f *SQLStore) CountRoleByNameQuery(ctx context.Context, query *model.CountRoleByNameQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var count int64
		whereObjs := make([]string, 0, len(query.Name))
		for _, name := range query.Name {
			whereObjs = append(whereObjs, name.String())
		}
		if err := sess.Where("tenant = ?", tenant).Where("name in ?", whereObjs).Count(&count).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		query.Result = count
		return nil

	})
}

func (f *SQLStore) AddRoleCommand(ctx context.Context, cmd *model.AddRoleCommand) error {
	cmds := model.AddRoleCommands{
		Organization: cmd.Organization,
		Roles:        []*model.AddRoleCommand{cmd},
	}
	err := f.AddRoleCommands(ctx, &cmds)
	if err != nil {
		return err
	}
	cmd.Result = cmds.Roles[0].Result
	return nil
}

func (f *SQLStore) AddRoleCommands(ctx context.Context, cmd *model.AddRoleCommands) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		org, err := getOrganizationByID(sess, plugins.ID(cmd.Organization.ID.String()), tenant)
		if err != nil {
			return err
		}
		entities := make([]entity.Role, 0, len(cmd.Roles))
		for _, role := range cmd.Roles {
			entity := entity.Role{
				Model: entity.Model{
					Tenant: tenant,
				},
				Name:        role.Name.String(),
				Description: role.Description,
			}
			entity.Generate()
			entities = append(entities, entity)
		}
		err = sess.Model(org).Association("Roles").Append(&entities)
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		for idx, entity := range entities {
			cmd.Roles[idx].Result = &model.Role{
				ID:          entity.ID,
				Name:        entity.Name,
				Description: entity.Description,
			}
		}
		return nil
	})
}

func (f *SQLStore) UpdateRoleCommand(ctx context.Context, cmd *model.UpdateRoleCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		value := entity.Role{
			Name:        cmd.Name.String(),
			Description: cmd.Description,
		}
		where := entity.Role{
			Model: entity.Model{
				ID:     plugins.ID(cmd.ID.String()),
				Tenant: tenant,
			},
		}
		return driver.NewDBErr(sess.DB, sess.Where(&where).Updates(&value).Error)
	})
}

func (f *SQLStore) DeleteRoleCommand(ctx context.Context, cmd *model.DeleteRoleCommand) error {
	cmds := model.DeleteRoleCommands{
		Organization: cmd.Organization,
		Roles:        []*model.DeleteRoleCommand{cmd},
	}
	return f.DeleteRoleCommands(ctx, &cmds)
}

func (f *SQLStore) DeleteRoleCommands(ctx context.Context, cmd *model.DeleteRoleCommands) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		org, err := getOrganizationByID(sess, plugins.ID(cmd.Organization.ID.String()), tenant)
		if err != nil {
			return err
		}
		roles := make([]entity.Role, 0, len(cmd.Roles))
		for _, cmd := range cmd.Roles {
			role := entity.Role{
				Model: entity.Model{
					ID:     plugins.ID(cmd.ID.String()),
					Tenant: tenant,
				},
			}
			roles = append(roles, role)
		}
		err = sess.Model(&org).Association("Roles").Delete(roles)
		return driver.NewDBErr(sess.DB, err)
	})
}

func getRoleByID(sess *DBSession, id plugins.ID, tenant string) (*entity.Role, error) {
	var role entity.Role
	err := sess.
		Where(&entity.Role{Model: entity.Model{ID: id, Tenant: tenant}}).
		First(&role).Error
	if err != nil {
		return nil, driver.NewDBErr(sess.DB, err)
	}
	return &role, nil
}

func (f *SQLStore) AddRolePermissionCommand(ctx context.Context, cmd *model.AddRolePermissionCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		role, err := getRoleByID(sess, plugins.ID(cmd.ID.String()), tenant)
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		permissions := make([]interface{}, 0, len(cmd.Permissions))
		for _, permission := range cmd.Permissions {
			permissions = append(permissions, &entity.Permission{
				Model: entity.Model{
					ID:     plugins.ID(permission.ID.String()),
					Tenant: role.Tenant,
				},
			})
		}
		return driver.NewDBErr(sess.DB, sess.Model(role).Association("Permissions").Append(permissions...))
	})
}

func (f *SQLStore) DeleteRolePermissionCommand(ctx context.Context, cmd *model.DeleteRolePermissionCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		role, err := getRoleByID(sess, plugins.ID(cmd.ID.String()), tenant)
		if err != nil {
			return err
		}
		permissions := make([]entity.Permission, 0, len(cmd.Permissions))
		for _, permission := range cmd.Permissions {
			permissions = append(permissions, entity.Permission{
				Model: entity.Model{
					ID:     plugins.ID(permission.ID.String()),
					Tenant: role.Tenant,
				},
			})
		}
		return driver.NewDBErr(sess.DB, sess.Model(&role).Association("Permissions").Delete(permissions))
	})
}
