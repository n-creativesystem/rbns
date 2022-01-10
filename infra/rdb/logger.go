package rdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/n-creativesystem/rbns/logger"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type gormLogger struct {
	skipErrRecordNotFound bool
	slowThreshold         time.Duration
	log                   logger.Logger
}

var _ glog.Interface = (*gormLogger)(nil)

func (l *gormLogger) LogMode(mode glog.LogLevel) glog.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.log.InfoWithContext(ctx, fmt.Sprintf(s, args...))
}

func (l *gormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.log.WarningWithContext(ctx, fmt.Sprintf(s, args...))
}

func (l *gormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	err := fmt.Errorf(s, args...)
	l.log.ErrorWithContext(ctx, err, err.Error())
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()
	fields := []interface{}{
		"rowsAffected", rowsAffected,
		"sourceField", utils.FileWithLineNum(),
		"elapsed", elapsed.String(),
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.skipErrRecordNotFound) {
		l.log.ErrorWithContext(ctx, err, sql, fields...)
		return
	}

	if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		l.log.WarningWithContext(ctx, sql, fields...)
		return
	}
	l.log.DebugWithContext(ctx, sql, fields...)
}
