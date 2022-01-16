package rdb

import (
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/mysql"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/postgres"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"github.com/n-creativesystem/rbns/storage"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"github.com/wader/gormstore/v2"
	str2duration "github.com/xhit/go-str2duration/v2"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// ContextSessionKey `context.Context`へセッションを格納するキー
type ContextSessionKey struct{}

type SQLStore struct {
	db   *gorm.DB
	log  logger.Logger
	quit chan struct{}
}

var (
	_ storage.Factory = (*SQLStore)(nil)
)

func (f *SQLStore) Initialize(settings *storage.Setting) error {
	log := &gormLogger{
		skipErrRecordNotFound: true,
		slowThreshold:         200 * time.Millisecond,
		log:                   f.log.WithSubModule("gorm"),
	}
	f.db.Logger = log
	sqlDB, err := f.db.DB()
	if err != nil {
		return err
	}
	var slaveDB *dbresolver.DBResolver
	if dsn := settings.Key("slave_dsn").String(); dsn != "" {
		dsn = os.ExpandEnv(dsn)
		var dialect gorm.Dialector
		switch f.db.Dialector.Name() {
		case "postgres":
			dialect = postgres.Open(dsn)
		case "mysql":
			dialect = mysql.Open(dsn)
		}
		slaveDB = dbresolver.Register(dbresolver.Config{Replicas: []gorm.Dialector{dialect}})
	}
	if maxIdleConns := settings.Key("key max_idle_conns").Int(); maxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(maxIdleConns)
		if slaveDB != nil {
			slaveDB.SetMaxIdleConns(maxIdleConns)
		}
	}
	if maxOpenConns := settings.Key("max_open_conns").Int(); maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(maxOpenConns)
		if slaveDB != nil {
			slaveDB.SetMaxOpenConns(maxOpenConns)
		}
	}
	if maxLifeTime := settings.Key("max_life_time").String(); maxLifeTime != "" {
		value, err := str2duration.ParseDuration(maxLifeTime)
		if err == nil {
			sqlDB.SetConnMaxLifetime(value)
			if slaveDB != nil {
				slaveDB.SetConnMaxLifetime(value)
			}
		}
	}
	if logSettings := settings.Key("log").StringMap(); logSettings != nil {
		if val, ok := logSettings["slow_threshold"].(string); ok {
			v, err := str2duration.ParseDuration(val)
			if err == nil {
				log.slowThreshold = v
			}
		}
	}
	if slaveDB != nil {
		if err := f.db.Use(slaveDB); err != nil {
			return err
		}
	}
	if settings.Key("migration").Bool() {
		if err := f.migration(); err != nil {
			return err
		}
	}
	return nil
}

func (f *SQLStore) SessionStore(keyPairs ...storage.KeyPairs) sessions.Store {
	bkp := [][]byte{}
	for _, kp := range keyPairs {
		bkp = append(bkp, kp)
	}
	f.quit = make(chan struct{})
	store := gormstore.NewOptions(f.db, gormstore.Options{SkipCreateTable: false}, bkp...)
	go store.PeriodicCleanup(1*time.Hour, f.quit)
	return store
}

func (f *SQLStore) Close() error {
	defer close(f.quit)
	db, _ := f.db.DB()
	return db.Close()
}

func NewFactory(db *gorm.DB, bus bus.Bus) *SQLStore {
	_ = db.Use(otelgorm.NewPlugin())
	factory := &SQLStore{
		db:  db.Session(&gorm.Session{}),
		log: logger.New("sql store"),
	}
	factory.bus()
	return factory
}
