package service

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/n-creativesystem/rbns/infra/rdb"
	"github.com/n-creativesystem/rbns/tests"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPermission(t *testing.T) {
	ctx := context.Background()
	cases := tests.MocksByPostgres{
		{
			Name: "create",
			Fn: func(db *gorm.DB, mock sqlmock.Sqlmock) func(t *testing.T) {
				expecteds := []struct {
					name, description string
				}{
					{
						name:        "create:user",
						description: "create user permission",
					},
					{
						name:        "update:user",
						description: "update user permisssion",
					},
					{
						name:        "read:user",
						description: "read user permission",
					},
					{
						name:        "delete:user",
						description: "delete user permission",
					},
				}
				repo := rdb.NewFactory(db)
				pSrv := NewPermissionService(repo.Reader(), repo.Writer())
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`INSERT INTO "permissions" ("id","created_at","updated_at","name","description") VALUES ($1,$2,$3,$4,$5),($6,$7,$8,$9,$10),($11,$12,$13,$14,$15),($16,$17,$18,$19,$20)`),
				).WillReturnResult(sqlmock.NewResult(0, 4))
				mock.ExpectCommit()
				return func(t *testing.T) {
					names := []string{"create:user", "update:user", "read:user", "delete:user"}
					descriptions := []string{"create user permission", "update user permisssion", "read user permission", "delete user permission"}
					out, err := pSrv.Create(ctx, names, descriptions)
					assert.NoError(t, err)
					for idx, entity := range out.Copy() {
						expected := expecteds[idx]
						assert.NotEmpty(t, *entity.GetID())
						assert.Equal(t, expected.name, entity.GetName())
						assert.Equal(t, expected.description, entity.GetDescription())
					}
				}
			},
		},
		{
			Name: "findById",
			Fn: func(db *gorm.DB, mock sqlmock.Sqlmock) func(t *testing.T) {
				repo := rdb.NewFactory(db)
				pSrv := NewPermissionService(repo.Reader(), repo.Writer())
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."id" = $1`),
				).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).AddRow("1", "create:user", "create user permission"))
				return func(t *testing.T) {
					out, err := pSrv.FindById(ctx, "1")
					assert.NoError(t, err)
					assert.Equal(t, "create:user", out.GetName())
					assert.Equal(t, "create user permission", out.GetDescription())
				}
			},
		},
		{
			Name: "findAll",
			Fn: func(db *gorm.DB, mock sqlmock.Sqlmock) func(t *testing.T) {
				expecteds := []struct {
					id                string
					name, description string
				}{
					{
						id:          "1",
						name:        "create:user",
						description: "create user permission",
					},
					{
						id:          "2",
						name:        "update:user",
						description: "update user permisssion",
					},
					{
						id:          "3",
						name:        "read:user",
						description: "read user permission",
					},
					{
						id:          "4",
						name:        "delete:user",
						description: "delete user permission",
					},
				}
				row := sqlmock.NewRows([]string{"id", "name", "description"})
				for _, e := range expecteds {
					row.AddRow(e.id, e.name, e.description)
				}
				repo := rdb.NewFactory(db)
				pSrv := NewPermissionService(repo.Reader(), repo.Writer())
				mock.ExpectQuery(
					regexp.QuoteMeta(`SELECT * FROM "permissions" ORDER BY id`),
				).WillReturnRows(row)
				return func(t *testing.T) {
					out, err := pSrv.FindAll(ctx)
					assert.NoError(t, err)
					for idx, entity := range out {
						expected := expecteds[idx]
						assert.Equal(t, expected.id, entity.GetID())
						assert.Equal(t, expected.name, entity.GetName())
						assert.Equal(t, expected.description, entity.GetDescription())
					}
				}
			},
		},
		{
			Name: "update",
			Fn: func(db *gorm.DB, mock sqlmock.Sqlmock) func(t *testing.T) {
				repo := rdb.NewFactory(db)
				pSrv := NewPermissionService(repo.Reader(), repo.Writer())
				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta(`UPDATE "permissions" SET "updated_at"=$1,"name"=$2,"description"=$3 WHERE "permissions"."id" = $4`),
				).WithArgs(sqlmock.AnyArg(), "read:permission", "read user permission", "1").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				return func(t *testing.T) {
					err := pSrv.Update(ctx, "1", "read:permission", "read user permission")
					assert.NoError(t, err)
				}
			},
		},
	}

	cases.Run(t)
}
