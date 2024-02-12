package grpc

import (
	"context"

	peggyTypes "bharvest.io/orchmon/client/grpc/protobuf/peggy"
	"google.golang.org/grpc"
)

type Client interface {
	Connect(ctx context.Context) error
	Terminate(_ context.Context) error
}

type (
	Injective struct {
		host string
		conn *grpc.ClientConn
		peggyClient peggyTypes.QueryClient
	}
)