package k8s

import (
	"context"
	"fmt"
	"github.com/dbunion/com/scheduler"
	"github.com/juju/errors"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"
	"time"
)

//PodClient ...
type PodClient struct {
	apiClient *Client
}

//newPodClient ...
func newPodClient(apiClient *Client) *PodClient {
	return &PodClient{
		apiClient: apiClient,
	}
}

func convertToFlexVolumeSource(v *v1.FlexVolumeSource) *scheduler.FlexVolumeSource {
	return &scheduler.FlexVolumeSource{
		Driver:   v.Driver,
		FSType:   v.FSType,
		ReadOnly: false,
		Options:  v.Options,
	}
}

func convertToVolume(v *v1.Volume) *scheduler.Volume {
	if v == nil {
		return nil
	}
	ret := &scheduler.Volume{
		Name:  v.Name,
		Value: make(map[string]interface{}),
	}

	if v.FlexVolume != nil {
		ret.Value["Type"] = "FlexVolume"
		ret.Value["FlexVolume"] = convertToFlexVolumeSource(v.FlexVolume)
	} else if v.HostPath != nil {
		ret.Value["Type"] = "HostPath"
		ret.Value["HostPath"] = v.HostPath
	} else if v.EmptyDir != nil {
		ret.Value["Type"] = "EmptyDir"
		ret.Value["EmptyDir"] = v.EmptyDir
	} else if v.GCEPersistentDisk != nil {
		ret.Value["Type"] = "GCEPersistentDisk"
		ret.Value["GCEPersistentDisk"] = v.GCEPersistentDisk
	} else if v.AWSElasticBlockStore != nil {
		ret.Value["Type"] = "AWSElasticBlockStore"
		ret.Value["AWSElasticBlockStore"] = v.AWSElasticBlockStore
	} else if v.Secret != nil {
		ret.Value["Type"] = "Secret"
		ret.Value["Secret"] = v.Secret
	} else if v.NFS != nil {
		ret.Value["Type"] = "NFS"
		ret.Value["NFS"] = v.NFS
	} else if v.ISCSI != nil {
		ret.Value["Type"] = "ISCSI"
		ret.Value["ISCSI"] = v.ISCSI
	} else if v.Glusterfs != nil {
		ret.Value["Type"] = "Glusterfs"
		ret.Value["Glusterfs"] = v.Glusterfs
	} else if v.PersistentVolumeClaim != nil {
		ret.Value["Type"] = "PersistentVolumeClaim"
		ret.Value["PersistentVolumeClaim"] = v.PersistentVolumeClaim
	} else if v.RBD != nil {
		ret.Value["Type"] = "RBD"
		ret.Value["RBD"] = v.RBD
	} else if v.Cinder != nil {
		ret.Value["Type"] = "Cinder"
		ret.Value["Cinder"] = v.Cinder
	} else if v.CephFS != nil {
		ret.Value["Type"] = "CephFS"
		ret.Value["CephFS"] = v.CephFS
	} else if v.Flocker != nil {
		ret.Value["Type"] = "Flocker"
		ret.Value["Flocker"] = v.Flocker
	} else if v.DownwardAPI != nil {
		ret.Value["Type"] = "DownwardAPI"
		ret.Value["DownwardAPI"] = v.DownwardAPI
	} else if v.FC != nil {
		ret.Value["Type"] = "FC"
		ret.Value["FC"] = v.FC
	} else if v.AzureFile != nil {
		ret.Value["Type"] = "AzureFile"
		ret.Value["AzureFile"] = v.AzureFile
	} else if v.ConfigMap != nil {
		ret.Value["Type"] = "ConfigMap"
		ret.Value["ConfigMap"] = v.ConfigMap
	} else if v.VsphereVolume != nil {
		ret.Value["Type"] = "VsphereVolume"
		ret.Value["VsphereVolume"] = v.VsphereVolume
	} else if v.Quobyte != nil {
		ret.Value["Type"] = "Quobyte"
		ret.Value["Quobyte"] = v.Quobyte
	} else if v.AzureDisk != nil {
		ret.Value["Type"] = "AzureDisk"
		ret.Value["AzureDisk"] = v.AzureDisk
	} else if v.PhotonPersistentDisk != nil {
		ret.Value["Type"] = "PhotonPersistentDisk"
		ret.Value["PhotonPersistentDisk"] = v.PhotonPersistentDisk
	} else if v.Projected != nil {
		ret.Value["Type"] = "Projected"
		ret.Value["Projected"] = v.Projected
	} else if v.PortworxVolume != nil {
		ret.Value["Type"] = "PortworxVolume"
		ret.Value["PortworxVolume"] = v.PortworxVolume
	} else if v.ScaleIO != nil {
		ret.Value["Type"] = "ScaleIO"
		ret.Value["ScaleIO"] = v.ScaleIO
	} else if v.StorageOS != nil {
		ret.Value["Type"] = "StorageOS"
		ret.Value["StorageOS"] = v.StorageOS
	} else if v.CSI != nil {
		ret.Value["Type"] = "CSI"
		ret.Value["CSI"] = v.CSI
	}
	return ret
}

