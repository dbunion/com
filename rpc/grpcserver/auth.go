package grpcserver

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Authenticator provides an interface to implement auth in Vitess in
// grpc server
type Authenticator interface {
	Authenticate(ctx context.Context, fullMethod string) (context.Context, error)
}

// authPlugins is a registry of AuthPlugin initializers.
var authPlugins = make(map[string]func([]byte) (Authenticator, error))

// RegisterAuthPlugin registers an implementation of AuthServer.
func RegisterAuthPlugin(name string, authPlugin func([]byte) (Authenticator, error)) error {
	if _, ok := authPlugins[name]; ok {
		return fmt.Errorf("AuthPlugin named %v already exists", name)
	}
	authPlugins[name] = authPlugin
	return nil
}

// GetAuthenticator returns an AuthPlugin by name, or log.Fatalf.
func GetAuthenticator(name string) (func([]byte) (Authenticator, error), error) {
	authPlugin, ok := authPlugins[name]
	if !ok {
		return nil, fmt.Errorf("no AuthPlugin name %v registered", name)
	}
	return authPlugin, nil
}

// FakeAuthStreamInterceptor fake interceptor to test plugin
func FakeAuthStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if fakeDummyAuthenticate(stream.Context()) {
		return handler(srv, stream)
	}
	return status.Errorf(codes.Unauthenticated, "username and password must be provided")
}

// FakeAuthUnaryInterceptor fake interceptor to test plugin
func FakeAuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if fakeDummyAuthenticate(ctx) {
		return handler(ctx, req)
	}
	return nil, status.Errorf(codes.Unauthenticated, "username and password must be provided")
}

func fakeDummyAuthenticate(ctx context.Context) bool {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md["username"]) == 0 || len(md["password"]) == 0 {
			return false
		}
		username := md["username"][0]
		password := md["password"][0]
		if username == "valid" && password == "valid" {
			return true
		}
		return false
	}
	return false
}
