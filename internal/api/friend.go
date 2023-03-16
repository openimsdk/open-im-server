package api

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/internal/api/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"

	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewFriend(c discoveryregistry.SvcDiscoveryRegistry) *Friend {
	return &Friend{c: c}
}

type Friend struct {
	c discoveryregistry.SvcDiscoveryRegistry
}

func (o *Friend) client() (friend.FriendClient, error) {
	conn, err := o.c.GetConn(config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		return nil, err
	}
	return friend.NewFriendClient(conn), nil
}

func (o *Friend) ApplyToAddFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.ApplyToAddFriend, o.client, c)
}

func (o *Friend) DeleteFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.DeleteFriend, o.client, c)
}

func (o *Friend) GetFriendApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyFrom, o.client, c)
}

func (o *Friend) GetSelfApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyFrom, o.client, c)
}

func (o *Friend) GetFriendList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetDesignatedFriends, o.client, c)
}

func (o *Friend) RespondFriendApply(c *gin.Context) {
	a2r.Call(friend.FriendClient.RespondFriendApply, o.client, c)
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
