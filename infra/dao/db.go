package dao

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/n-creativesystem/rbns/di"
	"github.com/n-creativesystem/rbns/infra/entity"
	"github.com/n-creativesystem/rbns/logger"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func Register() {
	di.MustRegister(New)
}

type driverType string

func (d driverType) String() string {
	return string(d)
}

const (
	postgreSQL = driverType("postgres")
	mySQL      = driverType("mysql")
)

var dialector driverType

type DataBase interface {
	Session(ctx context.Context) *gorm.DB
}

type database struct {
	db      *gorm.DB
	session *gorm.Session
}

var db DataBase
var once sync.Once

func New(opts ...Option) DataBase {
	conf := &Config{}
	*conf = *defaultConfig
	for _, opt := range opts {
		opt(conf)
	}
	once.Do(func() {
		gormOptions := []logger.GormOption{}
		if conf.debug {
			gormOptions = append(gormOptions, logger.WithDebug)
		}
		switch strings.ToLower(conf.dialector) {
		case "postgres":
			dialector = postgreSQL
		case "mysql":
			dialector = mySQL
		}
		fn := func(dsn string) *gorm.DB {
			gConf := &gorm.Config{
				Logger: logger.NewGorm(gormOptions...),
			}
			ormDB, err := newDriver(dsn, gConf)
			if err != nil {
				panic(err)
			}
			sqlDB, err := ormDB.DB()
			if err != nil {
				panic(err)
			}

			// Connection Pool
			sqlDB.SetMaxIdleConns(conf.maxIdleConns)
			sqlDB.SetMaxOpenConns(conf.maxOpenConns)
			sqlDB.SetConnMaxLifetime(conf.maxLifeTime)
			return ormDB
		}
		ormDB := fn(conf.masterDSN)
		if conf.slaveDSN != "" {
			ormDB.Use(dbresolver.Register(dbresolver.Config{Replicas: []gorm.Dialector{newDialector(conf.slaveDSN)}}).
				SetMaxIdleConns(conf.maxIdleConns).
				SetMaxOpenConns(conf.maxOpenConns).
				SetConnMaxLifetime(conf.maxLifeTime).
				SetConnMaxIdleTime(time.Hour))
		}
		db_ := &database{
			db:      fn(conf.masterDSN),
			session: &gorm.Session{},
		}
		if conf.session != nil {
			*db_.session = *conf.session
		}
		db = db_
	})
	if conf.migrationBack {
		if err := entity.MigrationBack(db.Session(context.Background())); err != nil {
			panic(err)
		}
	}
	if conf.migration {
		if err := entity.Migrations(db.Session(context.Background())); err != nil {
			panic(err)
		}
	}

	return db
}

func (db *database) Session(ctx context.Context) *gorm.DB {
	session := &gorm.Session{}
	*session = *db.session
	session.Context = ctx
	session.CreateBatchSize = 1000
	return db.db.Session(session)
}
