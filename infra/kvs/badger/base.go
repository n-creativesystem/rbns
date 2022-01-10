package badger

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/n-creativesystem/rbns/logger"
)

type CacheBadgerDB struct {
	*badger.DB
}

type badgerLogger struct {
	log logger.Logger
}

func (log *badgerLogger) Errorf(format string, a ...interface{}) {
	err := fmt.Errorf(format, a...)
	log.log.Error(err, err.Error())
}

func (log *badgerLogger) Warningf(format string, a ...interface{}) {
	log.log.Warning(fmt.Sprintf(format, a...))
}

func (log *badgerLogger) Infof(format string, a ...interface{}) {
	log.log.Info(fmt.Sprintf(format, a...))
}

func (log *badgerLogger) Debugf(format string, a ...interface{}) {
	log.log.Debug(fmt.Sprintf(format, a...))
}

func newBadgerDB(path string, inMemory bool) (*badger.DB, error) {
	opts := badger.DefaultOptions(path)
	opts = opts.WithLogger(&badgerLogger{logger.New("bagger")})
	opts = opts.WithInMemory(inMemory)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	go gc(db)
	return db, nil
}

func NewCacheDB(path string, inMemory bool) (*CacheBadgerDB, error) {
	db, err := newBadgerDB(path, inMemory)
	if err != nil {
		return nil, err
	}
	return &CacheBadgerDB{
		DB: db,
	}, nil
}

func gc(db *badger.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
	again:
		err := db.RunValueLogGC(0.7)
		if err == nil {
			goto again
		}
	}
}
