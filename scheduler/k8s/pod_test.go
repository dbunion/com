package k8s

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dbunion/com/scheduler"
	"testing"
	"text/template"
	"time"
)

func TestCreatePod(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetPodOperator().Create(context.Background(), &scheduler.Pod{
		Name:      fmt.Sprintf("scheduler-pod-test-%v", time.Now().UnixNano()),
		Namespace: defaultNamespace,
		Labels:    map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
		Spec: scheduler.PodSpec{
			Containers: []scheduler.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
			NodeSelector: map[string]string{"kubernetes.io/hostname": defaultDeploymentNode},
		},
	}, scheduler.Options{}); err != nil {
		t.Fatalf("create pod failure, err:%v", err)
	}
	t.Logf("create pod success")
}

func TestListPod(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetPodOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v", i, list[i].Name)
	}
}

func TestUpdatePod(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetPodOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetPodOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update pod err:%v", err)
		}
	}
}

func TestUpdatePodResource(t *testing.T) {
	if env == defaultEnv {
		return
	}

	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	pod, err := client.GetPodOperator().Get(context.Background(), defaultEnv, &scheduler.Pod{Name: "grafana-b4f886f5b-qcgvc"})
	if err != nil {
		t.Fatalf("%v", err)
	}

	pod.Labels["Update"] = time.Now().Format("20060102150405")
	pod.Spec.Containers[0].Resources.Limits["cpu"] = 5
	pod.Spec.Containers[0].Resources.Limits["memory"] = 5 * scheduler.QuantityG
	pod.Spec.Containers[0].Resources.Requests["memory"] = 5 * scheduler.QuantityG

	flexVolum, ok := pod.Spec.Volumes[1].Value["FlexVolume"].(*scheduler.FlexVolumeSource)
	if ok {
		flexVolum.Options["size"] = "101Gi"
		pod.Spec.Volumes[1].Value["FlexVolume"] = flexVolum
	}

	if err := client.GetPodOperator().Update(context.Background(), pod); err != nil {
		t.Fatalf("update pod err:%v", err)
	}
}

func TestDeletePod(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetPodOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetPodOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete pod err:%v", err)
		}
	}
}

func TestCreatePodByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := fmt.Sprintf("scheduler-pod-test-%v", time.Now().UnixNano())

	tpl, err := template.New("create_pod_yaml").Parse(podYamlTpl)
	if err != nil {
		t.Fatalf("parse yaml template error:%v", err)
	}

	var buffer bytes.Buffer
	m := map[string]interface{}{
		"Name":           name,
		"LabelApp":       defaultLabelApp,
		"LabelComponent": defaultLabelComponent,
		"DeploymentNode": defaultDeploymentNode,
	}
	if err := tpl.Execute(&buffer, m); err != nil {
		t.Fatalf("execute tmplate error:%v", err)
	}

	param := &scheduler.Pod{
		Name:      name,
		Namespace: defaultNamespace,
		YAML:      buffer.Bytes(),
	}
	if err := client.GetPodOperator().CreateWithYaml(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("Create Pod WithYaml error:%v", err)
	}

	pod, err := client.GetPodOperator().Get(context.Background(), defaultNamespace, param)
	if err != nil {
		t.Fatalf("get Pod error:%v", err)
	}

	t.Logf("pod content:%+v", pod)

	if err := client.GetPodOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("delete Pod err:%v", err)
	}
}

var podYamlTpl = `
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
  name: {{ .Name }}
spec:
  containers:
  - image: nginx
    imagePullPolicy: Always
    name: nginx
  nodeSelector:
    kubernetes.io/hostname: {{ .DeploymentNode }}
`
