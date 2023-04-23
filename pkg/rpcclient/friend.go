package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

type FriendClient struct {
	*MetaClient
}

func NewFriendClient(zk discoveryRegistry.SvcDiscoveryRegistry) *FriendClient {
	return &FriendClient{NewMetaClient(zk, config.Config.RpcRegisterName.OpenImFriendName)}
}

func (f *FriendClient) GetFriendsInfo(ctx context.Context, ownerUserID, friendUserID string) (resp *sdkws.FriendInfo, err error) {
	cc, err := f.getConn()
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

// possibleFriendUserID是否在userID的好友中
func (f *FriendClient) IsFriend(ctx context.Context, possibleFriendUserID, userID string) (bool, error) {
	cc, err := f.getConn()
	if err != nil {
		return false, err
	}
	resp, err := friend.NewFriendClient(cc).IsFriend(ctx, &friend.IsFriendReq{UserID1: userID, UserID2: possibleFriendUserID})
	if err != nil {
		return false, err
	}
	return resp.InUser1Friends, nil

}

func (f *FriendClient) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
	cc, err := f.getConn()
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
