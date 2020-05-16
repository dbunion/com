package v3

import (
	"github.com/dbunion/com/helm"
	"testing"
)

const (
	runTest                  = "test"
	runLocal                 = "local"
	env                      = runTest
	server                   = "https://127.0.0.1:6443"
	userName                 = "admin"
	password                 = "a813ccf47e92861b"
	token                    = "eyJhbGciOiJSUzI1NiIsImtpZCI6InVzT1hRbWF0SHdBeWVocDl1cTV0QlJlSGNvT3R3ZC1ySWFDdWc3VTBnTGMifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJjbHVzdGVycm9sZS1hZ2dyZWdhdGlvbi1jb250cm9sbGVyLXRva2VuLXg2bDluIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImNsdXN0ZXJyb2xlLWFnZ3JlZ2F0aW9uLWNvbnRyb2xsZXIiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJmZDdkNjhmMS02MWFlLTQ0MmMtOGVjNy0xNjEyYjI2YjgyYmIiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06Y2x1c3RlcnJvbGUtYWdncmVnYXRpb24tY29udHJvbGxlciJ9.vIms7qGlv2w7yIToL0tZMgrzR9cWU9i330cAKuFoWAIiOmgqvgBXTw1ZagqYbM5NllQ0u7HCJuaqRUMtNcJFpLtmYzU5LB3FDpNcQZR0grDdLrTnPoH2NMNoeOWrCh_2p6stqrrQytt26BZIKms6e_B7a1FNTVoY4FCBxyJ7uCkHqkcvxNqKLX9PgbWCuzcY827T1q3yv_CM_mM1A8SQmNSUTQHvd6YC-DqkPQssonfOtrbP0yDof-E-QldTLZvNdimkNWhkZ2glgD2jwAmvWlNGIER7PtQ7PwOfv3Z83YGClr5N6LonYskL8k4QZLFvX8yvGs7AfOR00ljPbQPSww%"
	certificateAuthorityData = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJd01ETXhNakV6TVRFd01Wb1hEVE13TURNeE1ERXpNVEV3TVZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTDNMCnlOa3lzbmJxTmZ2UlJJTlhPMDV4eHBIVnFleEN4ekM2SFNCQ09CUzBOQURla0FCaUlxZGM3aDJudmdXZExjQXYKWHpMa2hWbEFydE00WVJNWERiTUh1cmxMRUQ0NHN6NUpBd2VhNHZvNW9EdXdNMlZmL0g3aXlxcmFjZ0lia3c3YgpLeWNqMzlndk1sa09zNDFaS0twWnV1V1ZCNXRqbmlHMnlEbU5xVnJXbDRtMkdzeWhQTWJaYnIxQnJaT2I1UDdECldTWG9aYzBDQVBzN21DOWZKOVJqUkQycTlVTzl3SGFvSEczQkxSMnI3cnFkQkFIbDkyRkI3aGlhZjhRSUhiaU8KTEVxWFoxTC9zZVRwcVk1eCtxYWpobENNdUVVcjNlVUM0Vmt0eWc5WWJhYko3RUo0WHY0bWVvOHZHNkpGVC9FNwpDREJEOWFpV1hocXkvKytKd2xNQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFEbkRiRTdtL0JxS1dkWXNPeEhvQkRqenhITVYKQUFYVWZ0dWY4TlNCa050bDdBb01taklOMVJ0WUVqdW1JMUlwTzh1alJ5RCt4UWpETWZrS0xsVzJidG4wK0xtSApqRGcxYVd6bVNhRFUzcUdzWUdPUGZTOWx2WGdqTGV0UEtSeldHVzFVaVhGZFNGOWZnamxJNFlaYkNHV1FMYS9iCndNNUplTC84TkNmQ0N5R3NQT3dNY3VHdjlvQUJORUlsa2ZTNTdURWVTN0svaUFOVTdId1VXWnJ6MjFpRmxhejQKVGNSdkZGOFAvSXRFckcvRVNpZ1l4b3FFcU9taW1HR0JxMVRZMzBCZUptYm5RTnV0V2pLbHBGT3Vtd3JmbS9WaAo4eDdkUFlkMDZzenYvZFYwenRhOW9MZlBFckdrMkNqcnVhbDVnM1BsL0l0MitvNWtyL3FXaThya1c4bz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
	clientCertificateData    = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM5RENDQWR5Z0F3SUJBZ0lJY01McWRPMmRWQjh3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TURBek1USXhNekV4TURGYUZ3MHlNVEExTVRZd01EUXhOVE5hTURZeApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sc3dHUVlEVlFRREV4SmtiMk5yWlhJdFptOXlMV1JsCmMydDBiM0F3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRRFowcjRaTSs3WnFLME8KbVBvN1M5TWZSd1ZyblcxT1VRL0k4c20wVnhNVS9OWVNjQ3A2eGk1ZVdQWGZFR3Fsc0cxV0RHVFF4VmZxSlZrUgpNOFJTODNCbEhuTlZzbTdQNldTNlRoK05VS3MwUldTTGRpeVhMVjBoS3Y2Sm9EVnBjdk1rMFdjT1REUzUvNCtUCkpZMkt2OWVnSWlKMGFnZTR0aHBCaG9LcC96RlhXeDNBeVEzeVVUNUxzeXdhUEJubWo2dHhFKy9CMkVXaU93QTIKZFZ5VjBoYlhXOHF6SmFraTJTNDk5bVluRjBpT0hiRWlhS0JWL2E3K1ExMzVXdkg4T1RpclRlNlY0MUxMOElLdAozb3lUQlpZTXl2VlJQTVllMXZqK2NSVjI3U09rRk9sTHhlRkh0aCtvK2I1UHQ3YmxVNXlySVBqZVRqQXhpZzl2Cm5xd0kyc1J0QWdNQkFBR2pKekFsTUE0R0ExVWREd0VCL3dRRUF3SUZvREFUQmdOVkhTVUVEREFLQmdnckJnRUYKQlFjREFqQU5CZ2txaGtpRzl3MEJBUXNGQUFPQ0FRRUFYdzh1M0pJVFYrTnlIbDc3N1hEVStRQWtML3NPWkhUTQprMVNETThrL2ZnU0tXWHcrWkpQSUVHd3dUMklabXhwVWFpRVJnWmRMalVDY0h1NFZNMmx6bVk3S3cveWxNcXJ5Cmdqby80b3JaaTlPZm9Tdk11N1RKSW4yNkhZOUJ3b004NkZPZUExR3hzeHJjU2ZXcXF3WWlqRzlXSjVnRmRJRVYKZURyU0IvWEg5TTI0QVBLMm9penFDd29RQW9SNzk2TlJNaFJtU1JYYkJyZTdpRDd4WldMQVY0WGt4dTdrdWJRVQpsM3Nrdmg1ZWNpNlRTRmZ2TnBGOC96N0FQZDVkNENtUVAzQ2xwQUt0cWZjMWhYcm9tQnI5aDZHYktIMVQ2MG1LClR2bUgrdXZOSFdKZll1MWErREg1azhTeHJaQnlWVWl2UHJvbG1zZU9aZ29SY2xROHFxbjMwdz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
	clientKeyData            = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBMmRLK0dUUHUyYWl0RHBqNk8wdlRIMGNGYTUxdFRsRVB5UExKdEZjVEZQeldFbkFxCmVzWXVYbGoxM3hCcXBiQnRWZ3hrME1WWDZpVlpFVFBFVXZOd1pSNXpWYkp1eitsa3VrNGZqVkNyTkVWa2kzWXMKbHkxZElTcitpYUExYVhMekpORm5Ea3cwdWYrUGt5V05pci9Yb0NJaWRHb0h1TFlhUVlhQ3FmOHhWMXNkd01rTgo4bEUrUzdNc0dqd1o1bytyY1JQdndkaEZvanNBTm5WY2xkSVcxMXZLc3lXcEl0a3VQZlptSnhkSWpoMnhJbWlnClZmMnUva05kK1ZyeC9EazRxMDN1bGVOU3kvQ0NyZDZNa3dXV0RNcjFVVHpHSHRiNC9uRVZkdTBqcEJUcFM4WGgKUjdZZnFQbStUN2UyNVZPY3F5RDQzazR3TVlvUGI1NnNDTnJFYlFJREFRQUJBb0lCQVFDd2NMS25lNWU0RzVmQwp3VXZBaUZVRmV1UDdIZFFTb2pybERUVXhyVzY1aTJ4a0Z4Tnh5K0ozYmh3TGlpSzQvOFl3ODIwZVp4d2xnQWM3CmxmRXJPQ0lNNXJPSjhUVXRtT0tNdndkejBxbzdkeEpRblhMVktsdkxxQ1h5bVNGcXYwQVF6TVpyb0hMOVR2T2YKdjhWOVpKUm5FLzlReVNwa0JxNFI4Y0VacnFyckZGWVMzUGlOcUt3ZSs5SkpFQkNDNGdhNEhlcUtJTUFEMTBJQwpkSmpwc05kTEVBZ05mcWhIRXF4MmdRZkNyaEluYzErbGpKQ2t2SXE4SVNrN0V1VUtTQnIzWldzVlJpa1llZzF6CjdpZGJNL2ZsclRUOHJDdjQvZCtNakJQeUg2ODB3ZUcwOUZOYmZPVlVMYUZJM1NnbzNTZU1SN1pkU2k4SVBtNnEKb0FZL1MzQkpBb0dCQVB2bzFCUmdKZDAzMGRtVGUvLzdwa0dkemxzV01ZcU9TNEJnNWFVaWR2cEFxTit5MU43Lwo4WC9CNXU4K3FnRXFPaFdwcWl1RDFLdzZIclFWb21nc28yYUtLOXl0amhRSnVra0Vqck93QmpJNmwyelZYWTVOCjhUZkV2S2JuOWdxUG1Qc3BqemtxZk45TjhQRm5ORDVDNDQ4aFdmVWJLUUg0d21Wc1RrYzRNOHZ2QW9HQkFOMWMKT0QxeWtVSHRSY2tPU1R0bVNOY1pieUwvWHZheTMrQW9VZU1tUkV3Ujg3NlVkSVkvbCs1N09VaGlGaFlsZ2wrZgpKZlAyQ3VJYkhEZktrUm01VE9hUnFua1gyTGhsRmlDL3hzQTM5aFpNYWZsMEFZTDJUT3Y4OUZoWWZsdmVISTBmCmVXZ0duWHd2c2l5UDhhTlU1YnhoYXZQaDJmV3puaVlZenczR1VJbGpBb0dCQUpxcWN6d25aRmdVbzZXQTVndGYKU003VXd4UXQ1akQ1K09WVG1Pdnc0emYzWnpnaHBvMWNlN3hESXpVVk5uQU1xeWtOeU10NjM2TDI5RThCVVh2QQpuSHFpcnhlSS84alJ5Q3g0dmQwWllGU0tvTnBUam1PRysybVFRM3YrdzVQc1lyTWk2eXJnS1ZjNWxZdkNIOU55CllpRkpDdUJ3MHBiQlV5bk9lNmxDbXByRkFvR0FaNWgwSFpYVGswWUhCZHd3dTZMWDdnclNMMC9TVXFST0QrcnUKdjk2MTBlQUk4YVVxajNXTmxpZUhISEFESkNRenlxcUJxWlg1YSs0Q0c3NHFnQjV2ZGV5d0duSGxsQTZDOTVxbwoybWFXUGNOUFNWWllLc0U0S0swK2NXWWc1TCtqSHd5ZTFlZGFwcmJ2ck5sRTNncXFaYW8xMHFOZFRSYmRxYmlzCmxpYThwRDhDZ1lFQXBpVkE0WUUwRjBKZmlDTHFqSDZWTGZXLzNoVU9Sc015UVk4QkZOOUtuN3hMNWNuVGdsL2MKQTB4RWN1RXF6Vi85WWhha0ZMWE1XQWFJeW04TjlKQ3NGU0RYdXhqa2xLakdYWmJsa0xDc1hqYlVKU0IyWUJYQQpuWUk2YllHZjZ5NGJ3TnpoN0JHdDNMYXU1UXNGcXAyNzBHTmN4eXUzeVVZSWg3c0tUY1I0KzJrPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
	repoName                 = "stable"
	repoURL                  = "https://cetic.github.io/helm-charts"
)

