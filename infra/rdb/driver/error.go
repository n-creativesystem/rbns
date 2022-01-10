package driver

import (
	"gorm.io/gorm"

	"github.com/n-creativesystem/rbns/infra/rdb/driver/mysql"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/postgres"
)

func NewDBErr(db *gorm.DB, err error) error {
	if err == nil {
		return nil
	}

	switch db.Dialector.Name() {
	case PostgreSQL:
		err = postgres.NewDBErr(err)
	case MySQL:
		err = mysql.NewDBErr(err)
	}
	return err
}
