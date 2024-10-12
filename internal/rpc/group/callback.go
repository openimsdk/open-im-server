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
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
)

// CallbackBeforeCreateGroup callback before create group.
func (s *groupServer) webhookBeforeCreateGroup(ctx context.Context, before *config.BeforeConfig, req *group.CreateGroupReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
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

		if err := s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		datautil.NotNilReplace(&req.GroupInfo.GroupID, resp.GroupID)
		datautil.NotNilReplace(&req.GroupInfo.GroupName, resp.GroupName)
		datautil.NotNilReplace(&req.GroupInfo.Notification, resp.Notification)
		datautil.NotNilReplace(&req.GroupInfo.Introduction, resp.Introduction)
		datautil.NotNilReplace(&req.GroupInfo.FaceURL, resp.FaceURL)
		datautil.NotNilReplace(&req.GroupInfo.OwnerUserID, resp.OwnerUserID)
		datautil.NotNilReplace(&req.GroupInfo.Ex, resp.Ex)
		datautil.NotNilReplace(&req.GroupInfo.Status, resp.Status)
		datautil.NotNilReplace(&req.GroupInfo.CreatorUserID, resp.CreatorUserID)
		datautil.NotNilReplace(&req.GroupInfo.GroupType, resp.GroupType)
		datautil.NotNilReplace(&req.GroupInfo.NeedVerification, resp.NeedVerification)
		datautil.NotNilReplace(&req.GroupInfo.LookMemberInfo, resp.LookMemberInfo)
		return nil
	})
}

func (s *groupServer) webhookAfterCreateGroup(ctx context.Context, after *config.AfterConfig, req *group.CreateGroupReq) {
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
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterCreateGroupResp{}, after)
}

func (s *groupServer) webhookBeforeMembersJoinGroup(ctx context.Context, before *config.BeforeConfig, groupMembers []*model.GroupMember, groupID string, groupEx string) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		groupMembersMap := datautil.SliceToMap(groupMembers, func(e *model.GroupMember) string {
			return e.UserID
		})
		var groupMembersCallback []*callbackstruct.CallbackGroupMember

		for _, member := range groupMembers {
			groupMembersCallback = append(groupMembersCallback, &callbackstruct.CallbackGroupMember{
				UserID: member.UserID,
				Ex:     member.Ex,
			})
		}

		cbReq := &callbackstruct.CallbackBeforeMembersJoinGroupReq{
			CallbackCommand: callbackstruct.CallbackBeforeMembersJoinGroupCommand,
			GroupID:         groupID,
			MembersList:     groupMembersCallback,
			GroupEx:         groupEx,
		}
		resp := &callbackstruct.CallbackBeforeMembersJoinGroupResp{}

		if err := s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		for _, memberCallbackResp := range resp.MemberCallbackList {
			if _, ok := groupMembersMap[(*memberCallbackResp.UserID)]; ok {
				if memberCallbackResp.MuteEndTime != nil {
					groupMembersMap[(*memberCallbackResp.UserID)].MuteEndTime = time.UnixMilli(*memberCallbackResp.MuteEndTime)
				}

				datautil.NotNilReplace(&groupMembersMap[(*memberCallbackResp.UserID)].FaceURL, memberCallbackResp.FaceURL)
				datautil.NotNilReplace(&groupMembersMap[(*memberCallbackResp.UserID)].Ex, memberCallbackResp.Ex)
				datautil.NotNilReplace(&groupMembersMap[(*memberCallbackResp.UserID)].Nickname, memberCallbackResp.Nickname)
				datautil.NotNilReplace(&groupMembersMap[(*memberCallbackResp.UserID)].RoleLevel, memberCallbackResp.RoleLevel)
			}
		}

		return nil
	})
}

