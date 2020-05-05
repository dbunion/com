package v3

import (
	"context"
	"github.com/dbunion/com/vtctl"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	client, err := vtctl.NewClient(vtctl.TypeVtctlV3, vtctl.Config{
		RetryOption: nil,
		Servers: map[string]vtctl.ServerInfo{
			"test": {
				IP:   "127.0.0.1",
				Port: 15999,
			},
		},
		HealthCheck: 0,
	})

	if err != nil {
		t.Fatalf("create new vtctl error, err:%v", err)
	}

	values, err := client.RunCommand(context.Background(), []string{"ListAllTablets"}, time.Second*120)
	assert.NotNil(t, err)
	assert.Empty(t, values)
	t.Logf("values:%+v", values)
}
