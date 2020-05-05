package fun

import (
	"context"
	"github.com/zssky/tc/ssh"
	"time"
)

// ExecSSHCommand - exec ssh command
func (m *defaultFuncWrap) ExecSSHCommand(ctx context.Context, host string, port int, username, password string, cmd string) (string, error) {
	results, err := m.wrangler(ctx, func(ctx context.Context) ([]string, error) {
		client := ssh.NewSShClient(host, port, username, password, time.Second*5)
		data, err := client.ExecCmd(cmd)
		if err != nil {
			return nil, err
		}

		return []string{string(data)}, nil
	})

	if err != nil {
		return "", err
	}

	return results[0], nil
}
