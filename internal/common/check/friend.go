package check

import (
	"Open_IM/pkg/common/config"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	"Open_IM/pkg/proto/friend"
	sdkws "Open_IM/pkg/proto/sdkws"
	"context"
	"google.golang.org/grpc"
)

type FriendChecker struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func NewFriendChecker(zk discoveryRegistry.SvcDiscoveryRegistry) *FriendChecker {
	return &FriendChecker{
		zk: zk,
	}
}

func (f *FriendChecker) GetFriendsInfo(ctx context.Context, ownerUserID, friendUserID string) (resp *sdkws.FriendInfo, err error) {
	cc, err := f.getConn()
	if err != nil {
		return nil, err
	}
	r, err := friend.NewFriendClient(cc).GetPaginationFriends(ctx, &friend.GetPaginationFriendsReq{OwnerUserID: ownerUserID, FriendUserIDs: []string{friendUserID}})
	if err != nil {
		return nil, err
	}
	resp = r.FriendsInfo[0]
	return
}
func (f *FriendChecker) getConn() (*grpc.ClientConn, error) {
	return f.zk.GetConn(config.Config.RpcRegisterName.OpenImFriendName)
}

// possibleFriendUserID是否在userID的好友中
func (f *FriendChecker) IsFriend(ctx context.Context, possibleFriendUserID, userID string) (bool, error) {
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

func (f *FriendChecker) GetAllPageFriends(ctx context.Context, ownerUserID string) (resp []*sdkws.FriendInfo, err error) {

	cc, err := f.getConn()
	if err != nil {
		return nil, err
	}

	resp, err := friend.NewFriendClient(cc).GetPaginationFriends(ctx)
}
