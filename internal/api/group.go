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
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
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

// SetSendMessageSetting 设置群成员发消息权限：allowSendMsg 0=全员可发，1=仅群主/管理员可发（委托 SetGroupInfoEx）。
func (o *GroupApi) SetSendMessageSetting(c *gin.Context) {
	var req struct {
		GroupID      string `json:"groupID"`
		AllowSendMsg int32  `json:"allowSendMsg"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if req.GroupID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("groupID is empty"))
		return
	}
	if req.AllowSendMsg != 0 && req.AllowSendMsg != 1 {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("allowSendMsg must be 0 or 1"))
		return
	}
	resp, err := o.Client.SetGroupInfoEx(c.Request.Context(), &group.SetGroupInfoExReq{
		GroupID:      req.GroupID,
		AllowSendMsg: wrapperspb.Int32(req.AllowSendMsg),
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// GetSendMessageSetting 返回当前群的 allowSendMsg：0=全员可发，1=仅群主/管理员可发（与 get_groups_info 中字段一致）。
func (o *GroupApi) GetSendMessageSetting(c *gin.Context) {
	var req struct {
		GroupID string `json:"groupID"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if req.GroupID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("groupID is empty"))
		return
	}
	resp, err := o.Client.GetGroupsInfo(c.Request.Context(), &group.GetGroupsInfoReq{
		GroupIDs: []string{req.GroupID},
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(resp.GroupInfos) == 0 {
		apiresp.GinError(c, errs.ErrRecordNotFound.WrapMsg("group not found", "groupID", req.GroupID))
		return
	}
	gi := resp.GroupInfos[0]
	apiresp.GinSuccess(c, struct {
		GroupID      string `json:"groupID"`
		AllowSendMsg int32  `json:"allowSendMsg"`
	}{
		GroupID:      gi.GroupID,
		AllowSendMsg: gi.GetAllowSendMsg(),
	})
}

// SetInviteSetting 设置群成员邀请他人入群权限：allowAddMember 0=全员可邀请/拉人，1=仅群主/管理员（委托 SetGroupInfoEx；邀请走 InviteUserToGroup 时 RPC 已校验）。
func (o *GroupApi) SetInviteSetting(c *gin.Context) {
	var req struct {
		GroupID        string `json:"groupID"`
		AllowAddMember int32  `json:"allowAddMember"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if req.GroupID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("groupID is empty"))
		return
	}
	if req.AllowAddMember != 0 && req.AllowAddMember != 1 {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("allowAddMember must be 0 or 1"))
		return
	}
	resp, err := o.Client.SetGroupInfoEx(c.Request.Context(), &group.SetGroupInfoExReq{
		GroupID:        req.GroupID,
		AllowAddMember: wrapperspb.Int32(req.AllowAddMember),
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// GetInviteSetting 返回当前群的 allowAddMember（与 get_groups_info 中字段一致）。
func (o *GroupApi) GetInviteSetting(c *gin.Context) {
	var req struct {
		GroupID string `json:"groupID"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if req.GroupID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("groupID is empty"))
		return
	}
	resp, err := o.Client.GetGroupsInfo(c.Request.Context(), &group.GetGroupsInfoReq{
		GroupIDs: []string{req.GroupID},
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(resp.GroupInfos) == 0 {
		apiresp.GinError(c, errs.ErrRecordNotFound.WrapMsg("group not found", "groupID", req.GroupID))
		return
	}
	gi := resp.GroupInfos[0]
	apiresp.GinSuccess(c, struct {
		GroupID        string `json:"groupID"`
		AllowAddMember int32  `json:"allowAddMember"`
	}{
		GroupID:        gi.GroupID,
		AllowAddMember: gi.GetAllowAddMember(),
	})
}

// SetPinSetting 设置群成员置顶消息权限：allowPinMsg 0=全员可置顶，1=仅群主/管理员（委托 SetGroupInfoEx；置顶/取消置顶走 pin_group_message / unpin_group_message 时 RPC 已校验）。
func (o *GroupApi) SetPinSetting(c *gin.Context) {
	var req struct {
		GroupID     string `json:"groupID"`
		AllowPinMsg int32  `json:"allowPinMsg"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if req.GroupID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("groupID is empty"))
		return
	}
	if req.AllowPinMsg != 0 && req.AllowPinMsg != 1 {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("allowPinMsg must be 0 or 1"))
		return
	}
	resp, err := o.Client.SetGroupInfoEx(c.Request.Context(), &group.SetGroupInfoExReq{
		GroupID:     req.GroupID,
		AllowPinMsg: wrapperspb.Int32(req.AllowPinMsg),
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// GetPinSetting 返回当前群的 allowPinMsg（与 get_groups_info 中字段一致）。
func (o *GroupApi) GetPinSetting(c *gin.Context) {
	var req struct {
		GroupID string `json:"groupID"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if req.GroupID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("groupID is empty"))
		return
	}
	resp, err := o.Client.GetGroupsInfo(c.Request.Context(), &group.GetGroupsInfoReq{
		GroupIDs: []string{req.GroupID},
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(resp.GroupInfos) == 0 {
		apiresp.GinError(c, errs.ErrRecordNotFound.WrapMsg("group not found", "groupID", req.GroupID))
		return
	}
	gi := resp.GroupInfos[0]
	apiresp.GinSuccess(c, struct {
		GroupID     string `json:"groupID"`
		AllowPinMsg int32  `json:"allowPinMsg"`
	}{
		GroupID:     gi.GroupID,
		AllowPinMsg: gi.GetAllowPinMsg(),
	})
}

// SetEditSetting 设置群成员编辑群资料权限：allowEditGroupInfo 0=全员可编辑，1=仅群主/管理员（委托 SetGroupInfoEx；改群资料走 set_group_info / set_group_info_ex 时 RPC 已校验）。
func (o *GroupApi) SetEditSetting(c *gin.Context) {
	var req struct {
		GroupID            string `json:"groupID"`
		AllowEditGroupInfo int32  `json:"allowEditGroupInfo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if req.GroupID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("groupID is empty"))
		return
	}
	if req.AllowEditGroupInfo != 0 && req.AllowEditGroupInfo != 1 {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("allowEditGroupInfo must be 0 or 1"))
		return
	}
	resp, err := o.Client.SetGroupInfoEx(c.Request.Context(), &group.SetGroupInfoExReq{
		GroupID:            req.GroupID,
		AllowEditGroupInfo: wrapperspb.Int32(req.AllowEditGroupInfo),
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// GetEditSetting 返回当前群的 allowEditGroupInfo（与 get_groups_info 中字段一致）。
func (o *GroupApi) GetEditSetting(c *gin.Context) {
	var req struct {
		GroupID string `json:"groupID"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if req.GroupID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("groupID is empty"))
		return
	}
	resp, err := o.Client.GetGroupsInfo(c.Request.Context(), &group.GetGroupsInfoReq{
		GroupIDs: []string{req.GroupID},
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	if len(resp.GroupInfos) == 0 {
		apiresp.GinError(c, errs.ErrRecordNotFound.WrapMsg("group not found", "groupID", req.GroupID))
		return
	}
	gi := resp.GroupInfos[0]
	apiresp.GinSuccess(c, struct {
		GroupID            string `json:"groupID"`
		AllowEditGroupInfo int32  `json:"allowEditGroupInfo"`
	}{
		GroupID:            gi.GroupID,
		AllowEditGroupInfo: gi.GetAllowEditGroupInfo(),
	})
}

// SetMsgBurnDuration 设置群消息阅后即焚时长（秒）；burnDuration=0 表示关闭。
func (o *GroupApi) SetMsgBurnDuration(c *gin.Context) {
	var req struct {
		GroupID         string `json:"groupID"`
		BurnDuration    int32  `json:"burnDuration"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg(err.Error()))
		return
	}
	if req.GroupID == "" {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("groupID is empty"))
		return
	}
	if req.BurnDuration < 0 {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("burnDuration must be >= 0"))
		return
	}
	resp, err := o.Client.SetGroupInfoEx(c.Request.Context(), &group.SetGroupInfoExReq{
		GroupID:         req.GroupID,
		MsgBurnDuration: wrapperspb.Int32(req.BurnDuration),
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
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
	//a2r.Call(c, group.GroupClient.GetGroupsInfo, o.Client, c, a2r.NewNilReplaceOption(group.GroupClient.GetGroupsInfo))
}

func (o *GroupApi) KickGroupMember(c *gin.Context) {
	a2r.Call(c, group.GroupClient.KickGroupMember, o.Client)
}

func (o *GroupApi) GetGroupMembersInfo(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupMembersInfo, o.Client)
	//a2r.Call(c, group.GroupClient.GetGroupMembersInfo, o.Client, c, a2r.NewNilReplaceOption(group.GroupClient.GetGroupMembersInfo))
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

// func (g *Group) SetGroupMemberNickname(c *gin.Context) {
//	a2r.Call(c, group.GroupClient.SetGroupMemberNickname, g.userClient)
//}
//
// func (g *Group) GetGroupAllMemberList(c *gin.Context) {
//	a2r.Call(c, group.GroupClient.GetGroupAllMember, g.userClient)
//}

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

func (o *GroupApi) GetCommonGroupsWithFriend(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetCommonGroupsWithFriend, o.Client)
}

func (o *GroupApi) PinGroupMessage(c *gin.Context) {
	a2r.Call(c, group.GroupClient.PinGroupMessage, o.Client)
}

func (o *GroupApi) UnpinGroupMessage(c *gin.Context) {
	a2r.Call(c, group.GroupClient.UnpinGroupMessage, o.Client)
}

func (o *GroupApi) GetGroupPinnedMessages(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupPinnedMessages, o.Client)
}

func (o *GroupApi) SetGroupMute(c *gin.Context) {
	a2r.Call(c, group.GroupClient.SetGroupMute, o.Client)
}

func (o *GroupApi) GetGroupMute(c *gin.Context) {
	a2r.Call(c, group.GroupClient.GetGroupMute, o.Client)
}

func (o *GroupApi) PinGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.PinGroup, o.Client)
}

func (o *GroupApi) UnpinGroup(c *gin.Context) {
	a2r.Call(c, group.GroupClient.UnpinGroup, o.Client)
}
