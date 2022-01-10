package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/entity"
	"github.com/n-creativesystem/rbns/infra/rdb/driver"
)

func (f *SQLStore) addLoginUserBus() {
	bus.AddHandler("sql", f.AddLoginUserCommand)
	bus.AddHandler("sql", f.GetLoginUserByIDQuery)
	bus.AddEventListenerCtx(f.PublishTenantCommand)
}

func (f *SQLStore) AddLoginUserCommand(ctx context.Context, cmd *model.UpsertLoginUserCommand) error {
	return f.inTransactionWithDbSession(ctx, func(sess *DBSession) error {
		value := entity.LoginUser{
			Model: entity.Model{
				Tenant: dammyTenant,
			},
			OAuthId:       cmd.User.OAuthID,
			UserName:      cmd.User.UseName,
			Email:         cmd.User.Email,
			Role:          cmd.User.Role,
			SignupAllowed: cmd.SignupAllowed,
			OAuthName:     cmd.User.GetOAuthName(),
			Password:      cmd.User.GetPassword(),
			OAuthToken:    cmd.User.GetToken(),
		}
		value.Generate()
		return sess.Save(&value).Error
	})
}

func (f *SQLStore) GetLoginUserByIDQuery(ctx context.Context, query *model.GetLoginUserByIDQuery) error {
	return f.DbSessionWithTenant(ctx, func(sess *DBSession, tenant string) error {
		var result entity.LoginUser
		if err := sess.Preload("Tenants").Where("id = ?", query.ID).First(&result).Error; err != nil {
			return err
		}
		tenants := make([]model.Tenant, 0, len(result.Tenants))
		for _, t := range result.Tenants {
			tenants = append(tenants, model.Tenant{
				ID:   t.ID.String(),
				Name: t.Name,
			})
		}
		loginUser := (&model.LoginUser{
			ID:      result.ID.String(),
			UseName: result.UserName,
			Email:   result.Email,
			Role:    result.Role,
			Tenant:  tenant,
			Tenants: tenants,
		}).SetOAuthName(result.OAuthName).SetOAuthToken(result.OAuthToken).SetPassword(result.Password)
		query.Result = loginUser
		return nil
	})
}

func (f *SQLStore) PublishTenantCommand(ctx context.Context, cmd *model.AddTenantAndLoginUserCommand) error {
	return f.inTransactionWithDbSession(ctx, func(sess *DBSession) error {
		var loginUser entity.LoginUser
		if err := sess.Where("id = ?", cmd.LoginUser.ID).First(&loginUser).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		var tenant entity.Tenant
		if err := sess.Where("tenant = ?", cmd.Tenant.ID).First(&tenant).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		loginUser.Tenant = cmd.Tenant.ID
		if err := sess.Save(&loginUser).Model(&loginUser).Association("Tenant").Append(&tenant); err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		return nil
	})
}
