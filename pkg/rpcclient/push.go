package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/push"
	"google.golang.org/grpc"
)

type Push struct {
	conn   grpc.ClientConnInterface
	Client push.PushMsgServiceClient
	discov discoveryregistry.SvcDiscoveryRegistry
}

func NewPush(discov discoveryregistry.SvcDiscoveryRegistry) *Push {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImPushName)
	if err != nil {
		panic(err)
	}
	return &Push{
		discov: discov,
		conn:   conn,
		Client: push.NewPushMsgServiceClient(conn),
	}
}

type PushRpcClient Push

func NewPushRpcClient(discov discoveryregistry.SvcDiscoveryRegistry) PushRpcClient {
	return PushRpcClient(*NewPush(discov))
}

func (p *PushRpcClient) DelUserPushToken(ctx context.Context, req *push.DelUserPushTokenReq) (*push.DelUserPushTokenResp, error) {
	return p.Client.DelUserPushToken(ctx, req)
}
