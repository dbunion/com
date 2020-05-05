package helm

import "fmt"

const (
	// TypeHelm2 - type helm2
	TypeHelm2 = "helm2"
	// TypeHelm3 - type helm3
	TypeHelm3 = "helm3"
)

// Config - init config
type Config struct {
	Server   string `json:"server"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	RepoName string `json:"repo"`
	RepoURL  string `json:"repo_url"`
}

// Item - release item
type Item struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Revision   string `json:"revision"`
	Updated    string `json:"updated"`
	Status     string `json:"status"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
}

// RepoItem - repo item
type RepoItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// SearchItem - search repo item
type SearchItem struct {
	Name         string `json:"name"`
	ChartVersion string `json:"chart_version"`
	AppVersion   string `json:"app_version"`
	Description  string `json:"description"`
}

// BuildInfo - helm build info
type BuildInfo struct {
	Version      string `json:"version"`
	GitCommit    string `json:"git_commit"`
	GitTreeState string `json:"git_tree_state"`
	GoVersion    string `json:"go_version"`
}

// Helm interface contains all behaviors for Helm adapter.
// usage:
type Helm interface {
	// This command installs a chart archive.
	// flags： please read helm usage `helm install -h`
	Install(name, chart string, flags []string) error

	// list releases
	List(flags []string) ([]Item, error)

	// add chart repositories
	RepoAdd(name, url string, flags []string) error

	// list chart repositories
	RepoList(flags []string) ([]RepoItem, error)

	// remove chart repositories
	RepoRemove(name string, flags []string) error

	// update chart repositories
	RepoUpdate(flags []string) error

	// search repo
	SearchRepo(keyword string, flags []string) ([]SearchItem, error)

	// shows the status of a named release.
	Status(release string, flags []string) (string, error)

	// takes a release name and uninstalls the release.
	UnInstall(release string, flags []string) error

	// print the client version information
	Version(flags []string) (*BuildInfo, error)

	// start gc routine based on config string settings.
	StartAndGC(cfg Config) error
}

// Instance is a function create a new Helm Instance
type Instance func() Helm

var adapters = make(map[string]Instance)

// Register makes a Helm adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("Helm: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("Helm: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewHelm Create a new Helm driver by adapter name and config string.
// config need to be correct JSON as string: {"server": "xxx.xxx.com", "user": "xxxx", "password":"xxxxx"}.
// it will start gc automatically.
func NewHelm(adapterName string, cfg Config) (adapter Helm, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("Helm: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(cfg)
	if err != nil {
		adapter = nil
	}
	return
}
