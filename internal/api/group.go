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

	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type GroupApi rpcclient.Group

func NewGroupApi(client rpcclient.Group) GroupApi {
	return GroupApi(client)
}

func (o *GroupApi) CreateGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CreateGroup, o.Client, c)
}

func (o *GroupApi) SetGroupInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupInfo, o.Client, c)
}

func (o *GroupApi) SetGroupInfoEx(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupInfoEx, o.Client, c)
}

func (o *GroupApi) JoinGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.JoinGroup, o.Client, c)
}

func (o *GroupApi) QuitGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.QuitGroup, o.Client, c)
}

func (o *GroupApi) ApplicationGroupResponse(c *gin.Context) {
	a2r.Call(group.GroupClient.GroupApplicationResponse, o.Client, c)
}

func (o *GroupApi) TransferGroupOwner(c *gin.Context) {
	a2r.Call(group.GroupClient.TransferGroupOwner, o.Client, c)
}

func (o *GroupApi) GetRecvGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupApplicationList, o.Client, c)
}

func (o *GroupApi) GetUserReqGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetUserReqApplicationList, o.Client, c)
}

func (o *GroupApi) GetGroupUsersReqApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupUsersReqApplicationList, o.Client, c)
}

func (o *GroupApi) GetSpecifiedUserGroupRequestInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetSpecifiedUserGroupRequestInfo, o.Client, c)
}

func (o *GroupApi) GetGroupsInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupsInfo, o.Client, c)
	//a2r.Call(group.GroupClient.GetGroupsInfo, o.Client, c, a2r.NewNilReplaceOption(group.GroupClient.GetGroupsInfo))
}

func (o *GroupApi) KickGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.KickGroupMember, o.Client, c)
}

func (o *GroupApi) GetGroupMembersInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMembersInfo, o.Client, c)
	//a2r.Call(group.GroupClient.GetGroupMembersInfo, o.Client, c, a2r.NewNilReplaceOption(group.GroupClient.GetGroupMembersInfo))
}

func (o *GroupApi) GetGroupMemberList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMemberList, o.Client, c)
}

func (o *GroupApi) InviteUserToGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.InviteUserToGroup, o.Client, c)
}

func (o *GroupApi) GetJoinedGroupList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetJoinedGroupList, o.Client, c)
}

func (o *GroupApi) DismissGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.DismissGroup, o.Client, c)
}

func (o *GroupApi) MuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroupMember, o.Client, c)
}

func (o *GroupApi) CancelMuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroupMember, o.Client, c)
}

func (o *GroupApi) MuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroup, o.Client, c)
}

func (o *GroupApi) CancelMuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroup, o.Client, c)
}

func (o *GroupApi) SetGroupMemberInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupMemberInfo, o.Client, c)
}

func (o *GroupApi) GetGroupAbstractInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupAbstractInfo, o.Client, c)
}

// func (g *Group) SetGroupMemberNickname(c *gin.Context) {
//	a2r.Call(group.GroupClient.SetGroupMemberNickname, g.userClient, c)
//}
//
// func (g *Group) GetGroupAllMemberList(c *gin.Context) {
//	a2r.Call(group.GroupClient.GetGroupAllMember, g.userClient, c)
//}

func (o *GroupApi) GroupCreateCount(c *gin.Context) {
	a2r.Call(group.GroupClient.GroupCreateCount, o.Client, c)
}

func (o *GroupApi) GetGroups(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroups, o.Client, c)
}

func (o *GroupApi) GetGroupMemberUserIDs(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMemberUserIDs, o.Client, c)
}

func (o *GroupApi) GetIncrementalJoinGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.GetIncrementalJoinGroup, o.Client, c)
}

func (o *GroupApi) GetIncrementalGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.GetIncrementalGroupMember, o.Client, c)
}

func (o *GroupApi) GetIncrementalGroupMemberBatch(c *gin.Context) {
	a2r.Call(group.GroupClient.BatchGetIncrementalGroupMember, o.Client, c)
}

func (o *GroupApi) GetFullGroupMemberUserIDs(c *gin.Context) {
	a2r.Call(group.GroupClient.GetFullGroupMemberUserIDs, o.Client, c)
}

func (o *GroupApi) GetFullJoinGroupIDs(c *gin.Context) {
	a2r.Call(group.GroupClient.GetFullJoinGroupIDs, o.Client, c)
}
