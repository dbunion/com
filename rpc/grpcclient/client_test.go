package grpcclient

import (
	"fmt"
	"github.com/dbunion/com/rpc/grpcserver"
	"testing"
)

func TestClientConn(t *testing.T) {
	srv, err := grpcserver.NewRPCServer(&grpcserver.DefaultConfig)
	if err != nil {
		t.Fatalf("create new rpc server error, err:%v", err)
	}

	go func() {
		if err := srv.Run(); err != nil {
			panic(err)
		}
	}()

	target := fmt.Sprintf("%s:%d", "127.0.0.1", grpcserver.DefaultConfig.GRPCPort)
	client, err := NewConn(target, &DefaultConfig)
	if err != nil {
		t.Fatalf("create new conn err:%v", err)
	}

	stat := client.GetState()

	t.Logf("current stat:%v", stat)
	if err := client.Close(); err != nil {
		t.Fatalf("close conn err:%v", err)
	}

	t.Log("test client conn success")
}
