package memcache

import (
	_ "github.com/bradfitz/gomemcache/memcache"
	"github.com/dbunion/com/cache"
	"strconv"
	"testing"
	"time"
)

const (
	cacheKey  = "memoryCacheKey"
	cacheKey1 = "memoryCacheKey1"
)

func TestMemcachedCache(t *testing.T) {
	bm, err := cache.NewCache(cache.TypeMemoryCache, cache.Config{Server: "127.0.0.1", Port: 11211})
	if err != nil {
		t.Fatalf("create new cache error, err:%v", err)
	}

	timeoutDuration := 10 * time.Second
	if err = bm.Put(cacheKey, "1", timeoutDuration); err != nil {
		t.Fatalf("put key error, err:%v", err)
	}
	if !bm.IsExist(cacheKey) {
		t.Fatal("check key exist error")
	}

	time.Sleep(11 * time.Second)

	if bm.IsExist(cacheKey) {
		t.Fatal("check key exist error")
	}
	if err = bm.Put(cacheKey, "1", timeoutDuration); err != nil {
		t.Fatal("put key error", err)
	}

	if v, err := strconv.Atoi(string(bm.Get(cacheKey).([]byte))); err != nil || v != 1 {
		t.Fatal("get key error")
	}

	if err = bm.Incr(cacheKey); err != nil {
		t.Fatal("Incr Error", err)
	}

	if v, err := strconv.Atoi(string(bm.Get(cacheKey).([]byte))); err != nil || v != 2 {
		t.Fatal("get key error")
	}

	if err = bm.Decr(cacheKey); err != nil {
		t.Fatal("Decr Error", err)
	}

	if v, err := strconv.Atoi(string(bm.Get(cacheKey).([]byte))); err != nil || v != 1 {
		t.Fatal("get key error")
	}

	if err = bm.Delete(cacheKey); err != nil {
		t.Fatal("del Error", err)
	}

	if bm.IsExist(cacheKey) {
		t.Fatal("delete err")
	}

	//test string
	if err = bm.Put(cacheKey, "author", timeoutDuration); err != nil {
		t.Fatal("put key error", err)
	}
	if !bm.IsExist(cacheKey) {
		t.Fatal("check key exist error")
	}

	if v := bm.Get(cacheKey).([]byte); string(v) != "author" {
		t.Fatal("get key error")
	}

	//test GetMulti
	if err = bm.Put(cacheKey1, "author1", timeoutDuration); err != nil {
		t.Fatal("put key error", err)
	}
	if !bm.IsExist(cacheKey1) {
		t.Fatal("check key exist error")
	}

	vv := bm.GetMulti([]string{cacheKey, cacheKey1})
	if len(vv) != 2 {
		t.Fatal("get multi error")
	}
	if string(vv[0].([]byte)) != "author" && string(vv[0].([]byte)) != "author1" {
		t.Fatal("get multi error")
	}
	if string(vv[1].([]byte)) != "author1" && string(vv[1].([]byte)) != "author" {
		t.Fatal("get multi error")
	}

	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Fatal("clear all err")
	}
}
