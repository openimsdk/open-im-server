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

package group

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GroupEventCallbackConfig struct {
	CallbackUrl       string
	BeforeCreateGroup config.CallBackConfig
}

// CallbackBeforeCreateGroup callback before create group.
func CallbackBeforeCreateGroup(ctx context.Context, cfg *GroupEventCallbackConfig, req *group.CreateGroupReq) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}
	cbReq := &callbackstruct.CallbackBeforeCreateGroupReq{
		CallbackCommand: callbackstruct.CallbackBeforeCreateGroupCommand,
		OperationID:     mcontext.GetOperationID(ctx),
		GroupInfo:       req.GroupInfo,
	}
	cbReq.InitMemberList = append(cbReq.InitMemberList, &apistruct.GroupAddMemberInfo{
		UserID:    req.OwnerUserID,
		RoleLevel: constant.GroupOwner,
	})
	for _, userID := range req.AdminUserIDs {
		cbReq.InitMemberList = append(cbReq.InitMemberList, &apistruct.GroupAddMemberInfo{
			UserID:    userID,
			RoleLevel: constant.GroupAdmin,
		})
	}
	for _, userID := range req.MemberUserIDs {
		cbReq.InitMemberList = append(cbReq.InitMemberList, &apistruct.GroupAddMemberInfo{
			UserID:    userID,
			RoleLevel: constant.GroupOrdinaryUsers,
		})
	}
	resp := &callbackstruct.CallbackBeforeCreateGroupResp{}

	if err = http.CallBackPostReturn(ctx, cfg.CallbackUrl, cbReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}

	utils.NotNilReplace(&req.GroupInfo.GroupID, resp.GroupID)
	utils.NotNilReplace(&req.GroupInfo.GroupName, resp.GroupName)
	utils.NotNilReplace(&req.GroupInfo.Notification, resp.Notification)
	utils.NotNilReplace(&req.GroupInfo.Introduction, resp.Introduction)
	utils.NotNilReplace(&req.GroupInfo.FaceURL, resp.FaceURL)
	utils.NotNilReplace(&req.GroupInfo.OwnerUserID, resp.OwnerUserID)
	utils.NotNilReplace(&req.GroupInfo.Ex, resp.Ex)
	utils.NotNilReplace(&req.GroupInfo.Status, resp.Status)
	utils.NotNilReplace(&req.GroupInfo.CreatorUserID, resp.CreatorUserID)
	utils.NotNilReplace(&req.GroupInfo.GroupType, resp.GroupType)
	utils.NotNilReplace(&req.GroupInfo.NeedVerification, resp.NeedVerification)
	utils.NotNilReplace(&req.GroupInfo.LookMemberInfo, resp.LookMemberInfo)
	return nil
}

func CallbackAfterCreateGroup(ctx context.Context, cfg *GroupEventCallbackConfig, req *group.CreateGroupReq) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}

	cbReq := &callbackstruct.CallbackAfterCreateGroupReq{
		CallbackCommand: callbackstruct.CallbackAfterCreateGroupCommand,
		GroupInfo:       req.GroupInfo,
	}
	cbReq.InitMemberList = append(cbReq.InitMemberList, &apistruct.GroupAddMemberInfo{
		UserID:    req.OwnerUserID,
		RoleLevel: constant.GroupOwner,
	})
	for _, userID := range req.AdminUserIDs {
		cbReq.InitMemberList = append(cbReq.InitMemberList, &apistruct.GroupAddMemberInfo{
			UserID:    userID,
			RoleLevel: constant.GroupAdmin,
		})
	}
	for _, userID := range req.MemberUserIDs {
		cbReq.InitMemberList = append(cbReq.InitMemberList, &apistruct.GroupAddMemberInfo{
			UserID:    userID,
			RoleLevel: constant.GroupOrdinaryUsers,
		})
	}
	resp := &callbackstruct.CallbackAfterCreateGroupResp{}
	if err = http.CallBackPostReturn(ctx, cfg.CallbackUrl, cbReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}
	return nil
}

