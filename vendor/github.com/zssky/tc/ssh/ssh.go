package ssh

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

// SShClient - an SSh connection client
type SShClient struct {
	host   string
	port   int
	user   string
	passwd string
	config ssh.ClientConfig
}

// NewSShClient - create an new SShClient
func NewSShClient(host string, port int, user, passwd string, timeout time.Duration) *SShClient {
	return &SShClient{
		host:   host,
		port:   port,
		user:   user,
		passwd: passwd,
		config: ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(passwd),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: timeout,
		},
	}
}

// ExecCmd - exec shell command in remove server
func (s *SShClient) ExecCmd(cmd string) ([]byte, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), &s.config)
	if err != nil {
		return nil, fmt.Errorf("Cannot build SSH connection to database server, for %s", err)
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Cannot open session in SSH connection to database, err:%v", err.Error())
	}

	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var buff bytes.Buffer
	session.Stdout = &buff
	session.Stderr = &buff
	if err := session.Run(cmd); err != nil {
		return nil, fmt.Errorf("run cmd err:%v cmd:%v, output:%v", err, cmd, buff.String())
	}

	return buff.Bytes(), nil
}