func convertToVolumes(list []v1.Volume) []scheduler.Volume {
	ret := make([]scheduler.Volume, 0)
	for i := 0; i < len(list); i++ {
		ret = append(ret, *convertToVolume(&list[i]))
	}
	return ret
}

func convertToContainer(c *v1.Container) *scheduler.Container {
	ports := make([]scheduler.ContainerPort, 0)
	for i := 0; i < len(c.Ports); i++ {
		port := c.Ports[i]
		ports = append(ports, scheduler.ContainerPort{
			Name:          port.Name,
			HostPort:      port.HostPort,
			ContainerPort: port.ContainerPort,
			Protocol:      string(port.Protocol),
			HostIP:        port.HostIP,
		})
	}

	limits := make(scheduler.ResourceList)
	for key, value := range c.Resources.Limits {
		limits[key.String()], _ = value.AsInt64()
	}

	requests := make(scheduler.ResourceList)
	for key, value := range c.Resources.Requests {
		requests[key.String()], _ = value.AsInt64()
	}

	resource := scheduler.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}

	var volumeMount []scheduler.VolumeMount
	if c.VolumeMounts != nil {
		volumeMount = make([]scheduler.VolumeMount, len(c.VolumeMounts))
		for i := 0; i < len(c.VolumeMounts); i++ {
			mount := c.VolumeMounts[i]
			volumeMount[i] = scheduler.VolumeMount{
				Name:        mount.Name,
				ReadOnly:    mount.ReadOnly,
				MountPath:   mount.MountPath,
				SubPath:     mount.SubPath,
				SubPathExpr: mount.SubPathExpr,
			}
		}
	}

	ret := &scheduler.Container{
		Name:         c.Name,
		Image:        c.Image,
		Command:      c.Command,
		Args:         c.Args,
		WorkingDir:   c.WorkingDir,
		Ports:        ports,
		Resources:    resource,
		VolumeMounts: volumeMount,
	}

	return ret
}

func convertToContainers(list []v1.Container) []scheduler.Container {
	length := len(list)
	if length == 0 {
		return nil
	}

	ret := make([]scheduler.Container, length)
	for i := 0; i < length; i++ {
		ret[i] = *convertToContainer(&list[i])
	}
	return ret
}

func convertToPodTemplateSpec(p *v1.PodTemplateSpec) *scheduler.PodTemplateSpec {
	return &scheduler.PodTemplateSpec{
		Name:      p.Name,
		Namespace: p.Namespace,
		Labels:    p.Labels,
		Spec: scheduler.PodSpec{
			Volumes:                       convertToVolumes(p.Spec.Volumes),
			Containers:                    convertToContainers(p.Spec.Containers),
			RestartPolicy:                 string(p.Spec.RestartPolicy),
			TerminationGracePeriodSeconds: p.Spec.TerminationGracePeriodSeconds,
			ActiveDeadlineSeconds:         p.Spec.ActiveDeadlineSeconds,
			NodeSelector:                  p.Spec.NodeSelector,
			NodeName:                      p.Spec.NodeName,
			HostNetwork:                   p.Spec.HostNetwork,
			Hostname:                      p.Spec.Hostname,
			PriorityClassName:             p.Spec.PriorityClassName,
			Priority:                      p.Spec.Priority,
		},
	}
}

