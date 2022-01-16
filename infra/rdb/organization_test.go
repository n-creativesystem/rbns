package rdb_test

import (
	"context"
	"testing"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/rdb"
	"github.com/n-creativesystem/rbns/infra/rdb/mock"
	"github.com/n-creativesystem/rbns/ncsfw/tenants"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestOrganization(t *testing.T) {
	ctx := context.Background()
	mock.NewCase(mock.PostgreSQL, "org").Set(mock.Case{
		Name: "Organization query and command",
		Fn: func(store *rdb.SQLStore, db *gorm.DB) func(t *testing.T) {
			name, _ := model.NewName("test")
			cmd := model.AddOrganizationCommand{
				Name:        name,
				Description: "organization test",
			}
			return func(t *testing.T) {
				tenant, err := testAddTenant(ctx)
				if !assert.NoError(t, err) {
					return
				}
				ctx, _ = tenants.SetTenantWithContext(ctx, tenant.ID)
				t.Run("create", func(t *testing.T) {
					t.Helper()
					err := store.AddOrganizationCommand(ctx, &cmd)
					assert.NoError(t, err)
					assert.NotEmpty(t, cmd.Result.ID.String())
					assert.Equal(t, "test", cmd.Result.Name)
					assert.Equal(t, "organization test", cmd.Result.Description)
				})
				t.Run("find by id", func(t *testing.T) {
					t.Helper()
					query := model.GetOrganizationByIDQuery{
						PrimaryQuery: model.PrimaryQuery{
							ID: cmd.Result.ID,
						},
					}
					err := store.GetOrganizationByIDQuery(ctx, &query)
					r := query.Result
					assert.NoError(t, err)
					assert.Equal(t, cmd.Result.ID.String(), r.ID.String())
					assert.Equal(t, "test", r.Name)
					assert.Equal(t, "organization test", r.Description)
				})
				t.Run("find by name", func(t *testing.T) {
					t.Helper()
					query := model.GetOrganizationByNameQuery{
						Name: name,
					}
					err := store.GetOrganizationByNameQuery(ctx, &query)
					r := query.Result
					assert.NoError(t, err)
					assert.Equal(t, cmd.Result.ID.String(), r.ID.String())
					assert.Equal(t, "test", r.Name)
					assert.Equal(t, "organization test", r.Description)
				})
				t.Run("find all", func(t *testing.T) {
					t.Helper()
					query := model.GetOrganizationQuery{}
					err := store.GetOrganizationQuery(ctx, &query)
					if !assert.NoError(t, err) {
						return
					}
					for _, org := range query.Result {
						assert.Equal(t, cmd.Result.ID.String(), org.ID.String())
						assert.Equal(t, cmd.Result.Name, org.Name)
						assert.Equal(t, cmd.Result.Description, org.Description)
					}
				})
				t.Run("update", func(t *testing.T) {
					cmd := model.UpdateOrganizationCommand{
						PrimaryCommand: model.PrimaryCommand{
							ID: cmd.Result.ID,
						},
						Name:        name,
						Description: "test org",
					}
					err := store.UpdateOrganizationCommand(ctx, &cmd)
					assert.NoError(t, err)
				})
				t.Run("delete", func(t *testing.T) {
					t.Helper()
					cmd := model.DeleteOrganizationCommand{
						PrimaryCommand: model.PrimaryCommand{
							ID: cmd.Result.ID,
						},
					}
					err := store.DeleteOrganizationCommand(ctx, &cmd)
					assert.NoError(t, err)
				})
			}
		},
	}).Run(t)
}
