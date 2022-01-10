package rdb

import (
	"context"

	"gorm.io/gorm"
)

type DBSession struct {
	*gorm.DB
	events []interface{}
}

func startSession(ctx context.Context, db *gorm.DB, beginTx bool) *DBSession {
	value := ctx.Value(ContextSessionKey{})
	sess, ok := value.(*DBSession)
	if ok {
		sess.DB = sess.DB.Session(&gorm.Session{
			Context: ctx,
		})
		return sess
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
	}
	return sess
}
