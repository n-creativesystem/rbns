package rdb_test

import (
	"context"
	"testing"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/rdb"
	"github.com/n-creativesystem/rbns/infra/rdb/mock"
	"github.com/n-creativesystem/rbns/internal/contexts"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPermission(t *testing.T) {
	const tenant = "test"
	ctx := context.Background()
	ctx = contexts.ToTenantContext(ctx, tenant)
	mock.NewCase(mock.PostgreSQL, "permission").Set(mock.Case{
		Name: "Permission Query and Command",
		Fn: func(store *rdb.SQLStore, db *gorm.DB) func(t *testing.T) {
			return func(t *testing.T) {
				t.Helper()
				name, _ := model.NewName("create:permission")
				cmd := &model.AddPermissionCommand{
					Name:        name,
					Description: "test",
				}
				t.Run("create", func(t *testing.T) {
					t.Helper()
					err := store.AddPermissionCommand(ctx, cmd)
					if !assert.NoError(t, err) {
						return
					}
					assert.NotEmpty(t, cmd.Result.ID)
				})
				t.Run("find by id", func(t *testing.T) {
					t.Helper()
					query := model.GetPermissionByIDQuery{
						PrimaryCommand: model.PrimaryCommand{
							ID: cmd.Result.ID,
						},
					}
					err := store.GetPermissionByIDQuery(ctx, &query)
					if !assert.NoError(t, err) {
						return
					}
					assert.Equal(t, cmd.Result.ID.String(), query.Result.ID.String())
					assert.Equal(t, cmd.Result.Name, query.Result.Name)
					assert.Equal(t, cmd.Result.Description, query.Result.Description)
				})
				t.Run("find by name", func(t *testing.T) {
					t.Helper()
					query := model.GetPermissionByNameQuery{
						Name: name,
					}
					err := store.GetPermissionByNameQuery(ctx, &query)
					if !assert.NoError(t, err) {
						return
					}
					assert.Equal(t, cmd.Result.ID.String(), query.Result.ID.String())
					assert.Equal(t, cmd.Result.Name, query.Result.Name)
					assert.Equal(t, cmd.Result.Description, query.Result.Description)
				})
				t.Run("find all", func(t *testing.T) {
					t.Helper()
					query := model.GetPermissionQuery{}
					err := store.GetPermissionQuery(ctx, &query)
					if !assert.NoError(t, err) {
						return
					}
					res := query.Result
					for _, r := range res {
						assert.Equal(t, cmd.Result.ID.String(), r.ID.String())
						assert.Equal(t, cmd.Result.Name, r.Name)
						assert.Equal(t, cmd.Result.Description, r.Description)
					}
				})
				t.Run("update", func(t *testing.T) {
					t.Helper()
					updateCmd := model.UpdatePermissionCommand{
						PrimaryCommand: model.PrimaryCommand{
							ID: cmd.Result.ID,
						},
						Name:        name,
						Description: "test desc",
					}
					err := store.UpdatePermissionCommand(ctx, &updateCmd)
					assert.NoError(t, err)
					t.Run("updated find by id", func(t *testing.T) {
						t.Helper()
						query := model.GetPermissionByIDQuery{
							PrimaryCommand: model.PrimaryCommand{
								ID: cmd.Result.ID,
							},
						}
						err = store.GetPermissionByIDQuery(ctx, &query)
						if !assert.NoError(t, err) {
							return
						}
						assert.Equal(t, cmd.Result.ID.String(), query.Result.ID.String())
						assert.Equal(t, cmd.Result.Name, query.Result.Name)
						assert.NotEqual(t, cmd.Result.Description, query.Result.Description)
						assert.Equal(t, "test desc", query.Result.Description)
					})
				})
				t.Run("delete", func(t *testing.T) {
					deleteCmd := model.DeletePermissionCommand{
						PrimaryCommand: model.PrimaryCommand{
							ID: cmd.Result.ID,
						},
					}
					err := store.DeletePermissionCommand(ctx, &deleteCmd)
					assert.NoError(t, err)
					t.Run("deleted find by id", func(t *testing.T) {
						t.Helper()
						query := model.GetPermissionByIDQuery{
							PrimaryCommand: model.PrimaryCommand{
								ID: cmd.Result.ID,
							},
						}
						err = store.GetPermissionByIDQuery(ctx, &query)
						if !assert.Error(t, err) {
							return
						}
						assert.ErrorIs(t, err, model.ErrNoData)
					})
				})
			}
		},
	}).Run(t)
}
