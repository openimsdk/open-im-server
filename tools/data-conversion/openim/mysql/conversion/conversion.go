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

package conversion

import (
	"github.com/OpenIMSDK/protocol/constant"

	v2 "github.com/openimsdk/open-im-server/v3/tools/data-conversion/openim/mysql/v2"
	v3 "github.com/openimsdk/open-im-server/v3/tools/data-conversion/openim/mysql/v3"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/utils"
)

func Friend(v v2.Friend) (v3.FriendModel, bool) {
	utils.InitTime(&v.CreateTime)
	return v3.FriendModel{
		OwnerUserID:    v.OwnerUserID,
		FriendUserID:   v.FriendUserID,
		Remark:         v.Remark,
		CreateTime:     v.CreateTime,
		AddSource:      v.AddSource,
		OperatorUserID: v.OperatorUserID,
		Ex:             v.Ex,
	}, true
}

func FriendRequest(v v2.FriendRequest) (v3.FriendRequestModel, bool) {
	utils.InitTime(&v.CreateTime, &v.HandleTime)
	return v3.FriendRequestModel{
		FromUserID:    v.FromUserID,
		ToUserID:      v.ToUserID,
		HandleResult:  v.HandleResult,
		ReqMsg:        v.ReqMsg,
		CreateTime:    v.CreateTime,
		HandlerUserID: v.HandlerUserID,
		HandleMsg:     v.HandleMsg,
		HandleTime:    v.HandleTime,
		Ex:            v.Ex,
	}, true
}

func Group(v v2.Group) (v3.GroupModel, bool) {
	switch v.GroupType {
	case constant.WorkingGroup, constant.NormalGroup:
		v.GroupType = constant.WorkingGroup
	default:
		return v3.GroupModel{}, false
	}
	utils.InitTime(&v.CreateTime, &v.NotificationUpdateTime)
	return v3.GroupModel{
		GroupID:                v.GroupID,
		GroupName:              v.GroupName,
		Notification:           v.Notification,
		Introduction:           v.Introduction,
		FaceURL:                v.FaceURL,
		CreateTime:             v.CreateTime,
		Ex:                     v.Ex,
		Status:                 v.Status,
		CreatorUserID:          v.CreatorUserID,
		GroupType:              v.GroupType,
		NeedVerification:       v.NeedVerification,
		LookMemberInfo:         v.LookMemberInfo,
		ApplyMemberFriend:      v.ApplyMemberFriend,
		NotificationUpdateTime: v.NotificationUpdateTime,
		NotificationUserID:     v.NotificationUserID,
	}, true
}

func GroupMember(v v2.GroupMember) (v3.GroupMemberModel, bool) {
	utils.InitTime(&v.JoinTime, &v.MuteEndTime)
	return v3.GroupMemberModel{
		GroupID:        v.GroupID,
		UserID:         v.UserID,
		Nickname:       v.Nickname,
		FaceURL:        v.FaceURL,
		RoleLevel:      v.RoleLevel,
		JoinTime:       v.JoinTime,
		JoinSource:     v.JoinSource,
		InviterUserID:  v.InviterUserID,
		OperatorUserID: v.OperatorUserID,
		MuteEndTime:    v.MuteEndTime,
		Ex:             v.Ex,
	}, true
}

func GroupRequest(v v2.GroupRequest) (v3.GroupRequestModel, bool) {
	utils.InitTime(&v.ReqTime, &v.HandledTime)
	return v3.GroupRequestModel{
		UserID:        v.UserID,
		GroupID:       v.GroupID,
		HandleResult:  v.HandleResult,
		ReqMsg:        v.ReqMsg,
		HandledMsg:    v.HandledMsg,
		ReqTime:       v.ReqTime,
		HandleUserID:  v.HandleUserID,
		HandledTime:   v.HandledTime,
		JoinSource:    v.JoinSource,
		InviterUserID: v.InviterUserID,
		Ex:            v.Ex,
	}, true
}

func User(v v2.User) (v3.UserModel, bool) {
	utils.InitTime(&v.CreateTime)
	return v3.UserModel{
		UserID:           v.UserID,
		Nickname:         v.Nickname,
		FaceURL:          v.FaceURL,
		Ex:               v.Ex,
		CreateTime:       v.CreateTime,
		AppMangerLevel:   v.AppMangerLevel,
		GlobalRecvMsgOpt: v.GlobalRecvMsgOpt,
	}, true
}
