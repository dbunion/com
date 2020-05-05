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

func TestCreateService(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := client.GetServiceOperator().Create(context.Background(), &scheduler.Service{
		Name:      fmt.Sprintf("svc-test-%v", time.Now().UnixNano()),
		Namespace: defaultNamespace,
		Labels:    map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
		Spec: scheduler.ServiceSpec{
			Ports: []scheduler.ServicePort{
				{
					Name:       "web",
					Protocol:   "TCP",
					Port:       80,
					TargetPort: 80,
				},
			},
			Selector: map[string]string{"app": defaultLabelApp, "component": defaultLabelComponent},
		},
	}, scheduler.Options{}); err != nil {
		t.Fatalf("create service failure, err:%v", err)
	}
	t.Logf("create service success")
}

func TestListService(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetServiceOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(list); i++ {
		t.Logf("index:%v name:%v", i, list[i].Name)
	}
}

func TestUpdateService(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetServiceOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only update first one
	if len(list) > 0 {
		param := list[0]
		param.Labels["Update"] = time.Now().Format("20060102150405")
		if err := client.GetServiceOperator().Update(context.Background(), param); err != nil {
			t.Fatalf("update Service err:%v", err)
		}
	}
}

func TestDeleteService(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	list, err := client.GetServiceOperator().List(context.Background(), defaultNamespace, scheduler.Options{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	// only delete last one
	if len(list) > 0 {
		param := list[len(list)-1]
		if err := client.GetServiceOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
			t.Fatalf("delete Service err:%v", err)
		}
	}
}

func TestCreateServiceByYaml(t *testing.T) {
	if env == defaultEnv {
		return
	}
	client, err := newClient(&opt)
	if err != nil {
		t.Fatalf("%v", err)
	}

	name := fmt.Sprintf("svc-test-%v", time.Now().UnixNano())

	tpl, err := template.New("create_service_yaml").Parse(serviceYamlTpl)
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

	param := &scheduler.Service{
		Name:      name,
		Namespace: defaultNamespace,
		YAML:      buffer.Bytes(),
	}
	if err := client.GetServiceOperator().CreateWithYaml(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("Create Service WithYaml error:%v", err)
	}

	cfg, err := client.GetServiceOperator().Get(context.Background(), defaultNamespace, param)
	if err != nil {
		t.Fatalf("get service error:%v", err)
	}

	t.Logf("service content:%+v", cfg)

	if err := client.GetServiceOperator().Delete(context.Background(), param, scheduler.Options{}); err != nil {
		t.Fatalf("delete service err:%v", err)
	}
}

var serviceYamlTpl = `
apiVersion: v1
kind: Service
metadata:
  labels:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
  name: {{ .Name }}
spec:
  ports:
  - name: web
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: {{ .LabelApp }}
    component: {{ .LabelComponent }}
`
