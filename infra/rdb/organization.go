package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/entity"
	"github.com/n-creativesystem/rbns/infra/entity/plugins"
	"github.com/n-creativesystem/rbns/infra/rdb/driver"
)

func (f *SQLStore) addOrganizationBus() {
	bus.AddHandler("sql", f.GetOrganizationQuery)
	bus.AddHandler("sql", f.GetOrganizationByIDQuery)
	bus.AddHandler("sql", f.GetOrganizationByNameQuery)
	bus.AddHandler("sql", f.CountOrganizationByNameQuery)
	bus.AddHandler("sql", f.AddOrganizationCommand)
	bus.AddHandler("sql", f.UpdateOrganizationCommand)
	bus.AddHandler("sql", f.DeleteOrganizationCommand)
}

func (f *SQLStore) GetOrganizationQuery(ctx context.Context, query *model.GetOrganizationQuery) error {
	return f.DbSessionWithTenant(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		var organizations []entity.Organization
		where := entity.Organization{
			Model: entity.Model{
				Tenant: tenant,
			},
		}
		err := userPreload(sess).Where(&where).Order("id").Find(&organizations).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		if len(organizations) == 0 {
			return model.ErrNoDataFound
		}
		query.Result = make([]model.Organization, 0, len(organizations))
		for _, organization := range organizations {
			query.Result = append(query.Result, *organization.ConvertModel())
		}
		return nil
	})
}

func (f *SQLStore) GetOrganizationByIDQuery(ctx context.Context, query *model.GetOrganizationByIDQuery) error {
	return f.DbSessionWithTenant(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		organization := entity.Organization{
			Model: entity.Model{
				ID:     plugins.ID(query.ID.String()),
				Tenant: tenant,
			},
		}
		err := userPreload(sess).Order("id").Find(&organization).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		query.Result = organization.ConvertModel()
		return nil
	})
}

func (f *SQLStore) GetOrganizationByNameQuery(ctx context.Context, query *model.GetOrganizationByNameQuery) error {
	return f.DbSessionWithTenant(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		var org entity.Organization
		var err error

		where := entity.Organization{
			Model: entity.Model{
				Tenant: tenant,
			},
			Name: query.Name.String(),
		}
		if err = sess.Preload("Users").Preload("Roles").Where(&where).Find(&org).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		if org.ID == "" {
			return model.ErrNoDataFound
		}
		query.Result = org.ConvertModel()
		return nil
	})
}

func (f *SQLStore) CountOrganizationByNameQuery(ctx context.Context, query *model.CountOrganizationByNameQuery) error {
	return f.DbSessionWithTenant(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
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

func (f *SQLStore) AddOrganizationCommand(ctx context.Context, cmd *model.AddOrganizationCommand) error {
	return f.inTransactionWithToken(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		entity := entity.Organization{
			Model: entity.Model{
				Tenant: tenant,
			},
			Name:        cmd.Name.String(),
			Description: cmd.Description,
		}
		entity.Generate()
		err := sess.Create(&entity).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		cmd.Result = &model.Organization{
			ID:          entity.ID,
			Name:        entity.Name,
			Description: entity.Description,
		}
		return nil
	})
}

func (f *SQLStore) UpdateOrganizationCommand(ctx context.Context, cmd *model.UpdateOrganizationCommand) error {
	return f.inTransactionWithToken(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		value := entity.Organization{
			Model: entity.Model{
				Tenant: tenant,
			},
			Name:        cmd.Name.String(),
			Description: cmd.Description,
		}
		return driver.NewDBErr(sess.DB, sess.Where(&entity.Organization{Model: entity.Model{ID: plugins.ID(cmd.ID.String())}}).Updates(&value).Error)
	})
}

func (f *SQLStore) DeleteOrganizationCommand(ctx context.Context, cmd *model.DeleteOrganizationCommand) error {
	return f.inTransactionWithToken(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		entity := entity.Organization{
			Model: entity.Model{
				ID:     plugins.ID(cmd.ID.String()),
				Tenant: tenant,
			},
		}
		err := sess.Delete(&entity).Error
		return driver.NewDBErr(sess.DB, err)
	})
}

func (f *SQLStore) AddOrganizationUserCommand(ctx context.Context, cmd *model.AddOrganizationUserCommand) error {
	return f.inTransactionWithToken(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		org, err := getOrganizationByID(sess, plugins.ID(cmd.ID.String()), tenant)
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		users := make([]entity.User, 0, len(cmd.User))
		for _, u := range cmd.User {
			users = append(users, entity.User{
				Model: entity.Model{
					ID:     plugins.ID(u.ID),
					Tenant: org.Tenant,
				},
			})
		}
		return driver.NewDBErr(sess.DB, sess.Model(org).Association("Users").Append(&users))
	})
}

func (f *SQLStore) DeleteOrganizationUserCommand(ctx context.Context, cmd *model.DeleteOrganizationUserCommand) error {
	return f.inTransactionWithToken(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		org, err := getOrganizationByID(sess, plugins.ID(cmd.ID.String()), tenant)
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		users := make([]entity.User, 0, len(cmd.User))
		for _, u := range cmd.User {
			users = append(users, entity.User{
				Model: entity.Model{
					ID:     plugins.ID(u.ID),
					Tenant: org.Tenant,
				},
			})
		}
		return driver.NewDBErr(sess.DB, sess.Model(org).Association("Users").Delete(&users))
	})
}

func getOrganizationByID(sess *DBSession, id plugins.ID, tenant string) (*entity.Organization, error) {
	var organization entity.Organization
	err := sess.
		Where(&entity.Organization{Model: entity.Model{ID: id, Tenant: tenant}}).
		First(&organization).Error
	if err != nil {
		return nil, driver.NewDBErr(sess.DB, err)
	}
	return &organization, nil
}

func userPreload(sess *DBSession) *DBSession {
	return &DBSession{DB: sess.Preload("Users").Preload("Roles")}
}
