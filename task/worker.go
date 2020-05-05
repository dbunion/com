package task

import "fmt"

const (
	// TypeAsyncWorker - type async worker
	TypeAsyncWorker = "async_worker"
)

// Worker interface contains all behaviors for task worker
type Worker interface {
	// run worker
	Run() error

	// close worker
	Close() error

	// start gc routine based on config string settings.
	StartAndGC(config Config) error
}

// WorkerInstance is a function create a new Task worker Instance
type WorkerInstance func() Worker

var workerAdapters = make(map[string]WorkerInstance)

// RegisterWorker makes a Task worker adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
func RegisterWorker(name string, adapter WorkerInstance) {
	if adapter == nil {
		panic("Task: Register worker adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("Task: Register called twice for adapter " + name)
	}
	workerAdapters[name] = adapter
}

// NewWorker Create a new Task worker driver by adapter name and config string.
// config need to be correct JSON as string:
// it will start gc automatically.
func NewWorker(adapterName string, config Config) (adapter Worker, err error) {
	instanceFunc, ok := workerAdapters[adapterName]
	if !ok {
		err = fmt.Errorf("task: unknown worker adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