func convertToPodSpec(p *v1.Pod) *scheduler.PodSpec {
	return &scheduler.PodSpec{
		Volumes:                       convertToVolumes(p.Spec.Volumes),
		Containers:                    convertToContainers(p.Spec.Containers),
		RestartPolicy:                 string(p.Spec.RestartPolicy),
		TerminationGracePeriodSeconds: p.Spec.TerminationGracePeriodSeconds,
		ActiveDeadlineSeconds:         p.Spec.ActiveDeadlineSeconds,
		NodeSelector:                  p.Spec.NodeSelector,
		NodeName:                      p.Spec.NodeName,
		HostNetwork:                   p.Spec.HostNetwork,
		Hostname:                      p.Spec.Hostname,
		PriorityClassName:             p.Spec.PriorityClassName,
		Priority:                      p.Spec.Priority,
	}
}

func convertToPodCondition(p *v1.Pod) []scheduler.PodCondition {
	var conditions []scheduler.PodCondition
	if length := len(p.Status.Conditions); length > 0 {
		conditions = make([]scheduler.PodCondition, length)
		for i := 0; i < length; i++ {
			cond := p.Status.Conditions[i]
			conditions[i] = scheduler.PodCondition{
				Type:               string(cond.Type),
				Status:             string(cond.Status),
				LastProbeTime:      cond.LastProbeTime.String(),
				LastTransitionTime: cond.LastProbeTime.String(),
				Reason:             cond.Reason,
				Message:            cond.Message,
			}
		}

	}
	return conditions
}

func convertToPodStatus(p *v1.Pod) *scheduler.PodStatus {
	status := &scheduler.PodStatus{
		Phase:      string(p.Status.Phase),
		Conditions: convertToPodCondition(p),
		Message:    p.Status.Message,
		Reason:     p.Status.Reason,
		HostIP:     p.Status.HostIP,
		PodIP:      p.Status.PodIP,
	}

	if p.Status.StartTime != nil {
		status.StartTime = p.Status.StartTime.String()
	}

	return status
}

// convertToPod - convert k8s's Pod to Pod
func convertToPod(p *v1.Pod) *scheduler.Pod {
	if p == nil {
		return nil
	}

	pod := &scheduler.Pod{
		Name:      p.Name,
		Namespace: p.Namespace,
		Labels:    p.Labels,
		Spec:      *convertToPodSpec(p),
		Status:    *convertToPodStatus(p),
	}
	return pod
}

func convertToEvent(n *v1.Event) *scheduler.Event {
	if n == nil {
		return nil
	}

	event := &scheduler.Event{
		Reason:              n.Reason,
		Message:             n.Message,
		Source:              n.Source.String(),
		FirstTimestamp:      n.FirstTimestamp.String(),
		LastTimestamp:       n.LastTimestamp.String(),
		Count:               n.Count,
		Type:                n.Type,
		EventTime:           n.EventTime.String(),
		Action:              n.Action,
		ReportingController: n.ReportingController,
		ReportingInstance:   n.ReportingInstance,
	}

	return event
}

