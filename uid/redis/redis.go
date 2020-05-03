package redis

import (
	"fmt"
	"time"

	"github.com/dbunion/com/uid"
	"github.com/go-redis/redis/v7"
)

const (
	int32Key = "uid_int32_key"
	int64Key = "uid_int64_key"
	timeout  = time.Second * 60
)

// UID - Redis uuid gen
type UID struct {
	client *redis.Client
}

// NewRedisUID - create new uid with default collection name.
func NewRedisUID() uid.UID {
	return &UID{}
}

// HasInt32 - has int32 uid
func (r *UID) HasInt32() bool {
	return true
}

// NextUID32 - next int32 uid
func (r *UID) NextUID32() int32 {

	if _, err := r.client.Get(int32Key).Result(); err == redis.Nil {
		r.client.Set(int32Key, time.Now().Unix(), timeout)
	}

	val, err := r.client.IncrBy(int32Key, 1).Result()
	if err != nil {
		return -1
	}

	return int32(val)
}

// NextUID64 - next int64 uid
func (r *UID) NextUID64() int64 {

	if _, err := r.client.Get(int64Key).Result(); err == redis.Nil {
		r.client.Set(int64Key, time.Now().UnixNano(), timeout)
	}

	val, err := r.client.IncrBy(int64Key, 1).Result()
	if err != nil {
		return -1
	}

	return val
}

// Close - close connection
func (r *UID) Close() error {
	// do nothing
	return nil
}

// StartAndGC start uid adapter.
// config is like {"conn":"127.0.0.1:6379"}
// so no gc operation.
func (r *UID) StartAndGC(config uid.Config) error {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%v:%v", config.Server, config.Port),
		Password:        config.Password,
		DB:              0,
		MaxRetries:      3,
		MinRetryBackoff: time.Second * 1,
		MaxRetryBackoff: time.Second * 3,
		DialTimeout:     time.Second * 10,
		ReadTimeout:     time.Second * 10,
		WriteTimeout:    time.Second * 10,
		PoolSize:        5,
		MinIdleConns:    1,
	})

	if _, err := client.Ping().Result(); err != nil {
		return err
	}

	r.client = client

	return nil
}

func init() {
	uid.Register(uid.TypeRedis, NewRedisUID)
}
