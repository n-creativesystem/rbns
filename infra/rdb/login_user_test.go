package rdb_test

import (
	"context"
	"testing"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/rdb"
	"github.com/n-creativesystem/rbns/infra/rdb/mock"
	"github.com/n-creativesystem/rbns/ncsfw/tenants"
	"github.com/n-creativesystem/rbns/ncsfw/tracer"
	"github.com/n-creativesystem/rbns/version"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLoginUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx, _ = tenants.SetTenantWithContext(ctx, "dammy")
	_, _ = tracer.InitOpenTelemetryWithService(ctx, "test infra", tracer.Service{
		Name:    "rbns",
		Version: version.Version,
	})
	ctx, span := tracer.Start(ctx, "login user")
	defer span.End()
	mock.NewCase(mock.PostgreSQL, "login_user").Set(mock.Case{
		Name: "login user query and command",
		Fn: func(store *rdb.SQLStore, db *gorm.DB) func(t *testing.T) {
			return func(t *testing.T) {
				t.Run("save", func(t *testing.T) {
					t.Helper()
					cmd := &model.UpsertLoginUserCommand{
						User: &model.LoginUser{
							UserName: "user_name",
							Email:    "example@example.com",
							Role:     "Admin",
							Groups:   []string{},
						},
						SignupAllowed: true,
					}
					err := store.AddLoginUserCommand(ctx, cmd)
					assert.NoError(t, err)
				})
				t.Run("find", func(t *testing.T) {
					query := model.GetLoginUserByEmailQuery{
						Email: "example@example.com",
					}
					err := store.GetLoginUserByEmailQuery(ctx, &query)
					assert.NoError(t, err)
					assert.Equal(t, query.Result.UserName, "user_name")
				})
			}
		},
	}).Run(t)
}
