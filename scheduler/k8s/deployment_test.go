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

func TestCreateDeployment(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetDeploymentOperator().Create(context.Background(), &scheduler.Deployment{
		Version:   "apps/v1",
		Name:      fmt.Sprintf("dpl-test-%v", time.Now().UnixNano()),
		Namespace: defaultNamespace,
		Labels:    map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
		Spec: scheduler.DeploymentSpec{
			Replicas: 2,
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
		t.Fatalf("create deployment failure, err:%v", err)
	}
	t.Logf("create deployment success")
}

func TestListDeployment(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetDeploymentOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v", i, list[i].Name)
	}
}

func TestUpdateDeployment(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetDeploymentOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetDeploymentOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update Deployment err:%v", err)
		}
	}
}

func TestDeleteDeployment(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetDeploymentOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetDeploymentOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete Deployment err:%v", err)
		}
	}
}

func TestCreateDeploymentByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := fmt.Sprintf("dpl-test-%v", time.Now().UnixNano())

	tpl, err := template.New("create_deployment_yaml").Parse(deploymentYamlTpl)
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

	param := &scheduler.Deployment{
		Name:      name,
		Namespace: defaultNamespace,
		YAML:      buffer.Bytes(),
	}
	if err := client.GetDeploymentOperator().CreateWithYaml(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("CreateWithYaml error:%v", err)
	}

	cfg, err := client.GetDeploymentOperator().Get(context.Background(), defaultNamespace, param)
	if err != nil {
		t.Fatalf("get deployment error:%v", err)
	}

	t.Logf("deployment content:%+v", cfg)

	if err := client.GetDeploymentOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("delete deployment err:%v", err)
	}
}

var deploymentYamlTpl = `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
  name: {{ .Name }}
spec:
  replicas: 2
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: {{ .LabelApp }}
      component: {{ .LabelComponent }}
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
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
