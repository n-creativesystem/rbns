package cache

import (
	"fmt"
	"strings"
	"sync"

	"github.com/n-creativesystem/rbns/config"
)

type Cache interface {
	Get(key string) (value interface{}, ok bool)
	Set(key string, value interface{}) error
	SetTTL(key string, value interface{}) error
	Delete(key string) error
}

type Driver interface {
	Open(dsn string) (Cache, error)
}

var (
	drivers   = make(map[string]Driver)
	driversMu sync.Mutex
)

func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("cache: Register driver is nil")
	}

	if _, dup := drivers[name]; dup {
		panic("cache: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Open(driverName, dataSource string) (Cache, error) {
	driversMu.Lock()
	driver, ok := drivers[driverName]
	driversMu.Unlock()
	if !ok {
		return nil, fmt.Errorf("cache: unknown driver %q (forgotten import?)", driverName)
	}
	return driver.Open(dataSource)
}

func ParseDSN(dsn string) map[string]string {
	parameters := strings.Split(dsn, ";")
	mp := make(map[string]string, len(parameters))
	for _, parameter := range parameters {
		kv := strings.SplitN(parameter, "=", 2)
		if len(kv) == 2 {
			mp[kv[0]] = kv[1]
		}
	}
	return mp
}

func New(cfg *config.Config) (Cache, error) {
	return Open(cfg.CacheDriverName, cfg.CacheDataSource)
}
