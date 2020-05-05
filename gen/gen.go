package gen

import "fmt"

const (
	// TypeGormModel - type gorm model gen
	TypeGormModel = "gorm_model"
	// TypeBeegoModel - type beego model gen
	TypeBeegoModel = "beego_model"
)

// Item - gen item
type Item struct {
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

// Config - init config
type Config struct {
	Package      string `json:"package"`
	GenPath      string `json:"gen_path"`
	AllInOne     bool   `json:"all_in_one"`
	Items        []Item `json:"items"`
	MaxIdleConns int64  `json:"max_idle_conns"`
	MaxOpenConns int64  `json:"set_max_open_conns"`
}

// Generator interface contains all behaviors for model generator adapter.
type Generator interface {
	// Gen - gen file
	Gen() error

	// start gc routine based on config settings.
	StartAndGC(config Config) error
}

// Instance is a function create a new Generator Instance
type Instance func() Generator

var adapters = make(map[string]Instance)

// Register makes a Generator adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("Generator: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("Generator: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewGenerator Create a new Generator driver by adapter name and config setting.
// it will start gc automatically.
func NewGenerator(adapterName string, config Config) (adapter Generator, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("Generator: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
