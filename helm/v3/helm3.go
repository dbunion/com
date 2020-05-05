package v3

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/dbunion/com/helm"
	"github.com/zssky/tc/exec"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

// Helm3 is Helm center adapter and wrap helm3 command line
type Helm3 struct {
	cfg      helm.Config
	tempFile string
}

// NewHelm3 create new helm3 with default collection name.
func NewHelm3() helm.Helm {
	return &Helm3{}
}

// Install - This command installs a chart archive.
// flags： please read helm usage `helm install -h`
func (h *Helm3) Install(name, chart string, flags []string) error {
	params := []string{"helm", "--kubeconfig", h.tempFile, "install", name, chart}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	shell := strings.Join(params, " ")
	fmt.Printf("shell:%v\n", shell)

	_, err := exec.RunShellCommand(shell)
	if err != nil {
		return err
	}

	return nil
}

// List - list releases
func (h *Helm3) List(flags []string) ([]helm.Item, error) {
	params := []string{"helm", "--kubeconfig", h.tempFile, "list"}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	resp, err := exec.RunShellCommand(strings.Join(params, " "))
	if err != nil {
		return nil, err
	}

	list := make([]helm.Item, 0)
	lines := strings.Split(resp, "\n")
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "NAME") {
			continue
		}

		items := strings.Split(lines[i], "\t")
		if len(items) >= 7 {
			list = append(list, helm.Item{
				Name:       items[0],
				Namespace:  items[1],
				Revision:   items[2],
				Updated:    items[3],
				Status:     items[4],
				Chart:      items[5],
				AppVersion: items[6],
			})
		}
	}

	return list, nil
}

// RepoAdd - add chart repositories
func (h *Helm3) RepoAdd(name, url string, flags []string) error {
	params := []string{"helm", "--kubeconfig", h.tempFile, "repo", "add", name, url}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	_, err := exec.RunShellCommand(strings.Join(params, " "))
	if err != nil {
		return err
	}
	return nil
}

// RepoList - list chart repositories
func (h *Helm3) RepoList(flags []string) ([]helm.RepoItem, error) {
	params := []string{"helm", "--kubeconfig", h.tempFile, "repo", "list"}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	resp, err := exec.RunShellCommand(strings.Join(params, " "))
	if err != nil {
		return nil, err
	}

	list := make([]helm.RepoItem, 0)
	lines := strings.Split(resp, "\n")
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "NAME") {
			continue
		}

		items := strings.Split(lines[i], "\t")
		if len(items) >= 2 {
			list = append(list, helm.RepoItem{
				Name: items[0],
				URL:  items[1],
			})
		}
	}

	return list, nil
}

// RepoRemove - remove chart repositories
func (h *Helm3) RepoRemove(name string, flags []string) error {
	params := []string{"helm", "--kubeconfig", h.tempFile, "repo", "remove", name}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	_, err := exec.RunShellCommand(strings.Join(params, " "))
	if err != nil {
		return err
	}
	return nil
}

// RepoUpdate - update chart repositories
func (h *Helm3) RepoUpdate(flags []string) error {
	params := []string{"helm", "--kubeconfig", h.tempFile, "repo", "update"}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	_, err := exec.RunShellCommand(strings.Join(params, " "))
	if err != nil {
		return err
	}
	return nil
}

// SearchRepo - search repo
func (h *Helm3) SearchRepo(keyword string, flags []string) ([]helm.SearchItem, error) {
	params := []string{"helm", "--kubeconfig", h.tempFile, "search", "repo", keyword}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	resp, err := exec.RunShellCommand(strings.Join(params, " "))
	if err != nil {
		return nil, err
	}

	list := make([]helm.SearchItem, 0)
	lines := strings.Split(resp, "\n")
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "NAME") {
			continue
		}

		items := strings.Split(lines[i], "\t")
		if len(items) >= 4 {
			list = append(list, helm.SearchItem{
				Name:         items[0],
				ChartVersion: items[1],
				AppVersion:   items[2],
				Description:  items[3],
			})
		}
	}

	return list, nil
}

// Status - shows the status of a named release.
func (h *Helm3) Status(release string, flags []string) (string, error) {
	params := []string{"helm", "--kubeconfig", h.tempFile, "status", release}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	resp, err := exec.RunShellCommand(strings.Join(params, " "))
	if err != nil {
		return "", err
	}

	return resp, nil
}

// UnInstall - takes a release name and uninstalls the release.
func (h *Helm3) UnInstall(release string, flags []string) error {
	params := []string{"helm", "--kubeconfig", h.tempFile, "uninstall", release}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	_, err := exec.RunShellCommand(strings.Join(params, " "))
	if err != nil {
		return err
	}

	return nil
}

// Version - print the client version information
func (h *Helm3) Version(flags []string) (*helm.BuildInfo, error) {
	params := []string{"helm", "version"}

	if len(flags) > 0 {
		params = append(params, flags...)
	}

	resp, err := exec.RunShellCommand(strings.Join(params, " "))
	if err != nil {
		return nil, err
	}

	pos := strings.Index(resp, "{")
	end := strings.Index(resp, "}")

	substr := resp[pos+1 : end]
	substr = strings.Replace(substr, "\"", "", -1)
	items := strings.Split(substr, ",")
	if len(items) != 4 {
		return nil, fmt.Errorf("build info parse error")
	}

	buildInfo := &helm.BuildInfo{
		Version:      strings.Split(items[0], ":")[1],
		GitCommit:    strings.Split(items[1], ":")[1],
		GitTreeState: strings.Split(items[2], ":")[1],
		GoVersion:    strings.Split(items[3], ":")[1],
	}

	return buildInfo, nil
}

// StartAndGC start file Helm adapter.
func (h *Helm3) StartAndGC(cfg helm.Config) error {
	h.cfg = cfg

	// helm version check
	buildInfo, err := h.Version(nil)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(strings.ToLower(buildInfo.Version), "v3.") {
		return fmt.Errorf("invalid helm v3 version:%v", buildInfo.Version)
	}

	// init config file
	h.tempFile = fmt.Sprintf("/tmp/%x.config", md5.Sum([]byte(cfg.Server+cfg.UserName+cfg.Password)))
	fmt.Printf("tempFilePath: %v\n", h.tempFile)
	if _, err := os.Open(h.tempFile); err != nil {
		// build config content
		tpl, err := template.New("config").Parse(configTemplate)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, map[string]interface{}{"UserName": cfg.UserName, "Password": cfg.Password, "Server": cfg.Server}); err != nil {
			return err
		}

		if err := ioutil.WriteFile(h.tempFile, buf.Bytes(), os.ModePerm); err != nil {
			return err
		}
	}

	// add repo directly
	if err := h.RepoAdd(cfg.RepoName, cfg.RepoURL, nil); err != nil {
		return err
	}

	return nil
}

var configTemplate = `
apiVersion: v1
kind: Config
preferences: {}
users:
- name: kubecfg
  user:
    username: {{ .UserName }}
    password: {{ .Password }}
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: {{ .Server }}
  name: cluster-auth
contexts:
- context:
    cluster: cluster-auth
    user: kubecfg
  name: cluster-auth
current-context: cluster-auth
`

func init() {
	helm.Register(helm.TypeHelm3, NewHelm3)
}
