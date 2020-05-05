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

func TestCreateSTS(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetSTSOperator().Create(context.Background(), &scheduler.STS{
		Version:   "apps/v1",
		Name:      fmt.Sprintf("sts-test-%v", time.Now().UnixNano()),
		Namespace: defaultNamespace,
		Labels:    map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
		Spec: scheduler.STSSpec{
			ServiceName: "nginx",
			Replicas:    2,
			Selector:    map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
			Template: scheduler.PodTemplateSpec{
				Name:      "nginx",
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
			},
		},
	}, scheduler.Options{}); err != nil {
		t.Fatalf("create sts failure, err:%v", err)
	}
	t.Logf("create sts success")
}

func TestListSTS(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetSTSOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v Replicas:%v", i, list[i].Name, list[i].Status.Replicas)
	}
}

func TestUpdateSTS(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetSTSOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetSTSOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update STS err:%v", err)
		}
	}
}

func TestDeleteSTS(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetSTSOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetSTSOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete STS err:%v", err)
		}
	}
}

func TestCreateSTSByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := fmt.Sprintf("sts-test-%v", time.Now().UnixNano())

	tpl, err := template.New("create_sts_yaml").Parse(stsYamlTpl)
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

	param := &scheduler.STS{
		Name:      name,
		Namespace: defaultNamespace,
		YAML:      buffer.Bytes(),
	}
	if err := client.GetSTSOperator().CreateWithYaml(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("Create STS WithYaml error:%v", err)
	}

	cfg, err := client.GetSTSOperator().Get(context.Background(), defaultNamespace, param)
	if err != nil {
		t.Fatalf("get sts error:%v", err)
	}

	t.Logf("sts content:%+v", cfg)

	if err := client.GetSTSOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("delete sts err:%v", err)
	}
}

var stsYamlTpl = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
  name: {{ .Name }}
spec:
  replicas: 2
  selector:
    matchLabels:
      app: {{ .LabelApp }}
      component: {{ .LabelComponent }}
  serviceName: nginx
  template:
    metadata:
      labels:
        app: {{ .LabelApp }}
        component: {{ .LabelComponent }}
      name: nginx
    spec:
      containers:
      - image: nginx
        imagePullPolicy: Always
        name: nginx
      nodeSelector:
        kubernetes.io/hostname: {{ .DeploymentNode }}
`
