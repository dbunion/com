package k8s

import (
	"context"
	"github.com/dbunion/com/scheduler"
	v1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// ReplicaSetClient ...
type ReplicaSetClient struct {
	apiClient *Client
}

// newReplicaSetClient - new statefulSet client
func newReplicaSetClient(apiClient *Client) *ReplicaSetClient {
	return &ReplicaSetClient{
		apiClient: apiClient,
	}
}

// convertToReplicaSet - convert k8s's ReplicaSet to ReplicaSet
func convertToReplicaSet(n *v1.ReplicaSet) *scheduler.ReplicaSet {
	if n == nil {
		return nil
	}

	rs := &scheduler.ReplicaSet{
		Name:      n.Name,
		Namespace: n.Namespace,
		Labels:    n.Labels,
		Spec: scheduler.ReplicaSetSpec{
			Replicas:        *n.Spec.Replicas,
			MinReadySeconds: n.Spec.MinReadySeconds,
			Selector:        n.Spec.Selector.MatchLabels,
			Template:        *convertToPodTemplateSpec(&n.Spec.Template),
		},
		Status: scheduler.ReplicaSetStatus{
			Replicas:             n.Status.Replicas,
			FullyLabeledReplicas: n.Status.FullyLabeledReplicas,
			ReadyReplicas:        n.Status.ReadyReplicas,
			AvailableReplicas:    n.Status.AvailableReplicas,
		},
	}
	return rs
}

// Get - query ReplicaSets info
func (c *ReplicaSetClient) Get(ctx context.Context, namespace string, param *scheduler.ReplicaSet) (*scheduler.ReplicaSet, error) {
	n, err := c.apiClient.clientSet.AppsV1().ReplicaSets(namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToReplicaSet(n), nil
}

// List - query ReplicaSets list
func (c *ReplicaSetClient) List(ctx context.Context, namespace string, options scheduler.Options) ([]*scheduler.ReplicaSet, error) {
	list, err := c.apiClient.clientSet.AppsV1().ReplicaSets(namespace).List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}

	stsList := make([]*scheduler.ReplicaSet, 0)
	for i := 0; i < len(list.Items); i++ {
		stsList = append(stsList, convertToReplicaSet(&list.Items[i]))
	}

	return stsList, err
}

// Create - create new ReplicaSets
func (c *ReplicaSetClient) Create(ctx context.Context, param *scheduler.ReplicaSet, options scheduler.Options) error {
	req := &v1.ReplicaSet{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   param.Name,
			Labels: param.Labels,
		},
		Spec: v1.ReplicaSetSpec{
			Replicas:        &param.Spec.Replicas,
			MinReadySeconds: param.Spec.MinReadySeconds,
			Selector: &meta_v1.LabelSelector{
				MatchLabels: param.Spec.Selector,
			},
			Template: convertPodTemplateSpecToK8sPodTemplateSpec(param.Spec.Template),
		},
	}

	_, err := c.apiClient.clientSet.AppsV1().ReplicaSets(param.Namespace).Create(ctx, req, meta_v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new ReplicaSet with yaml
func (c *ReplicaSetClient) CreateWithYaml(ctx context.Context, param *scheduler.ReplicaSet, options scheduler.Options) error {
	var req v1.ReplicaSet
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.AppsV1().ReplicaSets(param.Namespace).Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update ReplicaSets content
func (c *ReplicaSetClient) Update(ctx context.Context, param *scheduler.ReplicaSet) error {
	req, err := c.apiClient.clientSet.AppsV1().ReplicaSets(param.Namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update fields
	req.Labels = param.Labels

	_, err = c.apiClient.clientSet.AppsV1().ReplicaSets(param.Namespace).Update(ctx, req, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete ReplicaSets map
func (c *ReplicaSetClient) Delete(ctx context.Context, param *scheduler.ReplicaSet, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.AppsV1().ReplicaSets(param.Namespace).Delete(ctx, param.Name, op)
}

// Watch - watch ReplicaSet change
func (c *ReplicaSetClient) Watch(ctx context.Context, param *scheduler.ReplicaSet, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.AppsV1().ReplicaSets(param.Namespace).Watch(ctx, op)
	if err != nil {
		return nil, err
	}

	return NewWatcher(w), nil
}
