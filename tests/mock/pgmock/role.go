package pgmock

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
)

func CreateRole(mock sqlmock.Sqlmock) (newId int) {
	newId = 1
	rows := sqlmock.NewRows([]string{"id"}).AddRow(newId)
	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "roles" ("created_at","updated_at","name","description") VALUES ($1,$2,$3,$4) RETURNING "id"`),
	).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "admin", "administrator").WillReturnRows(rows)
	mock.ExpectCommit()
	return newId
}

func FindByIdRole(mock sqlmock.Sqlmock) (id int, name string, description string) {
	id = 1
	name = "admin"
	description = "administrator"
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(id, name, description)
	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "roles" WHERE "roles"."id" = $1`),
	).WithArgs(1).WillReturnRows(rows)
	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE "role_permissions"."role_id" = $1`),
	).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"role_id", "permission_id"}))
	return
}

func FindAllRole(mock sqlmock.Sqlmock) {
	ids := []int{1, 2}
	names := []string{"admin", "guest"}
	descriptions := []string{"administrator", "guest"}
	rows := sqlmock.
		NewRows([]string{"id", "name", "description"}).
		AddRow(ids[0], names[0], descriptions[0]).
		AddRow(ids[1], names[1], descriptions[1])
	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "roles" ORDER BY id`),
	).WillReturnRows(rows)
	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "role_permissions" WHERE "role_permissions"."role_id" IN ($1,$2)`),
	).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"role_id", "permission_id"}))
}

func UpdateRole(mock sqlmock.Sqlmock) {
	id := 2
	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "roles" SET "updated_at"=$1,"name"=$2,"description"=$3 WHERE "roles"."id" = $4`),
	).WithArgs(sqlmock.AnyArg(), "view", "view", id).WillReturnResult(sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
}

func DeleteRole(mock sqlmock.Sqlmock) {
	id := 1
	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(`DELETE FROM "roles" WHERE "roles"."id" = $1`),
	).WithArgs(id).WillReturnResult(sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
}

func AddPermission(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO "role_permissions" ("role_id","permission_id") VALUES ($1,$2) ON CONFLICT DO NOTHING`),
	).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}
