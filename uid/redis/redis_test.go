package redis

import (
	"github.com/dbunion/com/uid"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestNewRedisInt32UID(t *testing.T) {
	s, err := uid.NewUID(uid.TypeRedis, uid.Config{Key: "uid", Server: "127.0.0.1", Port: 6379})
	if err != nil {
		t.Fatalf("start err:%v", err)
		return
	}

	// gen uid
	for i := 0; i < 10; i++ {
		if s.HasInt32() {
			sequence := s.NextUID32()
			assert.NotEqual(t, -1, sequence)
			t.Logf("Int32:%d", sequence)
		}
	}
	t.Logf("generate success")
}

func TestNewRedisInt64UID(t *testing.T) {
	s, err := uid.NewUID(uid.TypeRedis, uid.Config{Key: "uid", Server: "127.0.0.1", Port: 6379})
	if err != nil {
		t.Fatalf("start err:%v", err)
		return
	}

	assert.NotNil(t, s)

	// gen uid
	for i := 0; i < 10; i++ {
		sequence := s.NextUID64()

		assert.NotEqual(t, -1, sequence)
		t.Logf("Int64:%d", sequence)
	}
	t.Logf("generate success")
}
