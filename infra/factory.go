package infra

import (
	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/config"
	"github.com/n-creativesystem/rbns/infra/rdb"
	"github.com/n-creativesystem/rbns/infra/rdb/driver"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/customs"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/mysql"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/postgres"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/sqlite3"
	"github.com/n-creativesystem/rbns/storage"
	"gorm.io/gorm"
)

func NewFactory(conf *config.Config, bus bus.Bus) (*storage.FactorySet, error) {
	gormConfig := &gorm.Config{
		NamingStrategy: customs.NamingStrategy{},
	}
	v := conf.DatabaseRaw
	var factory storage.Factory
	dsn := v.Key("dsn").StringExpand()
	switch conf.StorageType {
	case driver.PostgreSQL:
		db, _ := postgres.New(dsn, gormConfig)
		factory = rdb.NewFactory(db, bus)
	case driver.MySQL:
		db, _ := mysql.New(dsn, gormConfig)
		factory = rdb.NewFactory(db, bus)
	case driver.SQLite3:
		db, _ := sqlite3.New(dsn, gormConfig)
		factory = rdb.NewFactory(db, bus)
	}
	return &storage.FactorySet{
		Factory: factory,
		Settings: &storage.Setting{
			Section: v,
		},
	}, nil
}
