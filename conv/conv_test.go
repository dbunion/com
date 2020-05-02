package conv

import (
	"testing"
)

func TestGetString(t *testing.T) {
	var t1 = "test1"
	if GetString(t1) != "test1" {
		t.Fatal("get string from string error")
	}
	var t2 = []byte("test2")
	if GetString(t2) != "test2" {
		t.Fatal("get string from byte array error")
	}
	var t3 = 1
	if GetString(t3) != "1" {
		t.Fatal("get string from int error")
	}
	var t4 int64 = 1
	if GetString(t4) != "1" {
		t.Fatal("get string from int64 error")
	}
	var t5 = 1.1
	if GetString(t5) != "1.1" {
		t.Fatal("get string from float64 error")
	}

	if GetString(nil) != "" {
		t.Fatal("get string from nil error")
	}
}

func TestGetInt(t *testing.T) {
	var t1 = 1
	if GetInt(t1) != 1 {
		t.Fatal("get int from int error")
	}
	var t2 int32 = 32
	if GetInt(t2) != 32 {
		t.Fatal("get int from int32 error")
	}
	var t3 int64 = 64
	if GetInt(t3) != 64 {
		t.Fatal("get int from int64 error")
	}
	var t4 = "128"
	if GetInt(t4) != 128 {
		t.Fatal("get int from num string error")
	}
	if GetInt(nil) != 0 {
		t.Fatal("get int from nil error")
	}
}

func TestGetInt64(t *testing.T) {
	var i int64 = 1
	var t1 = 1
	if i != GetInt64(t1) {
		t.Fatal("get int64 from int error")
	}
	var t2 int32 = 1
	if i != GetInt64(t2) {
		t.Fatal("get int64 from int32 error")
	}
	var t3 int64 = 1
	if i != GetInt64(t3) {
		t.Fatal("get int64 from int64 error")
	}
	var t4 = "1"
	if i != GetInt64(t4) {
		t.Fatal("get int64 from num string error")
	}
	if GetInt64(nil) != 0 {
		t.Fatal("get int64 from nil")
	}
}

func TestGetFloat64(t *testing.T) {
	var f = 1.11
	var t1 float32 = 1.11
	if f != GetFloat64(t1) {
		t.Fatal("get float64 from float32 error")
	}
	var t2 = 1.11
	if f != GetFloat64(t2) {
		t.Fatal("get float64 from float64 error")
	}
	var t3 = "1.11"
	if f != GetFloat64(t3) {
		t.Fatal("get float64 from string error")
	}

	var f2 float64 = 1
	var t4 = 1
	if f2 != GetFloat64(t4) {
		t.Fatal("get float64 from int error")
	}

	if GetFloat64(nil) != 0 {
		t.Fatal("get float64 from nil error")
	}
}

func TestGetBool(t *testing.T) {
	var t1 = true
	if !GetBool(t1) {
		t.Fatal("get bool from bool error")
	}
	var t2 = "true"
	if !GetBool(t2) {
		t.Fatal("get bool from string error")
	}
	if GetBool(nil) {
		t.Fatal("get bool from nil error")
	}
}
