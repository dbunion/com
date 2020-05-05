package k8s

import (
	"context"
	"github.com/dbunion/com/scheduler"
	v1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// DaemonSetClient ...
type DaemonSetClient struct {
	apiClient *Client
}

// newDaemonSetClient - new statefulSet client
func newDaemonSetClient(apiClient *Client) *DaemonSetClient {
	return &DaemonSetClient{
		apiClient: apiClient,
	}
}

// convertToDaemonSet - convert k8s's DaemonSet to DaemonSet
func convertToDaemonSet(n *v1.DaemonSet) *scheduler.DaemonSet {
	if n == nil {
		return nil
	}

	daeset := &scheduler.DaemonSet{
		Version:   n.APIVersion,
		Name:      n.Name,
		Namespace: n.Namespace,
		Labels:    n.Labels,
		Spec: scheduler.DaemonSetSpec{
			Selector:        n.Spec.Selector.MatchLabels,
			MinReadySeconds: n.Spec.MinReadySeconds,
		},
		Status: scheduler.DaemonSetStatus{
			CurrentNumberScheduled: n.Status.CurrentNumberScheduled,
			NumberMisscheduled:     n.Status.NumberMisscheduled,
			DesiredNumberScheduled: n.Status.DesiredNumberScheduled,
			NumberReady:            n.Status.NumberReady,
			NumberAvailable:        n.Status.NumberAvailable,
			NumberUnavailable:      n.Status.NumberUnavailable,
		},
	}
	return daeset
}

// Get - query DaemonSets info
func (c *DaemonSetClient) Get(ctx context.Context, namespace string, param *scheduler.DaemonSet) (*scheduler.DaemonSet, error) {
	n, err := c.apiClient.clientSet.AppsV1().DaemonSets(namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToDaemonSet(n), nil
}

// List - query DaemonSets list
func (c *DaemonSetClient) List(ctx context.Context, namespace string, options scheduler.Options) ([]*scheduler.DaemonSet, error) {
	list, err := c.apiClient.clientSet.AppsV1().DaemonSets(namespace).List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}

	stsList := make([]*scheduler.DaemonSet, 0)
	for i := 0; i < len(list.Items); i++ {
		stsList = append(stsList, convertToDaemonSet(&list.Items[i]))
	}

	return stsList, err
}

// Create - create new DaemonSets
func (c *DaemonSetClient) Create(ctx context.Context, param *scheduler.DaemonSet, options scheduler.Options) error {
	req := &v1.DaemonSet{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: param.Version,
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   param.Name,
			Labels: param.Labels,
		},
		Spec: v1.DaemonSetSpec{
			Selector: &meta_v1.LabelSelector{
				MatchLabels: param.Spec.Selector,
			},
			Template:             convertPodTemplateSpecToK8sPodTemplateSpec(param.Spec.Template),
			UpdateStrategy:       v1.DaemonSetUpdateStrategy{},
			MinReadySeconds:      0,
			RevisionHistoryLimit: nil,
		},
	}

	_, err := c.apiClient.clientSet.AppsV1().DaemonSets(param.Namespace).Create(ctx, req, meta_v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new daemonSet map with yalm
func (c *DaemonSetClient) CreateWithYaml(ctx context.Context, param *scheduler.DaemonSet, options scheduler.Options) error {
	var req v1.DaemonSet
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.AppsV1().DaemonSets(param.Namespace).Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update DaemonSets content
func (c *DaemonSetClient) Update(ctx context.Context, param *scheduler.DaemonSet) error {
	req, err := c.apiClient.clientSet.AppsV1().DaemonSets(param.Namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update fields
	req.Labels = param.Labels

	_, err = c.apiClient.clientSet.AppsV1().DaemonSets(param.Namespace).Update(ctx, req, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete DaemonSets map
func (c *DaemonSetClient) Delete(ctx context.Context, param *scheduler.DaemonSet, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.AppsV1().DaemonSets(param.Namespace).Delete(ctx, param.Name, op)
}

// Watch - watch DaemonSets change
func (c *DaemonSetClient) Watch(ctx context.Context, param *scheduler.DaemonSet, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.AppsV1().DaemonSets(param.Namespace).Watch(ctx, op)
	if err != nil {
		return nil, err
	}

	return NewWatcher(w), nil
}
