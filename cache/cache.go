package cache

import (
	"fmt"
	"time"
)

const (
	// TypeGoCache - type go cache
	TypeGoCache = "goCache"
	// TypeMemoryCache - type use memcached
	TypeMemoryCache = "memCache"
	// TypeRedisCache - type use redis
	TypeRedisCache = "redisCache"
)

// Config - define cache config
type Config struct {
	Expiration time.Duration `json:"expiration"`

	// public common
	Server   string `json:"server"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int64  `json:"port"`

	// redis
	DBNum     int64  `json:"db_num"`
	MaxIdle   int64  `json:"max_idle"`
	MaxActive int64  `json:"max_active"`
	Key       string `json:"key"`

	// Extend fields
	// Extended fields can be used if there is a special implementation
	Extend1 string `json:"extend_1"`
	Extend2 string `json:"extend_2"`
}

// Cache interface contains all behaviors for cache adapter.
type Cache interface {
	// get cached value by key.
	Get(key string) interface{}
	// GetMulti is a batch version of Get.
	GetMulti(keys []string) []interface{}
	// set cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// increase cached int value by key and return value.
	IncrBy(key string) (interface{}, error)
	// decrease cached int value by key and return value.
	DecrBy(key string) (interface{}, error)

	// check if cached value exists or not.
	IsExist(key string) bool
	// clear all cache.
	ClearAll() error
	// start gc routine based on config settings.
	StartAndGC(config Config) error

	// lock key
	TryLock(key string, val interface{}, timeout time.Duration) error
	// unlock key
	UnLock(key string, val interface{}) error
	// set key
	Set(key string, val interface{}) (bool, error)
	// expire key
	Expire(key string, timeout time.Duration) error
}

// Instance is a function create a new Cache Instance
type Instance func() Cache

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache Create a new cache driver by adapter name and config setting.
// it will start gc automatically.
func NewCache(adapterName string, config Config) (adapter Cache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
