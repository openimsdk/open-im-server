package check

import (
	"OpenIM/pkg/common/config"
	discoveryRegistry "OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/proto/friend"
	sdkws "OpenIM/pkg/proto/sdkws"
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
	r, err := friend.NewFriendClient(cc).GetDesignatedFriends(ctx, &friend.GetDesignatedFriendsReq{OwnerUserID: ownerUserID, FriendUserIDs: []string{friendUserID}})
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

func (f *FriendChecker) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
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
