package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/dbunion/com/log"
	"reflect"
	"time"
)

const (
	// TypeAsync - type async
	TypeAsync = "async"
)

// Config - task config
type Config struct {
	// async task
	BrokerType      string `json:"broker_type"`
	Broker          string `json:"broker"`
	DefaultQueue    string `json:"default_queue"`
	BrokerConfig    string `json:"broker_config"`
	BackendType     string `json:"backend_type"`
	BackendConfig   string `json:"backend_config"`
	ResultBackend   string `json:"result_backend"`
	ResultsExpireIn int    `json:"results_expire_in"`

	// worker
	Concurrency int `json:"concurrency"`

	FuncWraps map[string]FuncWrap `json:"func_wraps"`
	Logger    log.Logger          `json:"logger"`

	ErrorHandler    func(err error)
	PreTaskHandler  func(param *Param)
	PostTaskHandler func(param *Param)

	// Extend fields
	// Extended fields can be used if there is a special implementation
	Extend1 string `json:"extend_1"`
	Extend2 string `json:"extend_2"`
}

// ErrNotImpl - error for not impl
var ErrNotImpl = errors.New("method not impl")

// Option task options
type Option struct {
	ETA          *time.Time `json:"eta"`
	Priority     uint8      `json:"priority"`
	Immutable    bool       `json:"immutable"`
	RetryCount   int        `json:"retry_count"`
	RetryTimeout int        `json:"retry_timeout"`
}

// Arg represents a single argument passed to invocation fo a task
type Arg struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// Param task param
type Param struct {
	UUID        string        `json:"uuid"`
	Name        string        `json:"name"`
	Fun         string        `json:"fun"`
	Option      Option        `json:"option"`
	Args        []Arg         `json:"args"`
	WaitTimeOut time.Duration `json:"wait_time_out"`
}

// CallbackFunc - task call back function
type CallbackFunc func(param *Param, err error, message string) error

// FuncWrap - task worker func warp
type FuncWrap interface {
	// get all task func
	GetTasks() map[string]interface{}
	// stop task by uuid
	StopTask(uuid string) error
}

// ParamFromContext - convert context value to Param
type ParamFromContext func(ctx context.Context) *Param

// ParamFunc - point ParamFromContext impl
var ParamFunc ParamFromContext

// Task interface contains all behaviors for Task adapter.
type Task interface {
	// add new task
	AddTask(param *Param, onSuccess []*Param, onError []*Param, callbacks ...CallbackFunc) error

	// run all task
	// if chain is true, The tasks will be executed in turn, and the return value of the previous
	// task will be used as the parameter of the next task
	Run(chain bool) error

	// stop all task
	Stop() error

	// start gc routine based on config settings.
	StartAndGC(config Config) error
}

// Result - task result value
type Result interface {
	// Get - get result with sleep
	Get(sleepDuration time.Duration) ([]reflect.Value, error)

	// GetWithTimeout - get with timeout
	GetWithTimeout(timeoutDuration, sleepDuration time.Duration) ([]reflect.Value, error)
}

// Instance is a function create a new Task Instance
type Instance func() Task

var adapters = make(map[string]Instance)

// Register makes a Task adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("Task: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("Task: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewTask Create a new Task driver by adapter name and config string.
// config need to be correct JSON as string:
// {"server": "localhost:9092", "user": "xxxx", "password":"xxxxx"}.
// it will start gc automatically.
func NewTask(adapterName string, config Config) (adapter Task, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("task: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
