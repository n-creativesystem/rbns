package rdb

import (
	"context"

	"gorm.io/gorm"
)

type DBSession struct {
	*gorm.DB
	transaction bool
}

func (sess *DBSession) Rollback() error {
	if sess.transaction {
		return sess.DB.Rollback().Error
	}
	return nil
}

func (sess *DBSession) Commit() error {
	if sess.transaction {
		return sess.DB.Commit().Error
	}
	return nil
}

func startSession(ctx context.Context, db *gorm.DB, beginTx bool) (context.Context, *DBSession) {
	sess, ok := ctx.Value(ContextSessionKey{}).(*DBSession)
	if ok {
		return ctx, &DBSession{
			DB:          sess.DB,
			transaction: false,
		}
	}
	sess = &DBSession{
		DB: db.Session(&gorm.Session{
			Context: ctx,
		}),
	}
	if beginTx {
		tx := sess.Begin()
		tx.SkipDefaultTransaction = true
		sess.DB = tx
		sess.transaction = true
	}
	ctx = context.WithValue(ctx, ContextSessionKey{}, sess)
	return ctx, sess
}
