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

package convert

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"time"

	pbgroup "github.com/openimsdk/protocol/group"
	sdkws "github.com/openimsdk/protocol/sdkws"
)

func Db2PbGroupInfo(m *model.Group, ownerUserID string, memberCount uint32) *sdkws.GroupInfo {
	return &sdkws.GroupInfo{
		GroupID:                m.GroupID,
		GroupName:              m.GroupName,
		Notification:           m.Notification,
		Introduction:           m.Introduction,
		FaceURL:                m.FaceURL,
		OwnerUserID:            ownerUserID,
		CreateTime:             m.CreateTime.UnixMilli(),
		MemberCount:            memberCount,
		Ex:                     m.Ex,
		Status:                 m.Status,
		CreatorUserID:          m.CreatorUserID,
		GroupType:              m.GroupType,
		NeedVerification:       m.NeedVerification,
		LookMemberInfo:         m.LookMemberInfo,
		ApplyMemberFriend:      m.ApplyMemberFriend,
		NotificationUpdateTime: m.NotificationUpdateTime.UnixMilli(),
		NotificationUserID:     m.NotificationUserID,
	}
}

func Pb2DbGroupRequest(req *pbgroup.GroupApplicationResponseReq, handleUserID string) *model.GroupRequest {
	return &model.GroupRequest{
		UserID:       req.FromUserID,
		GroupID:      req.GroupID,
		HandleResult: req.HandleResult,
		HandledMsg:   req.HandledMsg,
		HandleUserID: handleUserID,
		HandledTime:  time.Now(),
	}
}

func Db2PbCMSGroup(m *model.Group, ownerUserID string, ownerUserName string, memberCount uint32) *pbgroup.CMSGroup {
	return &pbgroup.CMSGroup{
		GroupInfo:          Db2PbGroupInfo(m, ownerUserID, memberCount),
		GroupOwnerUserID:   ownerUserID,
		GroupOwnerUserName: ownerUserName,
	}
}

func Db2PbGroupMember(m *model.GroupMember) *sdkws.GroupMemberFullInfo {
	return &sdkws.GroupMemberFullInfo{
		GroupID:   m.GroupID,
		UserID:    m.UserID,
		RoleLevel: m.RoleLevel,
		JoinTime:  m.JoinTime.UnixMilli(),
		Nickname:  m.Nickname,
		FaceURL:   m.FaceURL,
		// AppMangerLevel: m.AppMangerLevel,
		JoinSource:     m.JoinSource,
		OperatorUserID: m.OperatorUserID,
		Ex:             m.Ex,
		MuteEndTime:    m.MuteEndTime.UnixMilli(),
		InviterUserID:  m.InviterUserID,
	}
}

func Db2PbGroupRequest(m *model.GroupRequest, user *sdkws.UserInfo, group *sdkws.GroupInfo) *sdkws.GroupRequest {
	var pu *sdkws.PublicUserInfo
	if user != nil {
		pu = &sdkws.PublicUserInfo{
			UserID:   user.UserID,
			Nickname: user.Nickname,
			FaceURL:  user.FaceURL,
			Ex:       user.Ex,
		}
	}
	return &sdkws.GroupRequest{
		UserInfo:      pu,
		GroupInfo:     group,
		HandleResult:  m.HandleResult,
		ReqMsg:        m.ReqMsg,
		HandleMsg:     m.HandledMsg,
		ReqTime:       m.ReqTime.UnixMilli(),
		HandleUserID:  m.HandleUserID,
		HandleTime:    m.HandledTime.UnixMilli(),
		Ex:            m.Ex,
		JoinSource:    m.JoinSource,
		InviterUserID: m.InviterUserID,
	}
}

func Db2PbGroupAbstractInfo(
	groupID string,
	groupMemberNumber uint32,
	groupMemberListHash uint64,
) *pbgroup.GroupAbstractInfo {
	return &pbgroup.GroupAbstractInfo{
		GroupID:             groupID,
		GroupMemberNumber:   groupMemberNumber,
		GroupMemberListHash: groupMemberListHash,
	}
}

func Pb2DBGroupInfo(m *sdkws.GroupInfo) *model.Group {
	return &model.Group{
		GroupID:                m.GroupID,
		GroupName:              m.GroupName,
		Notification:           m.Notification,
		Introduction:           m.Introduction,
		FaceURL:                m.FaceURL,
		CreateTime:             time.Now(),
		Ex:                     m.Ex,
		Status:                 m.Status,
		CreatorUserID:          m.CreatorUserID,
		GroupType:              m.GroupType,
		NeedVerification:       m.NeedVerification,
		LookMemberInfo:         m.LookMemberInfo,
		ApplyMemberFriend:      m.ApplyMemberFriend,
		NotificationUpdateTime: time.UnixMilli(m.NotificationUpdateTime),
		NotificationUserID:     m.NotificationUserID,
	}
}

// func Pb2DbGroupMember(m *sdkws.UserInfo) *relation.GroupMember {
//	return &relation.GroupMember{
//		UserID:   m.UserID,
//		Nickname: m.Nickname,
//		FaceURL:  m.FaceURL,
//		Ex:       m.Ex,
//	}
//}
