package file

import (
	"errors"
	"fmt"
	"github.com/dbunion/com/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

// Config is File config adapter.
type Config struct {
	fileType string
	path     string
	name     string
	file     string
	viper    *viper.Viper
}

// NewFileConfig create new file config with default collection name.
func NewFileConfig() config.Config {
	return &Config{}
}

// Get get config value by key.
func (fc *Config) Get(key string) interface{} {
	return fc.viper.Get(key)
}

// GetBool get bool config value by key
func (fc *Config) GetBool(key string) bool {
	return fc.viper.GetBool(key)
}

// GetFloat64 get float64 config value by key
func (fc *Config) GetFloat64(key string) float64 {
	return fc.viper.GetFloat64(key)
}

// GetInt get Int config value by key
func (fc *Config) GetInt(key string) int {
	return fc.viper.GetInt(key)
}

// GetInt32 get Int32 config value by key
func (fc *Config) GetInt32(key string) int32 {
	return fc.viper.GetInt32(key)
}

// GetInt64 get Int64 config value by key
func (fc *Config) GetInt64(key string) int64 {
	return fc.viper.GetInt64(key)
}

// GetString get string config value by key
func (fc *Config) GetString(key string) string {
	return fc.viper.GetString(key)
}

// GetStringMap get stringMap config value by key
func (fc *Config) GetStringMap(key string) map[string]interface{} {
	return fc.viper.GetStringMap(key)
}

// GetStringMapString get stringMapstring config value by key
func (fc *Config) GetStringMapString(key string) map[string]string {
	return fc.viper.GetStringMapString(key)
}

// GetStringSlice get stringSlice config value by key
func (fc *Config) GetStringSlice(key string) []string {
	return fc.viper.GetStringSlice(key)
}

// GetTime get time config value by key
func (fc *Config) GetTime(key string) time.Time {
	return fc.viper.GetTime(key)
}

// GetDuration get Duration config value by key
func (fc *Config) GetDuration(key string) time.Duration {
	return fc.viper.GetDuration(key)
}

// AllSettings get all config value
func (fc *Config) AllSettings() map[string]interface{} {
	return fc.viper.AllSettings()
}

// IsExist check config's existence in file.
func (fc *Config) IsExist(key string) bool {
	return fc.viper.IsSet(key)
}

// StartAndGC start file config adapter.
// config is like {"path":"/etc/niffer", "name":"xxxxx.json","config_type":"json"}
// the cache item in redis are stored forever,
// so no gc operation.
func (fc *Config) StartAndGC(cfg config.Param) error {
	if cfg.Path != "" {
		fc.path = cfg.Path
	}

	if cfg.Name != "" {
		fc.name = cfg.Name
	}

	if cfg.Type == "" {
		return errors.New("config type can not be empty, type:JSON/TOML/YAML/HCL/Java properties")
	}
	fc.fileType = cfg.Type

	vp := viper.New()

	if cfg.File != "" {
		fc.file = cfg.File
		vp.SetConfigFile(fc.file)
	}

	if fc.name != "" {
		vp.SetConfigName(fc.name)
	}

	if fc.path != "" {
		vp.AddConfigPath(fc.path)
	}

	if fc.fileType != "" {
		vp.SetConfigType(fc.fileType)
	}

	err := vp.ReadInConfig()
	if err != nil {
		panic(err)
	}

	fc.viper = vp

	// watch file change
	vp.WatchConfig()
	vp.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("config change, name:%v op:%s\n", in.Name, in.Op)
	})

	return nil
}

func init() {
	config.Register(config.TypeFile, NewFileConfig)
}