func (s *groupServer) webhookBeforeSetGroupMemberInfo(ctx context.Context, before *config.BeforeConfig, req *group.SetGroupMemberInfo) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := callbackstruct.CallbackBeforeSetGroupMemberInfoReq{
			CallbackCommand: callbackstruct.CallbackBeforeSetGroupMemberInfoCommand,
			GroupID:         req.GroupID,
			UserID:          req.UserID,
		}
		if req.Nickname != nil {
			cbReq.Nickname = &req.Nickname.Value
		}
		if req.FaceURL != nil {
			cbReq.FaceURL = &req.FaceURL.Value
		}
		if req.RoleLevel != nil {
			cbReq.RoleLevel = &req.RoleLevel.Value
		}
		if req.Ex != nil {
			cbReq.Ex = &req.Ex.Value
		}
		resp := &callbackstruct.CallbackBeforeSetGroupMemberInfoResp{}
		if err := s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
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
	})
}

func (s *groupServer) webhookAfterSetGroupMemberInfo(ctx context.Context, after *config.AfterConfig, req *group.SetGroupMemberInfo) {
	cbReq := callbackstruct.CallbackAfterSetGroupMemberInfoReq{
		CallbackCommand: callbackstruct.CallbackAfterSetGroupMemberInfoCommand,
		GroupID:         req.GroupID,
		UserID:          req.UserID,
	}
	if req.Nickname != nil {
		cbReq.Nickname = &req.Nickname.Value
	}
	if req.FaceURL != nil {
		cbReq.FaceURL = &req.FaceURL.Value
	}
	if req.RoleLevel != nil {
		cbReq.RoleLevel = &req.RoleLevel.Value
	}
	if req.Ex != nil {
		cbReq.Ex = &req.Ex.Value
	}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterSetGroupMemberInfoResp{}, after)
}

func (s *groupServer) webhookAfterQuitGroup(ctx context.Context, after *config.AfterConfig, req *group.QuitGroupReq) {
	cbReq := &callbackstruct.CallbackQuitGroupReq{
		CallbackCommand: callbackstruct.CallbackAfterQuitGroupCommand,
		GroupID:         req.GroupID,
		UserID:          req.UserID,
	}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackQuitGroupResp{}, after)
}

func (s *groupServer) webhookAfterKickGroupMember(ctx context.Context, after *config.AfterConfig, req *group.KickGroupMemberReq) {
	cbReq := &callbackstruct.CallbackKillGroupMemberReq{
		CallbackCommand: callbackstruct.CallbackAfterKickGroupCommand,
		GroupID:         req.GroupID,
		KickedUserIDs:   req.KickedUserIDs,
		Reason:          req.Reason,
	}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackKillGroupMemberResp{}, after)
}

func (s *groupServer) webhookAfterDismissGroup(ctx context.Context, after *config.AfterConfig, req *callbackstruct.CallbackDisMissGroupReq) {
	req.CallbackCommand = callbackstruct.CallbackAfterDisMissGroupCommand
	s.webhookClient.AsyncPost(ctx, req.GetCallbackCommand(), req, &callbackstruct.CallbackDisMissGroupResp{}, after)
}

func (s *groupServer) webhookBeforeApplyJoinGroup(ctx context.Context, before *config.BeforeConfig, req *callbackstruct.CallbackJoinGroupReq) (err error) {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		req.CallbackCommand = callbackstruct.CallbackBeforeJoinGroupCommand
		resp := &callbackstruct.CallbackJoinGroupResp{}
		if err := s.webhookClient.SyncPost(ctx, req.GetCallbackCommand(), req, resp, before); err != nil {
			return err
		}
		return nil
	})
}

func (s *groupServer) webhookAfterTransferGroupOwner(ctx context.Context, after *config.AfterConfig, req *group.TransferGroupOwnerReq) {
	cbReq := &callbackstruct.CallbackTransferGroupOwnerReq{
		CallbackCommand: callbackstruct.CallbackAfterTransferGroupOwnerCommand,
		GroupID:         req.GroupID,
		OldOwnerUserID:  req.OldOwnerUserID,
		NewOwnerUserID:  req.NewOwnerUserID,
	}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackTransferGroupOwnerResp{}, after)
}

