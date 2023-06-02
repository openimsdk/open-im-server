package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
)

type BlackClient struct {
	*MetaClient
}

func NewBlackClient(zk discoveryRegistry.SvcDiscoveryRegistry) *BlackClient {
	return &BlackClient{NewMetaClient(zk, config.Config.RpcRegisterName.OpenImFriendName)}
}

// possibleBlackUserID是否被userID拉黑，也就是是否在userID的黑名单中
func (b *BlackClient) IsBlocked(ctx context.Context, possibleBlackUserID, userID string) (bool, error) {
	cc, err := b.getConn(ctx)
	if err != nil {
		return false, err
	}
	r, err := friend.NewFriendClient(cc).IsBlack(ctx, &friend.IsBlackReq{UserID1: possibleBlackUserID, UserID2: userID})
	if err != nil {
		return false, err
	}
	return r.InUser2Blacks, nil
}
