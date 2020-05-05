package k8s

import (
	"context"
	"github.com/dbunion/com/scheduler"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// ReplicationControllerClient replication controller client
type ReplicationControllerClient struct {
	apiClient *Client
}

//newReplicationControllerClient ...
func newReplicationControllerClient(apiClient *Client) *ReplicationControllerClient {
	return &ReplicationControllerClient{
		apiClient: apiClient,
	}
}

// convertToRC - convert k8s's Pod to Pod
func convertToRC(n *v1.ReplicationController) *scheduler.RC {
	if n == nil {
		return nil
	}

	rc := &scheduler.RC{
		Version:   n.APIVersion,
		Name:      n.Name,
		Namespace: n.Namespace,
		Labels:    n.Labels,
		Spec: scheduler.RCSpec{
			Replicas:        *n.Spec.Replicas,
			MinReadySeconds: n.Spec.MinReadySeconds,
			Selector:        n.Spec.Selector,
			Template:        *convertToPodTemplateSpec(n.Spec.Template),
		},
		Status: scheduler.RCStatus{
			Replicas:             n.Status.ReadyReplicas,
			FullyLabeledReplicas: n.Status.FullyLabeledReplicas,
			ReadyReplicas:        n.Status.ReadyReplicas,
			AvailableReplicas:    n.Status.AvailableReplicas,
			ObservedGeneration:   n.Status.ObservedGeneration,
		},
		YAML: nil,
	}
	return rc
}

// Get - query ReplicationControllers info
func (c *ReplicationControllerClient) Get(ctx context.Context, namespace string, param *scheduler.RC) (*scheduler.RC, error) {
	n, err := c.apiClient.clientSet.CoreV1().ReplicationControllers(namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToRC(n), nil
}

// List - query ReplicationControllers list
func (c *ReplicationControllerClient) List(ctx context.Context, namespace string, options scheduler.Options) ([]*scheduler.RC, error) {
	list, err := c.apiClient.clientSet.CoreV1().ReplicationControllers(namespace).List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}
	podList := make([]*scheduler.RC, 0)
	for i := 0; i < len(list.Items); i++ {
		podList = append(podList, convertToRC(&list.Items[i]))
	}

	return podList, err
}

// Create - create new ReplicationControllers
func (c *ReplicationControllerClient) Create(ctx context.Context, param *scheduler.RC, options scheduler.Options) error {
	template := convertPodTemplateSpecToK8sPodTemplateSpec(param.Spec.Template)
	req := &v1.ReplicationController{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   param.Name,
			Labels: param.Labels,
		},
		Spec: v1.ReplicationControllerSpec{
			Replicas:        &param.Spec.Replicas,
			MinReadySeconds: param.Spec.MinReadySeconds,
			Selector:        param.Spec.Selector,
			Template:        &template,
		},
	}

	_, err := c.apiClient.clientSet.CoreV1().ReplicationControllers(param.Namespace).Create(ctx, req, meta_v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new ReplicationControllers with yaml
func (c *ReplicationControllerClient) CreateWithYaml(ctx context.Context, param *scheduler.RC, options scheduler.Options) error {
	var req v1.ReplicationController
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.CoreV1().ReplicationControllers(param.Namespace).Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update ReplicationControllers content
func (c *ReplicationControllerClient) Update(ctx context.Context, param *scheduler.RC) error {
	req, err := c.apiClient.clientSet.CoreV1().ReplicationControllers(param.Namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update fields
	req.Labels = param.Labels

	_, err = c.apiClient.clientSet.CoreV1().ReplicationControllers(param.Namespace).Update(ctx, req, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete ReplicationControllers map
func (c *ReplicationControllerClient) Delete(ctx context.Context, param *scheduler.RC, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.CoreV1().ReplicationControllers(param.Namespace).Delete(ctx, param.Name, op)
}

// Watch - watch RC change
func (c *ReplicationControllerClient) Watch(ctx context.Context, param *scheduler.RC, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.CoreV1().ReplicationControllers(param.Namespace).Watch(ctx, op)
	if err != nil {
		return nil, err
	}
	return NewWatcher(w), nil
}
