package k8s

const (
	defaultNode           = "127.0.0.1"
	defaultLabelApp       = "scheduler"
	defaultLabelComponent = "api"
	defaultNamespace      = "scheduler-test"
	defaultEnv            = "test"
)

var defaultDeploymentNode string
var opt ClientOpts
var env = defaultEnv // test local

func init() {
	defaultDeploymentNode = "127.0.0.1"
	opt = ClientOpts{
		URL:      "https://127.0.0.1:8080",
		Insecure: true,
		Username: "admin",
		Password: "123456",
		Data:     nil,
	}
}
