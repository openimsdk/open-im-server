package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/tools/a2r"
)

type GroupApi struct {
	Client group.GroupClient
}

func NewGroupApi(client group.GroupClient) GroupApi {
	return GroupApi{client}
}

func (o *GroupApi) CreateGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.CreateGroup, o.Client)
}

func (o *GroupApi) SetGroupInfo(c *gin.Context) {
	a2r.Call(c, group.GroupClient.SetGroupInfo, o.Client)
}

func (o *GroupApi) SetGroupInfoEx(c *gin.Context) {
	a2r.Call(c, group.GroupClient.SetGroupInfoEx, o.Client)
}

func (o *GroupApi) JoinGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.JoinGroup, o.Client)
}

func (o *GroupApi) QuitGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.QuitGroup, o.Client)
}

func (o *GroupApi) ApplicationGroupResponse(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GroupApplicationResponse, o.Client)
}

func (o *GroupApi) TransferGroupOwner(c *gin.Context) {
	a2r.Call(c, group.GroupClient.TransferGroupOwner, o.Client)
}

func (o *GroupApi) GetRecvGroupApplicationList(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupApplicationList, o.Client)
}

func (o *GroupApi) GetUserReqGroupApplicationList(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetUserReqApplicationList, o.Client)
}

func (o *GroupApi) GetGroupUsersReqApplicationList(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupUsersReqApplicationList, o.Client)
}

func (o *GroupApi) GetSpecifiedUserGroupRequestInfo(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetSpecifiedUserGroupRequestInfo, o.Client)
}

func (o *GroupApi) GetGroupsInfo(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupsInfo, o.Client)
}

func (o *GroupApi) KickGroupMember(c *gin.Context) {
	a2r.Call(c, group.GroupClient.KickGroupMember, o.Client)
}

func (o *GroupApi) GetGroupMembersInfo(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupMembersInfo, o.Client)
}

func (o *GroupApi) GetGroupMemberList(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupMemberList, o.Client)
}

func (o *GroupApi) InviteUserToGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.InviteUserToGroup, o.Client)
}

func (o *GroupApi) GetJoinedGroupList(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetJoinedGroupList, o.Client)
}

func (o *GroupApi) DismissGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.DismissGroup, o.Client)
}

func (o *GroupApi) MuteGroupMember(c *gin.Context) {
	a2r.Call(c, group.GroupClient.MuteGroupMember, o.Client)
}

func (o *GroupApi) CancelMuteGroupMember(c *gin.Context) {
	a2r.Call(c, group.GroupClient.CancelMuteGroupMember, o.Client)
}

func (o *GroupApi) MuteGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.MuteGroup, o.Client)
}

func (o *GroupApi) CancelMuteGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.CancelMuteGroup, o.Client)
}

func (o *GroupApi) SetGroupMemberInfo(c *gin.Context) {
	a2r.Call(c, group.GroupClient.SetGroupMemberInfo, o.Client)
}

func (o *GroupApi) GetGroupAbstractInfo(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupAbstractInfo, o.Client)
}

func (o *GroupApi) GroupCreateCount(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GroupCreateCount, o.Client)
}

func (o *GroupApi) GetGroups(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroups, o.Client)
}

func (o *GroupApi) GetGroupMemberUserIDs(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupMemberUserIDs, o.Client)
}

func (o *GroupApi) GetIncrementalJoinGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetIncrementalJoinGroup, o.Client)
}

func (o *GroupApi) GetIncrementalGroupMember(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetIncrementalGroupMember, o.Client)
}

func (o *GroupApi) GetIncrementalGroupMemberBatch(c *gin.Context) {
	a2r.Call(c, group.GroupClient.BatchGetIncrementalGroupMember, o.Client)
}

func (o *GroupApi) GetFullGroupMemberUserIDs(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetFullGroupMemberUserIDs, o.Client)
}

func (o *GroupApi) GetFullJoinGroupIDs(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetFullJoinGroupIDs, o.Client)
}

func (o *GroupApi) GetGroupApplicationUnhandledCount(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupApplicationUnhandledCount, o.Client)
}
