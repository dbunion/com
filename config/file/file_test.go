package file

import (
	"github.com/dbunion/com/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJsonFile(t *testing.T) {
	cfg, err := config.NewConfig(config.TypeFile, config.Param{File: "./test.json", Type: "json"})
	if err != nil {
		t.Fatalf("create config error, err:%v", err)
	}

	val := cfg.Get("key")
	assert.NotNil(t, val)
	assert.Equal(t, "test", val)
}

func TestJson(t *testing.T) {
	cfg, err := config.NewConfig(config.TypeFile, config.Param{Name: "test", Path: "./", Type: "json"})
	if err != nil {
		t.Fatalf("create config error, err:%v", err)
	}

	val := cfg.Get("key")
	assert.NotNil(t, val)
	assert.Equal(t, "test", val)

	bVal := cfg.GetBool("key1")
	assert.Equal(t, true, bVal)

	fVal := cfg.GetFloat64("key2")
	assert.Equal(t, 0.25, fVal)

	iVal := cfg.GetInt("key3")
	assert.Equal(t, 2334, iVal)

	sVal := cfg.GetString("key4")
	assert.Equal(t, "test4", sVal)

	sMap := cfg.GetStringMap("key5")
	assert.NotNil(t, sMap)
	assert.Contains(t, sMap, "name")
	assert.Contains(t, sMap, "age")
	assert.Contains(t, sMap, "sex")

	sMaps := cfg.GetStringMapString("key6")
	assert.NotNil(t, sMaps)
	assert.Contains(t, sMaps, "x1")
	assert.Contains(t, sMaps, "x2")
	assert.Contains(t, sMaps, "x3")

	ss := cfg.GetStringSlice("key7")
	assert.NotNil(t, ss)
	assert.Contains(t, ss, "t1")
	assert.Contains(t, ss, "t2")
	assert.Contains(t, ss, "t3")
	assert.NotContains(t, ss, "t4")
	assert.NotContains(t, ss, "aa")

	tm := cfg.GetTime("key8")
	assert.NotNil(t, tm)
	assert.IsType(t, time.Now(), tm)

	td := cfg.GetDuration("key9")
	assert.NotNil(t, td)
	assert.IsType(t, time.Second, td)

	tLgs := cfg.GetString("key10.key11.name")
	assert.NotNil(t, tLgs)
	assert.Equal(t, "testname", tLgs)
}

func TestYamlFile(t *testing.T) {
	cfg, err := config.NewConfig(config.TypeFile, config.Param{File: "./test.yaml", Type: "yaml"})
	if err != nil {
		t.Fatalf("create config error, err:%v", err)
	}

	val := cfg.Get("key")
	assert.NotNil(t, val)
	assert.Equal(t, "test", val)
}

func TestYaml(t *testing.T) {
	cfg, err := config.NewConfig(config.TypeFile, config.Param{Name: "test", Path: "./", Type: "yaml"})
	if err != nil {
		t.Fatalf("create config error, err:%v", err)
	}

	val := cfg.Get("key")
	assert.NotNil(t, val)
	assert.Equal(t, "test", val)

	bVal := cfg.GetBool("key1")
	assert.Equal(t, true, bVal)

	fVal := cfg.GetFloat64("key2")
	assert.Equal(t, 0.25, fVal)

	iVal := cfg.GetInt("key3")
	assert.Equal(t, 2334, iVal)

	sVal := cfg.GetString("key4")
	assert.Equal(t, "test4", sVal)

	sMap := cfg.GetStringMap("key5")
	assert.NotNil(t, sMap)
	assert.Contains(t, sMap, "name")
	assert.Contains(t, sMap, "age")
	assert.Contains(t, sMap, "sex")

	sMaps := cfg.GetStringMapString("key6")
	assert.NotNil(t, sMaps)
	assert.Contains(t, sMaps, "x1")
	assert.Contains(t, sMaps, "x2")
	assert.Contains(t, sMaps, "x3")

	ss := cfg.GetStringSlice("key7")
	assert.NotNil(t, ss)
	assert.Contains(t, ss, "t1")
	assert.Contains(t, ss, "t2")
	assert.Contains(t, ss, "t3")
	assert.NotContains(t, ss, "t4")
	assert.NotContains(t, ss, "aa")

	tm := cfg.GetTime("key8")
	assert.NotNil(t, tm)
	assert.IsType(t, time.Now(), tm)

	td := cfg.GetDuration("key9")
	assert.NotNil(t, td)
	assert.IsType(t, time.Second, td)

	tLgs := cfg.GetString("key10.key11.name")
	assert.NotNil(t, tLgs)
	assert.Equal(t, "testname", tLgs)

}
