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
}

func (f *SQLStore) AddTenantCommand(ctx context.Context, cmd *model.AddTenantCommand) error {
	return f.inTransactionWithDbSession(ctx, func(sess *DBSession) error {
		entity := entity.Tenant{
			Name: cmd.Name,
		}
		entity.ID.Generate()
		err := sess.Create(&entity).Error
		if err != nil {
			return err
		}
		cmd.Result = &model.Tenant{
			ID:   entity.ID.String(),
			Name: cmd.Name,
		}
		return nil
	})
}

func (f *SQLStore) GetTenantByName(ctx context.Context, cmd *model.GetTenantByNameQuery) error {
	return f.DbSessionFunc(ctx, func(sess *DBSession) error {
		where := entity.Tenant{
			Name: cmd.Name,
		}
		var value entity.Tenant
		err := sess.Where(&where).Model(&entity.Tenant{}).First(&value).Error
		if err != nil {
			return err
		}
		cmd.Result = &model.Tenant{
			ID:   value.ID.String(),
			Name: value.Name,
		}
		return nil
	})
}

func (f *SQLStore) DeleteTenantCommand(ctx context.Context, cmd *model.DeleteTenantCommand) error {
	return f.inTransactionWithDbSession(ctx, func(sess *DBSession) error {
		return driver.NewDBErr(sess.DB, sess.Where("id = ?", cmd.ID).Delete(&entity.Tenant{}).Error)
	})
}
