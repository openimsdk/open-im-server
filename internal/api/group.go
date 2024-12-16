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
	a2r.CallV2(group.CreateGroupCaller.Invoke, c)
}

func (o *GroupApi) SetGroupInfo(c *gin.Context) {
	a2r.CallV2(group.SetGroupInfoCaller.Invoke, c)
}

func (o *GroupApi) SetGroupInfoEx(c *gin.Context) {
	a2r.CallV2(group.SetGroupInfoExCaller.Invoke, c)
}

func (o *GroupApi) JoinGroup(c *gin.Context) {
	a2r.CallV2(group.JoinGroupCaller.Invoke, c)
}

func (o *GroupApi) QuitGroup(c *gin.Context) {
	a2r.CallV2(group.QuitGroupCaller.Invoke, c)
}

func (o *GroupApi) ApplicationGroupResponse(c *gin.Context) {
	a2r.CallV2(group.GroupApplicationResponseCaller.Invoke, c)
}

func (o *GroupApi) TransferGroupOwner(c *gin.Context) {
	a2r.CallV2(group.TransferGroupOwnerCaller.Invoke, c)
}

func (o *GroupApi) GetRecvGroupApplicationList(c *gin.Context) {
	a2r.CallV2(group.GetGroupApplicationListCaller.Invoke, c)
}

func (o *GroupApi) GetUserReqGroupApplicationList(c *gin.Context) {
	a2r.CallV2(group.GetUserReqApplicationListCaller.Invoke, c)
}

func (o *GroupApi) GetGroupUsersReqApplicationList(c *gin.Context) {
	a2r.CallV2(group.GetGroupUsersReqApplicationListCaller.Invoke, c)
}

func (o *GroupApi) GetSpecifiedUserGroupRequestInfo(c *gin.Context) {
	a2r.CallV2(group.GetSpecifiedUserGroupRequestInfoCaller.Invoke, c)
}

func (o *GroupApi) GetGroupsInfo(c *gin.Context) {
	a2r.CallV2(group.GetGroupsInfoCaller.Invoke, c)
	//a2r.Call(group.GroupClient.GetGroupsInfo, o.Client, c, a2r.NewNilReplaceOption(group.GroupClient.GetGroupsInfo))
}

func (o *GroupApi) KickGroupMember(c *gin.Context) {
	a2r.CallV2(group.KickGroupMemberCaller.Invoke, c)
}

func (o *GroupApi) GetGroupMembersInfo(c *gin.Context) {
	a2r.CallV2(group.GetGroupMembersInfoCaller.Invoke, c)
	//a2r.Call(group.GroupClient.GetGroupMembersInfo, o.Client, c, a2r.NewNilReplaceOption(group.GroupClient.GetGroupMembersInfo))
}

func (o *GroupApi) GetGroupMemberList(c *gin.Context) {
	a2r.CallV2(group.GetGroupMemberListCaller.Invoke, c)
}

func (o *GroupApi) InviteUserToGroup(c *gin.Context) {
	a2r.CallV2(group.InviteUserToGroupCaller.Invoke, c)
}

func (o *GroupApi) GetJoinedGroupList(c *gin.Context) {
	a2r.CallV2(group.GetJoinedGroupListCaller.Invoke, c)
}

func (o *GroupApi) DismissGroup(c *gin.Context) {
	a2r.CallV2(group.DismissGroupCaller.Invoke, c)
}

func (o *GroupApi) MuteGroupMember(c *gin.Context) {
	a2r.CallV2(group.MuteGroupMemberCaller.Invoke, c)
}

func (o *GroupApi) CancelMuteGroupMember(c *gin.Context) {
	a2r.CallV2(group.CancelMuteGroupMemberCaller.Invoke, c)
}

func (o *GroupApi) MuteGroup(c *gin.Context) {
	a2r.CallV2(group.MuteGroupCaller.Invoke, c)
}

func (o *GroupApi) CancelMuteGroup(c *gin.Context) {
	a2r.CallV2(group.CancelMuteGroupCaller.Invoke, c)
}

func (o *GroupApi) SetGroupMemberInfo(c *gin.Context) {
	a2r.CallV2(group.SetGroupMemberInfoCaller.Invoke, c)
}

func (o *GroupApi) GetGroupAbstractInfo(c *gin.Context) {
	a2r.CallV2(group.GetGroupAbstractInfoCaller.Invoke, c)
}

// func (g *Group) SetGroupMemberNickname(c *gin.Context) {
//	a2r.Call(group.GroupClient.SetGroupMemberNickname, g.userClient, c)
//}
//
// func (g *Group) GetGroupAllMemberList(c *gin.Context) {
//	a2r.Call(group.GroupClient.GetGroupAllMember, g.userClient, c)
//}

func (o *GroupApi) GroupCreateCount(c *gin.Context) {
	a2r.CallV2(group.GroupCreateCountCaller.Invoke, c)
}

func (o *GroupApi) GetGroups(c *gin.Context) {
	a2r.CallV2(group.GetGroupsCaller.Invoke, c)
}

func (o *GroupApi) GetGroupMemberUserIDs(c *gin.Context) {
	a2r.CallV2(group.GetGroupMemberUserIDsCaller.Invoke, c)
}

func (o *GroupApi) GetIncrementalJoinGroup(c *gin.Context) {
	a2r.CallV2(group.GetIncrementalJoinGroupCaller.Invoke, c)
}

func (o *GroupApi) GetIncrementalGroupMember(c *gin.Context) {
	a2r.CallV2(group.GetIncrementalGroupMemberCaller.Invoke, c)
}

func (o *GroupApi) GetIncrementalGroupMemberBatch(c *gin.Context) {
	a2r.CallV2(group.BatchGetIncrementalGroupMemberCaller.Invoke, c)
}

func (o *GroupApi) GetFullGroupMemberUserIDs(c *gin.Context) {
	a2r.CallV2(group.GetFullGroupMemberUserIDsCaller.Invoke, c)
}

func (o *GroupApi) GetFullJoinGroupIDs(c *gin.Context) {
	a2r.CallV2(group.GetFullJoinGroupIDsCaller.Invoke, c)
}
