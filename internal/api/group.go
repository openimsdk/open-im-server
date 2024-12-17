// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/tools/a2r"
)

type GroupApi struct{}

func NewGroupApi() GroupApi {
	return GroupApi{}
}

func (o *GroupApi) CreateGroup(c *gin.Context) {
	a2r.CallV2(c, group.CreateGroupCaller.Invoke)
}

func (o *GroupApi) SetGroupInfo(c *gin.Context) {
	a2r.CallV2(c, group.SetGroupInfoCaller.Invoke)
}

func (o *GroupApi) SetGroupInfoEx(c *gin.Context) {
	a2r.CallV2(c, group.SetGroupInfoExCaller.Invoke)
}

func (o *GroupApi) JoinGroup(c *gin.Context) {
	a2r.CallV2(c, group.JoinGroupCaller.Invoke)
}

func (o *GroupApi) QuitGroup(c *gin.Context) {
	a2r.CallV2(c, group.QuitGroupCaller.Invoke)
}

func (o *GroupApi) ApplicationGroupResponse(c *gin.Context) {
	a2r.CallV2(c, group.GroupApplicationResponseCaller.Invoke)
}

func (o *GroupApi) TransferGroupOwner(c *gin.Context) {
	a2r.CallV2(c, group.TransferGroupOwnerCaller.Invoke)
}

func (o *GroupApi) GetRecvGroupApplicationList(c *gin.Context) {
	a2r.CallV2(c, group.GetGroupApplicationListCaller.Invoke)
}

func (o *GroupApi) GetUserReqGroupApplicationList(c *gin.Context) {
	a2r.CallV2(c, group.GetUserReqApplicationListCaller.Invoke)
}

func (o *GroupApi) GetGroupUsersReqApplicationList(c *gin.Context) {
	a2r.CallV2(c, group.GetGroupUsersReqApplicationListCaller.Invoke)
}

func (o *GroupApi) GetSpecifiedUserGroupRequestInfo(c *gin.Context) {
	a2r.CallV2(c, group.GetSpecifiedUserGroupRequestInfoCaller.Invoke)
}

func (o *GroupApi) GetGroupsInfo(c *gin.Context) {
	a2r.CallV2(c, group.GetGroupsInfoCaller.Invoke)
	//a2r.Call(group.GroupClient.GetGroupsInfo, o.Client, c, a2r.NewNilReplaceOption(group.GroupClient.GetGroupsInfo))
}

func (o *GroupApi) KickGroupMember(c *gin.Context) {
	a2r.CallV2(c, group.KickGroupMemberCaller.Invoke)
}

func (o *GroupApi) GetGroupMembersInfo(c *gin.Context) {
	a2r.CallV2(c, group.GetGroupMembersInfoCaller.Invoke)
	//a2r.Call(group.GroupClient.GetGroupMembersInfo, o.Client, c, a2r.NewNilReplaceOption(group.GroupClient.GetGroupMembersInfo))
}

func (o *GroupApi) GetGroupMemberList(c *gin.Context) {
	a2r.CallV2(c, group.GetGroupMemberListCaller.Invoke)
}

func (o *GroupApi) InviteUserToGroup(c *gin.Context) {
	a2r.CallV2(c, group.InviteUserToGroupCaller.Invoke)
}

func (o *GroupApi) GetJoinedGroupList(c *gin.Context) {
	a2r.CallV2(c, group.GetJoinedGroupListCaller.Invoke)
}

func (o *GroupApi) DismissGroup(c *gin.Context) {
	a2r.CallV2(c, group.DismissGroupCaller.Invoke)
}

func (o *GroupApi) MuteGroupMember(c *gin.Context) {
	a2r.CallV2(c, group.MuteGroupMemberCaller.Invoke)
}

func (o *GroupApi) CancelMuteGroupMember(c *gin.Context) {
	a2r.CallV2(c, group.CancelMuteGroupMemberCaller.Invoke)
}

func (o *GroupApi) MuteGroup(c *gin.Context) {
	a2r.CallV2(c, group.MuteGroupCaller.Invoke)
}

func (o *GroupApi) CancelMuteGroup(c *gin.Context) {
	a2r.CallV2(c, group.CancelMuteGroupCaller.Invoke)
}

func (o *GroupApi) SetGroupMemberInfo(c *gin.Context) {
	a2r.CallV2(c, group.SetGroupMemberInfoCaller.Invoke)
}

func (o *GroupApi) GetGroupAbstractInfo(c *gin.Context) {
	a2r.CallV2(c, group.GetGroupAbstractInfoCaller.Invoke)
}

// func (g *Group) SetGroupMemberNickname(c *gin.Context) {
//	a2r.Call(group.GroupClient.SetGroupMemberNickname, g.userClient, c)
//}
//
// func (g *Group) GetGroupAllMemberList(c *gin.Context) {
//	a2r.Call(group.GroupClient.GetGroupAllMember, g.userClient, c)
//}

func (o *GroupApi) GroupCreateCount(c *gin.Context) {
	a2r.CallV2(c, group.GroupCreateCountCaller.Invoke)
}

func (o *GroupApi) GetGroups(c *gin.Context) {
	a2r.CallV2(c, group.GetGroupsCaller.Invoke)
}

func (o *GroupApi) GetGroupMemberUserIDs(c *gin.Context) {
	a2r.CallV2(c, group.GetGroupMemberUserIDsCaller.Invoke)
}

func (o *GroupApi) GetIncrementalJoinGroup(c *gin.Context) {
	a2r.CallV2(c, group.GetIncrementalJoinGroupCaller.Invoke)
}

func (o *GroupApi) GetIncrementalGroupMember(c *gin.Context) {
	a2r.CallV2(c, group.GetIncrementalGroupMemberCaller.Invoke)
}

func (o *GroupApi) GetIncrementalGroupMemberBatch(c *gin.Context) {
	a2r.CallV2(c, group.BatchGetIncrementalGroupMemberCaller.Invoke)
}

func (o *GroupApi) GetFullGroupMemberUserIDs(c *gin.Context) {
	a2r.CallV2(c, group.GetFullGroupMemberUserIDsCaller.Invoke)
}

func (o *GroupApi) GetFullJoinGroupIDs(c *gin.Context) {
	a2r.CallV2(c, group.GetFullJoinGroupIDsCaller.Invoke)
}