func (s *groupServer) webhookBeforeInviteUserToGroup(ctx context.Context, before *config.BeforeConfig, req *group.InviteUserToGroupReq) (err error) {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &callbackstruct.CallbackBeforeInviteUserToGroupReq{
			CallbackCommand: callbackstruct.CallbackBeforeInviteJoinGroupCommand,
			OperationID:     mcontext.GetOperationID(ctx),
			GroupID:         req.GroupID,
			Reason:          req.Reason,
			InvitedUserIDs:  req.InvitedUserIDs,
		}

		resp := &callbackstruct.CallbackBeforeInviteUserToGroupResp{}
		if err := s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		// Handle the scenario where certain members are refused
		// You might want to update the req.Members list or handle it as per your business logic

		// if len(resp.RefusedMembersAccount) > 0 {
		// implement members are refused
		// }

		return nil
	})
}

func (s *groupServer) webhookAfterJoinGroup(ctx context.Context, after *config.AfterConfig, req *group.JoinGroupReq) {
	cbReq := &callbackstruct.CallbackAfterJoinGroupReq{
		CallbackCommand: callbackstruct.CallbackAfterJoinGroupCommand,
		OperationID:     mcontext.GetOperationID(ctx),
		GroupID:         req.GroupID,
		ReqMessage:      req.ReqMessage,
		JoinSource:      req.JoinSource,
		InviterUserID:   req.InviterUserID,
	}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterJoinGroupResp{}, after)
}

func (s *groupServer) webhookBeforeSetGroupInfo(ctx context.Context, before *config.BeforeConfig, req *group.SetGroupInfoReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &callbackstruct.CallbackBeforeSetGroupInfoReq{
			CallbackCommand: callbackstruct.CallbackBeforeSetGroupInfoCommand,
			GroupID:         req.GroupInfoForSet.GroupID,
			Notification:    req.GroupInfoForSet.Notification,
			Introduction:    req.GroupInfoForSet.Introduction,
			FaceURL:         req.GroupInfoForSet.FaceURL,
			GroupName:       req.GroupInfoForSet.GroupName,
		}
		if req.GroupInfoForSet.Ex != nil {
			cbReq.Ex = req.GroupInfoForSet.Ex.Value
		}
		log.ZDebug(ctx, "debug CallbackBeforeSetGroupInfo", "ex", cbReq.Ex)
		if req.GroupInfoForSet.NeedVerification != nil {
			cbReq.NeedVerification = req.GroupInfoForSet.NeedVerification.Value
		}
		if req.GroupInfoForSet.LookMemberInfo != nil {
			cbReq.LookMemberInfo = req.GroupInfoForSet.LookMemberInfo.Value
		}
		if req.GroupInfoForSet.ApplyMemberFriend != nil {
			cbReq.ApplyMemberFriend = req.GroupInfoForSet.ApplyMemberFriend.Value
		}
		resp := &callbackstruct.CallbackBeforeSetGroupInfoResp{}

		if err := s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
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
		datautil.NotNilReplace(&req.GroupInfoForSet.GroupID, &resp.GroupID)
		datautil.NotNilReplace(&req.GroupInfoForSet.GroupName, &resp.GroupName)
		datautil.NotNilReplace(&req.GroupInfoForSet.FaceURL, &resp.FaceURL)
		datautil.NotNilReplace(&req.GroupInfoForSet.Introduction, &resp.Introduction)
		return nil
	})
}

