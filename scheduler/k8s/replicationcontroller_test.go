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

func TestCreateRC(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetRCOperator().Create(context.Background(), &scheduler.RC{
		Version:   "apps/v1",
		Name:      fmt.Sprintf("rc-test-%v", time.Now().UnixNano()),
		Namespace: defaultNamespace,
		Labels:    map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
		Spec: scheduler.RCSpec{
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
						},
					},
					NodeSelector: map[string]string{"kubernetes.io/hostname": defaultDeploymentNode},
				},
			},
		},
	}, scheduler.Options{}); err != nil {
		t.Fatalf("create rc failure, err:%v", err)
	}
	t.Logf("create rc success")
}

func TestListRC(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetRCOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v Replicas:%v", i, list[i].Name, list[i].Status.Replicas)
	}
}

func TestUpdateRC(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetRCOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetRCOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update RC err:%v", err)
		}
	}
}

func TestDeleteRC(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetRCOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetRCOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete RC err:%v", err)
		}
	}
}

func TestCreateRCByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := fmt.Sprintf("rc-test-%v", time.Now().UnixNano())

	tpl, err := template.New("create_rc_yaml").Parse(rcYamlTpl)
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

	param := &scheduler.RC{
		Name:      name,
		Namespace: defaultNamespace,
		YAML:      buffer.Bytes(),
	}
	if err := client.GetRCOperator().CreateWithYaml(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("Create RC WithYaml error:%v", err)
	}

	cfg, err := client.GetRCOperator().Get(context.Background(), defaultNamespace, param)
	if err != nil {
		t.Fatalf("get rc error:%v", err)
	}

	t.Logf("rc content:%+v", cfg)

	if err := client.GetRCOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("delete rc err:%v", err)
	}
}

var rcYamlTpl = `
apiVersion: v1
kind: ReplicationController
metadata:
  labels:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
  name: {{ .Name }}
spec:
  replicas: 2
  selector:
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
        imagePullPolicy: Always
        name: nginx
      nodeSelector:
        kubernetes.io/hostname: {{ .DeploymentNode }}
`
