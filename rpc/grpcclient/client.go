package grpcclient

import (
	"fmt"
	"time"

	"github.com/dbunion/com/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

// grpcDialOptions is a registry of functions that append grpcDialOption to use when dialing a service
var grpcDialOptions []func(opts []grpc.DialOption, grpcAuthStaticPassword []byte) ([]grpc.DialOption, error)

// RegisterGRPCDialOptions registers an implementation of AuthServer.
func RegisterGRPCDialOptions(grpcDialOptionsFunc func(opts []grpc.DialOption, grpcAuthStaticPassword []byte) ([]grpc.DialOption, error)) {
	grpcDialOptions = append(grpcDialOptions, grpcDialOptionsFunc)
}

// DefaultConfig - grpc client dial default config
var DefaultConfig rpc.Config = rpc.Config{
	MaxMessageSize:            rpc.DefaultMaxMessageSize,
	GRPCKeepAliveTime:         10 * time.Second,
	GRPCKeepAliveTimeout:      10 * time.Second,
	GRPCInitialConnWindowSize: 0,
	GRPCInitialWindowSize:     0,
}

// Conn - grpc conn
type Conn struct {
	*grpc.ClientConn
	cfg rpc.Config
}

// NewConn - create new grpc client conn
func NewConn(target string, cfg *rpc.Config) (*Conn, error) {
	conn := &Conn{
		cfg: *cfg,
	}

	return conn.Dial(target)
}

// Close - release conn
func (c *Conn) Close() error {
	return c.ClientConn.Close()
}

// Dial creates a grpc connection to the given target.
// failFast is a non-optional parameter because callers are required to specify
// what that should be.
func (c *Conn) Dial(target string) (*Conn, error) {
	newOpts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(int(c.cfg.MaxMessageSize)),
			grpc.MaxCallSendMsgSize(int(c.cfg.MaxMessageSize)),
		),
	}

	if c.cfg.GRPCKeepAliveTime != 0 || c.cfg.GRPCKeepAliveTimeout != 0 {
		kp := keepalive.ClientParameters{
			// After a duration of this time if the client doesn't see any activity it pings the server to see if the transport is still alive.
			Time: c.cfg.GRPCKeepAliveTime,
			// After having pinged for keepalive check, the client waits for a duration of Timeout and if no activity is seen even after that
			// the connection is closed. (This will eagerly fail inflight grpc requests even if they don't have timeouts.)
			Timeout:             c.cfg.GRPCKeepAliveTimeout,
			PermitWithoutStream: true,
		}
		newOpts = append(newOpts, grpc.WithKeepaliveParams(kp))
	}

	if c.cfg.GRPCInitialConnWindowSize != 0 {
		newOpts = append(newOpts, grpc.WithInitialConnWindowSize(c.cfg.GRPCInitialConnWindowSize))
	}

	if c.cfg.GRPCInitialWindowSize != 0 {
		newOpts = append(newOpts, grpc.WithInitialWindowSize(c.cfg.GRPCInitialWindowSize))
	}

	var err error
	for _, grpcDialOptionInitializer := range grpcDialOptions {
		newOpts, err = grpcDialOptionInitializer(newOpts, []byte(c.cfg.GRPCAuthStaticPassword))
		if err != nil {
			return nil, fmt.Errorf("there was an error initializing client grpc.DialOption: %v", err)
		}
	}

	if c.cfg.EnableGRPCPrometheus {
		newOpts = append(newOpts, grpc.WithUnaryInterceptor(grpcPrometheus.UnaryClientInterceptor))
		newOpts = append(newOpts, grpc.WithStreamInterceptor(grpcPrometheus.StreamClientInterceptor))
	}

	// secure dial opt init
	secOpt, err := secureDialOption(c.cfg)
	if err == nil {
		newOpts = append(newOpts, secOpt)
	}

	c.ClientConn, err = grpc.Dial(target, newOpts...)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// secureDialOption returns the gRPC dial option to use for the
// given client connection. It is either using TLS, or Insecure if
// nothing is set.
func secureDialOption(cfg rpc.Config) (grpc.DialOption, error) {
	// No security options set, just return.
	if (cfg.GRPCCert == "" || cfg.GRPCKey == "") && cfg.GRPCCA == "" {
		return grpc.WithInsecure(), nil
	}

	// Load the config.
	config, err := rpc.ClientConfig(cfg.GRPCCert, cfg.GRPCKey, cfg.GRPCCA, cfg.GRPCServerName)
	if err != nil {
		return nil, err
	}

	// Create the creds server options.
	creds := credentials.NewTLS(config)
	return grpc.WithTransportCredentials(creds), nil
}
