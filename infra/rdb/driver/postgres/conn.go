package postgres

import (
	"errors"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/n-creativesystem/rbns/domain/model"
	"google.golang.org/grpc/codes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetMigrator(db *gorm.DB) postgres.Migrator {
	return db.Migrator().(postgres.Migrator)
}

func Open(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

func New(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(Open(dsn), opts...)
}

func NewDBErr(err error) error {
	if err == nil {
		return nil
	}
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		switch pgError.SQLState() {
		case "23505":
			return model.ErrAlreadyExists
		}
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ErrNoData
	}
	return model.NewErrorStatus(uint32(codes.Internal), err.Error())
}

func HasUnique(db *gorm.DB, currentSchema, table interface{}, constraintName string) bool {
	var count int64
	db.Raw(
		"SELECT count(*) FROM INFORMATION_SCHEMA.table_constraints WHERE table_schema = ? AND table_name = ? AND constraint_name = ?",
		currentSchema, table, constraintName,
	).Row().Scan(&count)
	return count > 0
}

func AddUnique(db *gorm.DB, stmt *gorm.Statement, constraintName string, names ...string) error {
	migrator := GetMigrator(db)
	currentSchema, t := migrator.CurrentSchema(stmt, stmt.Schema.Table)
	if HasUnique(db, currentSchema, t, constraintName) {
		return nil
	}
	query := "alter table ? add CONSTRAINT ? unique("
	table := stmt.Schema.Table
	var vars = []interface{}{clause.Table{Name: table}, clause.Column{Name: constraintName}}
	if stmt.TableExpr != nil {
		vars[0] = stmt.TableExpr
	}
	for _, name := range names {
		f := stmt.Schema.LookUpField(name)
		query += "?,"
		vars = append(vars, clause.Column{Name: f.DBName})
	}
	query = strings.TrimSuffix(query, ",")
	query += ")"
	return db.Exec(query, vars...).Error
}