func convertVolumeToK8sVolume(volume scheduler.Volume) v1.Volume {
	return v1.Volume{
		Name: volume.Name,
		VolumeSource: v1.VolumeSource{
			HostPath:              nil,
			EmptyDir:              nil,
			GCEPersistentDisk:     nil,
			AWSElasticBlockStore:  nil,
			GitRepo:               nil,
			Secret:                nil,
			NFS:                   nil,
			ISCSI:                 nil,
			Glusterfs:             nil,
			PersistentVolumeClaim: nil,
			RBD:                   nil,
			FlexVolume:            nil,
			Cinder:                nil,
			CephFS:                nil,
			Flocker:               nil,
			DownwardAPI:           nil,
			FC:                    nil,
			AzureFile:             nil,
			ConfigMap:             nil,
			VsphereVolume:         nil,
			Quobyte:               nil,
			AzureDisk:             nil,
			PhotonPersistentDisk:  nil,
			Projected:             nil,
			PortworxVolume:        nil,
			ScaleIO:               nil,
			StorageOS:             nil,
			CSI:                   nil,
		},
	}
}

func convertVolumesToK8sVolumes(volumes []scheduler.Volume) []v1.Volume {
	vols := make([]v1.Volume, len(volumes))
	for i := 0; i < len(volumes); i++ {
		vols[i] = convertVolumeToK8sVolume(volumes[i])
	}
	return vols
}

func convertContainerPortToK8sContainerPort(port scheduler.ContainerPort) v1.ContainerPort {
	return v1.ContainerPort{
		Name:          port.Name,
		HostPort:      port.HostPort,
		ContainerPort: port.ContainerPort,
		Protocol:      v1.Protocol(port.Protocol),
		HostIP:        port.HostIP,
	}
}

func convertContainerPortsToK8sContainerPorts(ports []scheduler.ContainerPort) []v1.ContainerPort {
	cts := make([]v1.ContainerPort, len(ports))
	for i := 0; i < len(ports); i++ {
		cts[i] = convertContainerPortToK8sContainerPort(ports[i])
	}
	return cts
}

func convertResourceRequirementsToK8sResourceRequirements(r scheduler.ResourceRequirements) v1.ResourceRequirements {
	// requests
	requests := make(v1.ResourceList)
	for key, value := range r.Requests {
		if value >= scheduler.QuantityG {
			q, _ := resource.ParseQuantity(fmt.Sprintf("%vG", value/scheduler.QuantityG))
			requests[v1.ResourceName(key)] = q
		} else {
			requests[v1.ResourceName(key)] = *resource.NewQuantity(value, resource.DecimalExponent)
		}
	}

	// limits
	limits := make(v1.ResourceList)
	for key, value := range r.Limits {
		if value >= scheduler.QuantityG {
			q, _ := resource.ParseQuantity(fmt.Sprintf("%vG", value/scheduler.QuantityG))
			limits[v1.ResourceName(key)] = q
		} else {
			limits[v1.ResourceName(key)] = *resource.NewQuantity(value, resource.DecimalExponent)
		}
	}

	return v1.ResourceRequirements{
		Requests: requests,
		Limits:   limits,
	}
}

func convertVolumeMountToK8sVolumeMount(vol scheduler.VolumeMount) v1.VolumeMount {
	return v1.VolumeMount{
		Name:        vol.Name,
		ReadOnly:    vol.ReadOnly,
		MountPath:   vol.MountPath,
		SubPath:     vol.SubPath,
		SubPathExpr: vol.SubPathExpr,
	}
}

func convertVolumeMountsToK8sVolumeMounts(volumes []scheduler.VolumeMount) []v1.VolumeMount {
	vols := make([]v1.VolumeMount, len(volumes))
	for i := 0; i < len(volumes); i++ {
		vols[i] = convertVolumeMountToK8sVolumeMount(volumes[i])
	}
	return vols
}

func convertContainerToK8sContainer(container scheduler.Container) v1.Container {
	return v1.Container{
		Name:         container.Name,
		Image:        container.Image,
		Command:      container.Command,
		Args:         container.Args,
		WorkingDir:   container.WorkingDir,
		Ports:        convertContainerPortsToK8sContainerPorts(container.Ports),
		Resources:    convertResourceRequirementsToK8sResourceRequirements(container.Resources),
		VolumeMounts: convertVolumeMountsToK8sVolumeMounts(container.VolumeMounts),
	}
}