func TestHelm3Basic(t *testing.T) {
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		AuthType: helm.AuthTypeBasic,
		Cluster: &helm.Cluster{
			Server:                server,
			InsecureSkipTLSVerify: true,
		},
		AuthInfo: &helm.AuthInfo{
			Username: userName,
			Password: password,
		},
		RepoName: repoName,
		RepoURL:  repoURL,
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

func TestVersionBasic(t *testing.T) {
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		AuthType: helm.AuthTypeBasic,
		Cluster: &helm.Cluster{
			Server:                server,
			InsecureSkipTLSVerify: true,
		},
		AuthInfo: &helm.AuthInfo{
			Username: userName,
			Password: password,
		},
		RepoName: repoName,
		RepoURL:  repoURL,
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

func TestHelm3Token(t *testing.T) {
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		AuthType: helm.AuthTypeToken,
		Cluster: &helm.Cluster{
			Server:                server,
			InsecureSkipTLSVerify: true,
		},
		AuthInfo: &helm.AuthInfo{
			Token: token,
		},
		RepoName: repoName,
		RepoURL:  repoURL,
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

func TestVersionToken(t *testing.T) {
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		AuthType: helm.AuthTypeToken,
		Cluster: &helm.Cluster{
			Server:                server,
			InsecureSkipTLSVerify: true,
		},
		AuthInfo: &helm.AuthInfo{
			Token: token,
		},
		RepoName: repoName,
		RepoURL:  repoURL,
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

func TestHelm3Cert(t *testing.T) {
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		AuthType: helm.AuthTypeCert,
		Cluster: &helm.Cluster{
			Server:                   server,
			InsecureSkipTLSVerify:    true,
			CertificateAuthorityData: []byte(certificateAuthorityData),
		},
		AuthInfo: &helm.AuthInfo{
			ClientCertificateData: []byte(clientCertificateData),
			ClientKeyData:         []byte(clientKeyData),
		},
		RepoName: repoName,
		RepoURL:  repoURL,
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

func TestVersionCert(t *testing.T) {
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		AuthType: helm.AuthTypeCert,
		Cluster: &helm.Cluster{
			Server:                   server,
			InsecureSkipTLSVerify:    true,
			CertificateAuthorityData: []byte(certificateAuthorityData),
		},
		AuthInfo: &helm.AuthInfo{
			ClientCertificateData: []byte(clientCertificateData),
			ClientKeyData:         []byte(clientKeyData),
		},
		RepoName: repoName,
		RepoURL:  repoURL,
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

func TestListBasic(t *testing.T) {
	if env != runLocal {
		return
	}

	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		AuthType: helm.AuthTypeBasic,
		Cluster: &helm.Cluster{
			Server:                server,
			InsecureSkipTLSVerify: true,
		},
		AuthInfo: &helm.AuthInfo{
			Username: userName,
			Password: password,
		},
		RepoName: repoName,
		RepoURL:  repoURL,
	})

	if err != nil {
		t.Fatalf("create new helm error, err:%v", err)
	}

	items, err := h3.List([]string{})
	if err != nil {
		t.Fatalf("exec list command failure, err:%v", err)
	}

	t.Logf("items:%v", items)
}

func TestListToken(t *testing.T) {
	if env != runLocal {
		return
	}
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		AuthType: helm.AuthTypeToken,
		Cluster: &helm.Cluster{
			Server:                server,
			InsecureSkipTLSVerify: true,
		},
		AuthInfo: &helm.AuthInfo{
			Token:    token,
			Username: "ttt",
		},
		RepoName: repoName,
		RepoURL:  repoURL,
	})

	if err != nil {
		t.Fatalf("create new helm error, err:%v", err)
	}

	items, err := h3.List([]string{})
	if err != nil {
		t.Fatalf("exec list command failure, err:%v", err)
	}

	t.Logf("items:%v", items)
}

func TestListCert(t *testing.T) {
	if env != runLocal {
		return
	}
	h3, err := helm.NewHelm(helm.TypeHelm3, helm.Config{
		AuthType: helm.AuthTypeCert,
		Cluster: &helm.Cluster{
			Server:                   server,
			InsecureSkipTLSVerify:    true,
			CertificateAuthorityData: []byte(certificateAuthorityData),
		},
		AuthInfo: &helm.AuthInfo{
			ClientCertificateData: []byte(clientCertificateData),
			ClientKeyData:         []byte(clientKeyData),
		},
		RepoName: repoName,
		RepoURL:  repoURL,
	})

	if err != nil {
		t.Fatalf("create new helm error, err:%v", err)
	}

	items, err := h3.List([]string{})
	if err != nil {
		t.Fatalf("exec list command failure, err:%v", err)
	}

	t.Logf("items:%v", items)
}
