package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/entity"
	"github.com/n-creativesystem/rbns/infra/entity/plugins"
	"github.com/n-creativesystem/rbns/infra/rdb/driver"
)

func (f *SQLStore) addPermissionBus() {
	bus.AddHandler("sql", f.GetPermissionQuery)
	bus.AddHandler("sql", f.GetPermissionByIDQuery)
	bus.AddHandler("sql", f.GetPermissionByNameQuery)
	bus.AddHandler("sql", f.CountPermissionByNameQuery)

	bus.AddHandler("sql", f.AddPermissionCommand)
	bus.AddHandler("sql", f.AddPermissionCommands)
	bus.AddHandler("sql", f.UpdatePermissionCommand)
	bus.AddHandler("sql", f.DeletePermissionCommand)
}

func (f *SQLStore) GetPermissionQuery(ctx context.Context, query *model.GetPermissionQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var permissions []entity.Permission
		where := entity.Permission{
			Model: entity.Model{
				Tenant: tenant,
			},
		}
		err := sess.Where(&where).Order("id").Find(&permissions).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		if len(permissions) == 0 {
			return model.ErrNoData
		}
		query.Result = make([]model.Permission, 0, len(permissions))
		for _, permission := range permissions {
			query.Result = append(query.Result, *permission.ConvertModel())
		}
		return nil
	})
}

func (f *SQLStore) GetPermissionByIDQuery(ctx context.Context, query *model.GetPermissionByIDQuery) error {
	queries := model.GetPermissionByIDsQuery{
		Query: []model.PrimaryCommand{query.PrimaryCommand},
	}
	if err := f.GetPermissionByIDsQuery(ctx, &queries); err != nil {
		return err
	}
	p := queries.Result[0]
	query.Result = &p
	return nil
}

func (f *SQLStore) GetPermissionByIDsQuery(ctx context.Context, query *model.GetPermissionByIDsQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var permissions []entity.Permission
		ids := make([]string, 0, len(query.Query))
		for _, q := range query.Query {
			ids = append(ids, q.ID.String())
		}
		err := sess.Where("tenant = ?", tenant).Where("id in ?", ids).First(&permissions).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		query.Result = make([]model.Permission, 0, len(permissions))
		for _, ep := range permissions {
			query.Result = append(query.Result, *ep.ConvertModel())
		}
		return nil
	})
}

func (f *SQLStore) GetPermissionByNameQuery(ctx context.Context, query *model.GetPermissionByNameQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var permission entity.Permission

		where := entity.Permission{
			Model: entity.Model{
				Tenant: tenant,
			},
			Name: query.Name.String(),
		}
		if err := sess.Where(&where).Find(&permission).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		query.Result = permission.ConvertModel()
		return nil
	})
}

func (f *SQLStore) CountPermissionByNameQuery(ctx context.Context, query *model.CountPermissionByNameQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var count int64

		whereObjs := make([]string, 0, len(query.Name))
		for _, name := range query.Name {
			whereObjs = append(whereObjs, name.String())
		}
		if err := sess.Where("tenant = ?", tenant).Where("name in ?", whereObjs).Model(&entity.Permission{}).Count(&count).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		query.Result = count
		return nil
	})
}

func (f *SQLStore) AddPermissionCommand(ctx context.Context, cmd *model.AddPermissionCommand) error {
	cmds := model.AddPermissionCommands{
		AddPermissions: []model.AddPermissionCommand{*cmd},
	}
	err := f.AddPermissionCommands(ctx, &cmds)
	if err != nil {
		return err
	}
	cmd.Result = cmds.AddPermissions[0].Result
	return nil
}

func (f *SQLStore) AddPermissionCommands(ctx context.Context, cmd *model.AddPermissionCommands) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		entities := make([]entity.Permission, 0, len(cmd.AddPermissions))
		for _, c := range cmd.AddPermissions {
			p := entity.Permission{
				Model: entity.Model{
					Tenant: tenant,
				},
				Name:        c.Name.String(),
				Description: c.Description,
			}
			p.Generate()
			entities = append(entities, p)
		}
		err := sess.Create(&entities).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		for idx, entity := range entities {
			cmd.AddPermissions[idx].Result = &model.Permission{
				ID:          entity.ID,
				Name:        entity.Name,
				Description: entity.Description,
			}
		}
		return nil
	})
}

func (f *SQLStore) UpdatePermissionCommand(ctx context.Context, cmd *model.UpdatePermissionCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		value := entity.Permission{
			Name:        cmd.Name.String(),
			Description: cmd.Description,
		}
		where := entity.Permission{
			Model: entity.Model{
				ID:     plugins.ID(cmd.ID.String()),
				Tenant: tenant,
			},
		}
		return driver.NewDBErr(sess.DB, sess.Where(&where).Updates(&value).Error)
	})
}

func (f *SQLStore) DeletePermissionCommand(ctx context.Context, cmd *model.DeletePermissionCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		where := entity.Permission{
			Model: entity.Model{
				ID:     plugins.ID(cmd.ID.String()),
				Tenant: tenant,
			},
		}
		db := sess.Where(&where).Delete(&entity.Permission{})
		if db.RowsAffected == 0 {
			return model.ErrNoData
		}
		if err := driver.NewDBErr(sess.DB, sess.Error); err != nil {
			return err
		}
		sess.events = append(sess.events, &cmd)
		return nil
	})
}
