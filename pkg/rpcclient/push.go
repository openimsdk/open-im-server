package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/push"
	"google.golang.org/grpc"
)

type PushClient struct {
	conn *grpc.ClientConn
}

func NewPushClient(discov discoveryregistry.SvcDiscoveryRegistry) *PushClient {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImPushName)
	if err != nil {
		panic(err)
	}
	return &PushClient{conn: conn}
}

func (p *PushClient) DelUserPushToken(ctx context.Context, req *push.DelUserPushTokenReq) (*push.DelUserPushTokenResp, error) {
	resp, err := push.NewPushMsgServiceClient(p.conn).DelUserPushToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
