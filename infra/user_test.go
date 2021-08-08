package infra_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/infra"
	"github.com/n-creativesystem/rbns/infra/dao"
	"github.com/n-creativesystem/rbns/tests"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	ctx := context.Background()
	cases := tests.MocksByPostgres{
		{
			Name: "create",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "users" ("key","organization_id") VALUES ($1,$2)  ON CONFLICT DO NOTHING`),
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					tx := repo.NewConnection().Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						orgId, _ := model.NewID(tests.IDs[0])
						user, _ := model.NewUser("user1", nil, nil)
						_, err := tx.User().Create(orgId, user)
						return err
					})
					assert.NoError(t, err)
				}
			},
		},
		{
			Name: "add role",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "user_roles" ("user_key","role_id","organization_id") VALUES ($1,$2,$3),($4,$5,$6) ON CONFLICT DO NOTHING`),
				).WithArgs("user1", tests.IDs[1], tests.IDs[0], "user1", tests.IDs[2], tests.IDs[0]).WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					tx := repo.NewConnection().Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						orgId, _ := model.NewID(tests.IDs[0])
						userKey, _ := model.NewKey("user1")
						role, _ := model.NewRole(tests.IDs[1], "admin", "administrator", nil)
						role2, _ := model.NewRole(tests.IDs[2], "view", "viewer", nil)
						return tx.User().AddRole(orgId, userKey, *role, *role2)
					})
					assert.NoError(t, err)
				}
			},
		},
		{
			Name: "delete role",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`DELETE FROM "user_roles" WHERE "user_roles"."user_key" = $1 AND "user_roles"."role_id" = $2 AND "user_roles"."organization_id" = $3`),
				).WithArgs("user1", tests.IDs[1], tests.IDs[0]).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					tx := repo.NewConnection().Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						orgId, _ := model.NewID(tests.IDs[0])
						userKey, _ := model.NewKey("user1")
						role, _ := model.NewID(tests.IDs[1])
						return tx.User().DeleteRole(orgId, userKey, role)
					})
					assert.NoError(t, err)
				}
			},
		},
		{
			Name: "findByKey",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				usersRow := sqlmock.NewRows([]string{"organization_id", "key"}).AddRow(tests.IDs[0], "user1")
				orgRow := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(tests.IDs[0], "test", "organization test")
				userRoleRow := sqlmock.NewRows([]string{"user_key", "role_id"}).AddRow("user1", tests.IDs[1]).AddRow("user1", tests.IDs[2])
				roleRow := sqlmock.
					NewRows([]string{"id", "name", "description"}).
					AddRow(tests.IDs[1], "admin", "administrator").
					AddRow(tests.IDs[2], "view", "viewer")
				rolePermissionRow := sqlmock.NewRows([]string{"role_id", "permission_id"}).AddRow(tests.IDs[1], tests.IDs[3]).AddRow(tests.IDs[2], tests.IDs[4])
				permissionRow := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(tests.IDs[3], "create:permission", "create").AddRow(tests.IDs[4], "read:permission", "read")
				// users find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."key" = $1 AND "users"."organization_id" = $2`),
				).WithArgs("user1", tests.IDs[0]).WillReturnRows(usersRow)

				// organizations find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."id" = $1`),
				).WithArgs(tests.IDs[0]).WillReturnRows(orgRow)

				// user_roles find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "user_roles" WHERE "user_roles"."user_key" = $1`),
				).WithArgs("user1").WillReturnRows(userRoleRow)

				// roles find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "roles" WHERE "roles"."id" IN ($1,$2)`),
				).WithArgs(tests.IDs[1], tests.IDs[2]).WillReturnRows(roleRow)

				// role_permissions find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE "role_permissions"."role_id" IN ($1,$2)`),
				).WithArgs(tests.IDs[1], tests.IDs[2]).WillReturnRows(rolePermissionRow)

				// permissions find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."id" IN ($1,$2)`),
				).WithArgs(tests.IDs[3], tests.IDs[4]).WillReturnRows(permissionRow)

				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					orgId, _ := model.NewID(tests.IDs[0])
					userKey, _ := model.NewKey("user1")
					u, err := repo.NewConnection().User(ctx).FindByKey(orgId, userKey)
					assert.NoError(t, err)
					assert.Equal(t, *userKey.Value(), u.GetKey())
					assert.NotEmpty(t, u.GetRole())
					assert.NotEmpty(t, u.GetPermission())
					assert.Equal(t, tests.IDs[1], *u.GetRole()[0].GetID())
					assert.Equal(t, tests.IDs[2], *u.GetRole()[1].GetID())
					assert.Equal(t, tests.IDs[3], *u.GetPermission()[0].GetID())
					assert.Equal(t, tests.IDs[4], *u.GetPermission()[1].GetID())
				}
			},
		},
		{
			Name: "findAll",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				user1 := "user1"
				user2 := "user2"
				orgId := tests.IDs[0]
				admin := tests.IDs[1]
				view := tests.IDs[2]
				pCreate := tests.IDs[3]
				pRead := tests.IDs[4]
				usersRow := sqlmock.NewRows([]string{"organization_id", "key"}).
					AddRow(orgId, user1).
					AddRow(orgId, user2)
				orgRow := sqlmock.NewRows([]string{"id", "name", "description"}).
					AddRow(orgId, "test", "organization test")
				userRoleRow := sqlmock.NewRows([]string{"user_key", "role_id"}).
					AddRow(user1, admin).AddRow(user1, view).
					AddRow("user2", view)
				roleRow := sqlmock.
					NewRows([]string{"id", "name", "description"}).
					AddRow(admin, "admin", "administrator").
					AddRow(view, "view", "viewer")
				rolePermissionRow := sqlmock.NewRows([]string{"role_id", "permission_id"}).
					AddRow(admin, pCreate).
					AddRow(view, pRead)
				permissionRow := sqlmock.NewRows([]string{"id", "name", "description"}).
					AddRow(pCreate, "create:permission", "create").
					AddRow(pRead, "read:permission", "read")
				// users find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."organization_id" = $1`),
				).WithArgs(orgId).WillReturnRows(usersRow)

				// organizations find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."id" = $1`),
				).WithArgs(orgId).WillReturnRows(orgRow)

				// user_roles find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "user_roles" WHERE "user_roles"."user_key" IN ($1,$2)`),
				).WithArgs(user1, user2).WillReturnRows(userRoleRow)

				// roles find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "roles" WHERE "roles"."id" IN ($1,$2)`),
				).WithArgs(admin, view).WillReturnRows(roleRow)

				// role_permissions find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE "role_permissions"."role_id" IN ($1,$2)`),
				).WithArgs(admin, view).WillReturnRows(rolePermissionRow)

				// permissions find
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."id" IN ($1,$2)`),
				).WithArgs(pCreate, pRead).WillReturnRows(permissionRow)

				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					orgId, _ := model.NewID(tests.IDs[0])
					uRepo := repo.NewConnection().User(ctx)
					us, err := uRepo.FindAll(orgId)
					u := us[0]
					assert.NoError(t, err)
					assert.Equal(t, user1, u.GetKey())
					assert.NotEmpty(t, u.GetRole())
					assert.NotEmpty(t, u.GetPermission())
					assert.Equal(t, 2, len(u.GetRole()))
					assert.Equal(t, 2, len(u.GetPermission()))
					assert.Equal(t, admin, *u.GetRole()[0].GetID())
					assert.Equal(t, view, *u.GetRole()[1].GetID())
					assert.Equal(t, pCreate, *u.GetPermission()[0].GetID())
					assert.Equal(t, pRead, *u.GetPermission()[1].GetID())
					u = us[1]
					assert.Equal(t, user2, u.GetKey())
					assert.NotEmpty(t, u.GetRole())
					assert.NotEmpty(t, u.GetPermission())
					assert.Equal(t, 1, len(u.GetRole()))
					assert.Equal(t, 1, len(u.GetPermission()))
					assert.Equal(t, view, *u.GetRole()[0].GetID())
					assert.Equal(t, pRead, *u.GetPermission()[0].GetID())
				}
			},
		},
	}
	cases.Run(t)
}
