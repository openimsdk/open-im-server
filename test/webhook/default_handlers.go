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

package webhook

import (
	"github.com/gin-gonic/gin"
	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
)

// RegisterDefaultHandlers registers default success handlers for all supported OpenIM webhook commands.
func (s *Server) RegisterDefaultHandlers() {
	// User
	RegisterCallback[cbapi.CallbackBeforeUserRegisterReq, cbapi.CallbackBeforeUserRegisterResp](s, cbapi.CallbackBeforeUserRegisterCommand, func(c *gin.Context, req *cbapi.CallbackBeforeUserRegisterReq) (*cbapi.CallbackBeforeUserRegisterResp, error) {
		return &cbapi.CallbackBeforeUserRegisterResp{Users: req.Users}, nil
	})
	RegisterCallback[cbapi.CallbackAfterUserRegisterReq, cbapi.CallbackAfterUserRegisterResp](s, cbapi.CallbackAfterUserRegisterCommand, func(c *gin.Context, req *cbapi.CallbackAfterUserRegisterReq) (*cbapi.CallbackAfterUserRegisterResp, error) {
		return &cbapi.CallbackAfterUserRegisterResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeUpdateUserInfoReq, cbapi.CallbackBeforeUpdateUserInfoResp](s, cbapi.CallbackBeforeUpdateUserInfoCommand, func(c *gin.Context, req *cbapi.CallbackBeforeUpdateUserInfoReq) (*cbapi.CallbackBeforeUpdateUserInfoResp, error) {
		return &cbapi.CallbackBeforeUpdateUserInfoResp{Nickname: req.Nickname, FaceURL: req.FaceURL, Ex: req.Ex}, nil
	})
	RegisterCallback[cbapi.CallbackAfterUpdateUserInfoReq, cbapi.CallbackAfterUpdateUserInfoResp](s, cbapi.CallbackAfterUpdateUserInfoCommand, func(c *gin.Context, req *cbapi.CallbackAfterUpdateUserInfoReq) (*cbapi.CallbackAfterUpdateUserInfoResp, error) {
		return &cbapi.CallbackAfterUpdateUserInfoResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeUpdateUserInfoExReq, cbapi.CallbackBeforeUpdateUserInfoExResp](s, cbapi.CallbackBeforeUpdateUserInfoExCommand, func(c *gin.Context, req *cbapi.CallbackBeforeUpdateUserInfoExReq) (*cbapi.CallbackBeforeUpdateUserInfoExResp, error) {
		return &cbapi.CallbackBeforeUpdateUserInfoExResp{Nickname: req.Nickname, FaceURL: req.FaceURL, Ex: req.Ex}, nil
	})
	RegisterCallback[cbapi.CallbackAfterUpdateUserInfoExReq, cbapi.CallbackAfterUpdateUserInfoExResp](s, cbapi.CallbackAfterUpdateUserInfoExCommand, func(c *gin.Context, req *cbapi.CallbackAfterUpdateUserInfoExReq) (*cbapi.CallbackAfterUpdateUserInfoExResp, error) {
		return &cbapi.CallbackAfterUpdateUserInfoExResp{}, nil
	})

	// User Status Gateway
	RegisterCallback[cbapi.CallbackUserOnlineReq, cbapi.CallbackUserOnlineResp](s, cbapi.CallbackAfterUserOnlineCommand, func(c *gin.Context, req *cbapi.CallbackUserOnlineReq) (*cbapi.CallbackUserOnlineResp, error) {
		return &cbapi.CallbackUserOnlineResp{}, nil
	})
	RegisterCallback[cbapi.CallbackUserOfflineReq, cbapi.CallbackUserOfflineResp](s, cbapi.CallbackAfterUserOfflineCommand, func(c *gin.Context, req *cbapi.CallbackUserOfflineReq) (*cbapi.CallbackUserOfflineResp, error) {
		return &cbapi.CallbackUserOfflineResp{}, nil
	})
	RegisterCallback[cbapi.CallbackUserKickOffReq, cbapi.CallbackUserKickOffResp](s, cbapi.CallbackAfterUserKickOffCommand, func(c *gin.Context, req *cbapi.CallbackUserKickOffReq) (*cbapi.CallbackUserKickOffResp, error) {
		return &cbapi.CallbackUserKickOffResp{}, nil
	})

	// Friend
	RegisterCallback[cbapi.CallbackBeforeAddFriendReq, cbapi.CallbackBeforeAddFriendResp](s, cbapi.CallbackBeforeAddFriendCommand, func(c *gin.Context, req *cbapi.CallbackBeforeAddFriendReq) (*cbapi.CallbackBeforeAddFriendResp, error) {
		return &cbapi.CallbackBeforeAddFriendResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterAddFriendReq, cbapi.CallbackAfterAddFriendResp](s, cbapi.CallbackAfterAddFriendCommand, func(c *gin.Context, req *cbapi.CallbackAfterAddFriendReq) (*cbapi.CallbackAfterAddFriendResp, error) {
		return &cbapi.CallbackAfterAddFriendResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeAddFriendAgreeReq, cbapi.CallbackBeforeAddFriendAgreeResp](s, cbapi.CallbackBeforeAddFriendAgreeCommand, func(c *gin.Context, req *cbapi.CallbackBeforeAddFriendAgreeReq) (*cbapi.CallbackBeforeAddFriendAgreeResp, error) {
		return &cbapi.CallbackBeforeAddFriendAgreeResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterAddFriendAgreeReq, cbapi.CallbackAfterAddFriendAgreeResp](s, cbapi.CallbackAfterAddFriendAgreeCommand, func(c *gin.Context, req *cbapi.CallbackAfterAddFriendAgreeReq) (*cbapi.CallbackAfterAddFriendAgreeResp, error) {
		return &cbapi.CallbackAfterAddFriendAgreeResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterDeleteFriendReq, cbapi.CallbackAfterDeleteFriendResp](s, cbapi.CallbackAfterDeleteFriendCommand, func(c *gin.Context, req *cbapi.CallbackAfterDeleteFriendReq) (*cbapi.CallbackAfterDeleteFriendResp, error) {
		return &cbapi.CallbackAfterDeleteFriendResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeImportFriendsReq, cbapi.CallbackBeforeImportFriendsResp](s, cbapi.CallbackBeforeImportFriendsCommand, func(c *gin.Context, req *cbapi.CallbackBeforeImportFriendsReq) (*cbapi.CallbackBeforeImportFriendsResp, error) {
		return &cbapi.CallbackBeforeImportFriendsResp{FriendUserIDs: req.FriendUserIDs}, nil
	})
	RegisterCallback[cbapi.CallbackAfterImportFriendsReq, cbapi.CallbackAfterImportFriendsResp](s, cbapi.CallbackAfterImportFriendsCommand, func(c *gin.Context, req *cbapi.CallbackAfterImportFriendsReq) (*cbapi.CallbackAfterImportFriendsResp, error) {
		return &cbapi.CallbackAfterImportFriendsResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeSetFriendRemarkReq, cbapi.CallbackBeforeSetFriendRemarkResp](s, cbapi.CallbackBeforeSetFriendRemarkCommand, func(c *gin.Context, req *cbapi.CallbackBeforeSetFriendRemarkReq) (*cbapi.CallbackBeforeSetFriendRemarkResp, error) {
		return &cbapi.CallbackBeforeSetFriendRemarkResp{Remark: req.Remark}, nil
	})
	RegisterCallback[cbapi.CallbackAfterSetFriendRemarkReq, cbapi.CallbackAfterSetFriendRemarkResp](s, cbapi.CallbackAfterSetFriendRemarkCommand, func(c *gin.Context, req *cbapi.CallbackAfterSetFriendRemarkReq) (*cbapi.CallbackAfterSetFriendRemarkResp, error) {
		return &cbapi.CallbackAfterSetFriendRemarkResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeAddBlackReq, cbapi.CallbackBeforeAddBlackResp](s, cbapi.CallbackBeforeAddBlackCommand, func(c *gin.Context, req *cbapi.CallbackBeforeAddBlackReq) (*cbapi.CallbackBeforeAddBlackResp, error) {
		return &cbapi.CallbackBeforeAddBlackResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterRemoveBlackReq, cbapi.CallbackAfterRemoveBlackResp](s, cbapi.CallbackAfterRemoveBlackCommand, func(c *gin.Context, req *cbapi.CallbackAfterRemoveBlackReq) (*cbapi.CallbackAfterRemoveBlackResp, error) {
		return &cbapi.CallbackAfterRemoveBlackResp{}, nil
	})

	// Group
	RegisterCallback[cbapi.CallbackBeforeCreateGroupReq, cbapi.CallbackBeforeCreateGroupResp](s, cbapi.CallbackBeforeCreateGroupCommand, func(c *gin.Context, req *cbapi.CallbackBeforeCreateGroupReq) (*cbapi.CallbackBeforeCreateGroupResp, error) {
		return &cbapi.CallbackBeforeCreateGroupResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterCreateGroupReq, cbapi.CallbackAfterCreateGroupResp](s, cbapi.CallbackAfterCreateGroupCommand, func(c *gin.Context, req *cbapi.CallbackAfterCreateGroupReq) (*cbapi.CallbackAfterCreateGroupResp, error) {
		return &cbapi.CallbackAfterCreateGroupResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeMembersJoinGroupReq, cbapi.CallbackBeforeMembersJoinGroupResp](s, cbapi.CallbackBeforeMembersJoinGroupCommand, func(c *gin.Context, req *cbapi.CallbackBeforeMembersJoinGroupReq) (*cbapi.CallbackBeforeMembersJoinGroupResp, error) {
		return &cbapi.CallbackBeforeMembersJoinGroupResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeMembersJoinGroupReq, cbapi.CallbackBeforeMembersJoinGroupResp](s, cbapi.CallbackBeforeJoinGroupCommand, func(c *gin.Context, req *cbapi.CallbackBeforeMembersJoinGroupReq) (*cbapi.CallbackBeforeMembersJoinGroupResp, error) {
		return &cbapi.CallbackBeforeMembersJoinGroupResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterJoinGroupReq, cbapi.CallbackAfterJoinGroupResp](s, cbapi.CallbackAfterJoinGroupCommand, func(c *gin.Context, req *cbapi.CallbackAfterJoinGroupReq) (*cbapi.CallbackAfterJoinGroupResp, error) {
		return &cbapi.CallbackAfterJoinGroupResp{}, nil
	})
	RegisterCallback[cbapi.CallbackQuitGroupReq, cbapi.CallbackQuitGroupResp](s, cbapi.CallbackAfterQuitGroupCommand, func(c *gin.Context, req *cbapi.CallbackQuitGroupReq) (*cbapi.CallbackQuitGroupResp, error) {
		return &cbapi.CallbackQuitGroupResp{}, nil
	})
	RegisterCallback[cbapi.CallbackKillGroupMemberReq, cbapi.CallbackKillGroupMemberResp](s, cbapi.CallbackAfterKickGroupCommand, func(c *gin.Context, req *cbapi.CallbackKillGroupMemberReq) (*cbapi.CallbackKillGroupMemberResp, error) {
		return &cbapi.CallbackKillGroupMemberResp{}, nil
	})
	RegisterCallback[cbapi.CallbackDisMissGroupReq, cbapi.CallbackDisMissGroupResp](s, cbapi.CallbackAfterDisMissGroupCommand, func(c *gin.Context, req *cbapi.CallbackDisMissGroupReq) (*cbapi.CallbackDisMissGroupResp, error) {
		return &cbapi.CallbackDisMissGroupResp{}, nil
	})
	RegisterCallback[cbapi.CallbackTransferGroupOwnerReq, cbapi.CallbackTransferGroupOwnerResp](s, cbapi.CallbackAfterTransferGroupOwnerCommand, func(c *gin.Context, req *cbapi.CallbackTransferGroupOwnerReq) (*cbapi.CallbackTransferGroupOwnerResp, error) {
		return &cbapi.CallbackTransferGroupOwnerResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeInviteUserToGroupReq, cbapi.CallbackBeforeInviteUserToGroupResp](s, cbapi.CallbackBeforeInviteJoinGroupCommand, func(c *gin.Context, req *cbapi.CallbackBeforeInviteUserToGroupReq) (*cbapi.CallbackBeforeInviteUserToGroupResp, error) {
		return &cbapi.CallbackBeforeInviteUserToGroupResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeSetGroupInfoReq, cbapi.CallbackBeforeSetGroupInfoResp](s, cbapi.CallbackBeforeSetGroupInfoCommand, func(c *gin.Context, req *cbapi.CallbackBeforeSetGroupInfoReq) (*cbapi.CallbackBeforeSetGroupInfoResp, error) {
		return &cbapi.CallbackBeforeSetGroupInfoResp{GroupID: req.GroupID, GroupName: req.GroupName, Notification: req.Notification, Introduction: req.Introduction, FaceURL: req.FaceURL}, nil
	})
	RegisterCallback[cbapi.CallbackAfterSetGroupInfoReq, cbapi.CallbackAfterSetGroupInfoResp](s, cbapi.CallbackAfterSetGroupInfoCommand, func(c *gin.Context, req *cbapi.CallbackAfterSetGroupInfoReq) (*cbapi.CallbackAfterSetGroupInfoResp, error) {
		return &cbapi.CallbackAfterSetGroupInfoResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeSetGroupInfoExReq, cbapi.CallbackBeforeSetGroupInfoExResp](s, cbapi.CallbackBeforeSetGroupInfoExCommand, func(c *gin.Context, req *cbapi.CallbackBeforeSetGroupInfoExReq) (*cbapi.CallbackBeforeSetGroupInfoExResp, error) {
		return &cbapi.CallbackBeforeSetGroupInfoExResp{GroupID: req.GroupID, GroupName: req.GroupName, Notification: req.Notification, Introduction: req.Introduction, FaceURL: req.FaceURL, Ex: req.Ex, NeedVerification: req.NeedVerification, LookMemberInfo: req.LookMemberInfo, ApplyMemberFriend: req.ApplyMemberFriend}, nil
	})
	RegisterCallback[cbapi.CallbackAfterSetGroupInfoExReq, cbapi.CallbackAfterSetGroupInfoExResp](s, cbapi.CallbackAfterSetGroupInfoExCommand, func(c *gin.Context, req *cbapi.CallbackAfterSetGroupInfoExReq) (*cbapi.CallbackAfterSetGroupInfoExResp, error) {
		return &cbapi.CallbackAfterSetGroupInfoExResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeSetGroupMemberInfoReq, cbapi.CallbackBeforeSetGroupMemberInfoResp](s, cbapi.CallbackBeforeSetGroupMemberInfoCommand, func(c *gin.Context, req *cbapi.CallbackBeforeSetGroupMemberInfoReq) (*cbapi.CallbackBeforeSetGroupMemberInfoResp, error) {
		return &cbapi.CallbackBeforeSetGroupMemberInfoResp{Nickname: req.Nickname, FaceURL: req.FaceURL, RoleLevel: req.RoleLevel, Ex: req.Ex}, nil
	})
	RegisterCallback[cbapi.CallbackAfterSetGroupMemberInfoReq, cbapi.CallbackAfterSetGroupMemberInfoResp](s, cbapi.CallbackAfterSetGroupMemberInfoCommand, func(c *gin.Context, req *cbapi.CallbackAfterSetGroupMemberInfoReq) (*cbapi.CallbackAfterSetGroupMemberInfoResp, error) {
		return &cbapi.CallbackAfterSetGroupMemberInfoResp{}, nil
	})

	// Message & Push & Revoke
	RegisterCallback[cbapi.CallbackBeforeSendSingleMsgReq, cbapi.CallbackBeforeSendSingleMsgResp](s, cbapi.CallbackBeforeSendSingleMsgCommand, func(c *gin.Context, req *cbapi.CallbackBeforeSendSingleMsgReq) (*cbapi.CallbackBeforeSendSingleMsgResp, error) {
		return &cbapi.CallbackBeforeSendSingleMsgResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterSendSingleMsgReq, cbapi.CallbackAfterSendSingleMsgResp](s, cbapi.CallbackAfterSendSingleMsgCommand, func(c *gin.Context, req *cbapi.CallbackAfterSendSingleMsgReq) (*cbapi.CallbackAfterSendSingleMsgResp, error) {
		return &cbapi.CallbackAfterSendSingleMsgResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeSendGroupMsgReq, cbapi.CallbackBeforeSendGroupMsgResp](s, cbapi.CallbackBeforeSendGroupMsgCommand, func(c *gin.Context, req *cbapi.CallbackBeforeSendGroupMsgReq) (*cbapi.CallbackBeforeSendGroupMsgResp, error) {
		return &cbapi.CallbackBeforeSendGroupMsgResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterSendGroupMsgReq, cbapi.CallbackAfterSendGroupMsgResp](s, cbapi.CallbackAfterSendGroupMsgCommand, func(c *gin.Context, req *cbapi.CallbackAfterSendGroupMsgReq) (*cbapi.CallbackAfterSendGroupMsgResp, error) {
		return &cbapi.CallbackAfterSendGroupMsgResp{}, nil
	})
	RegisterCallback[cbapi.CallbackMsgModifyCommandReq, cbapi.CallbackMsgModifyCommandResp](s, cbapi.CallbackBeforeMsgModifyCommand, func(c *gin.Context, req *cbapi.CallbackMsgModifyCommandReq) (*cbapi.CallbackMsgModifyCommandResp, error) {
		return &cbapi.CallbackMsgModifyCommandResp{}, nil
	})
	RegisterCallback[cbapi.CallbackSingleMsgReadReq, cbapi.CallbackSingleMsgReadResp](s, cbapi.CallbackAfterSingleMsgReadCommand, func(c *gin.Context, req *cbapi.CallbackSingleMsgReadReq) (*cbapi.CallbackSingleMsgReadResp, error) {
		return &cbapi.CallbackSingleMsgReadResp{}, nil
	})
	RegisterCallback[cbapi.CallbackGroupMsgReadReq, cbapi.CallbackGroupMsgReadResp](s, cbapi.CallbackAfterGroupMsgReadCommand, func(c *gin.Context, req *cbapi.CallbackGroupMsgReadReq) (*cbapi.CallbackGroupMsgReadResp, error) {
		return &cbapi.CallbackGroupMsgReadResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterMsgSaveDBReq, cbapi.CallbackAfterMsgSaveDBResp](s, cbapi.CallbackAfterMsgSaveDBCommand, func(c *gin.Context, req *cbapi.CallbackAfterMsgSaveDBReq) (*cbapi.CallbackAfterMsgSaveDBResp, error) {
		return &cbapi.CallbackAfterMsgSaveDBResp{}, nil
	})
	RegisterCallback[cbapi.CallbackAfterRevokeMsgReq, cbapi.CallbackAfterRevokeMsgResp](s, cbapi.CallbackAfterRevokeMsgCommand, func(c *gin.Context, req *cbapi.CallbackAfterRevokeMsgReq) (*cbapi.CallbackAfterRevokeMsgResp, error) {
		return &cbapi.CallbackAfterRevokeMsgResp{}, nil
	})
	RegisterCallback[cbapi.CallbackBeforePushReq, cbapi.CallbackBeforePushResp](s, cbapi.CallbackBeforeOfflinePushCommand, func(c *gin.Context, req *cbapi.CallbackBeforePushReq) (*cbapi.CallbackBeforePushResp, error) {
		return &cbapi.CallbackBeforePushResp{UserIDs: req.UserIDList, OfflinePushInfo: req.OfflinePushInfo}, nil
	})
	RegisterCallback[cbapi.CallbackBeforePushReq, cbapi.CallbackBeforePushResp](s, cbapi.CallbackBeforeOnlinePushCommand, func(c *gin.Context, req *cbapi.CallbackBeforePushReq) (*cbapi.CallbackBeforePushResp, error) {
		return &cbapi.CallbackBeforePushResp{UserIDs: req.UserIDList}, nil
	})
	RegisterCallback[cbapi.CallbackBeforeSuperGroupOnlinePushReq, cbapi.CallbackBeforeSuperGroupOnlinePushResp](s, cbapi.CallbackBeforeGroupOnlinePushCommand, func(c *gin.Context, req *cbapi.CallbackBeforeSuperGroupOnlinePushReq) (*cbapi.CallbackBeforeSuperGroupOnlinePushResp, error) {
		return &cbapi.CallbackBeforeSuperGroupOnlinePushResp{}, nil
	})
}