func convertContainersToK8sContainers(containers []scheduler.Container) []v1.Container {
	cts := make([]v1.Container, len(containers))
	for i := 0; i < len(containers); i++ {
		cts[i] = convertContainerToK8sContainer(containers[i])
	}
	return cts
}

func updatePodSpec(src v1.PodSpec, change scheduler.PodSpec) v1.PodSpec {
	// update Resource
	for i := 0; i < len(src.Containers); i++ {
		if src.Containers[i].Name != change.Containers[i].Name {
			continue
		}

		srcResource := src.Containers[i].Resources
		changResource := change.Containers[i].Resources

		// update Limits
		for key, value := range srcResource.Limits {
			if v, found := changResource.Limits[key.String()]; found {
				srcInt64Value, ok := value.AsInt64()
				if !ok {
					continue
				}

				if v != srcInt64Value {
					if v >= scheduler.QuantityG {
						q, _ := resource.ParseQuantity(fmt.Sprintf("%vG", v/scheduler.QuantityG))
						srcResource.Limits[key] = q
					} else {
						srcResource.Limits[key] = *resource.NewQuantity(v, resource.DecimalExponent)
					}
				}
			}
		}

		// update Requests
		for key, value := range srcResource.Requests {
			if v, found := changResource.Requests[key.String()]; found {
				srcInt64Value, ok := value.AsInt64()
				if !ok {
					continue
				}
				if v != srcInt64Value {
					if v >= scheduler.QuantityG {
						q, _ := resource.ParseQuantity(fmt.Sprintf("%vG", v/scheduler.QuantityG))
						srcResource.Requests[key] = q
					} else {
						srcResource.Requests[key] = *resource.NewQuantity(v, resource.DecimalExponent)
					}
				}
			}
		}

		// update volume
		srcVolumes := src.Volumes
		changeVolumes := change.Volumes
		for i := 0; i < len(srcVolumes); i++ {
			if srcVolumes[i].Name != changeVolumes[i].Name {
				continue
			}

			// check FlexVolume
			if srcVolumes[i].FlexVolume != nil {
				if v, found := changeVolumes[i].Value["FlexVolume"]; found {
					flexVolumeSource, ok := v.(*scheduler.FlexVolumeSource)
					if ok && srcVolumes[i].FlexVolume.Options["size"] != flexVolumeSource.Options["size"] {
						src.Volumes[i].FlexVolume.Options["size"] = flexVolumeSource.Options["size"]
					}
				}
			}
		}

	}

	return src
}

func convertPodSpecToK8sPodSpec(spec scheduler.PodSpec) *v1.PodSpec {
	return &v1.PodSpec{
		Volumes:                       convertVolumesToK8sVolumes(spec.Volumes),
		InitContainers:                convertContainersToK8sContainers(spec.InitContainers),
		Containers:                    convertContainersToK8sContainers(spec.Containers),
		RestartPolicy:                 v1.RestartPolicy(spec.RestartPolicy),
		TerminationGracePeriodSeconds: spec.TerminationGracePeriodSeconds,
		ActiveDeadlineSeconds:         spec.ActiveDeadlineSeconds,
		NodeSelector:                  spec.NodeSelector,
		NodeName:                      spec.NodeName,
		HostNetwork:                   spec.HostNetwork,
		Hostname:                      spec.Hostname,
		PriorityClassName:             spec.PriorityClassName,
		Priority:                      spec.Priority,
	}
}

