package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"

	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewGroup(c discoveryregistry.SvcDiscoveryRegistry) *Group {
	return &Group{c: c}
}

type Group struct {
	c discoveryregistry.SvcDiscoveryRegistry
}

func (o *Group) client(ctx context.Context) (group.GroupClient, error) {
	conn, err := o.c.GetConn(ctx, config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "get conn", o.c.GetClientLocalConns())
	return group.NewGroupClient(conn), nil
}

func (o *Group) NewCreateGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CreateGroup, o.client, c)
}

func (o *Group) NewSetGroupInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupInfo, o.client, c)
}

func (o *Group) JoinGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.JoinGroup, o.client, c)
}

func (o *Group) QuitGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.QuitGroup, o.client, c)
}

func (o *Group) ApplicationGroupResponse(c *gin.Context) {
	a2r.Call(group.GroupClient.GroupApplicationResponse, o.client, c)
}

func (o *Group) TransferGroupOwner(c *gin.Context) {
	a2r.Call(group.GroupClient.TransferGroupOwner, o.client, c)
}

func (o *Group) GetRecvGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupApplicationList, o.client, c)
}

func (o *Group) GetUserReqGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetUserReqApplicationList, o.client, c)
}

func (o *Group) GetGroupsInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupsInfo, o.client, c)
}

func (o *Group) KickGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.KickGroupMember, o.client, c)
}

func (o *Group) GetGroupMembersInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMembersInfo, o.client, c)
}

func (o *Group) GetGroupMemberList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMemberList, o.client, c)
}

func (o *Group) InviteUserToGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.InviteUserToGroup, o.client, c)
}

func (o *Group) GetJoinedGroupList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetJoinedGroupList, o.client, c)
}

func (o *Group) DismissGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.DismissGroup, o.client, c)
}

func (o *Group) MuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroupMember, o.client, c)
}

func (o *Group) CancelMuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroupMember, o.client, c)
}

func (o *Group) MuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroup, o.client, c)
}

func (o *Group) CancelMuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroup, o.client, c)
}

func (o *Group) SetGroupMemberInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupMemberInfo, o.client, c)
}

func (o *Group) GetGroupAbstractInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupAbstractInfo, o.client, c)
}

//func (g *Group) SetGroupMemberNickname(c *gin.Context) {
//	a2r.Call(group.GroupClient.SetGroupMemberNickname, g.client, c)
//}
//
//func (g *Group) GetGroupAllMemberList(c *gin.Context) {
//	a2r.Call(group.GroupClient.GetGroupAllMember, g.client, c)
//}

func (o *Group) GetJoinedSuperGroupList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetJoinedSuperGroupList, o.client, c)
}

func (o *Group) GetSuperGroupsInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetSuperGroupsInfo, o.client, c)
}
