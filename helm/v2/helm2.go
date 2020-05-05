package v2

import (
	"github.com/dbunion/com/helm"
)

// Helm2 is Helm center adapter and wrap helm2 command line
type Helm2 struct {
	cfg helm.Config
}

// NewHelm2 create new helm2 with default collection name.
func NewHelm2() helm.Helm {
	return &Helm2{}
}

// Install - This command installs a chart archive.
// flags： please read helm usage `helm install -h`
func (h *Helm2) Install(name, chart string, flags []string) error {
	return nil
}

// List - list releases
func (h *Helm2) List(flags []string) ([]helm.Item, error) {
	return nil, nil
}

// RepoAdd - add chart repositories
func (h *Helm2) RepoAdd(name, url string, flags []string) error {
	return nil
}

// RepoList - list chart repositories
func (h *Helm2) RepoList(flags []string) ([]helm.RepoItem, error) {
	return nil, nil
}

// RepoRemove - remove chart repositories
func (h *Helm2) RepoRemove(name string, flags []string) error {
	return nil
}

// RepoUpdate - update chart repositories
func (h *Helm2) RepoUpdate(flags []string) error {
	return nil
}

// SearchRepo - search repo
func (h *Helm2) SearchRepo(keyword string, flags []string) ([]helm.SearchItem, error) {
	return nil, nil
}

// Status - shows the status of a named release.
func (h *Helm2) Status(release string, flags []string) (string, error) {
	return "", nil
}

// UnInstall - takes a release name and uninstalls the release.
func (h *Helm2) UnInstall(release string, flags []string) error {
	return nil
}

// Version - print the client version information
func (h *Helm2) Version(flags []string) (*helm.BuildInfo, error) {
	return nil, nil
}

// StartAndGC start file Helm adapter.
func (h *Helm2) StartAndGC(cfg helm.Config) error {
	h.cfg = cfg
	return nil
}

func init() {
	helm.Register(helm.TypeHelm2, NewHelm2)
}
