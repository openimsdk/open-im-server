package rpcclient

import (
	"context"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
)

type Third struct {
	conn   grpc.ClientConnInterface
	Client third.ThirdClient
	discov discoveryregistry.SvcDiscoveryRegistry
}

func NewThird(discov discoveryregistry.SvcDiscoveryRegistry) *Third {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImThirdName)
	if err != nil {
		panic(err)
	}
	client := third.NewThirdClient(conn)
	return &Third{discov: discov, Client: client, conn: conn}
}
