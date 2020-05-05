package vtctl

import (
	"context"
	"fmt"
	"time"

	"github.com/zssky/tc/retry"
)

const (
	// TypeVtctlV3 - vtctl v3
	TypeVtctlV3 = "vtctl_v3"
)

// ServerInfo - vtctld server info
type ServerInfo struct {
	IP     string `json:"name"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}

// Config - client config
type Config struct {
	RetryOption *retry.Options        `json:"option"`
	Servers     map[string]ServerInfo `json:"servers"`
	HealthCheck time.Duration         `json:"health_check"`
}

// Client interface contains all behaviors for vtctl client.
type Client interface {
	// run vtctl command
	RunCommand(ctx context.Context, args []string, timeout time.Duration) ([]string, error)

	// close connection
	Close() error

	// start gc routine based on config string settings.
	StartAndGC(config Config) error
}

// Instance is a function create a new Client Instance
type Instance func() Client

var adapters = make(map[string]Instance)

// Register makes a Client adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("Client: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("Client: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewClient Create a new Client driver by adapter name and config string.
// it will start gc automatically.
func NewClient(adapterName string, config Config) (adapter Client, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("newClient: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
