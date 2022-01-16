package badger

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	b3 "github.com/dgraph-io/badger/v3"
	"github.com/n-creativesystem/rbns/cache"
	"github.com/n-creativesystem/rbns/infra/kvs/badger"
	"github.com/n-creativesystem/rbns/ncsfw/logger"
)

func init() {
	cache.Register("badger", &badgerCache{
		ttl: 5 * time.Minute,
	})
}

type badgerCache struct {
	db  *badger.CacheBadgerDB
	ttl time.Duration
}

var _ cache.Driver = (*badgerCache)(nil)
var _ cache.Cache = (*badgerCache)(nil)

func (c *badgerCache) Open(dsn string) (cache.Cache, error) {
	var (
		path     string
		inMemory bool
		err      error
	)
	parameters := cache.ParseDSN(dsn)
	for key, value := range parameters {
		switch key {
		case "expiration":
			c.ttl, err = time.ParseDuration(value)
			if err != nil {
				c.ttl = 5 * time.Minute
				logger.Warning(fmt.Sprintf("cache: parse duration warning is expiration(%s)", value))
			}
		case "inmemory":
			inMemory, _ = strconv.ParseBool(value)
		case "path":
			path = value
		}
	}
	db, err := badger.NewCacheDB(path, inMemory)
	if err != nil {
		return nil, err
	}
	c.db = db
	return c, nil
}

func (c *badgerCache) prefix(key string) []byte {
	return []byte(fmt.Sprintf("cache:%s", key))
}

func (c *badgerCache) Set(key string, value interface{}) error {
	return c.db.Update(func(txn *b3.Txn) error {
		buf, err := json.Marshal(value)
		if err != nil {
			return err
		}
		e := b3.NewEntry(c.prefix(key), buf)
		return txn.SetEntry(e)
	})
}

func (c *badgerCache) SetTTL(key string, value interface{}) error {
	return c.db.Update(func(txn *b3.Txn) error {
		buf, err := json.Marshal(value)
		if err != nil {
			return err
		}
		e := b3.NewEntry(c.prefix(key), buf).WithTTL(c.ttl)
		return txn.SetEntry(e)
	})
}

func (c *badgerCache) Get(key string) (value interface{}, ok bool) {
	var valCopy []byte
	if err := c.db.View(func(txn *b3.Txn) error {
		item, err := txn.Get(c.prefix(key))
		if errors.Is(err, b3.ErrKeyNotFound) {
			return cache.ErrNotFound
		}
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
	}); err != nil {
		logger.Error(err, "cache get", "key", key)
		return
	}

	if err := json.Unmarshal(valCopy, &value); err != nil {
		logger.Error(err, "cache get json unmarshal", "key", key)
		value = nil
	}
	ok = true
	return
}

func (c *badgerCache) Delete(key string) error {
	return c.db.Update(func(txn *b3.Txn) error {
		return txn.Delete(c.prefix(key))
	})
}
