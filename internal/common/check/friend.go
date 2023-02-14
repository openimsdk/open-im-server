package check

import (
	"Open_IM/pkg/common/config"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	sdkws "Open_IM/pkg/proto/sdkws"
	"context"
	"errors"
	"google.golang.org/grpc"
)

type FriendChecker struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func (f *FriendChecker) GetFriendsInfo(ctx context.Context, ownerUserID, friendUserID string) (*sdkws.FriendInfo, error) {
	return nil, errors.New("TODO:GetUserInfo")
}
func (u *FriendChecker) getConn() (*grpc.ClientConn, error) {
	return u.zk.GetConn(config.Config.RpcRegisterName.OpenImFriendName)
}
