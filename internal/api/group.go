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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

type GroupApi rpcclient.Group

// newgroupapi creates a new group
func NewGroupApi(discov discoveryregistry.SvcDiscoveryRegistry) GroupApi {
	return GroupApi(*rpcclient.NewGroup(discov))
}

// create a new group
func (o *GroupApi) CreateGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CreateGroup, o.Client, c)
}

// set group info
func (o *GroupApi) SetGroupInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupInfo, o.Client, c)
}

// take into group
func (o *GroupApi) JoinGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.JoinGroup, o.Client, c)
}

// quit group
func (o *GroupApi) QuitGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.QuitGroup, o.Client, c)
}

// call group application response
func (o *GroupApi) ApplicationGroupResponse(c *gin.Context) {
	a2r.Call(group.GroupClient.GroupApplicationResponse, o.Client, c)
}

// transfer group owner
func (o *GroupApi) TransferGroupOwner(c *gin.Context) {
	a2r.Call(group.GroupClient.TransferGroupOwner, o.Client, c)
}

// get group application list
func (o *GroupApi) GetRecvGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupApplicationList, o.Client, c)
}

// get user group list request
func (o *GroupApi) GetUserReqGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetUserReqApplicationList, o.Client, c)
}

// get group infomation
func (o *GroupApi) GetGroupsInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupsInfo, o.Client, c)
}

// kick user out of group
func (o *GroupApi) KickGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.KickGroupMember, o.Client, c)
}

// get user info from group
func (o *GroupApi) GetGroupMembersInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMembersInfo, o.Client, c)
}

// get user list info from group
func (o *GroupApi) GetGroupMemberList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMemberList, o.Client, c)
}

// invite user to join group
func (o *GroupApi) InviteUserToGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.InviteUserToGroup, o.Client, c)
}

// get group list user joined
func (o *GroupApi) GetJoinedGroupList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetJoinedGroupList, o.Client, c)
}

// dismiss group
func (o *GroupApi) DismissGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.DismissGroup, o.Client, c)
}

// mute group member
func (o *GroupApi) MuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroupMember, o.Client, c)
}

// cancel mute group member
func (o *GroupApi) CancelMuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroupMember, o.Client, c)
}

// mute group
func (o *GroupApi) MuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroup, o.Client, c)
}

// cancel mute group
func (o *GroupApi) CancelMuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroup, o.Client, c)
}

// set group member info
func (o *GroupApi) SetGroupMemberInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupMemberInfo, o.Client, c)
}

// get group abstract info
func (o *GroupApi) GetGroupAbstractInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupAbstractInfo, o.Client, c)
}

//	func (g *Group) SetGroupMemberNickname(c *gin.Context) {
//		a2r.Call(group.GroupClient.SetGroupMemberNickname, g.userClient, c)
//	}
//
//	func (g *Group) GetGroupAllMemberList(c *gin.Context) {
//		a2r.Call(group.GroupClient.GetGroupAllMember, g.userClient, c)
//	}
//
// get joinde super group list
func (o *GroupApi) GetJoinedSuperGroupList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetJoinedSuperGroupList, o.Client, c)
}

// get super group info
func (o *GroupApi) GetSuperGroupsInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetSuperGroupsInfo, o.Client, c)
}
