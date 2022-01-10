package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/entity"
	"github.com/n-creativesystem/rbns/infra/entity/plugins"
	"github.com/n-creativesystem/rbns/infra/rdb/driver"
)

func (f *SQLStore) addUserBus() {
	bus.AddHandler("sql", f.GetUserQuery)
	bus.AddHandler("sql", f.GetUserByIDQuery)
	bus.AddHandler("sql", f.AddUserCommand)
	bus.AddHandler("sql", f.DeleteUserCommand)
	bus.AddHandler("sql", f.AddUserRoleCommand)
	bus.AddHandler("sql", f.DeleteUserRoleCommand)
	bus.AddHandler("sql", f.GetUserByIDsQuery)
}

func (f *SQLStore) GetUserQuery(ctx context.Context, query *model.GetUserQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var users []entity.User
		err := sess.
			Where(&entity.User{
				Model: entity.Model{
					Tenant: tenant,
				},
			}).
			Preload("Organization").
			Preload("Roles.Permissions").Find(&users).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		if len(users) == 0 {
			return model.ErrNoData
		}
		query.Result = make([]model.User, 0, len(users))
		for _, user := range users {
			query.Result = append(query.Result, *user.ConvertModel())
		}
		return nil
	})
}

func (f *SQLStore) GetUserByIDQuery(ctx context.Context, query *model.GetUserByIDQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var user entity.User
		err := sess.
			Where(&entity.User{
				Model: entity.Model{
					ID:     plugins.ID(query.ID.String()),
					Tenant: tenant,
				},
			}).
			Preload("Organization").
			Preload("Roles.Permissions").Find(&user).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		if user.ID.String() == "" {
			return model.ErrNoData
		}
		query.Result = user.ConvertModel()
		return nil
	})
}

func (f *SQLStore) GetUserByIDsQuery(ctx context.Context, query *model.GetUserByIDsQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		ids := make([]plugins.ID, 0, len(query.Query))
		for _, q := range query.Query {
			ids = append(ids, plugins.ID(q.ID.String()))
		}
		users, err := getUserByIDs(sess, ids, tenant)
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		query.Result = make([]model.User, 0, len(users))
		for _, u := range users {
			query.Result = append(query.Result, *u.ConvertModel())
		}
		return nil
	})
}

func (f *SQLStore) AddUserCommand(ctx context.Context, cmd *model.AddUserCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		eUser := entity.User{
			Model: entity.Model{
				ID:     plugins.ID(cmd.ID.String()),
				Tenant: tenant,
			},
		}
		err := sess.Create(&eUser).Error
		if err != nil {
			return err
		}
		cmd.Result = &model.User{
			ID:   eUser.ID.String(),
			Name: eUser.Name,
		}
		return nil
	})
}

func (f *SQLStore) DeleteUserCommand(ctx context.Context, cmd *model.DeleteUserCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		err := sess.
			Where(&entity.User{
				Model: entity.Model{
					ID:     plugins.ID(cmd.ID.String()),
					Tenant: tenant,
				},
			}).
			Delete(&entity.User{}).Error
		return driver.NewDBErr(sess.DB, err)
	})
}

func getUserByID(sess *DBSession, id plugins.ID, tenant string) (*entity.User, error) {
	var user entity.User
	err := sess.
		Where(&entity.User{
			Model: entity.Model{
				ID:     id,
				Tenant: tenant,
			},
		}).
		First(&user).Error
	if err != nil {
		return nil, driver.NewDBErr(sess.DB, err)
	}
	return &user, nil
}

func getUserByIDs(sess *DBSession, id []plugins.ID, tenant string) ([]entity.User, error) {
	var user []entity.User
	err := sess.
		Where("tenant = ?", tenant).
		Where("id in ?", id).
		Find(&user).Error
	if err != nil {
		return nil, driver.NewDBErr(sess.DB, err)
	}
	return user, nil
}

func (f *SQLStore) AddUserRoleCommand(ctx context.Context, cmd *model.AddUserRoleCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		org, err := getOrganizationByID(sess, plugins.ID(cmd.Organization.ID.String()), tenant)
		if err != nil {
			return err
		}
		user, err := getUserByID(sess, plugins.ID(cmd.ID.String()), tenant)
		if err != nil {
			return err
		}
		roles := make([]entity.Role, 0, len(cmd.Roles))
		for _, role := range cmd.Roles {
			roles = append(roles, entity.Role{
				Model: entity.Model{
					ID:     plugins.ID(role.ID.String()),
					Tenant: tenant,
				},
			})
		}
		// Organizationにも保存
		if err := driver.NewDBErr(sess.DB, sess.Model(&org).Association("Roles").Append(roles)); err != nil {
			return err
		}
		return driver.NewDBErr(sess.DB, sess.Model(&user).Association("Roles").Append(roles))
	})
}

func (f *SQLStore) DeleteUserRoleCommand(ctx context.Context, cmd *model.DeleteUserRoleCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		org, err := getOrganizationByID(sess, plugins.ID(cmd.Organization.ID.String()), tenant)
		if err != nil {
			return err
		}
		user, err := getUserByID(sess, plugins.ID(cmd.ID.String()), tenant)
		if err != nil {
			return err
		}
		roles := make([]entity.Role, 0, len(cmd.Roles))
		for _, role := range cmd.Roles {
			roles = append(roles, entity.Role{
				Model: entity.Model{
					ID:     plugins.ID(role.ID.String()),
					Tenant: tenant,
				},
			})
		}
		// Organizationからも削除
		if err := driver.NewDBErr(sess.DB, sess.Model(&org).Association("Roles").Delete(roles)); err != nil {
			return err
		}
		return driver.NewDBErr(sess.DB, sess.Model(&user).Association("Roles").Delete(roles))
	})
}
