package mysql

import (
	"github.com/dbunion/com/uid"
	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMyUID32(t *testing.T) {
	s, err := uid.NewUID(uid.TypeMySQL, uid.Config{
		Server:          "127.0.0.1",
		Port:            3306,
		User:            "test",
		Password:        "123456",
		InitValue:       time.Now().Unix(),
		Step:            10,
		DBName:          "test",
		TableName:       "int32_seq",
		AutoCreateTable: true,
	})

	if err != nil {
		t.Fatalf("start err:%v", err)
		return
	}

	assert.NotNil(t, s)

	for i := 0; i < 50; i++ {
		if s.HasInt32() {
			sequence := s.NextUID32()
			assert.NotEqual(t, -1, sequence)
			t.Logf("Int32:%v", sequence)
		}
	}

	if err := s.Close(); err != nil {
		t.Fatalf("close error:%v", err)
	}

	t.Logf("generate success")
}

func TestNewMyUID64(t *testing.T) {
	s, err := uid.NewUID(uid.TypeMySQL, uid.Config{
		Server:          "127.0.0.1",
		Port:            3306,
		User:            "test",
		Password:        "123456",
		InitValue:       time.Now().UnixNano(),
		Step:            1000,
		DBName:          "test",
		TableName:       "int64_seq",
		AutoCreateTable: true,
	})
	if err != nil {
		t.Fatalf("start err:%v", err)
		return
	}

	assert.NotNil(t, s)

	for i := 0; i < 5000; i++ {
		sequence := s.NextUID64()
		assert.NotEqual(t, -1, sequence)
		t.Logf("Int64:%v", sequence)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("close err:%v", err)
	}

	t.Logf("generate success")
}
