package scheduler

import (
	"context"
)

// Object - interface for base resource
type Object interface {
	GetName() string
}

// Options - resource options
type Options map[string]interface{}

// NodeStatus is information about the current status of a node.
type NodeStatus struct {
	Capacity    map[string]int64 `json:"capacity,omitempty" protobuf:"bytes,1,rep,name=capacity,casttype=ResourceList,castkey=ResourceName"`
	Allocatable map[string]int64 `json:"allocatable,omitempty" protobuf:"bytes,2,rep,name=allocatable,casttype=ResourceList,castkey=ResourceName"`
	Phase       string           `json:"phase,omitempty" protobuf:"bytes,3,opt,name=phase,casttype=NodePhase"`
}

// PodResource - node pods resource info
type PodResource struct {
	Resource
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// Resource - node Allocated resource
type Resource struct {
	CPURequest    int64 `json:"cpu_request"`
	CPULimit      int64 `json:"cpu_limit"`
	MemoryRequest int64 `json:"memory_request"`
	MemoryLimit   int64 `json:"memory_limit"`
}

// NodeDetail is node detail information
type NodeDetail struct {
	Capacity    map[string]int64 `json:"capacity,omitempty" protobuf:"bytes,1,rep,name=capacity,casttype=ResourceList,castkey=ResourceName"`
	Allocatable map[string]int64 `json:"allocatable,omitempty" protobuf:"bytes,2,rep,name=allocatable,casttype=ResourceList,castkey=ResourceName"`
	Pods        []PodResource    `json:"pods"`
	Resource    Resource         `json:"resource"`
}

// Node - cluster physical node
type Node struct {
	Name   string            `json:"name,omitempty" yaml:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Status NodeStatus        `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
	YAML   []byte            `json:"-"`
}

// GetName - object impl
func (n *Node) GetName() string {
	return n.Name
}

// NodeOperator - node Operator interface
type NodeOperator interface {
	Get(ctx context.Context, param *Node) (*Node, error)
	List(ctx context.Context, options Options) ([]*Node, error)
	Create(ctx context.Context, param *Node, options Options) error
	CreateWithYaml(ctx context.Context, param *Node, options Options) error
	Update(ctx context.Context, param *Node) error
	Delete(ctx context.Context, param *Node, options Options) error
	Watch(ctx context.Context, param *Node, options Options) (Interface, error)
	Describe(ctx context.Context, param *Node) (*NodeDetail, error)
}

// NamespaceStatus is information about the current status of a Namespace.
type NamespaceStatus struct {
	// Phase is the current lifecycle phase of the namespace.
	// More info: https://kubernetes.io/docs/tasks/administer-cluster/namespaces/
	// +optional
	Phase string `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=NamespacePhase"`
}

// Namespace - resource isolation unit
type Namespace struct {
	Name   string            `json:"name,omitempty" yaml:"name,omitempty"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Status NamespaceStatus   `json:"status"`
	YAML   []byte            `json:"-"`
}

// GetName - object impl
func (n *Namespace) GetName() string {
	return n.Name
}

// NamespaceOperator - namespace Operator define
type NamespaceOperator interface {
	Get(ctx context.Context, param *Namespace) (*Namespace, error)
	List(ctx context.Context, options Options) ([]*Namespace, error)
	Create(ctx context.Context, param *Namespace, options Options) error
	CreateWithYaml(ctx context.Context, param *Namespace, options Options) error
	Update(ctx context.Context, param *Namespace) error
	Delete(ctx context.Context, param *Namespace, options Options) error
	Watch(ctx context.Context, param *Namespace, options Options) (Interface, error)
}

// Config - common config file define
type Config struct {
	Name       string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace  string                 `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	BinaryData map[string][]byte      `json:"binaryData,omitempty" yaml:"binaryData,omitempty"`
	Data       map[string]string      `json:"data,omitempty" yaml:"data,omitempty"`
	Labels     map[string]string      `json:"labels,omitempty" yaml:"labels,omitempty"`
	Reserved   map[string]interface{} `json:"reserved,omitempty" yaml:"reserved,omitempty"`
	YAML       []byte                 `json:"-"`
}

// GetName - object impl
func (c *Config) GetName() string {
	return c.Name
}

// ConfigOperator - config Operator interface
type ConfigOperator interface {
	Get(ctx context.Context, namespace string, param *Config) (*Config, error)
	List(ctx context.Context, namespace string, options Options) ([]*Config, error)
	Create(ctx context.Context, param *Config, options Options) error
	CreateWithYaml(ctx context.Context, param *Config, options Options) error
	Update(ctx context.Context, param *Config) error
	Delete(ctx context.Context, param *Config, options Options) error
	Watch(ctx context.Context, param *Config, options Options) (Interface, error)
}

// ServicePort contains information on service's port.
type ServicePort struct {
	Name       string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	Protocol   string `json:"protocol,omitempty" protobuf:"bytes,2,opt,name=protocol,casttype=Protocol"`
	Port       int32  `json:"port" protobuf:"varint,3,opt,name=port"`
	TargetPort int32  `json:"targetPort,omitempty" protobuf:"bytes,4,opt,name=targetPort"`
}

// ServiceSpec - ServiceSpec describes the attributes that a user creates on a service.
type ServiceSpec struct {
	Ports           []ServicePort     `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"port" protobuf:"bytes,1,rep,name=ports"`
	Selector        map[string]string `json:"selector,omitempty" protobuf:"bytes,2,rep,name=selector"`
	ClusterIP       string            `json:"clusterIP,omitempty" protobuf:"bytes,3,opt,name=clusterIP"`
	Type            string            `json:"type,omitempty" protobuf:"bytes,4,opt,name=type,casttype=ServiceType"`
	ExternalIPs     []string          `json:"externalIPs,omitempty" protobuf:"bytes,5,rep,name=externalIPs"`
	SessionAffinity string            `json:"sessionAffinity,omitempty" protobuf:"bytes,7,opt,name=sessionAffinity,casttype=ServiceAffinity"`
	LoadBalancerIP  string            `json:"loadBalancerIP,omitempty" protobuf:"bytes,8,opt,name=loadBalancerIP"`
}

// Service - service struct define
type Service struct {
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Spec      ServiceSpec       `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	YAML      []byte            `json:"-"`
}

// GetName - object impl
func (s *Service) GetName() string {
	return s.Name
}

// ServiceOperator - service base Operator method
type ServiceOperator interface {
	Get(ctx context.Context, namespace string, param *Service) (*Service, error)
	List(ctx context.Context, namespace string, options Options) ([]*Service, error)
	Create(ctx context.Context, param *Service, options Options) error
	CreateWithYaml(ctx context.Context, param *Service, options Options) error
	Update(ctx context.Context, param *Service) error
	Delete(ctx context.Context, param *Service, options Options) error
	Watch(ctx context.Context, param *Service, options Options) (Interface, error)
}

// PodTemplateSpec describes the data a pod should have when created from a template
type PodTemplateSpec struct {
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Spec      PodSpec           `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// Pod - cluster pod
type Pod struct {
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Spec      PodSpec           `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status    PodStatus         `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
	YAML      []byte            `json:"-" yaml:"-"`
}

// GetName - object impl
func (p *Pod) GetName() string {
	return p.Name
}

// Volume represents a named volume in a pod that may be accessed by any container in the pod.
type Volume struct {
	Name  string                 `json:"name" protobuf:"bytes,1,opt,name=name"`
	Value map[string]interface{} `json:",inline" protobuf:"bytes,2,opt,name=value"`
}

// ContainerPort represents a network port in a single container.
type ContainerPort struct {
	Name          string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	HostPort      int32  `json:"hostPort,omitempty" protobuf:"varint,2,opt,name=hostPort"`
	ContainerPort int32  `json:"containerPort" protobuf:"varint,3,opt,name=containerPort"`
	Protocol      string `json:"protocol,omitempty" protobuf:"bytes,4,opt,name=protocol,casttype=Protocol"`
	HostIP        string `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
}

// ResourceList is a set of (resource name, quantity) pairs.
type ResourceList map[string]int64

// ResourceRequirements describes the compute resource requirements.
type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList,castkey=ResourceName"`
	Requests ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList,castkey=ResourceName"`
}

// VolumeMount describes a mounting of a Volume within a container.
type VolumeMount struct {
	Name        string `json:"name" protobuf:"bytes,1,opt,name=name"`
	ReadOnly    bool   `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`
	MountPath   string `json:"mountPath" protobuf:"bytes,3,opt,name=mountPath"`
	SubPath     string `json:"subPath,omitempty" protobuf:"bytes,4,opt,name=subPath"`
	SubPathExpr string `json:"subPathExpr,omitempty" protobuf:"bytes,6,opt,name=subPathExpr"`
}

// Container A single application container that you want to run within a pod.
type Container struct {
	Name         string               `json:"name" protobuf:"bytes,1,opt,name=name"`
	Image        string               `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	Command      []string             `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	Args         []string             `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
	WorkingDir   string               `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
	Ports        []ContainerPort      `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,6,rep,name=ports"`
	Resources    ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
	VolumeMounts []VolumeMount        `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,9,rep,name=volumeMounts"`
}

// PodSpec is a description of a pod.
type PodSpec struct {
	Volumes                       []Volume          `json:"volumes,omitempty" patchStrategy:"merge,retainKeys" patchMergeKey:"name" protobuf:"bytes,1,rep,name=volumes"`
	InitContainers                []Container       `json:"initContainers,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,20,rep,name=initContainers"`
	Containers                    []Container       `json:"containers" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=containers"`
	RestartPolicy                 string            `json:"restartPolicy,omitempty" protobuf:"bytes,3,opt,name=restartPolicy,casttype=RestartPolicy"`
	TerminationGracePeriodSeconds *int64            `json:"terminationGracePeriodSeconds,omitempty" protobuf:"varint,4,opt,name=terminationGracePeriodSeconds"`
	ActiveDeadlineSeconds         *int64            `json:"activeDeadlineSeconds,omitempty" protobuf:"varint,5,opt,name=activeDeadlineSeconds"`
	NodeSelector                  map[string]string `json:"nodeSelector,omitempty" protobuf:"bytes,7,rep,name=nodeSelector"`
	NodeName                      string            `json:"nodeName,omitempty" protobuf:"bytes,10,opt,name=nodeName"`
	HostNetwork                   bool              `json:"hostNetwork,omitempty" protobuf:"varint,11,opt,name=hostNetwork"`
	Hostname                      string            `json:"hostname,omitempty" protobuf:"bytes,16,opt,name=hostname"`
	PriorityClassName             string            `json:"priorityClassName,omitempty" protobuf:"bytes,24,opt,name=priorityClassName"`
	Priority                      *int32            `json:"priority,omitempty" protobuf:"bytes,25,opt,name=priority"`
}

// PodCondition contains details for the current condition of this pod.
type PodCondition struct {
	Type               string `json:"type" protobuf:"bytes,1,opt,name=type,casttype=PodConditionType"`
	Status             string `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	LastProbeTime      string `json:"lastProbeTime,omitempty" protobuf:"bytes,3,opt,name=lastProbeTime"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	Reason             string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	Message            string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// PodStatus represents information about the status of a pod. Status may trail the actual
type PodStatus struct {
	Phase      string         `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=PodPhase"`
	Conditions []PodCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
	Message    string         `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
	Reason     string         `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	HostIP     string         `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
	PodIP      string         `json:"podIP,omitempty" protobuf:"bytes,6,opt,name=podIP"`
	StartTime  string         `json:"startTime,omitempty" protobuf:"bytes,7,opt,name=startTime"`
}

// PodOperator - node Operator interface
type PodOperator interface {
	Get(ctx context.Context, namespace string, param *Pod) (*Pod, error)
	List(ctx context.Context, namespace string, options Options) ([]*Pod, error)
	Create(ctx context.Context, param *Pod, options Options) error
	CreateWithYaml(ctx context.Context, param *Pod, options Options) error
	Update(ctx context.Context, param *Pod) error
	Delete(ctx context.Context, param *Pod, options Options) error
	GetEvents(ctx context.Context, param *Pod) ([]*Event, error)
	GetLogs(ctx context.Context, namespace, name, container string) ([]byte, error)
	Watch(ctx context.Context, param *Pod, options Options) (Interface, error)
}

// Event is a report of an event somewhere in the cluster.
type Event struct {
	Reason              string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	Message             string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	Source              string `json:"source,omitempty" protobuf:"bytes,5,opt,name=source"`
	FirstTimestamp      string `json:"firstTimestamp,omitempty" protobuf:"bytes,6,opt,name=firstTimestamp"`
	LastTimestamp       string `json:"lastTimestamp,omitempty" protobuf:"bytes,7,opt,name=lastTimestamp"`
	Count               int32  `json:"count,omitempty" protobuf:"varint,8,opt,name=count"`
	Type                string `json:"type,omitempty" protobuf:"bytes,9,opt,name=type"`
	EventTime           string `json:"eventTime,omitempty" protobuf:"bytes,10,opt,name=eventTime"`
	Action              string `json:"action,omitempty" protobuf:"bytes,12,opt,name=action"`
	ReportingController string `json:"reportingComponent" protobuf:"bytes,14,opt,name=reportingComponent"`
	ReportingInstance   string `json:"reportingInstance" protobuf:"bytes,15,opt,name=reportingInstance"`
}

// RCSpec is the specification of a replication controller.
// As the internal representation of a replication controller, it may have either
// a TemplateRef or a Template set.
type RCSpec struct {
	Replicas        int32
	MinReadySeconds int32
	Selector        map[string]string
	Template        PodTemplateSpec
}

// RCStatus represents the current status of a replication
// controller.
type RCStatus struct {
	Replicas             int32
	FullyLabeledReplicas int32
	ReadyReplicas        int32
	AvailableReplicas    int32
	ObservedGeneration   int64
}

// RC - cluster RC
type RC struct {
	Version   string            `json:"version"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Spec      RCSpec
	Status    RCStatus
	YAML      []byte `json:"-"`
}

// GetName - object impl
func (r *RC) GetName() string {
	return r.Name
}

// RCOperator - rc Operator interface
type RCOperator interface {
	Get(ctx context.Context, namespace string, param *RC) (*RC, error)
	List(ctx context.Context, namespace string, options Options) ([]*RC, error)
	Create(ctx context.Context, param *RC, options Options) error
	CreateWithYaml(ctx context.Context, param *RC, options Options) error
	Update(ctx context.Context, param *RC) error
	Delete(ctx context.Context, param *RC, options Options) error
	Watch(ctx context.Context, param *RC, options Options) (Interface, error)
}

// A STSSpec is the specification of a StatefulSet.
type STSSpec struct {
	Replicas    int32             `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
	Selector    map[string]string `json:"selector" protobuf:"bytes,2,opt,name=selector"`
	Template    PodTemplateSpec   `json:"template" protobuf:"bytes,3,opt,name=template"`
	ServiceName string            `json:"serviceName" protobuf:"bytes,5,opt,name=serviceName"`
}

// STSStatus represents the current state of a StatefulSet.
type STSStatus struct {
	Replicas        int32  `json:"replicas" protobuf:"varint,2,opt,name=replicas"`
	ReadyReplicas   int32  `json:"readyReplicas,omitempty" protobuf:"varint,3,opt,name=readyReplicas"`
	CurrentReplicas int32  `json:"currentReplicas,omitempty" protobuf:"varint,4,opt,name=currentReplicas"`
	UpdatedReplicas int32  `json:"updatedReplicas,omitempty" protobuf:"varint,5,opt,name=updatedReplicas"`
	CurrentRevision string `json:"currentRevision,omitempty" protobuf:"bytes,6,opt,name=currentRevision"`
	UpdateRevision  string `json:"updateRevision,omitempty" protobuf:"bytes,7,opt,name=updateRevision"`
}

// STS - cluster Statefulset
type STS struct {
	Version   string            `json:"version"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Spec      STSSpec           `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status    STSStatus         `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
	YAML      []byte            `json:"-"`
}

// GetName - object impl
func (s *STS) GetName() string {
	return s.Name
}

// STSOperator - sts Operator interface
type STSOperator interface {
	Get(ctx context.Context, namespace string, param *STS) (*STS, error)
	List(ctx context.Context, namespace string, options Options) ([]*STS, error)
	Create(ctx context.Context, param *STS, options Options) error
	CreateWithYaml(ctx context.Context, param *STS, options Options) error
	Update(ctx context.Context, param *STS) error
	Delete(ctx context.Context, param *STS, options Options) error
	Watch(ctx context.Context, param *STS, options Options) (Interface, error)
}

// DaemonSetSpec is the specification of a daemon set.
type DaemonSetSpec struct {
	Selector        map[string]string `json:"selector" protobuf:"bytes,1,opt,name=selector"`
	Template        PodTemplateSpec   `json:"template" protobuf:"bytes,2,opt,name=template"`
	MinReadySeconds int32             `json:"minReadySeconds,omitempty" protobuf:"varint,4,opt,name=minReadySeconds"`
}

// DaemonSetStatus represents the current status of a daemon set.
type DaemonSetStatus struct {
	CurrentNumberScheduled int32 `json:"currentNumberScheduled" protobuf:"varint,1,opt,name=currentNumberScheduled"`
	NumberMisscheduled     int32 `json:"numberMisscheduled" protobuf:"varint,2,opt,name=numberMisscheduled"`
	DesiredNumberScheduled int32 `json:"desiredNumberScheduled" protobuf:"varint,3,opt,name=desiredNumberScheduled"`
	NumberReady            int32 `json:"numberReady" protobuf:"varint,4,opt,name=numberReady"`
	NumberAvailable        int32 `json:"numberAvailable,omitempty" protobuf:"varint,7,opt,name=numberAvailable"`
	NumberUnavailable      int32 `json:"numberUnavailable,omitempty" protobuf:"varint,8,opt,name=numberUnavailable"`
}

// DaemonSet - cluster DaemonSet
type DaemonSet struct {
	Version   string            `json:"version"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Spec      DaemonSetSpec     `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status    DaemonSetStatus   `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
	YAML      []byte            `json:"-"`
}

// GetName - object impl
func (d *DaemonSet) GetName() string {
	return d.Name
}

// DaemonSetOperator - DaemonSet Operator interface
type DaemonSetOperator interface {
	Get(ctx context.Context, namespace string, param *DaemonSet) (*DaemonSet, error)
	List(ctx context.Context, namespace string, options Options) ([]*DaemonSet, error)
	Create(ctx context.Context, param *DaemonSet, options Options) error
	CreateWithYaml(ctx context.Context, param *DaemonSet, options Options) error
	Update(ctx context.Context, param *DaemonSet) error
	Delete(ctx context.Context, param *DaemonSet, options Options) error
	Watch(ctx context.Context, param *DaemonSet, options Options) (Interface, error)
}

// DeploymentSpec specifies the state of a Deployment.
type DeploymentSpec struct {
	Replicas                int32             `json:"replicas"`
	Selector                map[string]string `json:"selector"`
	Template                PodTemplateSpec   `json:"template"`
	Strategy                string            `json:"strategy"`
	MinReadySeconds         int32             `json:"min_ready_seconds"`
	Paused                  bool              `json:"paused"`
	ProgressDeadlineSeconds *int32            `json:"progress_deadline_seconds"`
}

// DeploymentStatus holds information about the observed status of a deployment.
type DeploymentStatus struct {
	Replicas            int32 `json:"replicas"`
	UpdatedReplicas     int32 `json:"updated_replicas"`
	ReadyReplicas       int32 `json:"ready_replicas"`
	AvailableReplicas   int32 `json:"available_replicas"`
	UnavailableReplicas int32 `json:"unavailable_replicas"`
}

// Deployment - cluster Deployment
type Deployment struct {
	Version   string            `json:"version"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Spec      DeploymentSpec    `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status    DeploymentStatus  `json:"status,omitempty" yaml:"status,omitempty"`
	YAML      []byte            `json:"-" yaml:"-"`
}

// GetName - object impl
func (d *Deployment) GetName() string {
	return d.Name
}

// DeploymentOperator - Deployment Operator interface
type DeploymentOperator interface {
	Get(ctx context.Context, namespace string, param *Deployment) (*Deployment, error)
	List(ctx context.Context, namespace string, options Options) ([]*Deployment, error)
	Create(ctx context.Context, param *Deployment, options Options) error
	CreateWithYaml(ctx context.Context, param *Deployment, options Options) error
	Update(ctx context.Context, param *Deployment) error
	Delete(ctx context.Context, param *Deployment, options Options) error
	Watch(ctx context.Context, param *Deployment, options Options) (Interface, error)
}

// ReplicaSetSpec is the specification of a ReplicaSet.
// As the internal representation of a ReplicaSet, it must have
// a Template set.
type ReplicaSetSpec struct {
	Replicas        int32
	MinReadySeconds int32
	Selector        map[string]string
	Template        PodTemplateSpec
}

// ReplicaSetStatus represents the current status of a ReplicaSet.
type ReplicaSetStatus struct {
	Replicas             int32
	FullyLabeledReplicas int32
	ReadyReplicas        int32
	AvailableReplicas    int32
}

// ReplicaSet - cluster ReplicaSet
type ReplicaSet struct {
	Version   string            `json:"version"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Spec      ReplicaSetSpec
	Status    ReplicaSetStatus
	YAML      []byte `json:"-"`
}

// GetName - object impl
func (r *ReplicaSet) GetName() string {
	return r.Name
}

// ReplicaSetOperator - ReplicaSet Operator interface
type ReplicaSetOperator interface {
	Get(ctx context.Context, namespace string, param *ReplicaSet) (*ReplicaSet, error)
	List(ctx context.Context, namespace string, options Options) ([]*ReplicaSet, error)
	Create(ctx context.Context, param *ReplicaSet, options Options) error
	CreateWithYaml(ctx context.Context, param *ReplicaSet, options Options) error
	Update(ctx context.Context, param *ReplicaSet) error
	Delete(ctx context.Context, param *ReplicaSet, options Options) error
	Watch(ctx context.Context, param *ReplicaSet, options Options) (Interface, error)
}

// EventType defines the possible types of events.
type EventType string

// watch event
const (
	Added    EventType = "ADDED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
	Bookmark EventType = "BOOKMARK"
	Error    EventType = "ERROR"

	DefaultChanSize int32 = 100
)

// WatchEvent represents a single event to a watched resource.
type WatchEvent struct {
	Type EventType

	// Object is:
	//  * If Type is Added or Modified: the new state of the object.
	//  * If Type is Deleted: the state of the object immediately before deletion.
	//  * If Type is Bookmark: the object (instance of a type being watched) where
	//    only ResourceVersion field is set. On successful restart of watch from a
	//    bookmark resourceVersion, client is guaranteed to not get repeat event
	//    nor miss any events.
	//  * If Type is Error: *api.Status is recommended; other types may make sense
	//    depending on context.
	Object Object
}

// Interface can be implemented by anything that knows how to watch and report changes.
type Interface interface {
	// Stops watching. Will close the channel returned by ResultChan(). Releases
	// any resources used by the watch.
	Stop()

	// Returns a chan which will receive all the events. If an error occurs
	// or Stop() is called, this channel will be closed, in which case the
	// watch should be completely cleaned up.
	ResultChan() <-chan WatchEvent
}
