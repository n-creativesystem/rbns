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

func TestRole(t *testing.T) {
	ctx := context.Background()
	cases := tests.MocksByPostgres{
		{
			Name: "create",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "roles" ("id","created_at","updated_at","name","description") VALUES ($1,$2,$3,$4,$5)`),
				).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "admin", "administrator").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				con := infra.NewRepository(db).NewConnection()
				return func(t *testing.T) {
					tx := con.Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						name, _ := model.NewName("admin")
						pRepo := tx.Role()
						p, err := pRepo.Create(name, "administrator")
						assert.NoError(t, err)
						assert.NotEmpty(t, *p.GetID())
						return err
					})
					assert.NoError(t, err)
				}
			},
		},
		{
			Name: "findById",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				id := tests.IDs[0]
				name := "admin"
				description := "administrator"
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(id, name, description)
				// rolePermissions := sqlmock.NewRows([]string{"role_id", "permission_id"}).AddRow(id, tests.IDs[1])
				// permissionRow := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(tests.IDs[2], "create:permission", "test description")
				// userRoles := sqlmock.NewRows([]string{"role_id", "user_key", "organization_id"}).AddRow(id, "user1", tests.IDs[3])
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "roles" WHERE "roles"."id" = $1`),
				).WithArgs(id).WillReturnRows(rows)
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE "role_permissions"."role_id" = $1`),
				).WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"role_id", "permission_id"}))
				// mock.ExpectQuery(
				// 	regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."id" = $1`),
				// ).WithArgs(sqlmock.AnyArg()).WillReturnRows(permissionRow)
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "user_roles" WHERE "user_roles"."role_id" = $1`),
				).WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"role_id", "user_key", "organization_id"}))
				// mock.ExpectQuery(
				// 	regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."key" = $1`),
				// ).WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"user_key"}))
				pRepo := infra.NewRepository(db).NewConnection().Role(ctx)
				return func(t *testing.T) {
					mId, _ := model.NewID(id)
					p, err := pRepo.FindByID(mId)
					assert.NoError(t, err)
					assert.Equal(t, id, *p.GetID())
					assert.Equal(t, name, *p.GetName())
					assert.Equal(t, description, p.GetDescription())
				}
			},
		},
		{
			Name: "find all",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				names := []string{"admin", "guest"}
				descriptions := []string{"administrator", "guest"}
				rows := sqlmock.
					NewRows([]string{"id", "name", "description"}).
					AddRow(tests.IDs[0], names[0], descriptions[0]).
					AddRow(tests.IDs[1], names[1], descriptions[1])
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "roles" ORDER BY id`),
				).WillReturnRows(rows)
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE "role_permissions"."role_id" IN ($1,$2)`),
				).WithArgs(tests.IDs[0], tests.IDs[1]).WillReturnRows(sqlmock.NewRows([]string{"role_id", "permission_id"}))
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "user_roles" WHERE "user_roles"."role_id" IN ($1,$2)`),
				).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"role_id", "user_key", "organization_id"}))
				pRepo := infra.NewRepository(db).NewConnection().Role(ctx)
				return func(t *testing.T) {
					res, err := pRepo.FindAll()
					assert.NoError(t, err)
					for idx, r := range res {
						assert.Equal(t, tests.IDs[idx], *r.GetID())
						assert.Equal(t, names[idx], *r.GetName())
						assert.Equal(t, descriptions[idx], r.GetDescription())
					}
				}
			},
		},
		{
			Name: "update",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				id, name, desc := tests.IDs[0], "view", "view"
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`UPDATE "roles" SET "updated_at"=$1,"name"=$2,"description"=$3 WHERE "roles"."id" = $4`),
				).WithArgs(sqlmock.AnyArg(), "view", "view", id).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				con := infra.NewRepository(db).NewConnection()
				return func(t *testing.T) {
					tx := con.Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						p, _ := model.NewRole(id, name, desc, nil)
						return tx.Role().Update(p)
					})
					assert.NoError(t, err)
				}
			},
		},
		{
			Name: "delete",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				id := tests.IDs[0]
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`DELETE FROM "roles" WHERE "roles"."id" = $1`),
				).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				con := infra.NewRepository(db).NewConnection()
				return func(t *testing.T) {
					tx := con.Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						id, _ := model.NewID(id)
						return tx.Role().Delete(id)
					})
					assert.NoError(t, err)
				}
			},
		},
		{
			Name: "add permission",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "role_permissions" ("role_id","permission_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`),
				).WithArgs(tests.IDs[0], tests.IDs[1]).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				con := infra.NewRepository(db).NewConnection()
				return func(t *testing.T) {
					tx := con.Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						id, _ := model.NewID(tests.IDs[0])
						p, _ := model.NewPermission(tests.IDs[1], "aaa", "")
						return tx.Role().AddPermission(id, model.Permissions{
							*p,
						})
					})
					assert.NoError(t, err)
				}
			},
		},
		{
			Name: "delete role permission",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`DELETE FROM "role_permissions" WHERE "role_permissions"."role_id" = $1 AND "role_permissions"."permission_id" = $2`),
				).WithArgs(tests.IDs[0], tests.IDs[1]).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					tx := repo.NewConnection().Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						roleId, _ := model.NewID(tests.IDs[0])
						permissionId, _ := model.NewID(tests.IDs[1])
						return tx.Role().DeletePermission(roleId, permissionId)
					})
					assert.NoError(t, err)
				}
			},
		},
	}
	cases.Run(t)

}
