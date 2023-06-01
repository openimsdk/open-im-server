package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/push"
)

type PushClient struct {
	MetaClient
}

func NewPushClient(client discoveryregistry.SvcDiscoveryRegistry) *PushClient {
	return &PushClient{
		MetaClient: MetaClient{
			client:          client,
			rpcRegisterName: config.Config.RpcRegisterName.OpenImPushName,
		},
	}
}

func (p *PushClient) DelUserPushToken(ctx context.Context, req *push.DelUserPushTokenReq) (*push.DelUserPushTokenResp, error) {
	cc, err := p.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := push.NewPushMsgServiceClient(cc).DelUserPushToken(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
