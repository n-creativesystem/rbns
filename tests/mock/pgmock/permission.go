package pgmock

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
)

func CreatePermission(mock sqlmock.Sqlmock) {
	newId := 1
	rows := sqlmock.NewRows([]string{"id"}).AddRow(newId)
	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "permissions" ("created_at","updated_at","name","description") VALUES ($1,$2,$3,$4) RETURNING "id"`),
	).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "create:permission", "test").WillReturnRows(rows)
	mock.ExpectCommit()
}

func FindByIdPermission(mock sqlmock.Sqlmock) {
	id := 1
	name := "create:permission"
	description := "test description"
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(id, name, description)
	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "permissions" WHERE "permissions"."id" = $1`),
	).WithArgs(1).WillReturnRows(rows)
}

func FindAllPermission(mock sqlmock.Sqlmock) {
	ids := []int{1, 2}
	names := []string{"create:permission", "read:permission"}
	descriptions := []string{"permission desc1", "permission desc2"}
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(ids[0], names[0], descriptions[0]).AddRow(ids[1], names[1], descriptions[1])
	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "permissions" ORDER BY id`),
	).WillReturnRows(rows)
}

func UpdatePermission(mock sqlmock.Sqlmock) {
	id := 1
	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(`UPDATE "permissions" SET "updated_at"=$1,"name"=$2,"description"=$3 WHERE "permissions"."id" = $4`),
	).WithArgs(sqlmock.AnyArg(), "test", "test desc", id).WillReturnResult(sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
}

func DeletePermission(mock sqlmock.Sqlmock) {
	id := 1
	mock.ExpectBegin()
	mock.ExpectExec(
		regexp.QuoteMeta(`DELETE FROM "permissions" WHERE "permissions"."id" = $1`),
	).WithArgs(id).WillReturnResult(sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
}
