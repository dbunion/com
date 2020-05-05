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

func TestCreateReplicaSet(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetReplicaSetOperator().Create(context.Background(), &scheduler.ReplicaSet{
		Version:   "apps/v1",
		Name:      fmt.Sprintf("rs-test-%v", time.Now().UnixNano()),
		Namespace: defaultNamespace,
		Labels:    map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
		Spec: scheduler.ReplicaSetSpec{
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
		t.Fatalf("create replicaSet failure, err:%v", err)
	}
	t.Logf("create replicaSet success")
}

func TestListReplicaSet(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetReplicaSetOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v Replicas:%v", i, list[i].Name, list[i].Status.Replicas)
	}
}

func TestUpdateReplicaSet(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetReplicaSetOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetReplicaSetOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update ReplicaSet err:%v", err)
		}
	}
}

func TestDeleteReplicaSet(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetReplicaSetOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetReplicaSetOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete ReplicaSet err:%v", err)
		}
	}
}

func TestCreateReplicaSetByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := fmt.Sprintf("rs-test-%v", time.Now().UnixNano())

	tpl, err := template.New("create_replicaSet_yaml").Parse(replicaSetYamlTpl)
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

	param := &scheduler.ReplicaSet{
		Name:      name,
		Namespace: defaultNamespace,
		YAML:      buffer.Bytes(),
	}
	if err := client.GetReplicaSetOperator().CreateWithYaml(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("Create ReplicaSet WithYaml error:%v", err)
	}

	cfg, err := client.GetReplicaSetOperator().Get(context.Background(), defaultNamespace, param)
	if err != nil {
		t.Fatalf("get replicaSet error:%v", err)
	}

	t.Logf("replicaSet content:%+v", cfg)

	if err := client.GetReplicaSetOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("delete replicaSet err:%v", err)
	}
}

var replicaSetYamlTpl = `
apiVersion: extensions/v1beta1
kind: ReplicaSet
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
