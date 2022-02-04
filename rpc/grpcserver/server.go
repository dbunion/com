package grpcserver

import (
	"context"
	"fmt"
	"math"
	"net"
	"time"

	"github.com/dbunion/com/rpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// DefaultConfig - rpc default config
var DefaultConfig rpc.Config = rpc.Config{
	MaxMessageSize:            rpc.DefaultMaxMessageSize,
	GRPCPort:                  8264,
	GRPCInitialConnWindowSize: 0,
	GRPCInitialWindowSize:     0,
	EnableTracing:             false,
	EnableGRPCPrometheus:      false,
	GRPCKeepAliveEnforcementPolicyPermitWithoutStream: false,
	GRPCCert:                              "",
	GRPCKey:                               "",
	GRPCCA:                                "",
	GRPCAuth:                              "",
	GRPCMaxConnectionAge:                  time.Duration(math.MaxInt64),
	GRPCMaxConnectionAgeGrace:             time.Duration(math.MaxInt64),
	GRPCKeepAliveEnforcementPolicyMinTime: 5 * time.Minute,
}

var authPlugin Authenticator

// Server - rpc server
type Server struct {
	*grpc.Server
	cfg *rpc.Config
}

// NewRPCServer create the gRPC server we will be using.
// It has to be called after flags are parsed, but before
// services register themselves.
func NewRPCServer(cfg *rpc.Config, outOpts ...grpc.ServerOption) (*Server, error) {
	// skip if not registered
	if !cfg.IsGRPCEnabled() {
		return nil, fmt.Errorf("skipping gRPC server creation")
	}

	var opts []grpc.ServerOption
	if cfg.GRPCPort != 0 && cfg.GRPCCert != "" && cfg.GRPCKey != "" {
		config, err := rpc.ServerConfig(cfg.GRPCCert, cfg.GRPCKey, cfg.GRPCCA)
		if err != nil {
			return nil, fmt.Errorf("failed to log gRPC cert/key/ca: %v", err)
		}

		// create the creds server options
		creds := credentials.NewTLS(config)
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	// Override the default max message size for both send and receive
	// (which is 4 MiB in gRPC 1.0.0).
	// Large messages can occur when users try to insert or fetch very big
	// rows. If they hit the limit, they'll see the following error:
	// grpc: received message length XXXXXXX exceeding the max size 4194304
	// Note: For gRPC 1.0.0 it's sufficient to set the limit on the server only
	// because it's not enforced on the client side.
	opts = append(opts, grpc.MaxRecvMsgSize(int(cfg.MaxMessageSize)))
	opts = append(opts, grpc.MaxSendMsgSize(int(cfg.MaxMessageSize)))

	if cfg.GRPCInitialConnWindowSize != 0 {
		opts = append(opts, grpc.InitialConnWindowSize(cfg.GRPCInitialConnWindowSize))
	}

	if cfg.GRPCInitialWindowSize != 0 {
		opts = append(opts, grpc.InitialWindowSize(cfg.GRPCInitialWindowSize))
	}

	ep := keepalive.EnforcementPolicy{
		MinTime:             cfg.GRPCKeepAliveEnforcementPolicyMinTime,
		PermitWithoutStream: cfg.GRPCKeepAliveEnforcementPolicyPermitWithoutStream,
	}
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(ep))

	if cfg.GRPCMaxConnectionAge != 0 {
		ka := keepalive.ServerParameters{
			MaxConnectionAge: cfg.GRPCMaxConnectionAge,
		}
		if cfg.GRPCMaxConnectionAgeGrace != 0 {
			ka.MaxConnectionAgeGrace = cfg.GRPCMaxConnectionAgeGrace
		}
		opts = append(opts, grpc.KeepaliveParams(ka))
	}

	opts = append(opts, interceptors(cfg)...)

	if len(outOpts) > 0 {
		opts = append(opts, outOpts...)
	}

	return &Server{grpc.NewServer(opts...), cfg}, nil
}

// We can only set a ServerInterceptor once, so we chain multiple interceptors into one
func interceptors(cfg *rpc.Config) []grpc.ServerOption {
	interceptors := &InterceptorBuilder{}

	if cfg.GRPCAuth != "" {
		pluginInitializer, err := GetAuthenticator(cfg.GRPCAuth)
		if err != nil {
			fmt.Printf("GetAuthenticator err:%v\n", err)
		}
		authPluginImpl, err := pluginInitializer([]byte(cfg.GRPCAuthStaticPassword))
		if err != nil {
			fmt.Printf("Failed to load auth plugin: %v\n", err)
		}
		authPlugin = authPluginImpl
		interceptors.Add(authenticatingStreamInterceptor, authenticatingUnaryInterceptor)
	}

	if interceptors.NonEmpty() {
		return []grpc.ServerOption{
			grpc.StreamInterceptor(interceptors.StreamServerInterceptor),
			grpc.UnaryInterceptor(interceptors.UnaryStreamInterceptor)}
	}
	return []grpc.ServerOption{}
}

// Run - run rpc server
func (s *Server) Run() error {
	if s.cfg.EnableGRPCPrometheus {
		grpc_prometheus.Register(s.Server)
		grpc_prometheus.EnableHandlingTimeHistogram()
	}
	// skip if not registered
	if s.cfg.GRPCPort == 0 {
		return nil
	}

	// listen on the port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("cannot listen on port %v for gRPC: %v", s.cfg.GRPCPort, err)
	}

	// and serve on it
	go func() {
		panic(s.Server.Serve(listener))
	}()
	return nil
}

// Close - close rpc
func (s *Server) Close() {
	s.Server.GracefulStop()
}

func authenticatingStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	newCtx, err := authPlugin.Authenticate(stream.Context(), info.FullMethod)

	if err != nil {
		return err
	}

	wrapped := WrapServerStream(stream)
	wrapped.WrappedContext = newCtx
	return handler(srv, wrapped)
}

func authenticatingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	newCtx, err := authPlugin.Authenticate(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}

	return handler(newCtx, req)
}

// WrappedServerStream is based on the service stream wrapper from: https://github.com/grpc-ecosystem/go-grpc-middleware
type WrappedServerStream struct {
	grpc.ServerStream
	WrappedContext context.Context
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream returns a ServerStream that has the ability to overwrite context.
func WrapServerStream(stream grpc.ServerStream) *WrappedServerStream {
	if existing, ok := stream.(*WrappedServerStream); ok {
		return existing
	}
	return &WrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}

// InterceptorBuilder chains together multiple ServerInterceptors
type InterceptorBuilder struct {
	StreamServerInterceptor grpc.StreamServerInterceptor
	UnaryStreamInterceptor  grpc.UnaryServerInterceptor
}

// Add is used to add two ServerInterceptors.
func (collector *InterceptorBuilder) Add(s grpc.StreamServerInterceptor, u grpc.UnaryServerInterceptor) {
	if collector.StreamServerInterceptor == nil {
		collector.StreamServerInterceptor = s
		collector.UnaryStreamInterceptor = u
	} else {
		collector.StreamServerInterceptor = grpc_middleware.ChainStreamServer(collector.StreamServerInterceptor, s)
		collector.UnaryStreamInterceptor = grpc_middleware.ChainUnaryServer(collector.UnaryStreamInterceptor, u)
	}
}

// NonEmpty check if ServerInterceptor is empty.
func (collector *InterceptorBuilder) NonEmpty() bool {
	return collector.StreamServerInterceptor != nil
}
