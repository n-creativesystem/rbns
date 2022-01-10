package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/entity"
	"github.com/n-creativesystem/rbns/infra/entity/plugins"
)

func (f *SQLStore) addApiKeyBus() {
	bus.AddHandler("sql", f.AddApiKeyCommand)
	bus.AddHandler("sql", f.DeleteAPIKeyCommand)
	bus.AddHandler("sql", f.GetAPIKeyByNameQuery)
}

func (f *SQLStore) AddApiKeyCommand(ctx context.Context, cmd *model.AddApiKeyCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
		value := entity.ApiKey{
			Name: cmd.Name,
			Key:  cmd.HashedKey,
			Role: string(cmd.Role),
			Model: entity.Model{
				Tenant: tenant,
			},
		}
		value.Generate()
		if err := sess.Create(&value).Error; err != nil {
			return err
		}

		cmd.Result = &model.ApiKey{
			Id:   value.ID.String(),
			Name: value.Name,
			Key:  cmd.HashedKey,
			Role: model.RoleType(cmd.Role),
		}
		return nil
	})
}

func (f *SQLStore) DeleteAPIKeyCommand(ctx context.Context, cmd *model.DeleteAPIKeyCommand) error {
	return f.inTransactionWithToken(ctx, func(sess *DBSession, tenant string) error {
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
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		value := entity.ApiKey{
			Model: entity.Model{
				Tenant: tenant,
			},
			Key: cmd.KeyName,
		}
		var result entity.ApiKey
		err := sess.Where(&value).Model(&entity.ApiKey{}).First(&result).Error
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		cmd.Result = &model.ApiKey{
			Id:   result.ID.String(),
			Name: result.Name,
			Role: model.RoleType(result.Role),
			Key:  result.Key,
		}
		return nil
	})
}
