package gen

import "fmt"

const (
	// TypeGormModel - type gorm model gen
	TypeGormModel = "gorm_model"
	// TypeBeegoModel - type beego model gen
	TypeBeegoModel = "beego_model"
	// TypeService - type service gen
	TypeService = "service"
)

// Primary - primary key info
type Primary struct {
	Type    string
	Name    string
	Default interface{}
}

// IsEmpty - check primary is empty
func (p *Primary) IsEmpty() bool {
	return p.Type == "" && p.Name == "" && p.Default == nil
}

// Item - gen item
type Item struct {
	Name        string   `json:"name"`
	TableName   string   `json:"table_name"`
	Relations   []string `json:"relations"`
	Detail      string   `json:"detail"`
	Primary     Primary  `json:"primary"`
	DisableMode bool     `json:"disable_mode"`
}

// ModelGenConfig - model config
type ModelGenConfig struct {
	Items        []Item `json:"items"`
	MaxIdleConns int64  `json:"max_idle_conns"`
	MaxOpenConns int64  `json:"set_max_open_conns"`
}

// SItem - service config item
type SItem struct {
	Req        interface{} `json:"req"`
	Dst        interface{} `json:"dst"`
	Index      int64       `json:"index"`
	CheckApp   bool        `json:"check_app"`
	ViewCfg    ViewConfig  `json:"view_cfg"`
	DisableRPC bool        `json:"disable_rpc"`
}

// ServiceGenConfig - service config
type ServiceGenConfig struct {
	Items       []SItem         `json:"items"`
	ImportPaths []string        `json:"import_paths"`
	IgnoreName  map[string]bool `json:"ignore_name"`
	KeepName    map[string]bool `json:"keep_name"`
}

// ViewConfig - view gen config
type ViewConfig struct {
	Enable     bool              `json:"enable"`
	Mapping    map[string]string `json:"mapping"`
	RelMapping map[string]string `json:"relation_mapping"`
}

// Config - init config
type Config struct {
	Package    string           `json:"package"`
	GenPath    string           `json:"gen_path"`
	AllInOne   bool             `json:"all_in_one"`
	ModelCfg   ModelGenConfig   `json:"model_cfg"`
	ServiceCfg ServiceGenConfig `json:"service_cfg"`
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
