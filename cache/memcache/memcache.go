package memcache

import (
	"errors"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/dbunion/com/cache"
)

// Cache - Memcache adapter.
type Cache struct {
	conn   *memcache.Client
	server string
	port   int64
}

// NewMemCache create new memcache adapter.
func NewMemCache() cache.Cache {
	return &Cache{}
}

// Get get value from memcache.
func (rc *Cache) Get(key string) interface{} {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		return item.Value
	}
	return nil
}

// GetMulti get value from memcache.
func (rc *Cache) GetMulti(keys []string) []interface{} {
	size := len(keys)
	var rv []interface{}
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			for i := 0; i < size; i++ {
				rv = append(rv, err)
			}
			return rv
		}
	}
	mv, err := rc.conn.GetMulti(keys)
	if err == nil {
		for _, v := range mv {
			rv = append(rv, v.Value)
		}
		return rv
	}
	for i := 0; i < size; i++ {
		rv = append(rv, err)
	}
	return rv
}

// Put put value to memcache.
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	item := memcache.Item{Key: key, Expiration: int32(timeout / time.Second)}
	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		return errors.New("val only support string and []byte")
	}
	return rc.conn.Set(&item)
}

// Delete delete value in memcache.
func (rc *Cache) Delete(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.Delete(key)
}

// Incr increase counter.
func (rc *Cache) Incr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Increment(key, 1)
	return err
}

// Decr decrease counter.
func (rc *Cache) Decr(key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Decrement(key, 1)
	return err
}

// IncrBy increase counter.
func (rc *Cache) IncrBy(key string) (interface{}, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil, err
		}
	}
	return rc.conn.Increment(key, 1)
}

// DecrBy decrease counter.
func (rc *Cache) DecrBy(key string) (interface{}, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil, err
		}
	}
	return rc.conn.Decrement(key, 1)
}

// IsExist check value exists in memcache.
func (rc *Cache) IsExist(key string) bool {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false
		}
	}
	_, err := rc.conn.Get(key)
	return !(err != nil)
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

// ClearAll clear all cached in memcache.
func (rc *Cache) ClearAll() error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.FlushAll()
}

// StartAndGC start memcache adapter.
// config string is like {"conn":"connection info"}.
// if connecting error, return.
func (rc *Cache) StartAndGC(config cache.Config) error {
	rc.server = config.Server
	rc.port = config.Port
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
		if err := rc.conn.Ping(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcached and keep the connection.
func (rc *Cache) connectInit() error {
	rc.conn = memcache.New(fmt.Sprintf("%s:%d", rc.server, rc.port))
	return nil
}

func init() {
	cache.Register(cache.TypeMemoryCache, NewMemCache)
}
