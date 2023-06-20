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
		panic(err)
	}
	client := friend.NewFriendClient(conn)
	return &Friend{discov: discov, conn: conn, client: client}
}

type Friend struct {
	conn   *grpc.ClientConn
	client friend.FriendClient
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (o *Friend) Client() friend.FriendClient {
	return friend.NewFriendClient(o.conn)
}

func (o *Friend) ApplyToAddFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.ApplyToAddFriend, o.Client, c)
}

func (o *Friend) RespondFriendApply(c *gin.Context) {
	a2r.Call(friend.FriendClient.RespondFriendApply, o.Client, c)
}

func (o *Friend) DeleteFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.DeleteFriend, o.Client, c)
}

func (o *Friend) GetFriendApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyTo, o.Client, c)
}

func (o *Friend) GetSelfApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyFrom, o.Client, c)
}

func (o *Friend) GetFriendList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriends, o.Client, c)
}

func (o *Friend) SetFriendRemark(c *gin.Context) {
	a2r.Call(friend.FriendClient.SetFriendRemark, o.Client, c)
}

func (o *Friend) AddBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.AddBlack, o.Client, c)
}

func (o *Friend) GetPaginationBlacks(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationBlacks, o.Client, c)
}

func (o *Friend) RemoveBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.RemoveBlack, o.Client, c)
}

func (o *Friend) ImportFriends(c *gin.Context) {
	a2r.Call(friend.FriendClient.ImportFriends, o.Client, c)
}

func (o *Friend) IsFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.IsFriend, o.Client, c)
}
