package snowflake

import (
	"github.com/dbunion/com/uid"
	"testing"
)

func TestNewSnowflake(t *testing.T) {
	s := NewSnowflake()
	config := uid.Config{NodeID: 1}
	if err := s.StartAndGC(config); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		t.Logf("%d", s.NextUID64())
	}

	t.Logf("generate success")
}

func BenchmarkNewSnowflake(b *testing.B) {
	s1 := NewSnowflake()
	if err := s1.StartAndGC(uid.Config{NodeID: 1}); err != nil {
		panic(err)
	}

	s2 := NewSnowflake()
	if err := s2.StartAndGC(uid.Config{NodeID: 2}); err != nil {
		panic(err)
	}

	s3 := NewSnowflake()
	if err := s3.StartAndGC(uid.Config{NodeID: 3}); err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		b.Logf("Node:1 UID:%d", s1.NextUID64())
		b.Logf("Node:2 UID:%d", s2.NextUID64())
		b.Logf("Node:3 UID:%d", s3.NextUID64())
	}
}
