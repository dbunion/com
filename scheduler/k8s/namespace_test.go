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

func TestCreateNamespace(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetNamespaceOperator().Create(context.Background(), &scheduler.Namespace{
		Name:   fmt.Sprintf("%s-%v", defaultNamespace, time.Now().UnixNano()),
		Labels: map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
	}, scheduler.Options{}); err != nil {
		t.Fatalf("create namespace failure, err:%v", err)
	}
	t.Logf("create namespace success")
}

func TestListNamespace(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	options := make(scheduler.Options)
	options[listOptionsKeyLimit] = 10
	options[listOptionsKeyLabelSelector] = fmt.Sprintf("app=%s", defaultLabelApp)

	list, err := client.GetNamespaceOperator().List(context.Background(), options)
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v", i, list[i].Name)
	}
}

func TestUpdateNamespace(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	options := make(scheduler.Options)
	options[listOptionsKeyLabelSelector] = fmt.Sprintf("app=%s", defaultLabelApp)

	list, err := client.GetNamespaceOperator().List(context.Background(), options)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetNamespaceOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update namespace err:%v", err)
		}
	}
}

func TestDeleteNamespace(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	options := make(scheduler.Options)
	options[listOptionsKeyLabelSelector] = fmt.Sprintf("app=%s", defaultLabelApp)

	list, err := client.GetNamespaceOperator().List(context.Background(), options)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetNamespaceOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete namespace err:%v", err)
		}
		t.Logf("success delete namespace %v", param.Name)
	}
}

func TestCreateNamespaceByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := fmt.Sprintf("%s-%v", defaultNamespace, time.Now().UnixNano())
	tpl, err := template.New("create_namespace_yaml").Parse(namespaceYamlTpl)
	if err != nil {
		t.Fatalf("parse yalm template error:%v", err)
	}

	var buffer bytes.Buffer
	m := map[string]interface{}{"Name": name, "LabelApp": defaultLabelApp, "LabelComponent": defaultLabelComponent}
	if err := tpl.Execute(&buffer, m); err != nil {
		t.Fatalf("execute tmplate error:%v", err)
	}

	param := &scheduler.Config{
		Name:      name,
		Namespace: defaultNamespace,
		YAML:      buffer.Bytes(),
	}
	if err := client.GetConfigOperator().CreateWithYaml(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("CreateWithYaml error:%v", err)
	}

	cfg, err := client.GetConfigOperator().Get(context.Background(), defaultNamespace, param)
	if err != nil {
		t.Fatalf("get config error:%v", err)
	}

	t.Logf("config content:%+v", cfg)

	if err := client.GetConfigOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("delete config err:%v", err)
	}
}

var namespaceYamlTpl = `
apiVersion: v1
kind: Namespace
metadata:
  labels:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
  name: {{ .Name }}
`
