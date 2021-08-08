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

func TestPermission(t *testing.T) {
	ctx := context.Background()
	cases := tests.MocksByPostgres{
		{
			Name: "create",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "permissions" ("id","created_at","updated_at","name","description") VALUES ($1,$2,$3,$4,$5)`),
				).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "create:permission", "test").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				con := infra.NewRepository(db).NewConnection()
				return func(t *testing.T) {
					tx := con.Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						name, _ := model.NewName("create:permission")
						pRepo := tx.Permission()
						p, err := pRepo.Create(name, "test")
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
				name := "create:permission"
				description := "test description"
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(id, name, description)
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."id" = $1`),
				).WithArgs(id).WillReturnRows(rows)
				pRepo := infra.NewRepository(db).NewConnection().Permission(ctx)
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
				names := []string{"create:permission", "read:permission"}
				descriptions := []string{"permission desc1", "permission desc2"}
				rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(tests.IDs[0], names[0], descriptions[0]).AddRow(tests.IDs[1], names[1], descriptions[1])
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "permissions" ORDER BY id`),
				).WillReturnRows(rows)
				pRepo := infra.NewRepository(db).NewConnection().Permission(ctx)
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
				id, name, desc := tests.IDs[0], "test", "test desc"
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`UPDATE "permissions" SET "updated_at"=$1,"name"=$2,"description"=$3 WHERE "permissions"."id" = $4`),
				).WithArgs(sqlmock.AnyArg(), name, desc, id).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				con := infra.NewRepository(db).NewConnection()
				return func(t *testing.T) {
					tx := con.Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						p, _ := model.NewPermission(id, name, desc)
						return tx.Permission().Update(p)
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
					regexp.QuoteMeta(`DELETE FROM "permissions" WHERE "permissions"."id" = $1`),
				).WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				con := infra.NewRepository(db).NewConnection()
				return func(t *testing.T) {
					tx := con.Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						id, _ := model.NewID(tests.IDs[0])
						return tx.Permission().Delete(id)
					})
					assert.NoError(t, err)
				}
			},
		},
	}
	cases.Run(t)
}
