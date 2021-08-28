package mysql

import (
	"github.com/n-creativesystem/rbns/domain/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}

func New(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), opts...)
}

func NewDBErr(err error) error {
	return model.NewDBErr(err)
}
