package driver

import (
	"errors"

	"github.com/n-creativesystem/rbns/infra/rdb/driver/customs"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/mysql"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/postgres"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/sqlite3"
	"gorm.io/gorm"
)

const (
	PostgreSQL = "postgres"
	MySQL      = "mysql"
	SQLite3    = "sqlite3"
)

func Open(driverName, dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	opts = append([]gorm.Option{&gorm.Config{
		NamingStrategy: customs.NamingStrategy{},
	}}, opts...)
	switch driverName {
	case PostgreSQL:
		return postgres.New(dsn, opts...)
	case MySQL:
		return mysql.New(dsn, opts...)
	case SQLite3:
		return sqlite3.New(dsn, opts...)
	}
	return nil, errors.New("No support driver")
}
