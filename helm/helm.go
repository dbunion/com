package helm

import (
	"fmt"
)

const (
	// TypeHelm2 - type helm2
	TypeHelm2 = "helm2"
	// TypeHelm3 - type helm3
	TypeHelm3 = "helm3"
)

// AuthType - helm auth type
type AuthType string

const (
	// AuthTypeBasic - auth type base
	AuthTypeBasic AuthType = "basic"
	// AuthTypeToken - auth type token
	AuthTypeToken AuthType = "token"
	// AuthTypeCert - auth type cert
	AuthTypeCert AuthType = "cert"
)

// Cluster contains information about how to communicate with a kubernetes cluster
type Cluster struct {
	// Server is the address of the kubernetes cluster (https://hostname:port).
	Server string `json:"server"`
	// InsecureSkipTLSVerify skips the validity check for the server's certificate. This will make your HTTPS connections insecure.
	// +optional
	InsecureSkipTLSVerify bool `json:"insecure-skip-tls-verify,omitempty"`
	// CertificateAuthority is the path to a cert file for the certificate authority.
	// +optional
	CertificateAuthority string `json:"certificate-authority,omitempty"`
	// CertificateAuthorityData contains PEM-encoded certificate authority certificates. Overrides CertificateAuthority
	// +optional
	CertificateAuthorityData []byte `json:"certificate-authority-data,omitempty"`
}

// AuthInfo contains information that describes identity information.  This is use to tell the kubernetes cluster who you are.
type AuthInfo struct {
	// ClientCertificate is the path to a client cert file for TLS.
	// +optional
	ClientCertificate string `json:"client-certificate,omitempty"`
	// ClientCertificateData contains PEM-encoded data from a client cert file for TLS. Overrides ClientCertificate
	// +optional
	ClientCertificateData []byte `json:"client-certificate-data,omitempty"`
	// ClientKey is the path to a client key file for TLS.
	// +optional
	ClientKey string `json:"client-key,omitempty"`
	// ClientKeyData contains PEM-encoded data from a client key file for TLS. Overrides ClientKey
	// +optional
	ClientKeyData []byte `json:"client-key-data,omitempty"`
	// Token is the bearer token for authentication to the kubernetes cluster.
	// +optional
	Token string `json:"token,omitempty"`
	// TokenFile is a pointer to a file that contains a bearer token (as described above).  If both Token and TokenFile are present, Token takes precedence.
	// +optional
	TokenFile string `json:"tokenFile,omitempty"`
	// Username is the username for basic authentication to the kubernetes cluster.
	// +optional
	Username string `json:"username,omitempty"`
	// Password is the password for basic authentication to the kubernetes cluster.
	// +optional
	Password string `json:"password,omitempty"`
}

// Config - init config
type Config struct {
	AuthType AuthType `json:"auth_type"`
	// Cluster is a object of referencable names to cluster config
	Cluster *Cluster `json:"cluster"`
	// AuthInfo is a object of referencable names to user config
	AuthInfo *AuthInfo `json:"users"`
	// RepoName is a helm chart repo name
	RepoName string `json:"repo"`
	// RepoURL is a helm chart repo url
	RepoURL string `json:"repo_url"`
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
