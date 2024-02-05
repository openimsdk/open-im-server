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

package pkg

import (
	"time"

	mongomodel "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	mysqlmodel "github.com/openimsdk/open-im-server/v3/tools/data-conversion/openim/mysql/v3"
	mongomodelrtc "github.com/openimsdk/open-im-server/v3/tools/up35/pkg/internal/rtc/mongo/table"
	mysqlmodelrtc "github.com/openimsdk/open-im-server/v3/tools/up35/pkg/internal/rtc/mysql"
)

type convert struct{}

func (convert) User(v mysqlmodel.UserModel) mongomodel.UserModel {
	return mongomodel.UserModel{
		UserID:           v.UserID,
		Nickname:         v.Nickname,
		FaceURL:          v.FaceURL,
		Ex:               v.Ex,
		AppMangerLevel:   v.AppMangerLevel,
		GlobalRecvMsgOpt: v.GlobalRecvMsgOpt,
		CreateTime:       v.CreateTime,
	}
}

func (convert) Friend(v mysqlmodel.FriendModel) mongomodel.FriendModel {
	return mongomodel.FriendModel{
		OwnerUserID:    v.OwnerUserID,
		FriendUserID:   v.FriendUserID,
		Remark:         v.Remark,
		CreateTime:     v.CreateTime,
		AddSource:      v.AddSource,
		OperatorUserID: v.OperatorUserID,
		Ex:             v.Ex,
	}
}

func (convert) FriendRequest(v mysqlmodel.FriendRequestModel) mongomodel.FriendRequestModel {
	return mongomodel.FriendRequestModel{
		FromUserID:    v.FromUserID,
		ToUserID:      v.ToUserID,
		HandleResult:  v.HandleResult,
		ReqMsg:        v.ReqMsg,
		CreateTime:    v.CreateTime,
		HandlerUserID: v.HandlerUserID,
		HandleMsg:     v.HandleMsg,
		HandleTime:    v.HandleTime,
		Ex:            v.Ex,
	}
}

func (convert) Black(v mysqlmodel.BlackModel) mongomodel.BlackModel {
	return mongomodel.BlackModel{
		OwnerUserID:    v.OwnerUserID,
		BlockUserID:    v.BlockUserID,
		CreateTime:     v.CreateTime,
		AddSource:      v.AddSource,
		OperatorUserID: v.OperatorUserID,
		Ex:             v.Ex,
	}
}

func (convert) Group(v mysqlmodel.GroupModel) mongomodel.GroupModel {
	return mongomodel.GroupModel{
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
	}
}

func (convert) GroupMember(v mysqlmodel.GroupMemberModel) mongomodel.GroupMemberModel {
	return mongomodel.GroupMemberModel{
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
	}
}

func (convert) GroupRequest(v mysqlmodel.GroupRequestModel) mongomodel.GroupRequestModel {
	return mongomodel.GroupRequestModel{
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
	}
}

func (convert) Conversation(v mysqlmodel.ConversationModel) mongomodel.ConversationModel {
	return mongomodel.ConversationModel{
		OwnerUserID:           v.OwnerUserID,
		ConversationID:        v.ConversationID,
		ConversationType:      v.ConversationType,
		UserID:                v.UserID,
		GroupID:               v.GroupID,
		RecvMsgOpt:            v.RecvMsgOpt,
		IsPinned:              v.IsPinned,
		IsPrivateChat:         v.IsPrivateChat,
		BurnDuration:          v.BurnDuration,
		GroupAtType:           v.GroupAtType,
		AttachedInfo:          v.AttachedInfo,
		Ex:                    v.Ex,
		MaxSeq:                v.MaxSeq,
		MinSeq:                v.MinSeq,
		CreateTime:            v.CreateTime,
		IsMsgDestruct:         v.IsMsgDestruct,
		MsgDestructTime:       v.MsgDestructTime,
		LatestMsgDestructTime: v.LatestMsgDestructTime,
	}
}

func (convert) Object(engine string) func(v mysqlmodel.ObjectModel) mongomodel.ObjectModel {
	return func(v mysqlmodel.ObjectModel) mongomodel.ObjectModel {
		return mongomodel.ObjectModel{
			Name:        v.Name,
			UserID:      v.UserID,
			Hash:        v.Hash,
			Engine:      engine,
			Key:         v.Key,
			Size:        v.Size,
			ContentType: v.ContentType,
			Group:       v.Cause,
			CreateTime:  v.CreateTime,
		}
	}
}

func (convert) Log(v mysqlmodel.Log) mongomodel.LogModel {
	return mongomodel.LogModel{
		LogID:      v.LogID,
		Platform:   v.Platform,
		UserID:     v.UserID,
		CreateTime: v.CreateTime,
		Url:        v.Url,
		FileName:   v.FileName,
		SystemType: v.SystemType,
		Version:    v.Version,
		Ex:         v.Ex,
	}
}

func (convert) SignalModel(v mysqlmodelrtc.SignalModel) mongomodelrtc.SignalModel {
	return mongomodelrtc.SignalModel{
		SID:           v.SID,
		InviterUserID: v.InviterUserID,
		CustomData:    v.CustomData,
		GroupID:       v.GroupID,
		RoomID:        v.RoomID,
		Timeout:       v.Timeout,
		MediaType:     v.MediaType,
		PlatformID:    v.PlatformID,
		SessionType:   v.SessionType,
		InitiateTime:  v.InitiateTime,
		EndTime:       v.EndTime,
		FileURL:       v.FileURL,
		Title:         v.Title,
		Desc:          v.Desc,
		Ex:            v.Ex,
		IOSPushSound:  v.IOSPushSound,
		IOSBadgeCount: v.IOSBadgeCount,
		SignalInfo:    v.SignalInfo,
	}
}

func (convert) SignalInvitationModel(v mysqlmodelrtc.SignalInvitationModel) mongomodelrtc.SignalInvitationModel {
	return mongomodelrtc.SignalInvitationModel{
		SID:          v.SID,
		UserID:       v.UserID,
		Status:       v.Status,
		InitiateTime: v.InitiateTime,
		HandleTime:   v.HandleTime,
	}
}

func (convert) Meeting(v mysqlmodelrtc.MeetingInfo) mongomodelrtc.MeetingInfo {
	return mongomodelrtc.MeetingInfo{
		RoomID:      v.RoomID,
		MeetingName: v.MeetingName,
		HostUserID:  v.HostUserID,
		Status:      v.Status,
		StartTime:   time.Unix(v.StartTime, 0),
		EndTime:     time.Unix(v.EndTime, 0),
		CreateTime:  v.CreateTime,
		Ex:          v.Ex,
	}
}

func (convert) MeetingInvitationInfo(v mysqlmodelrtc.MeetingInvitationInfo) mongomodelrtc.MeetingInvitationInfo {
	return mongomodelrtc.MeetingInvitationInfo{
		RoomID:     v.RoomID,
		UserID:     v.UserID,
		CreateTime: v.CreateTime,
	}
}

func (convert) MeetingVideoRecord(v mysqlmodelrtc.MeetingVideoRecord) mongomodelrtc.MeetingVideoRecord {
	return mongomodelrtc.MeetingVideoRecord{
		RoomID:     v.RoomID,
		FileURL:    v.FileURL,
		CreateTime: v.CreateTime,
	}
}
