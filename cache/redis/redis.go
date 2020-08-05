package redis

import (
	"errors"
	"fmt"
	"github.com/dbunion/com/cache"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	defaultKey = "cache"
)

// Cache is Redis cache adapter.
type Cache struct {
	p         *redis.Pool // redis connection pool
	server    string
	connInfo  string
	port      int64
	dbNum     int64
	key       string
	password  string
	maxIdle   int64
	maxActive int64
}

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() cache.Cache {
	return &Cache{key: defaultKey}
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer func() {
		err = c.Close()
	}()

	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

// Get cache from redis.
func (rc *Cache) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// GetMulti get cache from redis.
func (rc *Cache) GetMulti(keys []string) []interface{} {
	c := rc.p.Get()
	defer func() {
		if err := c.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Put put cache to redis.
func (rc *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Delete delete cache in redis.
func (rc *Cache) Delete(key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *Cache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

// Incr increase counter in redis.
func (rc *Cache) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (rc *Cache) Decr(key string) error {
	_, err := redis.Bool(rc.do("DECRBY", key, 1))
	return err
}

// IncrBy increase counter in redis.
func (rc *Cache) IncrBy(key string) (interface{}, error) {
	return rc.do("INCRBY", key, 1)
}

// DecrBy decrease counter in redis.
func (rc *Cache) DecrBy(key string) (interface{}, error) {
	return rc.do("INCRBY", key, -1)
}

// TryLock ...
func (rc *Cache) TryLock(key string, val interface{}, timeout time.Duration) error {
	luaLock := `
if redis.call('setNx',KEYS[1],ARGV[1])  then
   if redis.call('get',KEYS[1])==ARGV[1] then
      return redis.call('expire',KEYS[1],ARGV[2])
   else
      return 0
   end
end`
	lua := redis.NewScript(1, luaLock)
	conn := rc.p.Get()
	_, err := lua.Do(conn, fmt.Sprintf("%v", key), fmt.Sprintf("%v", val), int64(timeout/time.Second))
	if err != nil {
		return err
	}
	return nil
}

// UnLock ...
func (rc *Cache) UnLock(key string, val interface{}) error {
	luaLock := `
if redis.call('get', KEYS[1]) == ARGV[1]
then
return redis.call('del', KEYS[1])
else
return false
end
`
	lua := redis.NewScript(1, luaLock)
	conn := rc.p.Get()
	_, err := lua.Do(conn, key, val)
	if err != nil {
		return err
	}
	return nil
}

// Set put cache to redis.
func (rc *Cache) Set(key string, val interface{}) (bool, error) {
	v, err := redis.Bool(rc.do("SETNX", key, val))
	if err != nil {
		return false, err
	}

	return v, nil
}

// Expire key.
func (rc *Cache) Expire(key string, timeout time.Duration) error {
	ttl := int(timeout / time.Second)
	if ttl < 0 {
		return nil
	}

	_, err := rc.do("expire", key, ttl)
	if err != nil {
		return err
	}

	return nil
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *Cache) ClearAll() error {
	c := rc.p.Get()
	defer func() {
		if err := c.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	cachedKeys, err := redis.Strings(c.Do("KEYS", rc.key+":*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *Cache) StartAndGC(config cache.Config) error {
	rc.key = defaultKey
	if config.Key != "" {
		rc.key = config.Key
	}

	if config.Server == "" {
		return errors.New("server config is empty")
	}

	rc.server = config.Server
	rc.port = config.Port
	rc.password = config.Password
	rc.dbNum = config.DBNum

	connInfo := fmt.Sprintf("%s:%d", rc.server, rc.port)
	rc.connInfo = connInfo
	rc.maxIdle = config.MaxIdle
	rc.maxActive = config.MaxActive

	rc.connectInit()

	c := rc.p.Get()
	defer func() {
		if err := c.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	return c.Err()
}

// connect to redis.
func (rc *Cache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.connInfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err = c.Do("AUTH", rc.password); err != nil {
				if cErr := c.Close(); cErr != nil {
					return nil, fmt.Errorf("authError:%v CloseError:%v", err, cErr)
				}
				return nil, err
			}
		}

		_, selectErr := c.Do("SELECT", rc.dbNum)
		if selectErr != nil {
			if cErr := c.Close(); cErr != nil {
				return nil, fmt.Errorf("selectError:%v CloseError:%v", selectErr, cErr)
			}
			return nil, selectErr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     int(rc.maxIdle),
		MaxActive:   int(rc.maxActive),
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func init() {
	cache.Register(cache.TypeRedisCache, NewRedisCache)
}