func (s *groupServer) webhookAfterSetGroupInfo(ctx context.Context, after *config.AfterConfig, req *group.SetGroupInfoReq) {
	cbReq := &callbackstruct.CallbackAfterSetGroupInfoReq{
		CallbackCommand: callbackstruct.CallbackAfterSetGroupInfoCommand,
		GroupID:         req.GroupInfoForSet.GroupID,
		Notification:    req.GroupInfoForSet.Notification,
		Introduction:    req.GroupInfoForSet.Introduction,
		FaceURL:         req.GroupInfoForSet.FaceURL,
		GroupName:       req.GroupInfoForSet.GroupName,
	}
	if req.GroupInfoForSet.Ex != nil {
		cbReq.Ex = &req.GroupInfoForSet.Ex.Value
	}
	if req.GroupInfoForSet.NeedVerification != nil {
		cbReq.NeedVerification = &req.GroupInfoForSet.NeedVerification.Value
	}
	if req.GroupInfoForSet.LookMemberInfo != nil {
		cbReq.LookMemberInfo = &req.GroupInfoForSet.LookMemberInfo.Value
	}
	if req.GroupInfoForSet.ApplyMemberFriend != nil {
		cbReq.ApplyMemberFriend = &req.GroupInfoForSet.ApplyMemberFriend.Value
	}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterSetGroupInfoResp{}, after)
}

func (s *groupServer) webhookBeforeSetGroupInfoEx(ctx context.Context, before *config.BeforeConfig, req *group.SetGroupInfoExReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &callbackstruct.CallbackBeforeSetGroupInfoExReq{
			CallbackCommand: callbackstruct.CallbackBeforeSetGroupInfoExCommand,
			GroupID:         req.GroupID,
			GroupName:       req.GroupName,
			Notification:    req.Notification,
			Introduction:    req.Introduction,
			FaceURL:         req.FaceURL,
		}

		if req.Ex != nil {
			cbReq.Ex = req.Ex
		}
		log.ZDebug(ctx, "debug CallbackBeforeSetGroupInfoEx", "ex", cbReq.Ex)

		if req.NeedVerification != nil {
			cbReq.NeedVerification = req.NeedVerification
		}
		if req.LookMemberInfo != nil {
			cbReq.LookMemberInfo = req.LookMemberInfo
		}
		if req.ApplyMemberFriend != nil {
			cbReq.ApplyMemberFriend = req.ApplyMemberFriend
		}

		resp := &callbackstruct.CallbackBeforeSetGroupInfoExResp{}

		if err := s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		datautil.NotNilReplace(&req.GroupID, &resp.GroupID)
		datautil.NotNilReplace(&req.GroupName, &resp.GroupName)
		datautil.NotNilReplace(&req.FaceURL, &resp.FaceURL)
		datautil.NotNilReplace(&req.Introduction, &resp.Introduction)
		datautil.NotNilReplace(&req.Ex, &resp.Ex)
		datautil.NotNilReplace(&req.NeedVerification, &resp.NeedVerification)
		datautil.NotNilReplace(&req.LookMemberInfo, &resp.LookMemberInfo)
		datautil.NotNilReplace(&req.ApplyMemberFriend, &resp.ApplyMemberFriend)

		return nil
	})
}

func (s *groupServer) webhookAfterSetGroupInfoEx(ctx context.Context, after *config.AfterConfig, req *group.SetGroupInfoExReq) {
	cbReq := &callbackstruct.CallbackAfterSetGroupInfoExReq{
		CallbackCommand: callbackstruct.CallbackAfterSetGroupInfoExCommand,
		GroupID:         req.GroupID,
		GroupName:       req.GroupName,
		Notification:    req.Notification,
		Introduction:    req.Introduction,
		FaceURL:         req.FaceURL,
	}

	if req.Ex != nil {
		cbReq.Ex = req.Ex
	}
	if req.NeedVerification != nil {
		cbReq.NeedVerification = req.NeedVerification
	}
	if req.LookMemberInfo != nil {
		cbReq.LookMemberInfo = req.LookMemberInfo
	}
	if req.ApplyMemberFriend != nil {
		cbReq.ApplyMemberFriend = req.ApplyMemberFriend
	}

	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterSetGroupInfoExResp{}, after)
}
