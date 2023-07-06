package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"

	"github.com/gin-gonic/gin"
)

type FriendApi rpcclient.Friend

func NewFriendApi(discov discoveryregistry.SvcDiscoveryRegistry) FriendApi {
	return FriendApi(*rpcclient.NewFriend(discov))
}

func (o *FriendApi) ApplyToAddFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.ApplyToAddFriend, o.Client, c)
}

func (o *FriendApi) RespondFriendApply(c *gin.Context) {
	a2r.Call(friend.FriendClient.RespondFriendApply, o.Client, c)
}

func (o *FriendApi) DeleteFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.DeleteFriend, o.Client, c)
}

func (o *FriendApi) GetFriendApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyTo, o.Client, c)
}

func (o *FriendApi) GetSelfApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyFrom, o.Client, c)
}

func (o *FriendApi) GetFriendList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriends, o.Client, c)
}

func (o *FriendApi) SetFriendRemark(c *gin.Context) {
	a2r.Call(friend.FriendClient.SetFriendRemark, o.Client, c)
}

func (o *FriendApi) AddBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.AddBlack, o.Client, c)
}

func (o *FriendApi) GetPaginationBlacks(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationBlacks, o.Client, c)
}

func (o *FriendApi) RemoveBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.RemoveBlack, o.Client, c)
}

func (o *FriendApi) ImportFriends(c *gin.Context) {
	a2r.Call(friend.FriendClient.ImportFriends, o.Client, c)
}

func (o *FriendApi) IsFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.IsFriend, o.Client, c)
}
