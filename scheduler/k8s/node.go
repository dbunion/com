package k8s

import (
	"context"
	"github.com/dbunion/com/scheduler"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/yaml"
)

//NodeClient ...
type NodeClient struct {
	apiClient *Client
}

// newNodeClient - create new node client
func newNodeClient(apiClient *Client) *NodeClient {
	return &NodeClient{
		apiClient: apiClient,
	}
}

// convertToNode - convert k8s'sNode to Node
func convertToNode(n *v1.Node) *scheduler.Node {
	if n == nil {
		return nil
	}

	node := &scheduler.Node{
		Name:   n.Name,
		Labels: n.Labels,
		Status: *convertToSchedulerNodeStatus(&n.Status),
	}
	return node
}

func convertToSchedulerNodeStatus(c *v1.NodeStatus) *scheduler.NodeStatus {
	if c == nil {
		return nil
	}

	cap := make(map[string]int64)
	for k, v := range c.Capacity {
		cap[k.String()], _ = v.AsInt64()
	}

	all := make(map[string]int64)
	for k, v := range c.Allocatable {
		all[k.String()], _ = v.AsInt64()
	}

	s := &scheduler.NodeStatus{
		Capacity:    cap,
		Allocatable: all,
		Phase:       string(c.Phase),
	}

	return s
}

// Get - query node info
func (c *NodeClient) Get(ctx context.Context, param *scheduler.Node) (*scheduler.Node, error) {
	n, err := c.apiClient.clientSet.CoreV1().Nodes().Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToNode(n), nil
}

// List - query node list
func (c *NodeClient) List(ctx context.Context, options scheduler.Options) ([]*scheduler.Node, error) {
	list, err := c.apiClient.clientSet.CoreV1().Nodes().List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}

	slist := make([]*scheduler.Node, 0)
	for i := 0; i < len(list.Items); i++ {
		slist = append(slist, convertToNode(&list.Items[i]))
	}

	return slist, err
}

// Create - create new node
func (c *NodeClient) Create(ctx context.Context, param *scheduler.Node, options scheduler.Options) error {
	req := &v1.Node{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   param.Name,
			Labels: param.Labels,
		},
	}

	_, err := c.apiClient.clientSet.CoreV1().Nodes().Create(ctx, req, meta_v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new node with yaml
func (c *NodeClient) CreateWithYaml(ctx context.Context, param *scheduler.Node, options scheduler.Options) error {
	var req v1.Node
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.CoreV1().Nodes().Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update node content
func (c *NodeClient) Update(ctx context.Context, param *scheduler.Node) error {
	req, err := c.apiClient.clientSet.CoreV1().Nodes().Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update fields
	req.Labels = param.Labels

	_, err = c.apiClient.clientSet.CoreV1().Nodes().Update(ctx, req, convertToUpdateOptions(scheduler.Options{}))
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete node map
func (c *NodeClient) Delete(ctx context.Context, param *scheduler.Node, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.CoreV1().Nodes().Delete(ctx, param.Name, op)
}

// Watch - watch Node change
func (c *NodeClient) Watch(ctx context.Context, param *scheduler.Node, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.CoreV1().Nodes().Watch(ctx, op)
	if err != nil {
		return nil, err
	}
	return NewWatcher(w), nil
}

// Describe - describe node resource info
func (c *NodeClient) Describe(ctx context.Context, param *scheduler.Node) (*scheduler.NodeDetail, error) {
	node, err := c.Get(ctx, param)
	if err != nil {
		return nil, err
	}

	name := param.Name
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + name + ",status.phase!=" + string(v1.PodSucceeded) + ",status.phase!=" + string(v1.PodFailed))
	if err != nil {
		return nil, err
	}

	canViewPods := true
	nodeNonTerminatedPodsList, err := c.apiClient.clientSet.CoreV1().Pods("").List(ctx, meta_v1.ListOptions{FieldSelector: fieldSelector.String()})
	if err != nil {
		if !errors.IsForbidden(err) {
			return nil, err
		}
		canViewPods = false
	}

	pods := make([]scheduler.PodResource, 0)
	if canViewPods && nodeNonTerminatedPodsList != nil {
		for _, pod := range nodeNonTerminatedPodsList.Items {
			req, limit := PodRequestsAndLimits(&pod)
			cpuReq, cpuLimit, memoryReq, memoryLimit := req[v1.ResourceCPU], limit[v1.ResourceCPU], req[v1.ResourceMemory], limit[v1.ResourceMemory]
			pods = append(pods, scheduler.PodResource{
				Resource: scheduler.Resource{
					CPURequest:    cpuReq.MilliValue(),
					CPULimit:      cpuLimit.MilliValue(),
					MemoryRequest: memoryReq.MilliValue(),
					MemoryLimit:   memoryLimit.MilliValue(),
				},
				Namespace: pod.Namespace,
				Name:      pod.Name,
			})
		}
	}

	nodeDetail := &scheduler.NodeDetail{
		Capacity:    node.Status.Capacity,
		Allocatable: node.Status.Allocatable,
		Pods:        pods,
	}

	return nodeDetail, nil
}

// PodRequestsAndLimits returns a dictionary of all defined resources summed up for all
// containers of the pod. If pod overhead is non-nil, the pod overhead is added to the
// total container resource requests and to the total container limits which have a
// non-zero quantity.
func PodRequestsAndLimits(pod *v1.Pod) (reqs, limits v1.ResourceList) {
	reqs, limits = v1.ResourceList{}, v1.ResourceList{}
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}
	// init containers define the minimum of any resource
	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}

	// Add overhead for running a pod to the sum of requests and to non-zero limits:
	if pod.Spec.Overhead != nil {
		addResourceList(reqs, pod.Spec.Overhead)

		for name, quantity := range pod.Spec.Overhead {
			if value, ok := limits[name]; ok && !value.IsZero() {
				value.Add(quantity)
				limits[name] = value
			}
		}
	}
	return
}

// addResourceList adds the resources in newList to list
func addResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}

// maxResourceList sets list to the greater of list/newList for every resource
// either list
func maxResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = quantity.DeepCopy()
			continue
		} else {
			if quantity.Cmp(value) > 0 {
				list[name] = quantity.DeepCopy()
			}
		}
	}
}
