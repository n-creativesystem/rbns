package rdb

import (
	"github.com/n-creativesystem/rbns/infra/rdb/driver/postgres"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func IsDuplication(db *gorm.DB, err error) error {
	switch db.Dialector.(type) {
	case *pgdriver.Dialector:
		if postgres.IsDuplication(err) {
			postgres.NewDBErr(err)
		}
		return err
	default:
		return err
	}
}
