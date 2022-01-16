package driver

import (
	"errors"

	"gorm.io/gorm"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/mysql"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/postgres"
)

func NewDBErr(db *gorm.DB, err error) error {
	if err == nil {
		return nil
	}
	var errStatus model.ErrorStatus
	if errors.As(err, &errStatus) {
		return err
	}
	switch db.Dialector.Name() {
	case PostgreSQL:
		err = postgres.NewDBErr(err)
	case MySQL:
		err = mysql.NewDBErr(err)
	}
	return err
}
