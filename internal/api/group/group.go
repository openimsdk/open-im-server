package group

import (
	"OpenIM/internal/a2r"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/proto/group"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewGroup(zk *openKeeper.ZkClient) *Group {
	return &Group{zk: zk}
}

type Group struct {
	zk *openKeeper.ZkClient
}

func (g *Group) getGroupClient() (group.GroupClient, error) {
	conn, err := g.zk.GetConn(config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		return nil, err
	}
	return group.NewGroupClient(conn), nil
}

func (g *Group) NewCreateGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CreateGroup, g.getGroupClient, c)
}

func (g *Group) NewSetGroupInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupInfo, g.getGroupClient, c)
}

func (g *Group) JoinGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.JoinGroup, g.getGroupClient, c)
}

func (g *Group) QuitGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.QuitGroup, g.getGroupClient, c)
}

func (g *Group) ApplicationGroupResponse(c *gin.Context) {
	a2r.Call(group.GroupClient.GroupApplicationResponse, g.getGroupClient, c)
}

func (g *Group) TransferGroupOwner(c *gin.Context) {
	a2r.Call(group.GroupClient.TransferGroupOwner, g.getGroupClient, c)
}

func (g *Group) GetRecvGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupApplicationList, g.getGroupClient, c)
}

func (g *Group) GetUserReqGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetUserReqApplicationList, g.getGroupClient, c)
}

func (g *Group) GetGroupsInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupsInfo, g.getGroupClient, c)
}

func (g *Group) KickGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.KickGroupMember, g.getGroupClient, c)
}

func (g *Group) GetGroupMembersInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMembersInfo, g.getGroupClient, c)
}

func (g *Group) InviteUserToGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.InviteUserToGroup, g.getGroupClient, c)
}

func (g *Group) GetJoinedGroupList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetJoinedGroupList, g.getGroupClient, c)
}

func (g *Group) DismissGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.DismissGroup, g.getGroupClient, c)
}

func (g *Group) MuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroupMember, g.getGroupClient, c)
}

func (g *Group) CancelMuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroupMember, g.getGroupClient, c)
}

func (g *Group) MuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroup, g.getGroupClient, c)
}

func (g *Group) CancelMuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroup, g.getGroupClient, c)
}

func (g *Group) SetGroupMemberInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupMemberInfo, g.getGroupClient, c)
}

func (g *Group) GetGroupAbstractInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupAbstractInfo, g.getGroupClient, c)
}

//func (g *Group) SetGroupMemberNickname(c *gin.Context) {
//	a2r.Call(group.GroupClient.SetGroupMemberNickname, g.getGroupClient, c)
//}
//
//func (g *Group) GetGroupAllMemberList(c *gin.Context) {
//	a2r.Call(group.GroupClient.GetGroupAllMember, g.getGroupClient, c)
//}

func (g *Group) GetJoinedSuperGroupList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetJoinedSuperGroupList, g.getGroupClient, c)
}

func (g *Group) GetSuperGroupsInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetSuperGroupsInfo, g.getGroupClient, c)
}
