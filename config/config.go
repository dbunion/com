package config

import (
	"fmt"
	"time"
)

const (
	// TypeFile - config type file
	TypeFile = "file"
	// TypeETCD - config type etcd
	TypeETCD = "etcd"
	// TypeConsul - config type consul
	TypeConsul = "consul"
)

// Param - config param
type Param struct {
	// public param
	Type string `json:"type"`

	// for key/value store
	Server   string `json:"server"`
	User     string `json:"user"`
	Password string `json:"password"`
	Prefix   string `json:"prefix"`
	Port     int64  `json:"port"`

	// file param
	Path string `json:"path"`
	Name string `json:"name"`
	File string `json:"file"`

	// Extend fields
	// Extended fields can be used if there is a special implementation
	Extend1 string `json:"extend_1"`
	Extend2 string `json:"extend_2"`
}

// Config interface contains all behaviors for config adapter.
type Config interface {
	// get config value by key.
	Get(key string) interface{}
	// get bool config value by key
	GetBool(key string) bool
	// get float64 config value by key
	GetFloat64(key string) float64
	// get Int config value by key
	GetInt(key string) int
	// get Int32 config value by key
	GetInt32(key string) int32
	// get Int64 config value by key
	GetInt64(key string) int64
	// get string config value by key
	GetString(key string) string
	// get stringMap config value by key
	GetStringMap(key string) map[string]interface{}
	// get stringMapstring config value by key
	GetStringMapString(key string) map[string]string
	// get stringSlice config value by key
	GetStringSlice(key string) []string
	// get time config value by key
	GetTime(key string) time.Time
	// get Duration config value by key
	GetDuration(key string) time.Duration
	// check if config value exists or not.
	IsExist(key string) bool
	// get all config value
	AllSettings() map[string]interface{}
	// start gc routine based on config param settings.
	StartAndGC(param Param) error
}

// Instance is a function create a new Config Instance
type Instance func() Config

var adapters = make(map[string]Instance)

// Register makes a Config adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("config: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("config: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewConfig Create a new config driver by adapter name and config string.
func NewConfig(adapterName string, config Param) (adapter Config, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("config: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
