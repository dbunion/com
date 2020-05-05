package scheduler

import (
	"fmt"
	"github.com/dbunion/com/log"
)

const (
	// TypeK8s - type k8s
	TypeK8s = "k8s"
	// TypeK3s = type k3s
	TypeK3s = "k3s"
	// TypeNomad - type nomad
	TypeNomad = "nomad"
	// TypeDockerCompose - type docker compose
	TypeDockerCompose = "docker_compose"
)

// Param - scheduler config
type Param struct {
	// k8s setting
	Server   string `json:"server"`
	Token    string `json:"token"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Insecure bool   `json:"insecure"`

	// nomad setting

	// public param
	Logger log.Logger `json:"logger"`

	// Extend fields
	// Extended fields can be used if there is a special implementation
	Extend1 string `json:"extend_1"`
	Extend2 string `json:"extend_2"`
}

// Scheduler interface contains all behaviors for scheduler adapter.
type Scheduler interface {
	// GetNodeOperator - get node Operator
	GetNodeOperator() NodeOperator

	// GetNamespaceOperator - get namespace Operator
	GetNamespaceOperator() NamespaceOperator

	// GetConfigOperator - get config Operator
	GetConfigOperator() ConfigOperator

	// GetServiceOperator - get service Operator
	GetServiceOperator() ServiceOperator

	// GetPodOperator - get pod operator
	GetPodOperator() PodOperator

	// GetRCOperator - get rc Operator
	GetRCOperator() RCOperator

	// GetSTSOperator - get sts Operator
	GetSTSOperator() STSOperator

	// GetDaemonSetOperator - get DaemonSet Operator
	GetDaemonSetOperator() DaemonSetOperator

	// GetDeploymentOperator - get Deployment Operator
	GetDeploymentOperator() DeploymentOperator

	// GetReplicaSetOperator - get replicaset Operator
	GetReplicaSetOperator() ReplicaSetOperator

	// close connection
	Close() error

	// start gc routine based on config string settings.
	StartAndGC(config Param) error
}

// Instance is a function create a new Scheduler Instance
type Instance func() Scheduler

var adapters = make(map[string]Instance)

// Register makes a Scheduler adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("scheduler: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("scheduler: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewScheduler Create a new Scheduler driver by adapter name and config setting.
func NewScheduler(adapterName string, config Param) (Scheduler, error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("newScheduler: unknown adapter name %q (forgot to import?)", adapterName)
	}
	adapter := instanceFunc()
	if err := adapter.StartAndGC(config); err != nil {
		return nil, err
	}
	return adapter, nil
}
