package uid

import (
	"fmt"
)

// Config - uid config
type Config struct {
	// Redis config
	Key string `json:"key"`

	// public param
	Server   string `json:"server"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int64  `json:"port"`

	// Snow flake node id
	NodeID int64 `json:"node_id"`

	// Mysql config
	InitValue       int64  `json:"init_value"`
	Step            int64  `json:"step"`
	DBName          string `json:"db_name"`
	TableName       string `json:"table_name"`
	AutoCreateTable bool   `json:"auto_create_table"`
}

// UID interface contains all behaviors for UID adapter.
// usage:
//	uid.Register("uid", uid.NewUID)
type UID interface {
	// has int32 uid
	HasInt32() bool

	// next int32 uid
	NextUID32() int32

	// next int64 uid
	NextUID64() int64

	// close connection
	Close() error

	// start gc routine based on config settings.
	StartAndGC(config Config) error
}

// Instance is a function create a new UID Instance
type Instance func() UID

var adapters = make(map[string]Instance)

// Register makes a UID adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("UID: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("UID: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewUID Create a new UID driver by adapter name and config string.
// config need to be correct JSON as string:
// {"server": "localhost:9092", "user": "xxxx", "password":"xxxxx"}.
// it will start gc automatically.
func NewUID(adapterName string, config Config) (adapter UID, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("UID: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
