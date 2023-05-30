package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"google.golang.org/grpc"
)

type FriendClient struct {
	conn *grpc.ClientConn
}

func NewFriendClient(discov discoveryRegistry.SvcDiscoveryRegistry) *FriendClient {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		panic(err)
	}
	return &FriendClient{conn: conn}
}

func (f *FriendClient) GetFriendsInfo(ctx context.Context, ownerUserID, friendUserID string) (resp *sdkws.FriendInfo, err error) {
	r, err := friend.NewFriendClient(f.conn).GetDesignatedFriends(ctx, &friend.GetDesignatedFriendsReq{OwnerUserID: ownerUserID, FriendUserIDs: []string{friendUserID}})
	if err != nil {
		return nil, err
	}
	resp = r.FriendsInfo[0]
	return
}

// possibleFriendUserID是否在userID的好友中
func (f *FriendClient) IsFriend(ctx context.Context, possibleFriendUserID, userID string) (bool, error) {
	resp, err := friend.NewFriendClient(f.conn).IsFriend(ctx, &friend.IsFriendReq{UserID1: userID, UserID2: possibleFriendUserID})
	if err != nil {
		return false, err
	}
	return resp.InUser1Friends, nil

}

func (f *FriendClient) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
	req := friend.GetFriendIDsReq{UserID: ownerUserID}
	resp, err := friend.NewFriendClient(f.conn).GetFriendIDs(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp.FriendIDs, err
}
