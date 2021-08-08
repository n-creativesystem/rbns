package dao

import (
	"errors"

	"github.com/n-creativesystem/rbns/infra/dao/driver/mysql"
	"github.com/n-creativesystem/rbns/infra/dao/driver/postgres"
	"gorm.io/gorm"
)

func newDriver(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	if d := newDialector(dsn); d != nil {
		return gorm.Open(d, opts...)
	} else {
		return nil, errors.New("can not find this dialector")
	}
}

func newDialector(dsn string) gorm.Dialector {
	switch dialector {
	case postgreSQL:
		return postgres.Open(dsn)
	case mySQL:
		return mysql.Open(dsn)
	default:
		return nil
	}
}
