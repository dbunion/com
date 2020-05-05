package k8s

import (
	"context"
	"github.com/dbunion/com/scheduler"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// ConfigMapClient - config map client wrap
type ConfigMapClient struct {
	apiClient *Client
}

//newConfigMapClient...
func newConfigMapClient(apiClient *Client) *ConfigMapClient {
	return &ConfigMapClient{
		apiClient: apiClient,
	}
}

// convertToConfig - convert configmap to config
func convertToConfig(c *v1.ConfigMap) *scheduler.Config {
	if c == nil {
		return nil
	}

	config := &scheduler.Config{
		Name:       c.Name,
		Namespace:  c.Namespace,
		BinaryData: c.BinaryData,
		Data:       c.Data,
		Labels:     c.Labels,
		Reserved:   nil,
	}
	return config
}

// Get - query namespace list
func (c *ConfigMapClient) Get(ctx context.Context, namespace string, param *scheduler.Config) (*scheduler.Config, error) {
	ns, err := c.apiClient.clientSet.CoreV1().ConfigMaps(namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToConfig(ns), nil
}

// List - query config map list
func (c *ConfigMapClient) List(ctx context.Context, namespace string, options scheduler.Options) ([]*scheduler.Config, error) {
	list, err := c.apiClient.clientSet.CoreV1().ConfigMaps(namespace).List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}

	slist := make([]*scheduler.Config, 0)
	for i := 0; i < len(list.Items); i++ {
		slist = append(slist, convertToConfig(&list.Items[i]))

	}

	return slist, err
}

// Create - create new config map
func (c *ConfigMapClient) Create(ctx context.Context, param *scheduler.Config, options scheduler.Options) error {
	req := &v1.ConfigMap{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      param.Name,
			Namespace: param.Namespace,
			Labels:    param.Labels,
		},
		Data:       param.Data,
		BinaryData: param.BinaryData,
	}

	_, err := c.apiClient.clientSet.CoreV1().ConfigMaps(param.Namespace).Create(ctx, req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new config map with yaml
func (c *ConfigMapClient) CreateWithYaml(ctx context.Context, param *scheduler.Config, options scheduler.Options) error {
	var req v1.ConfigMap
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.CoreV1().ConfigMaps(param.Namespace).Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update config content
func (c *ConfigMapClient) Update(ctx context.Context, param *scheduler.Config) error {
	req, err := c.apiClient.clientSet.CoreV1().ConfigMaps(param.Namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update fields
	req.Labels = param.Labels
	req.Data = param.Data
	req.BinaryData = param.BinaryData

	_, err = c.apiClient.clientSet.CoreV1().ConfigMaps(param.Namespace).Update(ctx, req, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete config map
func (c *ConfigMapClient) Delete(ctx context.Context, param *scheduler.Config, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.CoreV1().ConfigMaps(param.Namespace).Delete(ctx, param.Name, op)
}

// Watch - watch Config change
func (c *ConfigMapClient) Watch(ctx context.Context, param *scheduler.Config, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.CoreV1().ConfigMaps(param.Namespace).Watch(ctx, op)
	if err != nil {
		return nil, err
	}

	return NewWatcher(w), nil
}