// Get - query pod info
func (c *PodClient) Get(ctx context.Context, namespace string, param *scheduler.Pod) (*scheduler.Pod, error) {
	n, err := c.apiClient.clientSet.CoreV1().Pods(namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return convertToPod(n), nil
}

// List - query pod list
func (c *PodClient) List(ctx context.Context, namespace string, options scheduler.Options) ([]*scheduler.Pod, error) {
	list, err := c.apiClient.clientSet.CoreV1().Pods(namespace).List(ctx, convertToListOptions(options))
	if err != nil {
		return nil, err
	}

	podList := make([]*scheduler.Pod, 0)
	for i := 0; i < len(list.Items); i++ {
		podList = append(podList, convertToPod(&list.Items[i]))
	}

	return podList, err
}

// Create - create new pod
func (c *PodClient) Create(ctx context.Context, param *scheduler.Pod, options scheduler.Options) error {
	req := &v1.Pod{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   param.Name,
			Labels: param.Labels,
		},
		Spec: *convertPodSpecToK8sPodSpec(param.Spec),
	}

	_, err := c.apiClient.clientSet.CoreV1().Pods(param.Namespace).Create(ctx, req, meta_v1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// CreateWithYaml - create new pod with yaml
func (c *PodClient) CreateWithYaml(ctx context.Context, param *scheduler.Pod, options scheduler.Options) error {
	var req v1.Pod
	if err := yaml.Unmarshal(param.YAML, &req); err != nil {
		return err
	}

	_, err := c.apiClient.clientSet.CoreV1().Pods(param.Namespace).Create(ctx, &req, convertToCreateOptions(options))
	if err != nil {
		return err
	}

	return nil
}

// Update - update pod content
func (c *PodClient) Update(ctx context.Context, param *scheduler.Pod) error {
	req, err := c.apiClient.clientSet.CoreV1().Pods(param.Namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return err
	}

	// update Labels
	req.Labels = param.Labels

	// update PodSpec
	req.Spec = updatePodSpec(req.Spec, param.Spec)

	_, err = c.apiClient.clientSet.CoreV1().Pods(param.Namespace).Update(ctx, req, meta_v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Delete - delete pod map
func (c *PodClient) Delete(ctx context.Context, param *scheduler.Pod, options scheduler.Options) error {
	op := convertToDeleteOptions(options)
	return c.apiClient.clientSet.CoreV1().Pods(param.Namespace).Delete(ctx, param.Name, op)
}

// GetEvents -  query pod event
func (c *PodClient) GetEvents(ctx context.Context, param *scheduler.Pod) ([]*scheduler.Event, error) {
	pod, err := c.apiClient.clientSet.CoreV1().Pods(param.Namespace).Get(ctx, param.Name, meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	events, err := c.apiClient.clientSet.CoreV1().Events(param.Namespace).Search(scheme.Scheme, pod)
	if err != nil {
		return nil, errors.Trace(err)
	}

	list := make([]*scheduler.Event, 0)
	for i := 0; i < len(events.Items); i++ {
		list = append(list, convertToEvent(&events.Items[i]))
	}

	return list, nil
}

// GetLogs - get pod logs
func (c *PodClient) GetLogs(ctx context.Context, namespace, name, container string) ([]byte, error) {
	sinceTime := meta_v1.NewTime(time.Now().Add(time.Duration(-2 * time.Hour)))

	lines := int64(2000)

	logOptions := &v1.PodLogOptions{
		SinceTime:  &sinceTime,
		TailLines:  &lines,
		Container:  container,
		Follow:     false,
		Previous:   false,
		Timestamps: false,
	}

	return c.getRawPodLogs(ctx, namespace, name, logOptions)
}

// getRawPodLogs - get raw pod logs
func (c *PodClient) getRawPodLogs(ctx context.Context, namespace, podID string, logOptions *v1.PodLogOptions) ([]byte, error) {
	req := c.apiClient.clientSet.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(podID).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec)

	readCloser, err := req.Stream(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer func() {
		_ = readCloser.Close()
	}()

	result, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return result, nil
}

// Watch - watch Pod change
func (c *PodClient) Watch(ctx context.Context, param *scheduler.Pod, options scheduler.Options) (scheduler.Interface, error) {
	op := convertToListOptions(options)
	w, err := c.apiClient.clientSet.CoreV1().Pods(param.Namespace).Watch(ctx, op)
	if err != nil {
		return nil, err
	}
	return NewWatcher(w), nil
}
