package grpcserver

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	// StaticAuthPlugin implements AuthPlugin interface
	_ Authenticator = (*StaticAuthPlugin)(nil)
)

// StaticAuthConfigEntry is the container for server side credentials. Current implementation matches the
// the one from the client but this will change in the future as we hooked this pluging into ACL
// features.
type StaticAuthConfigEntry struct {
	Username string
	Password string
	// TODO (@rafael) Add authorization parameters
}

// StaticAuthPlugin  implements static username/password authentication for grpc. It contains an array of username/passwords
// that will be authorized to connect to the grpc server.
type StaticAuthPlugin struct {
	entries []StaticAuthConfigEntry
}

// Authenticate implements AuthPlugin interface. This method will be used inside a middleware in grpc_server to authenticate
// incoming requests.
func (sa *StaticAuthPlugin) Authenticate(ctx context.Context, fullMethod string) (context.Context, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md["username"]) == 0 || len(md["password"]) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "username and password must be provided")
		}
		username := md["username"][0]
		password := md["password"][0]
		for _, authEntry := range sa.entries {
			if username == authEntry.Username && password == authEntry.Password {
				return ctx, nil
			}
		}
		return nil, status.Errorf(codes.PermissionDenied, "auth failure: caller %q provided invalid credentials", username)
	}
	return nil, status.Errorf(codes.Unauthenticated, "username and password must be provided")
}

func staticAuthPluginInitializer(grpcAuthStaticPassword []byte) (Authenticator, error) {
	entries := make([]StaticAuthConfigEntry, 0)
	if grpcAuthStaticPassword == nil {
		err := fmt.Errorf("failed to load static auth plugin. Plugin configured but grpc_auth_static_password_file not provided")
		return nil, err
	}

	if err := json.Unmarshal(grpcAuthStaticPassword, &entries); err != nil {
		err := fmt.Errorf("fail to load static auth plugin: %v", err)
		return nil, err
	}

	staticAuthPlugin := &StaticAuthPlugin{
		entries: entries,
	}
	return staticAuthPlugin, nil
}

func init() {
	if err := RegisterAuthPlugin("static", staticAuthPluginInitializer); err != nil {
		panic(err)
	}
}
