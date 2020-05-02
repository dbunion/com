package rpc

import "time"

const (
	// DefaultMaxMessageSize - Default Max message size
	DefaultMaxMessageSize = 512 * 1024 * 1024
)

// Config - rpc base config
type Config struct {
	// Maximum allowed RPC message size. Larger messages will be rejected by gRPC with the error 'exceeding the max size'.")
	MaxMessageSize int32

	// GRPCPort is the port to listen on for gRPC. If not set or zero, don't listen.
	GRPCPort int32

	// GRPCInitialConnWindowSize ServerOption that sets window size for a connection.
	// The lower bound for window size is 64K and any value smaller than that will be ignored.
	GRPCInitialConnWindowSize int32

	// GRPCInitialWindowSize ServerOption that sets window size for stream.
	// The lower bound for window size is 64K and any value smaller than that will be ignored.
	GRPCInitialWindowSize int32

	// EnableTracing sets a flag to enable grpc client/server tracing.
	EnableTracing bool

	// EnableGRPCPrometheus sets a flag to enable grpc client/server grpc monitoring.
	EnableGRPCPrometheus bool

	// EnforcementPolicy PermitWithoutStream - If true, server allows keepalive pings
	// even when there are no active streams (RPCs). If false, and client sends ping when
	// there are no active streams, server will send GOAWAY and close the connection.
	GRPCKeepAliveEnforcementPolicyPermitWithoutStream bool

	// GRPCCert is the cert to use if TLS is enabled
	GRPCCert string

	// GRPCKey is the key to use if TLS is enabled
	GRPCKey string

	// GRPCCA is the CA to use if TLS is enabled
	GRPCCA string

	// GRPCAuth which auth plugin to use (at the moment now only static is supported)
	GRPCAuth string

	// the server name to use to validate server certificate
	GRPCServerName string

	// GRPCAuthStaticPassword JSON File to read the users/passwords from
	GRPCAuthStaticPassword string

	// GRPCMaxConnectionAge is the maximum age of a client connection, before GoAway is sent.
	// This is useful for L4 loadbalancing to ensure rebalancing after scaling.
	GRPCMaxConnectionAge time.Duration

	// GRPCMaxConnectionAgeGrace is an additional grace period after GRPCMaxConnectionAge, after which
	// connections are forcibly closed.
	GRPCMaxConnectionAgeGrace time.Duration

	// EnforcementPolicy MinTime that sets the keepalive enforcement policy on the server.
	// This is the minimum amount of time a client should wait before sending a keepalive ping.
	GRPCKeepAliveEnforcementPolicyMinTime time.Duration

	// After a duration of this time if the client doesn't see any activity it pings the server to see if the transport is still alive.
	GRPCKeepAliveTime time.Duration

	// After having pinged for keepalive check, the client waits for a duration of Timeout and if no activity is seen even after that the connection is closed.
	GRPCKeepAliveTimeout time.Duration
}

// IsGRPCEnabled returns true if gRPC server is set
func (c *Config) IsGRPCEnabled() bool {
	return c.GRPCPort != 0
}
