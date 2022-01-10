package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/bus"
	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/internal/contexts"
	"gorm.io/gorm"
)

const (
	dammyTenant = "dammy"
)

type dbSessionFunc func(sess *DBSession) error

type dbSessionFuncWithTenant func(sess *DBSession, tenant string) error

type dbTransactionFunc func(sess *DBSession) error

type dbTransactionFuncWithTenant func(sess *DBSession, tenant string) error

func (f *SQLStore) inTransactionWithToken(ctx context.Context, callback dbTransactionFuncWithTenant) error {
	return inTransactionWithTenantAndContext(ctx, f.db, callback)
}

func (f *SQLStore) inTransactionWithDbSession(ctx context.Context, callback dbTransactionFunc) error {
	return inTransactionWithDbSession(ctx, f.db, callback)
}

func (f *SQLStore) DbSessionFunc(ctx context.Context, callback dbSessionFunc) error {
	return dbSssion(ctx, f.db, callback)
}

func (f *SQLStore) DbSessionWithTenant(ctx context.Context, callback dbSessionFuncWithTenant) error {
	return dbSssionWithToken(ctx, f.db, callback)
}

func (f *SQLStore) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return f.inTransaction(ctx, fn)
}

func (f *SQLStore) inTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return inTransactionWithTenantAndContext(ctx, f.db, func(sess *DBSession, tenant string) error {
		withValue := context.WithValue(ctx, ContextSessionKey{}, sess)
		return fn(withValue)
	})
}

func dbSssion(ctx context.Context, db *gorm.DB, callback dbSessionFunc) error {
	sess := startSession(ctx, db, false)
	err := callback(sess)
	if err != nil {
		return err
	}
	if len(sess.events) > 0 {
		for _, e := range sess.events {
			if err = bus.PublishCtx(ctx, e); err != nil {
				GetLogger(db).ErrorWithContext(ctx, err, "Failed to publish event after commit.")
			}
		}
	}
	return nil
}

func dbSssionWithToken(ctx context.Context, db *gorm.DB, callback dbSessionFuncWithTenant) error {
	tenant := contexts.FromTenantContext(ctx)
	if tenant == "" {
		return model.ErrTenantRequired
	} else if tenant == dammyTenant {
		tenant = ""
	}
	sess := startSession(ctx, db, false)
	err := callback(sess, tenant)
	if err != nil {
		return err
	}
	if len(sess.events) > 0 {
		for _, e := range sess.events {
			if err = bus.PublishCtx(ctx, e); err != nil {
				GetLogger(db).ErrorWithContext(ctx, err, "Failed to publish event after commit.")
			}
		}
	}
	return nil
}

func inTransactionWithDbSession(ctx context.Context, db *gorm.DB, callback dbTransactionFunc) error {
	var err error
	defer func() {
		if err != nil {
			GetLogger(db).ErrorWithContext(ctx, err, "inTransactionWithDbSession")
		}
	}()
	sess := startSession(ctx, db, true)
	err = callback(sess)
	if err != nil {
		sess.Rollback()
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			sess.Rollback()
			GetLogger(db).PanicWithContext(ctx, err.(error), "inTransactionWithDbSession")
		}
	}()
	if err = sess.Commit().Error; err != nil {
		return err
	}
	if len(sess.events) > 0 {
		for _, e := range sess.events {
			if err = bus.PublishCtx(ctx, e); err != nil {
				GetLogger(db).ErrorWithContext(ctx, err, "Failed to publish event after commit.")
			}
		}
	}
	return nil
}

func inTransactionWithTenantAndContext(ctx context.Context, db *gorm.DB, callback dbTransactionFuncWithTenant) error {
	var err error
	tenant := contexts.FromTenantContext(ctx)
	if tenant == "" {
		return model.ErrTenantRequired
	}
	defer func() {
		if err != nil {
			GetLogger(db).ErrorWithContext(ctx, err, "inTransactionWithTenantAndContext")
		}
	}()
	sess := startSession(ctx, db, true)
	err = callback(sess, tenant)
	if err != nil {
		sess.Rollback()
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			sess.Rollback()
			GetLogger(db).PanicWithContext(ctx, err.(error), "inTransactionWithCtx")
		}
	}()
	if err = sess.Commit().Error; err != nil {
		return err
	}
	if len(sess.events) > 0 {
		for _, e := range sess.events {
			if err = bus.PublishCtx(ctx, e); err != nil {
				GetLogger(db).ErrorWithContext(ctx, err, "Failed to publish event after commit.")
			}
		}
	}
	return nil
}