func CallbackBeforeMemberJoinGroup(ctx context.Context, cfg *GroupEventCallbackConfig, groupMember *relation.GroupMemberModel, groupEx string) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}
	callbackReq := &callbackstruct.CallbackBeforeMemberJoinGroupReq{
		CallbackCommand: callbackstruct.CallbackBeforeMemberJoinGroupCommand,
		GroupID:         groupMember.GroupID,
		UserID:          groupMember.UserID,
		Ex:              groupMember.Ex,
		GroupEx:         groupEx,
	}
	resp := &callbackstruct.CallbackBeforeMemberJoinGroupResp{}

	if err = http.CallBackPostReturn(ctx, cfg.CallbackUrl, callbackReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}

	if resp.MuteEndTime != nil {
		groupMember.MuteEndTime = time.UnixMilli(*resp.MuteEndTime)
	}

	utils.NotNilReplace(&groupMember.FaceURL, resp.FaceURL)
	utils.NotNilReplace(&groupMember.Ex, resp.Ex)
	utils.NotNilReplace(&groupMember.Nickname, resp.Nickname)
	utils.NotNilReplace(&groupMember.RoleLevel, resp.RoleLevel)
	return nil
}

func CallbackBeforeSetGroupMemberInfo(ctx context.Context, cfg *GroupEventCallbackConfig, req *group.SetGroupMemberInfo) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}

	callbackReq := callbackstruct.CallbackBeforeSetGroupMemberInfoReq{
		CallbackCommand: callbackstruct.CallbackBeforeSetGroupMemberInfoCommand,
		GroupID:         req.GroupID,
		UserID:          req.UserID,
	}
	if req.Nickname != nil {
		callbackReq.Nickname = &req.Nickname.Value
	}
	if req.FaceURL != nil {
		callbackReq.FaceURL = &req.FaceURL.Value
	}
	if req.RoleLevel != nil {
		callbackReq.RoleLevel = &req.RoleLevel.Value
	}
	if req.Ex != nil {
		callbackReq.Ex = &req.Ex.Value
	}
	resp := &callbackstruct.CallbackBeforeSetGroupMemberInfoResp{}
	err = http.CallBackPostReturn(
		ctx,
		cfg.CallbackUrl,
		callbackReq,
		resp,
		cfg.BeforeCreateGroup,
	)
	if err != nil {
		return err
	}
	if resp.FaceURL != nil {
		req.FaceURL = wrapperspb.String(*resp.FaceURL)
	}
	if resp.Nickname != nil {
		req.Nickname = wrapperspb.String(*resp.Nickname)
	}
	if resp.RoleLevel != nil {
		req.RoleLevel = wrapperspb.Int32(*resp.RoleLevel)
	}
	if resp.Ex != nil {
		req.Ex = wrapperspb.String(*resp.Ex)
	}
	return nil
}

func CallbackAfterSetGroupMemberInfo(ctx context.Context, cfg *GroupEventCallbackConfig, req *group.SetGroupMemberInfo) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}
	callbackReq := callbackstruct.CallbackAfterSetGroupMemberInfoReq{
		CallbackCommand: callbackstruct.CallbackAfterSetGroupMemberInfoCommand,
		GroupID:         req.GroupID,
		UserID:          req.UserID,
	}
	if req.Nickname != nil {
		callbackReq.Nickname = &req.Nickname.Value
	}
	if req.FaceURL != nil {
		callbackReq.FaceURL = &req.FaceURL.Value
	}
	if req.RoleLevel != nil {
		callbackReq.RoleLevel = &req.RoleLevel.Value
	}
	if req.Ex != nil {
		callbackReq.Ex = &req.Ex.Value
	}
	resp := &callbackstruct.CallbackAfterSetGroupMemberInfoResp{}
	if err = http.CallBackPostReturn(ctx, cfg.CallbackUrl, callbackReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}
	return nil
}

func CallbackQuitGroup(ctx context.Context, cfg *GroupEventCallbackConfig, req *group.QuitGroupReq) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}
	cbReq := &callbackstruct.CallbackQuitGroupReq{
		CallbackCommand: callbackstruct.CallbackQuitGroupCommand,
		GroupID:         req.GroupID,
		UserID:          req.UserID,
	}
	resp := &callbackstruct.CallbackQuitGroupResp{}
	if err = http.CallBackPostReturn(ctx, cfg.CallbackUrl, cbReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}
	return nil
}

