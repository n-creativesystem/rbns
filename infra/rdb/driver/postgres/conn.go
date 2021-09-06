package postgres

import (
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/n-creativesystem/rbns/domain/model"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

func New(dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	return gorm.Open(Open(dsn), opts...)
}

func IsDuplication(err error) bool {
	if pgError, ok := err.(*pgconn.PgError); ok {
		switch pgError.SQLState() {
		case "23505":
			return true
		default:
			return false
		}
	}
	return false
}

func NewDBErr(err error) error {
	if err == nil {
		return nil
	}
	if pgError, ok := err.(*pgconn.PgError); ok {
		switch pgError.SQLState() {
		case "23505":
			return status.Error(http.StatusConflict, model.ErrAlreadyExists.Error())
		}
	}
	if err == gorm.ErrRecordNotFound {
		return status.Error(http.StatusNotFound, model.ErrNoData.Error())
	}
	return model.NewDBErr(err)
}
