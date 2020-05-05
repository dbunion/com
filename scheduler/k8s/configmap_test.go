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

func TestCreateConfig(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetConfigOperator().Create(context.Background(), &scheduler.Config{
		Name:       fmt.Sprintf("test-configmap-%v", time.Now().UnixNano()),
		Namespace:  defaultNamespace,
		BinaryData: nil,
		Data:       map[string]string{"config.json": `{"key":"value"}`},
		Labels:     map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
		Reserved:   nil,
	}, scheduler.Options{}); err != nil {
		t.Fatalf("create config failure, err:%v", err)
	}
	t.Logf("create config success")
}

func TestListConfig(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetConfigOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v", i, list[i].Name)
	}
}

func TestUpdateConfig(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetConfigOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetConfigOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update config err:%v", err)
		}
	}
}

func TestDeleteConfig(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetConfigOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetConfigOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete config err:%v", err)
		}
	}
}

func TestCreateByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := fmt.Sprintf("test-configmap-%v", time.Now().UnixNano())

	tpl, err := template.New("create_yaml").Parse(configYamlTpl)
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

var configYamlTpl = `
apiVersion: v1
data:
  config.json: '{"key":"value"}'
kind: ConfigMap
metadata:
  labels:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
  name: {{ .Name }}
`
