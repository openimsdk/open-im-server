package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"google.golang.org/grpc"
)

type BlackClient struct {
	conn *grpc.ClientConn
}

func NewBlackClient(discov discoveryRegistry.SvcDiscoveryRegistry) *BlackClient {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		panic(err)
	}
	return &BlackClient{conn: conn}
}

// possibleBlackUserID是否被userID拉黑，也就是是否在userID的黑名单中
func (b *BlackClient) IsBlocked(ctx context.Context, possibleBlackUserID, userID string) (bool, error) {
	r, err := friend.NewFriendClient(b.conn).IsBlack(ctx, &friend.IsBlackReq{UserID1: possibleBlackUserID, UserID2: userID})
	if err != nil {
		return false, err
	}
	return r.InUser2Blacks, nil
}
