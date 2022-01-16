package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/entity"
	"github.com/n-creativesystem/rbns/infra/rdb/driver"
)

func (f *SQLStore) addTenantBas() {
	bus.AddHandler("sql", f.AddTenantCommand)
	bus.AddHandler("sql", f.GetTenantByName)
	bus.AddHandler("sql", f.DeleteTenantCommand)
}

func (f *SQLStore) AddTenantCommand(ctx context.Context, cmd *model.AddTenantCommand) error {
	return f.inTransactionWithDbSession(ctx, func(ctx context.Context, sess *DBSession) error {
		var loginUser entity.LoginUser
		if err := sess.Where("email = ?", cmd.OwnerEmail).First(&loginUser).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}

		entity := entity.Tenant{
			Name:        cmd.Name.String(),
			Description: cmd.Description,
		}
		entity.ID.Generate()
		err := sess.Create(&entity).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		loginUser.Tenant = entity.ID.String()
		if err := sess.Model(&loginUser).Association("Tenants").Append(&entity); err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		cmd.Result = &model.Tenant{
			ID:          entity.ID.String(),
			Name:        cmd.Name.String(),
			Description: cmd.Description,
		}
		return nil
	})
}

func (f *SQLStore) GetTenantByName(ctx context.Context, cmd *model.GetTenantByNameQuery) error {
	return f.DbSessionFunc(ctx, func(ctx context.Context, sess *DBSession) error {
		where := entity.Tenant{
			Name: cmd.Name,
		}
		var value entity.Tenant
		err := sess.Where(&where).Model(&entity.Tenant{}).First(&value).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		cmd.Result = &model.Tenant{
			ID:   value.ID.String(),
			Name: value.Name,
		}
		return nil
	})
}

func (f *SQLStore) DeleteTenantCommand(ctx context.Context, cmd *model.DeleteTenantCommand) error {
	return f.inTransactionWithDbSession(ctx, func(ctx context.Context, sess *DBSession) error {
		if err := sess.Where("id = ?", cmd.ID).Delete(&entity.Tenant{}).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		return nil
	})
}

func getTenant(sess *DBSession, tenant string) (*entity.Tenant, error) {
	var result entity.Tenant
	if err := sess.Where("id = ?", tenant).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
