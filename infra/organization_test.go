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

func TestOrganization(t *testing.T) {
	ctx := context.Background()
	cases := tests.MocksByPostgres{
		{
			Name: "create",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "organizations" ("id","created_at","updated_at","name","description") VALUES ($1,$2,$3,$4,$5)`),
				).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					tx := repo.NewConnection().Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						name, _ := model.NewName("test")
						e, err := tx.Organization().Create(name, "organization test")
						assert.NoError(t, err)
						assert.NotEmpty(t, *e.GetID())
						assert.Equal(t, "test", *e.GetName())
						assert.Equal(t, "organization test", e.GetDescription())
						return err
					})
					assert.NoError(t, err)
				}
			},
		},
		{
			Name: "findById",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				row := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(tests.IDs[0], "test", "organization test")
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."id" = $1 ORDER BY id`),
				).WithArgs(tests.IDs[0]).WillReturnRows(row)
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."organization_id" = $1`),
				).WithArgs(tests.IDs[0]).WillReturnRows(sqlmock.NewRows([]string{""}))
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					id, _ := model.NewID(tests.IDs[0])
					r, err := repo.NewConnection().Organization(ctx).FindByID(id)
					assert.NoError(t, err)
					assert.Equal(t, tests.IDs[0], *r.GetID())
					assert.Equal(t, "test", *r.GetName())
					assert.Equal(t, "organization test", r.GetDescription())
				}
			},
		},
		{
			Name: "findByName",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				row := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(tests.IDs[0], "test", "organization test")
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "organizations" WHERE "organizations"."name" = $1`),
				).WithArgs("organization test").WillReturnRows(row)
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."organization_id" = $1`),
				).WithArgs(tests.IDs[0]).WillReturnRows(sqlmock.NewRows([]string{""}))
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					name, _ := model.NewName("organization test")
					r, err := repo.NewConnection().Organization(ctx).FindByName(name)
					assert.NoError(t, err)
					assert.Equal(t, tests.IDs[0], *r.GetID())
					assert.Equal(t, "test", *r.GetName())
					assert.Equal(t, "organization test", r.GetDescription())
				}
			},
		},
		{
			Name: "findAll",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				expecteds := []struct {
					id, name, description string
				}{
					{
						id:          tests.IDs[0],
						name:        "test",
						description: "organization test",
					},
					{
						id:          tests.IDs[1],
						name:        "test2",
						description: "org test2",
					},
				}
				row := sqlmock.NewRows([]string{"id", "name", "description"})
				for _, expected := range expecteds {
					row.AddRow(expected.id, expected.name, expected.description)
				}
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "organizations" ORDER BY id`),
				).WillReturnRows(row)
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."organization_id" IN ($1,$2)`),
				).WithArgs(tests.IDs[0], tests.IDs[1]).WillReturnRows(sqlmock.NewRows([]string{""}))
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					orgs, err := repo.NewConnection().Organization(ctx).FindAll()
					assert.NoError(t, err)
					for idx, org := range orgs {
						assert.Equal(t, expecteds[idx].id, *org.GetID())
						assert.Equal(t, expecteds[idx].name, *org.GetName())
						assert.Equal(t, expecteds[idx].description, org.GetDescription())
					}
				}
			},
		},
		{
			Name: "update",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`UPDATE "organizations" SET "updated_at"=$1,"name"=$2,"description"=$3 WHERE "organizations"."id" = $4`),
				).WithArgs(sqlmock.AnyArg(), "test", "test org", tests.IDs[0]).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					tx := repo.NewConnection().Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						org, _ := model.NewOrganization(tests.IDs[0], "test", "test org")
						return tx.Organization().Update(org)
					})
					assert.NoError(t, err)
				}
			},
		},
		{
			Name: "delete",
			Fn: func(db dao.DataBase, mock sqlmock.Sqlmock) func(t *testing.T) {
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`DELETE FROM "organizations" WHERE "organizations"."id" = $1`),
				).WithArgs(tests.IDs[0]).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				var repo repository.Repository = infra.NewRepository(db)
				return func(t *testing.T) {
					tx := repo.NewConnection().Transaction(ctx)
					err := tx.Do(func(tx repository.Transaction) error {
						id, _ := model.NewID(tests.IDs[0])
						return tx.Organization().Delete(id)
					})
					assert.NoError(t, err)
				}
			},
		},
	}
	cases.Run(t)
}
