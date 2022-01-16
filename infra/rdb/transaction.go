package rdb

import (
	"context"

	"github.com/n-creativesystem/rbns/domain/model"
	"github.com/n-creativesystem/rbns/ncsfw/tenants"
	"gorm.io/gorm"
)

const (
	dammyTenant = "dammy"
)

type dbSessionFunc func(ctx context.Context, sess *DBSession) error

type dbSessionFuncWithTenant func(ctx context.Context, sess *DBSession, tenant string) error

type dbTransactionFunc func(ctx context.Context, sess *DBSession) error

type dbTransactionFuncWithTenant func(ctx context.Context, sess *DBSession, tenant string) error

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
	return inTransactionWithTenantAndContext(ctx, f.db, func(ctx context.Context, sess *DBSession, tenant string) error {
		withValue := context.WithValue(ctx, ContextSessionKey{}, sess)
		return fn(withValue)
	})
}

func dbSssion(ctx context.Context, db *gorm.DB, callback dbSessionFunc) error {
	ctx, sess := startSession(ctx, db, false)
	err := callback(ctx, sess)
	if err != nil {
		GetLogger(db).ErrorWithContext(ctx, err, "dbSssion")
		return err
	}
	return nil
}

func dbSssionWithToken(ctx context.Context, db *gorm.DB, callback dbSessionFuncWithTenant) error {
	tenant := tenants.FromTenantContext(ctx)
	if tenant == "" {
		return model.ErrTenantEmpty
	} else if tenant == dammyTenant {
		tenant = ""
	}
	ctx, sess := startSession(ctx, db, false)
	err := callback(ctx, sess, tenant)
	if err != nil {
		GetLogger(db).ErrorWithContext(ctx, err, "dbSssionWithToken")
		return err
	}
	return nil
}

func inTransactionWithDbSession(ctx context.Context, db *gorm.DB, callback dbTransactionFunc) error {
	ctx, sess := startSession(ctx, db, true)
	err := callback(ctx, sess)
	if err != nil {
		_ = sess.Rollback()
		GetLogger(db).ErrorWithContext(ctx, err, "inTransactionWithDbSession")
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			_ = sess.Rollback()
			GetLogger(db).ErrorWithContext(ctx, err.(error), "inTransactionWithDbSession")
		}
	}()
	if err := sess.Commit(); err != nil {
		GetLogger(db).ErrorWithContext(ctx, err, "inTransactionWithDbSession")
		return err
	}
	return nil
}

func inTransactionWithTenantAndContext(ctx context.Context, db *gorm.DB, callback dbTransactionFuncWithTenant) error {
	tenant := tenants.FromTenantContext(ctx)
	if tenant == "" {
		return model.ErrTenantEmpty
	}
	ctx, sess := startSession(ctx, db, true)
	err := callback(ctx, sess, tenant)
	if err != nil {
		_ = sess.Rollback()
		GetLogger(db).ErrorWithContext(ctx, err, "inTransactionWithTenantAndContext")
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			_ = sess.Rollback()
			GetLogger(db).ErrorWithContext(ctx, err.(error), "inTransactionWithTenantAndContext")
		}
	}()
	if err := sess.Commit(); err != nil {
		GetLogger(db).ErrorWithContext(ctx, err, "inTransactionWithTenantAndContext")
		return err
	}
	return nil
}
