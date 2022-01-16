package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/entity"
	"github.com/n-creativesystem/rbns/infra/entity/plugins"
	"github.com/n-creativesystem/rbns/infra/rdb/driver"
)

func (f *SQLStore) addApiKeyBus() {
	bus.AddHandler("sql", f.AddApiKeyCommand)
	bus.AddHandler("sql", f.DeleteAPIKeyCommand)
	bus.AddHandler("sql", f.GetAPIKeyByNameQuery)
}

func (f *SQLStore) AddApiKeyCommand(ctx context.Context, cmd *model.AddApiKeyCommand) error {
	return f.inTransactionWithToken(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		value := entity.ApiKey{
			Name: cmd.Name,
			Key:  cmd.HashedKey,
			Role: uint(cmd.Role),
			Model: entity.Model{
				Tenant: tenant,
			},
		}
		value.Generate()
		if err := sess.Create(&value).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}

		cmd.Result = &model.ApiKey{
			Id:   value.ID.String(),
			Name: value.Name,
			Key:  cmd.HashedKey,
			Role: cmd.Role,
		}
		return nil
	})
}

func (f *SQLStore) DeleteAPIKeyCommand(ctx context.Context, cmd *model.DeleteAPIKeyCommand) error {
	return f.inTransactionWithToken(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		value := entity.ApiKey{
			Model: entity.Model{
				ID:     plugins.ID(cmd.ID.String()),
				Tenant: tenant,
			},
		}
		return sess.Delete(&value).Error
	})
}

func (f *SQLStore) GetAPIKeyByNameQuery(ctx context.Context, cmd *model.GetAPIKeyByNameQuery) error {
	return f.DbSessionWithTenant(ctx, func(ctx context.Context, sess *DBSession, tenant string) error {
		value := entity.ApiKey{
			Model: entity.Model{
				Tenant: tenant,
			},
			Key: cmd.KeyName,
		}
		var result entity.ApiKey
		err := sess.Where(&value).Model(&entity.ApiKey{}).First(&result).Error
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		if err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		cmd.Result = &model.ApiKey{
			Id:   result.ID.String(),
			Name: result.Name,
			Role: model.RoleLevel(result.Role),
			Key:  result.Key,
		}
		return nil
	})
}
