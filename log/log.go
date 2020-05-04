package log

import (
	"fmt"
)

// Level - log level
type Level string

const (
	// LevelInfo - info log
	LevelInfo Level = "info"
	// LevelDebug - debug log
	LevelDebug Level = "debug"
	// LevelWarning - warning log
	LevelWarning Level = "warn"
	// LevelError - error log
	LevelError Level = "error"
	// LevelFatal - fatal log
	LevelFatal Level = "fatal"
)

const (
	// TypeZsskyLog - use zssky log
	TypeZsskyLog = "zssky"
	// TypeLogrus - use logrus
	TypeLogrus = "logrus"
)

// Config - log base config
type Config struct {
	Level         Level  `json:"level"`
	FilePath      string `json:"file_path"`
	HighLighting  bool   `json:"highlighting"`
	RotateByDay   bool   `json:"rotate_by_day"`
	RotateByHour  bool   `json:"rotate_by_hour"`
	JSONFormatter bool   `json:"json_formatter"`

	// Extend fields
	// Extended fields can be used if there is a special implementation
	Extend1 string `json:"extend_1"`
	Extend2 string `json:"extend_2"`
}

// Logger interface contains all behaviors for log adapter.
type Logger interface {
	// Infof - info format log
	Infof(format string, v ...interface{})

	// Info - info format log
	Info(v ...interface{})

	// Debugf - debug format log
	Debugf(format string, v ...interface{})

	// Debug - debug format log
	Debug(v ...interface{})

	// Warnf - warn format log
	Warnf(format string, v ...interface{})

	// Warn - warn log
	Warn(v ...interface{})

	// Warningf - Warning format log
	Warningf(format string, v ...interface{})

	// Warning - Warning log
	Warning(v ...interface{})

	// Errorf - error format log
	Errorf(format string, v ...interface{})

	// Error - error format log
	Error(v ...interface{})

	// Fatalf - fatal format log
	Fatalf(format string, v ...interface{})

	// Fatal - fatal format log
	Fatal(v ...interface{})

	// Fatalln - fatal log
	Fatalln(v ...interface{})

	// Printf - printf format value
	Printf(format string, v ...interface{})

	// Print - printf value
	Print(v ...interface{})

	// Println - print value
	Println(v ...interface{})

	// Panic - panic
	Panic(v ...interface{})

	// Panicf - panic format value
	Panicf(format string, v ...interface{})

	// Panicln - panic
	Panicln(v ...interface{})

	// close connection
	Close() error

	// start gc routine based on config settings.
	StartAndGC(config Config) error
}

// Instance is a function create a new Logger Instance
type Instance func() Logger

var adapters = make(map[string]Instance)

// Register makes a Logger adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("Logger: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("Logger: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewLogger Create a new Logger driver by adapter name and config string.
// it will start gc automatically.
func NewLogger(adapterName string, config Config) (adapter Logger, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("NewLogger: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
