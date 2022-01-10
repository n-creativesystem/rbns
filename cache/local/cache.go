package local

import (
	"fmt"
	"strings"
	"time"

	"github.com/n-creativesystem/rbns/cache"
	"github.com/n-creativesystem/rbns/logger"
	goCache "github.com/patrickmn/go-cache"
)

func init() {
	cache.Register("local", &localCache{})
}

type localCache struct {
	cache *goCache.Cache
}

var _ cache.Driver = (*localCache)(nil)
var _ cache.Cache = (*localCache)(nil)

func (c *localCache) Open(dsn string) (cache.Cache, error) {
	var (
		defaultExpiration time.Duration
		cleanupInterval   time.Duration
		err               error
	)
	parameters := cache.ParseDSN(dsn)
	for key, value := range parameters {
		switch strings.ToLower(key) {
		case "expiration":
			defaultExpiration, err = time.ParseDuration(value)
			if err != nil {
				defaultExpiration = 5 * time.Minute
				logger.Warning(fmt.Sprintf("cache: parse duration warning is expiration(%s)", value))
			}
		case "cleanupinterval":
			cleanupInterval, err = time.ParseDuration(value)
			if err != nil {
				cleanupInterval = 10 * time.Minute
				logger.Warning(fmt.Sprintf("cache: parse duration warning is cleanup interval(%s)", value))
			}
		}
	}
	c.cache = goCache.New(defaultExpiration, cleanupInterval)
	return c, nil
}

func (c *localCache) SetTTL(key string, value interface{}) error {
	c.cache.SetDefault(key, value)
	return nil
}

func (c *localCache) Set(key string, value interface{}) error {
	c.cache.Set(key, value, -1)
	return nil
}

func (c *localCache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *localCache) Delete(key string) error {
	c.cache.Delete(key)
	return nil
}
