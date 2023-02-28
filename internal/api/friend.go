package api

import (
	"OpenIM/internal/api/a2r"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/proto/friend"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewFriend(zk *openKeeper.ZkClient) *Friend {
	return &Friend{zk: zk}
}

type Friend struct {
	zk *openKeeper.ZkClient
}

func (f *Friend) getGroupClient() (friend.FriendClient, error) {
	conn, err := f.zk.GetConn(config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		return nil, err
	}
	return friend.NewFriendClient(conn), nil
}

func (f *Friend) AddFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.AddFriend, f.getGroupClient, c)
}

func (f *Friend) DeleteFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.DeleteFriend, f.getGroupClient, c)
}

func (f *Friend) GetFriendApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetFriendApplyList, f.getGroupClient, c)
}

func (f *Friend) GetSelfApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetSelfApplyList, f.getGroupClient, c)
}

func (f *Friend) GetFriendList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetFriendList, f.getGroupClient, c)
}

func (f *Friend) AddFriendResponse(c *gin.Context) {
	a2r.Call(friend.FriendClient.AddFriendResponse, f.getGroupClient, c)
}

func (f *Friend) SetFriendRemark(c *gin.Context) {
	a2r.Call(friend.FriendClient.SetFriendRemark, f.getGroupClient, c)
}

func (f *Friend) AddBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.AddBlack, f.getGroupClient, c)
}

func (f *Friend) GetBlacklist(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetBlacklist, f.getGroupClient, c)
}

func (f *Friend) RemoveBlacklist(c *gin.Context) {
	a2r.Call(friend.FriendClient.RemoveBlacklist, f.getGroupClient, c)
}

func (f *Friend) ImportFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.ImportFriend, f.getGroupClient, c)
}

func (f *Friend) IsFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.IsFriend, f.getGroupClient, c)
}
