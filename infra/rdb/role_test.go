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

func TestRole(t *testing.T) {
	tenant := "test"
	ctx := context.Background()
	ctx, _ = tenants.SetTenantWithContext(ctx, tenant)
	mock.NewCase(mock.PostgreSQL, "role").Set(mock.Case{
		Name: "role query and command",
		Fn: func(store *rdb.SQLStore, db *gorm.DB) func(t *testing.T) {
			orgName, _ := model.NewName("test organization")
			permissionName, _ := model.NewName("create:permission")
			addPermissions := model.AddPermissionCommand{
				Name:        permissionName,
				Description: "test",
			}
			orgCmd := model.AddOrganizationCommand{
				Name:        orgName,
				Description: "organization test desc",
			}
			name, _ := model.NewName("admin")
			cmd := model.AddRoleCommand{
				Name:        name,
				Description: "administrator",
			}
			return func(t *testing.T) {
				t.Helper()
				if err := store.AddOrganizationCommand(ctx, &orgCmd); !assert.NoError(t, err) {
					return
				} else {
					cmd.Organization = orgCmd.Result
				}
				if err := store.AddPermissionCommand(ctx, &addPermissions); !assert.NoError(t, err) {
					return
				}
				t.Run("create", func(t *testing.T) {
					t.Helper()
					err := store.AddRoleCommand(ctx, &cmd)
					if !assert.NoError(t, err) {
						return
					}
					assert.NotEmpty(t, cmd.Result)
				})
				t.Run("find by id", func(t *testing.T) {
					t.Helper()
					query := model.GetRoleByIDQuery{
						Organization: cmd.Organization,
						PrimaryQuery: model.PrimaryQuery{
							ID: cmd.Result.ID,
						},
					}
					err := store.GetRoleByIDQuery(ctx, &query)
					if !assert.NoError(t, err) {
						return
					}
					p := query.Result
					assert.Equal(t, cmd.Result.ID.String(), p.ID.String())
					assert.Equal(t, cmd.Result.Name, p.Name)
					assert.Equal(t, cmd.Result.Description, p.Description)
				})
				t.Run("find all", func(t *testing.T) {
					t.Helper()
					query := model.GetRoleQuery{
						Organization: cmd.Organization,
					}
					err := store.GetRoleQuery(ctx, &query)
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
					updateCmd := model.UpdateRoleCommand{
						Organization: cmd.Organization,
						PrimaryCommand: model.PrimaryCommand{
							ID: cmd.Result.ID,
						},
						Name:        name,
						Description: "desc",
					}
					err := store.UpdateRoleCommand(ctx, &updateCmd)
					assert.NoError(t, err)
					t.Run("updated find by id", func(t *testing.T) {
						query := model.GetRoleByIDQuery{
							Organization: cmd.Organization,
							PrimaryQuery: model.PrimaryQuery{
								ID: cmd.Result.ID,
							},
						}
						err := store.GetRoleByIDQuery(ctx, &query)
						if !assert.NoError(t, err) {
							return
						}
						p := query.Result
						assert.Equal(t, cmd.Result.ID.String(), p.ID.String())
						assert.Equal(t, cmd.Result.Name, p.Name)
						assert.NotEqual(t, cmd.Result.Description, p.Description)
						assert.Equal(t, "desc", p.Description)
					})
				})
				t.Run("add permission", func(t *testing.T) {
					t.Helper()
					cmd := model.AddRolePermissionCommand{
						Role: &model.Role{
							ID: cmd.Result.ID,
						},
						Permissions: []model.Permission{*addPermissions.Result},
					}
					err := store.AddRolePermissionCommand(ctx, &cmd)
					assert.NoError(t, err)
				})
				t.Run("delete permission", func(t *testing.T) {
					t.Helper()
					cmd := model.DeleteRolePermissionCommand{
						Role: &model.Role{
							ID: cmd.Result.ID,
						},
						Permissions: []model.Permission{*addPermissions.Result},
					}
					err := store.DeleteRolePermissionCommand(ctx, &cmd)
					assert.NoError(t, err)
				})
				t.Run("delete", func(t *testing.T) {
					t.Helper()
					cmd := model.DeleteRoleCommand{
						Organization: cmd.Organization,
						PrimaryCommand: model.PrimaryCommand{
							ID: cmd.Result.ID,
						},
					}
					err := store.DeleteRoleCommand(ctx, &cmd)
					assert.NoError(t, err)
				})
			}
		},
	}).Run(t)
}