func CallbackKillGroupMember(ctx context.Context, cfg *GroupEventCallbackConfig, req *pbgroup.KickGroupMemberReq) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}
	cbReq := &callbackstruct.CallbackKillGroupMemberReq{
		CallbackCommand: callbackstruct.CallbackKillGroupCommand,
		GroupID:         req.GroupID,
		KickedUserIDs:   req.KickedUserIDs,
	}
	resp := &callbackstruct.CallbackKillGroupMemberResp{}
	if err = http.CallBackPostReturn(ctx, cfg.CallbackUrl, cbReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}
	return nil
}

func CallbackDismissGroup(ctx context.Context, cfg *GroupEventCallbackConfig, req *callbackstruct.CallbackDisMissGroupReq) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}
	req.CallbackCommand = callbackstruct.CallbackDisMissGroupCommand
	resp := &callbackstruct.CallbackDisMissGroupResp{}
	if err = http.CallBackPostReturn(ctx, cfg.CallbackUrl, req, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}
	return nil
}

func CallbackApplyJoinGroupBefore(ctx context.Context, cfg *GroupEventCallbackConfig, req *callbackstruct.CallbackJoinGroupReq) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}

	req.CallbackCommand = callbackstruct.CallbackBeforeJoinGroupCommand

	resp := &callbackstruct.CallbackJoinGroupResp{}
	if err = http.CallBackPostReturn(ctx, cfg.CallbackUrl, req, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}

	return nil
}

func CallbackAfterTransferGroupOwner(ctx context.Context, cfg *GroupEventCallbackConfig, req *pbgroup.TransferGroupOwnerReq) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}

	cbReq := &callbackstruct.CallbackTransferGroupOwnerReq{
		CallbackCommand: callbackstruct.CallbackAfterTransferGroupOwner,
		GroupID:         req.GroupID,
		OldOwnerUserID:  req.OldOwnerUserID,
		NewOwnerUserID:  req.NewOwnerUserID,
	}

	resp := &callbackstruct.CallbackTransferGroupOwnerResp{}
	if err = http.CallBackPostReturn(ctx, cfg.CallbackUrl, cbReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}
	return nil
}

func CallbackBeforeInviteUserToGroup(ctx context.Context, cfg *GroupEventCallbackConfig, req *group.InviteUserToGroupReq) (err error) {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}

	callbackReq := &callbackstruct.CallbackBeforeInviteUserToGroupReq{
		CallbackCommand: callbackstruct.CallbackBeforeInviteJoinGroupCommand,
		OperationID:     mcontext.GetOperationID(ctx),
		GroupID:         req.GroupID,
		Reason:          req.Reason,
		InvitedUserIDs:  req.InvitedUserIDs,
	}

	resp := &callbackstruct.CallbackBeforeInviteUserToGroupResp{}
	err = http.CallBackPostReturn(
		ctx,
		cfg.CallbackUrl,
		callbackReq,
		resp,
		cfg.BeforeCreateGroup,
	)

	if err != nil {
		return err
	}

	if len(resp.RefusedMembersAccount) > 0 {
		// Handle the scenario where certain members are refused
		// You might want to update the req.Members list or handle it as per your business logic
	}
	return nil
}

func CallbackAfterJoinGroup(ctx context.Context, cfg *GroupEventCallbackConfig, req *group.JoinGroupReq) error {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}
	callbackReq := &callbackstruct.CallbackAfterJoinGroupReq{
		CallbackCommand: callbackstruct.CallbackAfterJoinGroupCommand,
		OperationID:     mcontext.GetOperationID(ctx),
		GroupID:         req.GroupID,
		ReqMessage:      req.ReqMessage,
		JoinSource:      req.JoinSource,
		InviterUserID:   req.InviterUserID,
	}
	resp := &callbackstruct.CallbackAfterJoinGroupResp{}
	if err := http.CallBackPostReturn(ctx, cfg.CallbackUrl, callbackReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}
	return nil
}

