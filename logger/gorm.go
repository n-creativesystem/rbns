package logger

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var defaultGormConfig = &config{
	debug:      true,
	timeFormat: TimestampFormat,
}

type config struct {
	debug      bool
	timeFormat string
}

type GormOption func(conf *config)

func WithDebug(conf *config) {
	conf.debug = true
}

func WithTimeFormat(timeFmt string) GormOption {
	return func(conf *config) {
		conf.timeFormat = timeFmt
	}
}

type gormLogger struct {
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
	log                   *logrus.Logger
}

var _ gLogger.Interface = (*gormLogger)(nil)

func NewGorm(opts ...GormOption) *gormLogger {
	conf := &config{}
	*conf = *defaultGormConfig
	for _, opt := range opts {
		opt(conf)
	}
	log := logrus.New()
	if conf.debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: conf.timeFormat,
	})
	return &gormLogger{
		SkipErrRecordNotFound: true,
		SlowThreshold:         200 * time.Millisecond,
		log:                   log,
	}
}

func (l *gormLogger) LogMode(mode gLogger.LogLevel) gLogger.Interface {
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
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		l.log.WithContext(ctx).WithFields(fields).Errorf("%s", sql)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.log.WithContext(ctx).WithFields(fields).Warnf("%s", sql)
		return
	}
	l.log.WithContext(ctx).WithFields(fields).Debugf("%s", sql)
}
