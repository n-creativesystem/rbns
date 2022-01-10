package rdb_test

import (
	"context"
	"testing"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/rdb"
	"github.com/n-creativesystem/rbns/infra/rdb/mock"
	"github.com/n-creativesystem/rbns/internal/contexts"
	"github.com/n-creativesystem/rbns/tracer"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLoginUser(t *testing.T) {
	ctx := context.Background()
	ctx = contexts.ToTenantContext(ctx, "dammy")
	tr, _ := tracer.InitOpenTelemetry("test infra")
	defer tr.Cleanup(ctx)
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
							ID:      "test",
							UseName: "user_name",
							Email:   "example@example.com",
							Role:    "Admin",
							Groups:  []string{},
						},
						SignupAllowed: true,
					}
					err := store.AddLoginUserCommand(ctx, cmd)
					assert.NoError(t, err)
				})
				t.Run("find", func(t *testing.T) {
					query := model.GetLoginUserByIDQuery{
						ID: "test",
					}
					err := store.GetLoginUserByIDQuery(ctx, &query)
					assert.NoError(t, err)
					assert.Equal(t, query.Result.UseName, "user_name")
				})
			}
		},
	}).Run(t)
}
