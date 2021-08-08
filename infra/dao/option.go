package dao

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var defaultConfig = &Config{
	debug:         false,
	migration:     false,
	migrationBack: false,
	maxIdleConns:  10,
	maxOpenConns:  100,
	maxLifeTime:   1,
	session:       &gorm.Session{},
}

type Config struct {
	session       *gorm.Session
	debug         bool
	migration     bool
	migrationBack bool
	maxIdleConns  int
	maxOpenConns  int
	maxLifeTime   time.Duration
	dialector     string
	masterDSN     string
	slaveDSN      string
}

type Option func(conf *Config)

func WithMigration(conf *Config) {
	conf.migration = true
}

func WithMigrationBack(conf *Config) {
	conf.migrationBack = true
}

func Debug(conf *Config) {
	conf.debug = true
}

func DryRun(conf *Config) {
	conf.session.DryRun = true
}

func Session(session *gorm.Session) Option {
	return func(conf *Config) {
		if session != nil {
			*conf.session = *session
		}
	}
}

func WithMaxIdleConn(count int) Option {
	if count <= 0 {
		count = 10
	}
	return func(conf *Config) {
		conf.maxIdleConns = count
	}
}

func WithMaxOpenConns(count int) Option {
	if count <= 0 {
		count = 100
	}
	return func(conf *Config) {
		conf.maxOpenConns = count
	}
}

func WithMaxLifeTime(hour int) Option {
	if hour <= 0 {
		hour = 1
	}
	return func(conf *Config) {
		conf.maxLifeTime = time.Duration(hour) * time.Hour
	}
}

func WithDialector(name string) Option {
	if name == "" {
		name = postgreSQL.String()
	}
	return func(conf *Config) {
		conf.dialector = name
	}
}

func WithMasterDSN(dsn string) Option {
	if dsn == "" {
		logrus.Fatalln("master dsn is empty")
	}
	return func(conf *Config) {
		conf.masterDSN = os.ExpandEnv(dsn)
	}
}

func WithSlaveDSN(dsn string) Option {
	return func(conf *Config) {
		if dsn == "" {
			dsn = conf.masterDSN
		}
		conf.slaveDSN = os.ExpandEnv(dsn)
	}
}
