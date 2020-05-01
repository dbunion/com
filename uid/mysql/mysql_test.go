package mysql

import (
	"github.com/dbunion/com/uid"
	"testing"
	"time"
)

func TestNewMyUID32(t *testing.T) {
	s := NewMyUID()
	if err := s.StartAndGC(uid.Config{
		Server:          "127.0.0.1",
		Port:            3306,
		User:            "test",
		Password:        "123456",
		InitValue:       time.Now().Unix(),
		Step:            10,
		DBName:          "test",
		TableName:       "int32_seq",
		AutoCreateTable: true,
	}); err != nil {
		t.Fatalf("start err:%v", err)
		return
	}

	for i := 0; i < 50; i++ {
		if s.HasInt32() {
			t.Logf("Int32:%v", s.NextUID32())
		}
	}

	if err := s.Close(); err != nil {
		t.Fatalf("close error:%v", err)
	}

	t.Logf("generate success")
}

func TestNewMyUID64(t *testing.T) {
	s := NewMyUID()
	if err := s.StartAndGC(uid.Config{
		Server:          "127.0.0.1",
		Port:            3306,
		User:            "test",
		Password:        "123456",
		InitValue:       time.Now().UnixNano(),
		Step:            1000,
		DBName:          "test",
		TableName:       "int64_seq",
		AutoCreateTable: true,
	}); err != nil {
		t.Fatalf("start err:%v", err)
		return
	}

	for i := 0; i < 5000; i++ {
		t.Logf("Int64:%v", s.NextUID64())
	}
	if err := s.Close(); err != nil {
		t.Fatalf("close err:%v", err)
	}

	t.Logf("generate success")
}
