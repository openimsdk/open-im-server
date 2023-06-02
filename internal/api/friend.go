package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
)

func NewFriend(discov discoveryregistry.SvcDiscoveryRegistry) *Friend {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		// panic(err)
	}
	return &Friend{conn: conn, discov: discov}
}

type Friend struct {
	conn   *grpc.ClientConn
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (o *Friend) client(ctx context.Context) (friend.FriendClient, error) {
	c, err := o.discov.GetConn(ctx, config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		return nil, err
	}
	return friend.NewFriendClient(c), nil
}

func (o *Friend) ApplyToAddFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.ApplyToAddFriend, o.client, c)
}

func (o *Friend) RespondFriendApply(c *gin.Context) {
	a2r.Call(friend.FriendClient.RespondFriendApply, o.client, c)
}

func (o *Friend) DeleteFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.DeleteFriend, o.client, c)
}

func (o *Friend) GetFriendApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyTo, o.client, c)
}

func (o *Friend) GetSelfApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyFrom, o.client, c)
}

func (o *Friend) GetFriendList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriends, o.client, c)
}

func (o *Friend) SetFriendRemark(c *gin.Context) {
	a2r.Call(friend.FriendClient.SetFriendRemark, o.client, c)
}

func (o *Friend) AddBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.AddBlack, o.client, c)
}

func (o *Friend) GetPaginationBlacks(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationBlacks, o.client, c)
}

func (o *Friend) RemoveBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.RemoveBlack, o.client, c)
}

func (o *Friend) ImportFriends(c *gin.Context) {
	a2r.Call(friend.FriendClient.ImportFriends, o.client, c)
}

func (o *Friend) IsFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.IsFriend, o.client, c)
}