func CallbackBeforeSetGroupInfo(ctx context.Context, cfg *GroupEventCallbackConfig, req *group.SetGroupInfoReq) error {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}
	callbackReq := &callbackstruct.CallbackBeforeSetGroupInfoReq{
		CallbackCommand: callbackstruct.CallbackBeforeSetGroupInfoCommand,
		GroupID:         req.GroupInfoForSet.GroupID,
		Notification:    req.GroupInfoForSet.Notification,
		Introduction:    req.GroupInfoForSet.Introduction,
		FaceURL:         req.GroupInfoForSet.FaceURL,
		GroupName:       req.GroupInfoForSet.GroupName,
	}

	if req.GroupInfoForSet.Ex != nil {
		callbackReq.Ex = req.GroupInfoForSet.Ex.Value
	}
	log.ZDebug(ctx, "debug CallbackBeforeSetGroupInfo", callbackReq.Ex)
	if req.GroupInfoForSet.NeedVerification != nil {
		callbackReq.NeedVerification = req.GroupInfoForSet.NeedVerification.Value
	}
	if req.GroupInfoForSet.LookMemberInfo != nil {
		callbackReq.LookMemberInfo = req.GroupInfoForSet.LookMemberInfo.Value
	}
	if req.GroupInfoForSet.ApplyMemberFriend != nil {
		callbackReq.ApplyMemberFriend = req.GroupInfoForSet.ApplyMemberFriend.Value
	}
	resp := &callbackstruct.CallbackBeforeSetGroupInfoResp{}

	if err := http.CallBackPostReturn(ctx, cfg.CallbackUrl, callbackReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}

	if resp.Ex != nil {
		req.GroupInfoForSet.Ex = wrapperspb.String(*resp.Ex)
	}
	if resp.NeedVerification != nil {
		req.GroupInfoForSet.NeedVerification = wrapperspb.Int32(*resp.NeedVerification)
	}
	if resp.LookMemberInfo != nil {
		req.GroupInfoForSet.LookMemberInfo = wrapperspb.Int32(*resp.LookMemberInfo)
	}
	if resp.ApplyMemberFriend != nil {
		req.GroupInfoForSet.ApplyMemberFriend = wrapperspb.Int32(*resp.ApplyMemberFriend)
	}
	utils.NotNilReplace(&req.GroupInfoForSet.GroupID, &resp.GroupID)
	utils.NotNilReplace(&req.GroupInfoForSet.GroupName, &resp.GroupName)
	utils.NotNilReplace(&req.GroupInfoForSet.FaceURL, &resp.FaceURL)
	utils.NotNilReplace(&req.GroupInfoForSet.Introduction, &resp.Introduction)
	return nil
}

func CallbackAfterSetGroupInfo(ctx context.Context, cfg *GroupEventCallbackConfig, req *group.SetGroupInfoReq) error {
	if !cfg.BeforeCreateGroup.Enable {
		return nil
	}
	callbackReq := &callbackstruct.CallbackAfterSetGroupInfoReq{
		CallbackCommand: callbackstruct.CallbackAfterSetGroupInfoCommand,
		GroupID:         req.GroupInfoForSet.GroupID,
		Notification:    req.GroupInfoForSet.Notification,
		Introduction:    req.GroupInfoForSet.Introduction,
		FaceURL:         req.GroupInfoForSet.FaceURL,
		GroupName:       req.GroupInfoForSet.GroupName,
	}
	if req.GroupInfoForSet.Ex != nil {
		callbackReq.Ex = &req.GroupInfoForSet.Ex.Value
	}
	if req.GroupInfoForSet.NeedVerification != nil {
		callbackReq.NeedVerification = &req.GroupInfoForSet.NeedVerification.Value
	}
	if req.GroupInfoForSet.LookMemberInfo != nil {
		callbackReq.LookMemberInfo = &req.GroupInfoForSet.LookMemberInfo.Value
	}
	if req.GroupInfoForSet.ApplyMemberFriend != nil {
		callbackReq.ApplyMemberFriend = &req.GroupInfoForSet.ApplyMemberFriend.Value
	}
	resp := &callbackstruct.CallbackAfterSetGroupInfoResp{}
	if err := http.CallBackPostReturn(ctx, cfg.CallbackUrl, callbackReq, resp, cfg.BeforeCreateGroup); err != nil {
		return err
	}
	return nil
}
