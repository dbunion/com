package grpcclient

import (
	"encoding/json"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	// StaticAuthClientCreds implements client interface to be able to WithPerRPCCredentials
	_ credentials.PerRPCCredentials = (*StaticAuthClientCreds)(nil)
)

// StaticAuthClientCreds holder for client credentials
type StaticAuthClientCreds struct {
	Username string
	Password string
}

// GetRequestMetadata  gets the request metadata as a map from StaticAuthClientCreds
func (c *StaticAuthClientCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"username": c.Username,
		"password": c.Password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
// Given that people can use this with or without TLS, at the moment we are not enforcing
// transport security
func (c *StaticAuthClientCreds) RequireTransportSecurity() bool {
	return false
}

// AppendStaticAuth optionally appends static auth credentials if provided.
func AppendStaticAuth(opts []grpc.DialOption, grpcAuthStaticPassword []byte) ([]grpc.DialOption, error) {
	if len(grpcAuthStaticPassword) == 0 {
		return opts, nil
	}
	clientCreds := &StaticAuthClientCreds{}
	if err := json.Unmarshal(grpcAuthStaticPassword, clientCreds); err != nil {
		return nil, err
	}
	creds := grpc.WithPerRPCCredentials(clientCreds)
	opts = append(opts, creds)
	return opts, nil
}

func init() {
	RegisterGRPCDialOptions(AppendStaticAuth)
}
