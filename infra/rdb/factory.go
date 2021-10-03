package rdb

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/n-creativesystem/rbns/domain/repository"
	"github.com/n-creativesystem/rbns/infra/rdb/driver/postgres"
	"github.com/n-creativesystem/rbns/infra/rdb/entity"
	"github.com/n-creativesystem/rbns/storage"
	"github.com/sirupsen/logrus"
	str2duration "github.com/xhit/go-str2duration/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"gorm.io/plugin/dbresolver"
)

type gormLogger struct {
	skipErrRecordNotFound bool
	slowThreshold         time.Duration
	log                   *logrus.Logger
}

var _ logger.Interface = (*gormLogger)(nil)

func (l *gormLogger) LogMode(mode logger.LogLevel) logger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Infof(s, args...)
}

func (l *gormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Warnf(s, args...)
}

func (l *gormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.log.WithContext(ctx).Errorf(s, args...)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()
	fields := logrus.Fields{
		"rowsAffected": rowsAffected,
		"sourceField":  utils.FileWithLineNum(),
		"elapsed":      elapsed,
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.skipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.log.WithContext(ctx).WithFields(fields).Errorf("%s", sql)
		return
	}

	if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		l.log.WithContext(ctx).WithFields(fields).Warnf("%s", sql)
		return
	}
	l.log.WithContext(ctx).WithFields(fields).Debugf("%s", sql)
}

type rdbFactory struct {
	db *gorm.DB
}

var _ storage.Factory = (*rdbFactory)(nil)

func (f *rdbFactory) Initialize(mp map[string]interface{}, logger *logrus.Logger) error {
	log := &gormLogger{
		skipErrRecordNotFound: true,
		slowThreshold:         200 * time.Millisecond,
		log:                   logger,
	}
	f.db.Logger = log
	sqlDB, err := f.db.DB()
	if err != nil {
		return err
	}
	var slaveDB *dbresolver.DBResolver
	if dsn, ok := mp["slave_dsn"].(string); ok {
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
	if maxIdleConns, ok := mp["max_idle_conns"].(int); ok {
		sqlDB.SetMaxIdleConns(maxIdleConns)
		if slaveDB != nil {
			slaveDB.SetMaxIdleConns(maxIdleConns)
		}
	}
	if maxOpenConns, ok := mp["max_open_conns"].(int); ok {
		sqlDB.SetMaxOpenConns(maxOpenConns)
		if slaveDB != nil {
			slaveDB.SetMaxOpenConns(maxOpenConns)
		}
	}
	if maxLifeTime, ok := mp["max_life_time"].(string); ok {
		value, err := str2duration.ParseDuration(maxLifeTime)
		if err == nil {
			sqlDB.SetConnMaxLifetime(value)
			if slaveDB != nil {
				slaveDB.SetConnMaxLifetime(value)
			}
		}
	}
	if logSettings, ok := mp["log"].(map[string]interface{}); ok {
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
	migration, _ := mp["migration"].(bool)
	if migration {
		if err := f.db.AutoMigrate(entity.Permission{}, entity.Role{}, entity.RolePermission{}, entity.Organization{}, entity.User{}, entity.UserRole{}, entity.Resource{}, entity.ResourcePermissions{}); err != nil {
			return err
		}
	}
	return nil
}

func (f *rdbFactory) Reader() repository.Reader {
	return &reader{
		driver: f.db,
	}
}

func (f *rdbFactory) Writer() repository.Writer {
	return &tx{
		driver: f.db,
	}
}

func (f *rdbFactory) Close() error {
	db, _ := f.db.DB()
	return db.Close()
}

func NewFactory(db *gorm.DB) storage.Factory {
	return &rdbFactory{
		db: db,
	}
}

type reader struct {
	driver *gorm.DB
}

var _ repository.Reader = (*reader)(nil)

func (r *reader) Permission(ctx context.Context) repository.Permission {
	return &permission{
		db: r.driver.Session(&gorm.Session{Context: ctx}),
	}
}

func (r *reader) Role(ctx context.Context) repository.Role {
	return &role{
		db: r.driver.Session(&gorm.Session{Context: ctx}),
	}
}

func (r *reader) Organization(ctx context.Context) repository.Organization {
	return &organization{
		db: r.driver.Session(&gorm.Session{Context: ctx}),
	}
}

func (r *reader) User(ctx context.Context) repository.User {
	return &user{
		db: r.driver.Session(&gorm.Session{Context: ctx}),
	}
}

func (r *reader) Resource(ctx context.Context) repository.Resource {
	return &resource{
		db: r.driver.Session(&gorm.Session{Context: ctx}),
	}
}

type tx struct {
	driver *gorm.DB
}

var _ repository.Writer = (*tx)(nil)

func (t *tx) Do(ctx context.Context, fn func(tx repository.Transaction) error) error {
	var err error
	defer func() {
		if err != nil {
			logrus.Println(err)
		}
	}()
	tx := t.driver.Session(&gorm.Session{
		Context: ctx,
	}).Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()
	tx.SkipDefaultTransaction = true
	err = fn(&writer{
		db: tx,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

type writer struct {
	db *gorm.DB
}

var _ repository.Transaction = (*writer)(nil)

func (tx *writer) Permission() repository.PermissionCommand {
	return &permission{
		db: tx.db,
	}
}

func (tx *writer) Role() repository.RoleCommand {
	return &role{
		db: tx.db,
	}
}

func (tx *writer) Organization() repository.OrganizationCommand {
	return &organization{
		db: tx.db,
	}
}

func (tx *writer) User() repository.UserCommand {
	return &user{
		db: tx.db,
	}
}

func (tx *writer) Resource() repository.ResourceCommand {
	return &resource{
		db: tx.db,
	}
}
