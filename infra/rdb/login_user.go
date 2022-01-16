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
	// bus.AddHandler("sql", f.GetLoginUserByIDQuery)
	bus.AddHandler("sql", f.GetLoginUserByEmailQuery)
}

func (f *SQLStore) AddLoginUserCommand(ctx context.Context, cmd *model.UpsertLoginUserCommand) error {
	return f.inTransactionWithDbSession(ctx, func(ctx context.Context, sess *DBSession) error {
		var count int64
		if err := sess.Model(&entity.LoginUser{}).Where("email = ?", cmd.User.Email).Count(&count).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		if count == 0 {
			value := entity.LoginUser{
				OAuthId:       cmd.User.OAuthID,
				UserName:      cmd.User.UserName,
				Email:         cmd.User.Email,
				Role:          cmd.User.Role,
				SignupAllowed: cmd.SignupAllowed,
				OAuthName:     cmd.User.GetOAuthName(),
				Password:      cmd.User.GetPassword(),
				OAuthToken:    cmd.User.GetToken(),
			}
			return driver.NewDBErr(sess.DB, sess.Save(&value).Error)
		}
		return nil
	})
}

func (f *SQLStore) GetLoginUserByEmailQuery(ctx context.Context, query *model.GetLoginUserByEmailQuery) error {
	return f.DbSessionFunc(ctx, func(ctx context.Context, sess *DBSession) error {
		var result entity.LoginUser
		if err := sess.Where("email = ?", query.Email).Preload("Tenants").First(&result).Error; err != nil {
			return driver.NewDBErr(sess.DB, err)
		}
		loginUser := (&model.LoginUser{
			// ID:      result.ID.String(),
			OAuthID:  result.OAuthId,
			UserName: result.UserName,
			Email:    result.Email,
			Role:     result.Role,
			Tenants:  make([]model.Tenant, 0, len(result.Tenants)),
		}).SetOAuthName(result.OAuthName).SetOAuthToken(result.OAuthToken).SetPassword(result.Password).SetTenant(result.Tenant)
		for _, tenant := range result.Tenants {
			loginUser.Tenants = append(loginUser.Tenants, model.Tenant{
				ID:          tenant.ID.String(),
				Name:        tenant.Name,
				Description: tenant.Description,
			})
		}
		query.Result = loginUser
		return nil
	})
}
