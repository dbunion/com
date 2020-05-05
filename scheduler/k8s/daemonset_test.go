package k8s

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"text/template"
	"time"

	"github.com/dbunion/com/scheduler"
)

func TestCreateDaemonSet(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetDaemonSetOperator().Create(context.Background(), &scheduler.DaemonSet{
		Version:   "apps/v1",
		Name:      fmt.Sprintf("daemonset-test-%v", time.Now().UnixNano()),
		Namespace: defaultNamespace,
		Labels:    map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
		Spec: scheduler.DaemonSetSpec{
			Selector: map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
			Template: scheduler.PodTemplateSpec{
				Name:      "nginx",
				Namespace: defaultNamespace,
				Labels:    map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
				Spec: scheduler.PodSpec{
					Containers: []scheduler.Container{
						{
							Name:  "nginx",
							Image: "nginx",
							Resources: scheduler.ResourceRequirements{
								Limits:   scheduler.ResourceList{"cpu": 0, "memory": 0},
								Requests: scheduler.ResourceList{"cpu": 0, "memory": 0},
							},
						},
					},
					NodeSelector: map[string]string{"kubernetes.io/hostname": defaultDeploymentNode},
				},
			},
		},
	}, scheduler.Options{}); err != nil {
		t.Fatalf("create daemonSet failure, err:%v", err)
	}
	t.Logf("create daemonSet success")
}

func TestListDaemonSet(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetDaemonSetOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v NumberReady:%v", i, list[i].Name, list[i].Status.NumberReady)
	}
}

func TestUpdateDaemonSet(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetDaemonSetOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetDaemonSetOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update DaemonSet err:%v", err)
		}
	}
}

func TestDeleteDaemonSet(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetDaemonSetOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetDaemonSetOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete DaemonSet err:%v", err)
		}
	}
}

func TestCreateDaemonSetByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := fmt.Sprintf("daemonset-test-%v", time.Now().UnixNano())

	tpl, err := template.New("create_daemonSet_yaml").Parse(daemonSetYamlTpl)
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

	param := &scheduler.DaemonSet{
		Name:      name,
		Namespace: defaultNamespace,
		YAML:      buffer.Bytes(),
	}
	if err := client.GetDaemonSetOperator().CreateWithYaml(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("Create DaemonSet WithYaml error:%v", err)
	}

	cfg, err := client.GetDaemonSetOperator().Get(context.Background(), defaultNamespace, param)
	if err != nil {
		t.Fatalf("get daemonSet error:%v", err)
	}

	t.Logf("daemonSet content:%+v", cfg)

	if err := client.GetDaemonSetOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("delete daemonSet err:%v", err)
	}
}

var daemonSetYamlTpl = `
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  labels:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
  name: {{ .Name }}
spec:
  selector:
    matchLabels:
      app: {{ .LabelApp }}
      component: {{ .LabelComponent }}
  template:
    metadata:
      labels:
        app: {{ .LabelApp }}
        component: {{ .LabelComponent }}
      name: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
      nodeSelector:
        kubernetes.io/hostname: {{ .DeploymentNode }}
`
