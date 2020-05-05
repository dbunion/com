package k8s

import (
	"context"
	"github.com/dbunion/com/scheduler"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// NameSpaceClient - namespace operator client
type NameSpaceClient struct {
	apiClient *Client
}

// newNameSpaceClient - create new namespace
func newNameSpaceClient(apiClient *Client) *NameSpaceClient {
	return &NameSpaceClient{
		apiClient: apiClient,
	}
}

// convertToNamespace - convert Namespace to Namespace
func convertToNamespace(c *v1.Namespace) *scheduler.Namespace {
	if c == nil {
		return nil
	}

	ns := scheduler.Namespace{
		Name:   c.Name,
		Labels: c.Labels,
		Status: scheduler.NamespaceStatus{
			Phase: string(c.Status.Phase),
		},
	}
	return &ns
}

// Get - query namespace list
func (c *NameSpaceClient) Get(ctx context.Context, param *scheduler.Namespace) (*scheduler.Namespace, error) {
	ns, err := c.apiClient.clientSet.CoreV1().Namespaces().Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToNamespace(ns), nil
}

// List - query namespace list
func (c *NameSpaceClient) List(ctx context.Context, options scheduler.Options) ([]*scheduler.Namespace, error) {
	list, err := c.apiClient.clientSet.CoreV1().Namespaces().List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}

	slist := make([]*scheduler.Namespace, 0)
	for i := 0; i < len(list.Items); i++ {
		slist = append(slist, convertToNamespace(&list.Items[i]))
	}

	return slist, nil
}

// Create - create new namespace
func (c *NameSpaceClient) Create(ctx context.Context, param *scheduler.Namespace, options scheduler.Options) error {
	req := &v1.Namespace{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   param.Name,
			Labels: param.Labels,
		},
	}

	_, err := c.apiClient.clientSet.CoreV1().Namespaces().Create(ctx, req, meta_v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new namespace with yaml
func (c *NameSpaceClient) CreateWithYaml(ctx context.Context, param *scheduler.Namespace, options scheduler.Options) error {
	var req v1.Namespace
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.CoreV1().Namespaces().Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update namespace
func (c *NameSpaceClient) Update(ctx context.Context, param *scheduler.Namespace) error {
	req, err := c.apiClient.clientSet.CoreV1().Namespaces().Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update fields
	req.Labels = param.Labels

	_, err = c.apiClient.clientSet.CoreV1().Namespaces().Update(ctx, req, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete namespace
func (c *NameSpaceClient) Delete(ctx context.Context, param *scheduler.Namespace, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.CoreV1().Namespaces().Delete(ctx, param.Name, op)
}

// Watch - watch Namespace change
func (c *NameSpaceClient) Watch(ctx context.Context, param *scheduler.Namespace, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.CoreV1().ConfigMaps(param.Name).Watch(ctx, op)
	if err != nil {
		return nil, err
	}
	return NewWatcher(w), nil
}
