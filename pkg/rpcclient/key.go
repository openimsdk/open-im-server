package rpcclient

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/key"
	"google.golang.org/grpc"
)

func NewKey(discov discoveryregistry.SvcDiscoveryRegistry) *Key {
	conn, err := discov.GetConn(context.Background(), "key")
	if err != nil {
		panic(err)
	}
	client := key.NewKeyClient(conn)
	return &Key{discov: discov, conn: conn, Client: client}
}

type Key struct {
	conn   grpc.ClientConnInterface
	Client key.KeyClient
	discov discoveryregistry.SvcDiscoveryRegistry
}
