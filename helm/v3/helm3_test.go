package v3

import (
	"github.com/dbunion/com/helm"
	"testing"
)

func TestHelm3(t *testing.T) {
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		Server:   "127.0.0.1:8080",
		UserName: "admin",
		Password: "password",
		RepoName: "cetic",
		RepoURL:  "https://cetic.github.io/helm-charts",
	})

	if err != nil {
		t.Fatalf("create new helm error, err:%v", err)
	}

	repoList, err := h3.RepoList([]string{})
	if err != nil {
		t.Fatalf("%v", err)
	}

	for i := 0; i < len(repoList); i++ {
		t.Logf("%v", repoList[i])
	}

	if err := h3.RepoUpdate([]string{}); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestVersion(t *testing.T) {
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		Server:   "http://127.0.0.1:8080",
		UserName: "admin",
		Password: "123456",
		RepoName: "cetic",
		RepoURL:  "https://cetic.github.io/helm-charts",
	})

	if err != nil {
		t.Fatalf("create new helm error, err:%v", err)
	}

	buildInfo, err := h3.Version(nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("buildInfo:%v", buildInfo)
}
