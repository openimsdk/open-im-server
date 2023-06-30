package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"google.golang.org/grpc"
)

type Friend struct {
	conn   grpc.ClientConnInterface
	Client friend.FriendClient
	discov discoveryregistry.SvcDiscoveryRegistry
}

func NewFriend(discov discoveryregistry.SvcDiscoveryRegistry) *Friend {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		panic(err)
	}
	client := friend.NewFriendClient(conn)
	return &Friend{discov: discov, conn: conn, Client: client}
}

type FriendRpcClient Friend

func NewFriendRpcClient(discov discoveryregistry.SvcDiscoveryRegistry) FriendRpcClient {
	return FriendRpcClient(*NewFriend(discov))
}

func (f *FriendRpcClient) GetFriendsInfo(ctx context.Context, ownerUserID, friendUserID string) (resp *sdkws.FriendInfo, err error) {
	r, err := f.Client.GetDesignatedFriends(ctx, &friend.GetDesignatedFriendsReq{OwnerUserID: ownerUserID, FriendUserIDs: []string{friendUserID}})
	if err != nil {
		return nil, err
	}
	resp = r.FriendsInfo[0]
	return
}

// possibleFriendUserID是否在userID的好友中
func (f *FriendRpcClient) IsFriend(ctx context.Context, possibleFriendUserID, userID string) (bool, error) {
	resp, err := f.Client.IsFriend(ctx, &friend.IsFriendReq{UserID1: userID, UserID2: possibleFriendUserID})
	if err != nil {
		return false, err
	}
	return resp.InUser1Friends, nil

}

func (f *FriendRpcClient) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
	req := friend.GetFriendIDsReq{UserID: ownerUserID}
	resp, err := f.Client.GetFriendIDs(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp.FriendIDs, nil
}

func (b *FriendRpcClient) IsBlocked(ctx context.Context, possibleBlackUserID, userID string) (bool, error) {
	r, err := b.Client.IsBlack(ctx, &friend.IsBlackReq{UserID1: possibleBlackUserID, UserID2: userID})
	if err != nil {
		return false, err
	}
	return r.InUser2Blacks, nil
}
