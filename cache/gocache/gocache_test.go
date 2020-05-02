package gocache

import (
	"testing"
	"time"

	"github.com/dbunion/com/cache"
	"github.com/dbunion/com/conv"
	_ "github.com/patrickmn/go-cache"
)

const (
	cacheKey  = "goCacheKey"
	cacheKey1 = "goCacheKey1"
)

func TestGoCache(t *testing.T) {
	bm, err := cache.NewCache(cache.TypeGoCache, cache.Config{Expiration: 0})
	if err != nil {
		t.Fatalf("init err:%v", err)
	}

	timeoutDuration := 10 * time.Second
	if err = bm.Put(cacheKey, 1, timeoutDuration); err != nil {
		t.Fatalf("put key error, err:%v", err)
	}

	if !bm.IsExist(cacheKey) {
		t.Fatal("check key exist error")
	}

	time.Sleep(11 * time.Second)
	if bm.IsExist(cacheKey) {
		t.Fatal("check key exist error")
	}

	if err = bm.Put(cacheKey, 1, timeoutDuration); err != nil {
		t.Fatalf("put key error, err:%v", err)
	}

	if v := conv.GetInt(bm.Get(cacheKey)); v != 1 {
		t.Fatalf("get key error, err:%v", v)
	}

	if err = bm.Incr(cacheKey); err != nil {
		t.Fatalf("incr key error, err:%v", err)
	}

	if v := conv.GetInt(bm.Get(cacheKey)); v != 2 {
		t.Fatalf("get key error, err:%v", v)
	}

	if err = bm.Decr(cacheKey); err != nil {
		t.Fatalf("decr key error, err:%v", err)
	}

	if v := conv.GetInt(bm.Get(cacheKey)); v != 1 {
		t.Fatal("get key error")
	}

	if err = bm.Delete(cacheKey); err != nil {
		t.Fatalf("del key error, err:%v", err)
	}

	if bm.IsExist(cacheKey) {
		t.Fatalf("check key exist error")
	}

	// test string
	if err = bm.Put(cacheKey, "author", timeoutDuration); err != nil {
		t.Fatalf("put key error:%v", err)
	}

	if !bm.IsExist(cacheKey) {
		t.Fatalf("check key exist error")
	}

	if v := conv.GetString(bm.Get(cacheKey)); v != "author" {
		t.Fatalf("get key error")
	}

	// test GetMulti
	if err = bm.Put(cacheKey1, "author1", timeoutDuration); err != nil {
		t.Fatalf("put key error:%v", err)
	}

	if !bm.IsExist(cacheKey1) {
		t.Fatal("check key exist error")
	}

	vv := bm.GetMulti([]string{cacheKey, cacheKey1})
	if len(vv) != 2 {
		t.Fatalf("get multi error")
	}

	if v := conv.GetString(vv[0]); v != "author" {
		t.Fatalf("get multi error")
	}

	if v := conv.GetString(vv[1]); v != "author1" {
		t.Fatalf("get multi error")
	}

	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Fatalf("clear all error, err:%v", err)
	}
}
