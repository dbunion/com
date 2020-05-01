package redis

import (
	"github.com/dbunion/com/uid"

	"testing"
)

func TestNewRedisInt32UID(t *testing.T) {
	s := NewRedisUID()
	if err := s.StartAndGC(uid.Config{Key: "uid", Server: "127.0.0.1", Port: 6379}); err != nil {
		t.Fatalf("start err:%v", err)
		return
	}

	// gen uid
	for i := 0; i < 10; i++ {
		if s.HasInt32() {
			t.Logf("Int32:%d", s.NextUID32())
		}
	}
	t.Logf("generate success")
}

func TestNewRedisInt64UID(t *testing.T) {
	s := NewRedisUID()
	if err := s.StartAndGC(uid.Config{Key: "uid", Server: "127.0.0.1", Port: 6379}); err != nil {
		t.Fatalf("start err:%v", err)
		return
	}

	// gen uid
	for i := 0; i < 10; i++ {
		t.Logf("Int64:%d", s.NextUID64())
	}
	t.Logf("generate success")
}
