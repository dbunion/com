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

func TestCreateNode(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetNodeOperator().Create(context.Background(), &scheduler.Node{
		Name:   defaultNode,
		Labels: map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
	}, scheduler.Options{}); err != nil {
		t.Fatalf("create node failure, err:%v", err)
	}
	t.Logf("create node success")
}

func TestDescribeNode(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	nodeDesc, err := client.GetNodeOperator().Describe(context.Background(), &scheduler.Node{
		Name: defaultNode,
	})

	if err != nil {
		t.Fatalf("Describe err:%v", err)
	}

	t.Logf("node desc:%+v", nodeDesc)
}

func TestGetNode(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	node, err := client.GetNodeOperator().Get(context.Background(), &scheduler.Node{
		Name: defaultNode,
	})

	if err != nil {
		t.Fatalf("get node info failure")
	}

	t.Logf("get node success, node:%v", node)
}

func TestListNode(t *testing.T) {
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

	list, err := client.GetNodeOperator().List(context.Background(), options)
	if err != nil {
		t.Fatalf("list node failure, err:%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v", i, list[i].Name)
	}
}

func TestUpdateNode(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	options := make(scheduler.Options)
	options[listOptionsKeyLabelSelector] = fmt.Sprintf("app=%s", defaultLabelApp)

	list, err := client.GetNodeOperator().List(context.Background(), options)
	if err != nil {
		t.Fatalf("list node failure, err:%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetNodeOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update node err:%v", err)
		}
	}
}

func TestDeleteNode(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	options := make(scheduler.Options)
	options[listOptionsKeyLabelSelector] = fmt.Sprintf("app=%s", defaultLabelApp)

	list, err := client.GetNodeOperator().List(context.Background(), options)
	if err != nil {
		t.Fatalf("list nodes err:%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetNodeOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete node err:%v", err)
		}
		t.Logf("success delete node %v", param.Name)
	}
}

func TestCreateNodeByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := defaultNode

	tpl, err := template.New("create_node_yaml").Parse(nodeYamlTpl)
	if err != nil {
		t.Fatalf("parse yaml template error:%v", err)
	}

	var buffer bytes.Buffer
	m := map[string]interface{}{
		"Name":           name,
		"LabelApp":       defaultLabelApp,
		"LabelComponent": defaultLabelComponent,
	}
	if err := tpl.Execute(&buffer, m); err != nil {
		t.Fatalf("execute tmplate error:%v", err)
	}

	param := &scheduler.Node{
		Name: name,
		YAML: buffer.Bytes(),
	}
	if err := client.GetNodeOperator().CreateWithYaml(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("Create node WithYaml error:%v", err)
	}

	node, err := client.GetNodeOperator().Get(context.Background(), param)
	if err != nil {
		t.Fatalf("get node error:%v", err)
	}

	t.Logf("node content:%+v", node)

	if err := client.GetNodeOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("delete node err:%v", err)
	}
}

var nodeYamlTpl = `
apiVersion: v1
kind: Node
metadata:
  labels:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
  name: {{ .Name }}
`
