package k8s

import (
	"context"
	"github.com/dbunion/com/scheduler"
	v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// DeploymentClient ...
type DeploymentClient struct {
	apiClient *Client
}

// newDeploymentClient - new statefulSet client
func newDeploymentClient(apiClient *Client) *DeploymentClient {
	return &DeploymentClient{
		apiClient: apiClient,
	}
}

// convertToDeployment - convert k8s's Deployment to Deployment
func convertToDeployment(n *v1.Deployment) *scheduler.Deployment {
	if n == nil {
		return nil
	}

	dp := &scheduler.Deployment{
		Version:   n.APIVersion,
		Name:      n.Name,
		Namespace: n.Namespace,
		Labels:    n.Labels,
		Spec: scheduler.DeploymentSpec{
			Replicas:                *n.Spec.Replicas,
			Selector:                n.Spec.Selector.MatchLabels,
			Template:                *convertToPodTemplateSpec(&n.Spec.Template),
			Strategy:                n.Spec.Strategy.String(),
			MinReadySeconds:         n.Spec.MinReadySeconds,
			Paused:                  n.Spec.Paused,
			ProgressDeadlineSeconds: n.Spec.ProgressDeadlineSeconds,
		},
		Status: scheduler.DeploymentStatus{
			Replicas:            n.Status.Replicas,
			UpdatedReplicas:     n.Status.UpdatedReplicas,
			ReadyReplicas:       n.Status.ReadyReplicas,
			AvailableReplicas:   n.Status.AvailableReplicas,
			UnavailableReplicas: n.Status.UnavailableReplicas,
		},
	}
	return dp
}

func convertPodTemplateSpecToK8sPodTemplateSpec(template scheduler.PodTemplateSpec) core_v1.PodTemplateSpec {
	return core_v1.PodTemplateSpec{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      template.Name,
			Namespace: template.Namespace,
			Labels:    template.Labels,
		},
		Spec: *convertPodSpecToK8sPodSpec(template.Spec),
	}
}

// Get - query Deployments info
func (c *DeploymentClient) Get(ctx context.Context, namespace string, param *scheduler.Deployment) (*scheduler.Deployment, error) {
	n, err := c.apiClient.clientSet.AppsV1().Deployments(namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToDeployment(n), nil
}

// List - query Deployments list
func (c *DeploymentClient) List(ctx context.Context, namespace string, options scheduler.Options) ([]*scheduler.Deployment, error) {
	list, err := c.apiClient.clientSet.AppsV1().Deployments(namespace).List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}

	stsList := make([]*scheduler.Deployment, 0)
	for i := 0; i < len(list.Items); i++ {
		stsList = append(stsList, convertToDeployment(&list.Items[i]))
	}

	return stsList, err
}

// Create - create new Deployments
func (c *DeploymentClient) Create(ctx context.Context, param *scheduler.Deployment, options scheduler.Options) error {
	req := &v1.Deployment{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: param.Version,
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   param.Name,
			Labels: param.Labels,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &param.Spec.Replicas,
			Selector: &meta_v1.LabelSelector{
				MatchLabels: param.Spec.Selector,
			},
			Template:                convertPodTemplateSpecToK8sPodTemplateSpec(param.Spec.Template),
			Strategy:                v1.DeploymentStrategy{Type: v1.DeploymentStrategyType(param.Spec.Strategy)},
			MinReadySeconds:         param.Spec.MinReadySeconds,
			Paused:                  param.Spec.Paused,
			ProgressDeadlineSeconds: param.Spec.ProgressDeadlineSeconds,
		},
	}

	_, err := c.apiClient.clientSet.AppsV1().Deployments(param.Namespace).Create(ctx, req, meta_v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new Deployments with yaml
func (c *DeploymentClient) CreateWithYaml(ctx context.Context, param *scheduler.Deployment, options scheduler.Options) error {
	var req v1.Deployment
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.AppsV1().Deployments(param.Namespace).Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update Deployments content
func (c *DeploymentClient) Update(ctx context.Context, param *scheduler.Deployment) error {
	req, err := c.apiClient.clientSet.AppsV1().Deployments(param.Namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update fields
	req.Labels = param.Labels

	_, err = c.apiClient.clientSet.AppsV1().Deployments(param.Namespace).Update(ctx, req, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete Deployments map
func (c *DeploymentClient) Delete(ctx context.Context, param *scheduler.Deployment, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.AppsV1().Deployments(param.Namespace).Delete(ctx, param.Name, op)
}

// Watch - watch Deployments change
func (c *DeploymentClient) Watch(ctx context.Context, param *scheduler.Deployment, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.AppsV1().Deployments(param.Namespace).Watch(ctx, op)
	if err != nil {
		return nil, err
	}

	return NewWatcher(w), nil
}
