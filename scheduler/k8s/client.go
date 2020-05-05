package k8s

import (
	"github.com/dbunion/com/log"
	"github.com/dbunion/com/scheduler"
	"github.com/juju/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	defaultLogger log.Logger
)

// ClientOpts - client options
type ClientOpts struct {
	URL         string
	Insecure    bool
	Username    string
	Password    string
	BearerToken string
	Data        interface{}
}

// APIClient - api client
type APIClient struct {
	clientSet *kubernetes.Clientset
}

// NewAPIClient - create new api client
func newAPIClient(opts *ClientOpts) (*APIClient, error) {
	config := &rest.Config{
		Host:     opts.URL,
		Username: opts.Username,
		Password: opts.Password,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: opts.Insecure,
		},
		BearerToken: opts.BearerToken,
	}

	// Create the ClientSet
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &APIClient{
		clientSet: client,
	}, nil
}

// Client - k8s base client
type Client struct {
	APIClient
	cfg                   scheduler.Param
	logger                log.Logger
	ConfigMap             scheduler.ConfigOperator
	Namespace             scheduler.NamespaceOperator
	Service               scheduler.ServiceOperator
	Pod                   scheduler.PodOperator
	Node                  scheduler.NodeOperator
	ReplicationController scheduler.RCOperator
	StatefulSet           scheduler.STSOperator
	DaemonSet             scheduler.DaemonSetOperator
	Deployment            scheduler.DeploymentOperator
	ReplicaSet            scheduler.ReplicaSetOperator
}

// NewK8sClient - create new scheduler client
func NewK8sClient() scheduler.Scheduler {
	return &Client{
		logger: defaultLogger,
	}
}

// newClient - create new api client
func newClient(opts *ClientOpts) (*Client, error) {
	baseClient, err := newAPIClient(opts)
	if err != nil {
		return nil, err
	}

	client := &Client{
		APIClient: *baseClient,
	}

	// init base operator
	client.Node = newNodeClient(client)
	client.Namespace = newNameSpaceClient(client)
	client.ConfigMap = newConfigMapClient(client)
	client.Service = newServiceClient(client)
	client.Pod = newPodClient(client)
	client.ReplicationController = newReplicationControllerClient(client)
	client.StatefulSet = newStatefulSetClient(client)
	client.DaemonSet = newDaemonSetClient(client)
	client.Deployment = newDeploymentClient(client)
	client.ReplicaSet = newReplicaSetClient(client)

	return client, nil
}

// GetNodeOperator - get node Operator
func (c *Client) GetNodeOperator() scheduler.NodeOperator {
	return c.Node
}

// GetNamespaceOperator - get namespace Operator
func (c *Client) GetNamespaceOperator() scheduler.NamespaceOperator {
	return c.Namespace
}

// GetConfigOperator - get config Operator
func (c *Client) GetConfigOperator() scheduler.ConfigOperator {
	return c.ConfigMap
}

// GetServiceOperator - get service Operator
func (c *Client) GetServiceOperator() scheduler.ServiceOperator {
	return c.Service
}

// GetPodOperator - get pod Operator
func (c *Client) GetPodOperator() scheduler.PodOperator {
	return c.Pod
}

// GetRCOperator - get rc Operator
func (c *Client) GetRCOperator() scheduler.RCOperator {
	return c.ReplicationController
}

// GetSTSOperator - get sts Operator
func (c *Client) GetSTSOperator() scheduler.STSOperator {
	return c.StatefulSet
}

// GetDaemonSetOperator - get DaemonSet Operator
func (c *Client) GetDaemonSetOperator() scheduler.DaemonSetOperator {
	return c.DaemonSet
}

// GetDeploymentOperator - get Deployment Operator
func (c *Client) GetDeploymentOperator() scheduler.DeploymentOperator {
	return c.Deployment
}

// GetReplicaSetOperator - get ReplicaSet Operator
func (c *Client) GetReplicaSetOperator() scheduler.ReplicaSetOperator {
	return c.ReplicaSet
}

// Close - release resource
func (c *Client) Close() error {
	return nil
}

// StartAndGC - init base object
func (c *Client) StartAndGC(config scheduler.Param) error {
	c.cfg = config

	if config.Logger != nil {
		c.logger = config.Logger
	}

	opts := &ClientOpts{}
	opts.URL = config.Server
	opts.Username = config.User
	opts.Password = config.Password
	opts.Insecure = config.Insecure
	if config.Token != "" {
		opts.BearerToken = config.Token
	}

	cli, err := newClient(opts)
	if err != nil {
		return err
	}

	*c = *cli
	return nil
}

// init - init env and logger
func init() {
	scheduler.Register(scheduler.TypeK8s, NewK8sClient)
	defaultLogger, _ = log.NewLogger(log.TypeZsskyLog, log.Config{
		Level:    log.LevelInfo,
		FilePath: "/tmp/scheduler.log",
	})
}
