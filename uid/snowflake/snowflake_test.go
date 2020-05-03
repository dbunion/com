package snowflake

import (
	"github.com/dbunion/com/uid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSnowflake(t *testing.T) {
	s, err := uid.NewUID(uid.TypeSnowFlake, uid.Config{NodeID: 1})
	if err != nil {
		t.Fatalf("create new uid error, err:%v", err)
	}

	assert.NotNil(t, s)

	for i := 0; i < 10; i++ {
		sequence := s.NextUID64()
		assert.NotEqual(t, -1, sequence)
		t.Logf("%d", sequence)
	}

	t.Logf("generate success")
}

func BenchmarkNewSnowflake(b *testing.B) {
	s1, err := uid.NewUID(uid.TypeSnowFlake, uid.Config{NodeID: 1})
	if err != nil {
		b.Fatalf("create new uid error, err:%v", err)
	}

	s2, err := uid.NewUID(uid.TypeSnowFlake, uid.Config{NodeID: 2})
	if err != nil {
		b.Fatalf("create new uid error, err:%v", err)
	}

	s3, err := uid.NewUID(uid.TypeSnowFlake, uid.Config{NodeID: 3})
	if err != nil {
		b.Fatalf("create new uid error, err:%v", err)
	}

	for i := 0; i < b.N; i++ {
		b.Logf("Node:1 UID:%d", s1.NextUID64())
		b.Logf("Node:2 UID:%d", s2.NextUID64())
		b.Logf("Node:3 UID:%d", s3.NextUID64())
	}
}
