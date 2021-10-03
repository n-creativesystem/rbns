package sqlite3

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}

func New(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(Open(dsn), opts...)
}
