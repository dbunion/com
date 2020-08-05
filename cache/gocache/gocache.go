package gocache

import (
	"fmt"
	"time"

	"github.com/dbunion/com/cache"
	goCache "github.com/patrickmn/go-cache"
)

// Cache is go-cache adapter.
type Cache struct {
	p                 *goCache.Cache
	defaultExpiration time.Duration
}

// NewGoCache create new go-cache with default collection name.
func NewGoCache() cache.Cache {
	return &Cache{defaultExpiration: goCache.DefaultExpiration}
}

// Get cache from go-cache.
func (rc *Cache) Get(key string) interface{} {
	if v, found := rc.p.Get(key); found {
		return v
	}
	return nil
}

// GetMulti get cache from go-cache.
func (rc *Cache) GetMulti(keys []string) []interface{} {
	values := make([]interface{}, 0)
	for _, key := range keys {
		v, found := rc.p.Get(key)
		if found {
			values = append(values, v)
			continue
		}
		values = append(values, nil)
	}
	return values
}

// Put put cache to go-cache.
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	rc.p.Set(key, val, timeout)
	return nil
}

// Delete delete cache in go-cache.
func (rc *Cache) Delete(key string) error {
	rc.p.Delete(key)
	return nil
}

// IsExist check cache's existence in go-cache.
func (rc *Cache) IsExist(key string) bool {
	_, found := rc.p.Get(key)
	return found
}

// Incr increase counter in go-cache.
func (rc *Cache) Incr(key string) error {
	return rc.p.Increment(key, 1)
}

// Decr decrease counter in go-cache.
func (rc *Cache) Decr(key string) error {
	return rc.p.Decrement(key, 1)
}

// IncrBy increase counter in go-cache.
func (rc *Cache) IncrBy(key string) (interface{}, error) {
	if err := rc.p.Increment(key, 1); err != nil {
		return nil, err
	}

	if v, found := rc.p.Get(key); found {
		return v, nil
	}

	return nil, fmt.Errorf("key not exist")
}

// DecrBy decrease counter in go-cache.
func (rc *Cache) DecrBy(key string) (interface{}, error) {
	if err := rc.p.Decrement(key, -1); err != nil {
		return nil, err
	}

	if v, found := rc.p.Get(key); found {
		return v, nil
	}
	return nil, fmt.Errorf("key not exist")
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *Cache) ClearAll() error {
	items := rc.p.Items()
	for key := range items {
		if err := rc.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// TryLock ...
func (rc *Cache) TryLock(key string, val interface{}, timeout time.Duration) error {
	return nil
}

// UnLock ...
func (rc *Cache) UnLock(key string, val interface{}) error {
	return nil
}

// Set put cache to redis.
func (rc *Cache) Set(key string, val interface{}) (bool, error) {

	return false, nil
}

// Expire key.
func (rc *Cache) Expire(key string, timeout time.Duration) error {
	return nil
}

// StartAndGC start go-cache adapter.
// the cache item in go-cache are only in memory,
// so no gc operation.
func (rc *Cache) StartAndGC(config cache.Config) error {
	rc.defaultExpiration = config.Expiration
	rc.p = goCache.New(rc.defaultExpiration, rc.defaultExpiration)

	return nil
}

func init() {
	cache.Register(cache.TypeGoCache, NewGoCache)
}
