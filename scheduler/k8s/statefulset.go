package k8s

import (
	"context"
	"github.com/dbunion/com/scheduler"
	v1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// StatefulSetClient ...
type StatefulSetClient struct {
	apiClient *Client
}

// newStatefulSetClient - new statefulSet client
func newStatefulSetClient(apiClient *Client) *StatefulSetClient {
	return &StatefulSetClient{
		apiClient: apiClient,
	}
}

// convertToSTS - convert k8s's StatefulSet to STS
func convertToSTS(n *v1.StatefulSet) *scheduler.STS {
	if n == nil {
		return nil
	}

	sts := &scheduler.STS{
		Version:   n.APIVersion,
		Name:      n.Name,
		Namespace: n.Namespace,
		Labels:    n.Labels,
		Spec: scheduler.STSSpec{
			Replicas:    *n.Spec.Replicas,
			Selector:    n.Spec.Selector.MatchLabels,
			Template:    *convertToPodTemplateSpec(&n.Spec.Template),
			ServiceName: n.Spec.ServiceName,
		},
		Status: scheduler.STSStatus{
			Replicas:        n.Status.Replicas,
			ReadyReplicas:   n.Status.ReadyReplicas,
			CurrentReplicas: n.Status.CurrentReplicas,
			UpdatedReplicas: n.Status.UpdatedReplicas,
			CurrentRevision: n.Status.CurrentRevision,
			UpdateRevision:  n.Status.UpdateRevision,
		},
	}
	return sts
}

// Get - query StatefulSets info
func (c *StatefulSetClient) Get(ctx context.Context, namespace string, param *scheduler.STS) (*scheduler.STS, error) {
	n, err := c.apiClient.clientSet.AppsV1().StatefulSets(namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToSTS(n), nil
}

// List - query StatefulSets list
func (c *StatefulSetClient) List(ctx context.Context, namespace string, options scheduler.Options) ([]*scheduler.STS, error) {
	list, err := c.apiClient.clientSet.AppsV1().StatefulSets(namespace).List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}

	stsList := make([]*scheduler.STS, 0)
	for i := 0; i < len(list.Items); i++ {
		stsList = append(stsList, convertToSTS(&list.Items[i]))
	}

	return stsList, err
}

// Create - create new StatefulSets
func (c *StatefulSetClient) Create(ctx context.Context, param *scheduler.STS, options scheduler.Options) error {
	req := &v1.StatefulSet{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   param.Name,
			Labels: param.Labels,
		},
		Spec: v1.StatefulSetSpec{
			Replicas:    &param.Spec.Replicas,
			Selector:    &meta_v1.LabelSelector{MatchLabels: param.Spec.Selector},
			Template:    convertPodTemplateSpecToK8sPodTemplateSpec(param.Spec.Template),
			ServiceName: param.Spec.ServiceName,
		},
	}

	_, err := c.apiClient.clientSet.AppsV1().StatefulSets(param.Namespace).Create(ctx, req, meta_v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new StatefulSet with yaml
func (c *StatefulSetClient) CreateWithYaml(ctx context.Context, param *scheduler.STS, options scheduler.Options) error {
	var req v1.StatefulSet
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.AppsV1().StatefulSets(param.Namespace).Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update StatefulSets content
func (c *StatefulSetClient) Update(ctx context.Context, param *scheduler.STS) error {
	req, err := c.apiClient.clientSet.AppsV1().StatefulSets(param.Namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update fields
	req.Labels = param.Labels

	_, err = c.apiClient.clientSet.AppsV1().StatefulSets(param.Namespace).Update(ctx, req, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete StatefulSets map
func (c *StatefulSetClient) Delete(ctx context.Context, param *scheduler.STS, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.AppsV1().StatefulSets(param.Namespace).Delete(ctx, param.Name, op)
}

// Watch - watch StatefulSets change
func (c *StatefulSetClient) Watch(ctx context.Context, param *scheduler.STS, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.AppsV1().StatefulSets(param.Namespace).Watch(ctx, op)
	if err != nil {
		return nil, err
	}
	return NewWatcher(w), nil
}
