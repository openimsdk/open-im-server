package check

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"google.golang.org/grpc"
)

type MessageGateWayRpcClient struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func NewMessageGateWayRpcClient(zk discoveryRegistry.SvcDiscoveryRegistry) *MessageGateWayRpcClient {
	return &MessageGateWayRpcClient{
		zk: zk,
	}
}

func (m *MessageGateWayRpcClient) GetFriendsInfo(ctx context.Context, ownerUserID, friendUserID string) (resp *sdkws.FriendInfo, err error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	r, err := friend.NewFriendClient(cc).GetDesignatedFriends(ctx, &friend.GetDesignatedFriendsReq{OwnerUserID: ownerUserID, FriendUserIDs: []string{friendUserID}})
	if err != nil {
		return nil, err
	}
	resp = r.FriendsInfo[0]
	return
}
func (m *MessageGateWayRpcClient) getConn() (*grpc.ClientConn, error) {
	return m.zk.GetConn(config.Config.RpcRegisterName.OpenImMessageGatewayName)
}

// possibleFriendUserID是否在userID的好友中
func (m *MessageGateWayRpcClient) IsFriend(ctx context.Context, possibleFriendUserID, userID string) (bool, error) {
	cc, err := m.getConn()
	if err != nil {
		return false, err
	}
	resp, err := friend.NewFriendClient(cc).IsFriend(ctx, &friend.IsFriendReq{UserID1: userID, UserID2: possibleFriendUserID})
	if err != nil {
		return false, err
	}
	return resp.InUser1Friends, nil

}

func (m *MessageGateWayRpcClient) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	req := friend.GetFriendIDsReq{UserID: ownerUserID}
	resp, err := friend.NewFriendClient(cc).GetFriendIDs(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp.FriendIDs, err
}
