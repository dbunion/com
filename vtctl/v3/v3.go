package v3

import (
	"context"
	"fmt"
	"github.com/zssky/tc/retry"
	"io"
	"time"

	"github.com/dbunion/com/rpc/grpcclient"
	"github.com/dbunion/com/vtctl"
	logutilpb "vitess.io/vitess/go/vt/proto/logutil"
	vtctldatapb "vitess.io/vitess/go/vt/proto/vtctldata"
	vtctlservicepb "vitess.io/vitess/go/vt/proto/vtctlservice"
)

var defaultRetryOption = &retry.Options{
	InitialBackoff: time.Second * 1,
	MaxBackoff:     time.Second * 3,
	Multiplier:     2,
	MaxRetries:     3,
}

type client struct {
	config      vtctl.Config
	cc          []*grpcclient.Conn
	c           []vtctlservicepb.VtctlClient
	retryOption *retry.Options
	done        bool
	cWeight     map[int]int    // client index => weight
	cMapping    map[string]int // server key => client index
}

// NewClient - create new vtctl client
func NewClient() vtctl.Client {
	return &client{
		cc:       make([]*grpcclient.Conn, 0),
		c:        make([]vtctlservicepb.VtctlClient, 0),
		cWeight:  make(map[int]int),
		cMapping: make(map[string]int),
	}
}

func (c *client) formatServer(srv vtctl.ServerInfo) string {
	return fmt.Sprintf("%v:%v", srv.IP, srv.Port)
}

func (c *client) withRetryFunc(ctx context.Context, action string, fun func() error) error {
	if c.done {
		return fmt.Errorf("context is done, action:%v", action)
	}

	retryAttempts := 0
	for r := retry.Start(*c.retryOption); r.Next(); retryAttempts++ {
		if err := fun(); err != nil {
			continue
		}
		break
	}
	if retryAttempts == c.retryOption.MaxRetries+1 {
		return fmt.Errorf("with retry func reatch max retry count, action:%v", action)
	}

	return nil
}

// run vtctl command
func (c *client) RunCommand(ctx context.Context, args []string, timeout time.Duration) ([]string, error) {
	query := &vtctldatapb.ExecuteVtctlCommandRequest{
		Args:          args,
		ActionTimeout: int64(timeout.Nanoseconds()),
	}

	var values []string
	if err := c.withRetryFunc(ctx, "RunCommand", func() error {
		values = make([]string, 0)

		// choose health client
		client := c.c[0]

		stream, err := client.ExecuteVtctlCommand(ctx, query)
		if err != nil {
			return err
		}

		streamAdapter := &eventStreamAdapter{stream}
		recv := func(e *logutilpb.Event) {
			values = append(values, e.Value)
		}

		// stream the result
		for {
			e, err := streamAdapter.Recv()
			switch err {
			case nil:
				recv(e)
			case io.EOF:
				return nil
			default:
				return fmt.Errorf("remote error: %v", err)
			}
		}

	}); err != nil {
		return nil, err
	}
	return values, nil
}

// close connection
func (c *client) Close() error {
	c.done = true
	return nil
}

// start gc routine based on config string settings.
func (c *client) StartAndGC(config vtctl.Config) error {
	if len(config.Servers) <= 0 {
		return fmt.Errorf("invalid server param, server can not be empty")
	}

	c.config = config
	index := 0
	for key := range config.Servers {
		srv := config.Servers[key]

		cc, err := grpcclient.NewConn(c.formatServer(srv), &grpcclient.DefaultConfig)
		if err != nil {
			return err
		}

		c.retryOption = config.RetryOption
		if config.RetryOption == nil {
			c.retryOption = defaultRetryOption
		}

		c.cc = append(c.cc, cc)
		c.c = append(c.c, vtctlservicepb.NewVtctlClient(cc.ClientConn))

		c.cWeight[index] = 1
		c.cMapping[key] = index

		index++
	}
	return nil
}

type eventStreamAdapter struct {
	stream vtctlservicepb.Vtctl_ExecuteVtctlCommandClient
}

func (e *eventStreamAdapter) Recv() (*logutilpb.Event, error) {
	le, err := e.stream.Recv()
	if err != nil {
		return nil, err
	}
	return le.Event, nil
}

func init() {
	vtctl.Register(vtctl.TypeVtctlV3, NewClient)
}
