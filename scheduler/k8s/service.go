package k8s

import (
	"context"
	"github.com/dbunion/com/scheduler"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/yaml"
	"strconv"
)

// ServiceClient ...
type ServiceClient struct {
	apiClient *Client
}

//newServiceClient ...
func newServiceClient(apiClient *Client) *ServiceClient {
	return &ServiceClient{
		apiClient: apiClient,
	}
}

// convertToService - convert k8s'Service to Service
func convertToService(c *v1.Service) *scheduler.Service {
	if c == nil {
		return nil
	}

	s := &scheduler.Service{
		Name:      c.Name,
		Namespace: c.Namespace,
		Labels:    c.Labels,
		Spec: scheduler.ServiceSpec{
			Ports:           convertToSchedulerServicePorts(c.Spec.Ports),
			Selector:        c.Spec.Selector,
			ClusterIP:       c.Spec.ClusterIP,
			Type:            string(c.Spec.Type),
			ExternalIPs:     c.Spec.ExternalIPs,
			SessionAffinity: string(c.Spec.SessionAffinity),
			LoadBalancerIP:  c.Spec.LoadBalancerIP,
		},
	}
	return s
}

// convertToServicePort - convert ServicePort to k8s'ServicePort
func convertToServicePort(c *scheduler.ServicePort) *v1.ServicePort {
	if c == nil {
		return nil
	}

	s := &v1.ServicePort{
		Name:       c.Name,
		Protocol:   v1.Protocol(c.Protocol),
		Port:       c.Port,
		TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: c.TargetPort},
	}

	return s
}

func convertToServicePorts(list []scheduler.ServicePort) []v1.ServicePort {
	if len(list) == 0 {
		return nil
	}

	ports := make([]v1.ServicePort, 0)
	for i := 0; i < len(list); i++ {
		ports = append(ports, *convertToServicePort(&list[i]))
	}
	return ports
}

// convertToSchedulerServicePort - convert ServicePort to k8s'ServicePort
func convertToSchedulerServicePort(c *v1.ServicePort) *scheduler.ServicePort {
	if c == nil {
		return nil
	}

	s := &scheduler.ServicePort{
		Name:     c.Name,
		Protocol: string(c.Protocol),
		Port:     c.Port,
	}

	if c.TargetPort.Type == intstr.Int {
		s.TargetPort = c.TargetPort.IntVal
	} else if c.TargetPort.Type == intstr.String {
		port, _ := strconv.Atoi(c.TargetPort.StrVal)
		s.TargetPort = int32(port)
	}

	return s
}

func convertToSchedulerServicePorts(list []v1.ServicePort) []scheduler.ServicePort {
	if len(list) == 0 {
		return nil
	}

	ports := make([]scheduler.ServicePort, 0)
	for i := 0; i < len(list); i++ {
		ports = append(ports, *convertToSchedulerServicePort(&list[i]))
	}
	return ports
}

// Get - query service list
func (c *ServiceClient) Get(ctx context.Context, namespace string, param *scheduler.Service) (*scheduler.Service, error) {
	s, err := c.apiClient.clientSet.CoreV1().Services(namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToService(s), nil
}

// List - query service map list
func (c *ServiceClient) List(ctx context.Context, namespace string, options scheduler.Options) ([]*scheduler.Service, error) {
	list, err := c.apiClient.clientSet.CoreV1().Services(namespace).List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}

	slist := make([]*scheduler.Service, 0)
	for i := 0; i < len(list.Items); i++ {
		slist = append(slist, convertToService(&list.Items[i]))

	}

	return slist, err
}

// Create - create new service map
func (c *ServiceClient) Create(ctx context.Context, param *scheduler.Service, options scheduler.Options) error {
	req := &v1.Service{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      param.Name,
			Namespace: param.Namespace,
			Labels:    param.Labels,
		},
		Spec: v1.ServiceSpec{
			Ports:     convertToServicePorts(param.Spec.Ports),
			Selector:  param.Spec.Selector,
			ClusterIP: param.Spec.ClusterIP,
			Type:      v1.ServiceType(param.Spec.Type),
		},
	}

	_, err := c.apiClient.clientSet.CoreV1().Services(param.Namespace).Create(ctx, req, meta_v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new service with yaml
func (c *ServiceClient) CreateWithYaml(ctx context.Context, param *scheduler.Service, options scheduler.Options) error {
	var req v1.Service
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.CoreV1().Services(param.Namespace).Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update service content
func (c *ServiceClient) Update(ctx context.Context, param *scheduler.Service) error {
	req, err := c.apiClient.clientSet.CoreV1().Services(param.Namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update fields
	req.Labels = param.Labels

	_, err = c.apiClient.clientSet.CoreV1().Services(param.Namespace).Update(ctx, req, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete service map
func (c *ServiceClient) Delete(ctx context.Context, param *scheduler.Service, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.CoreV1().Services(param.Namespace).Delete(ctx, param.Name, op)
}

// Watch - watch Service change
func (c *ServiceClient) Watch(ctx context.Context, param *scheduler.Service, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.CoreV1().Services(param.Namespace).Watch(ctx, op)
	if err != nil {
		return nil, err
	}
	return NewWatcher(w), nil
}
