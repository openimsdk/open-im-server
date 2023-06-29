package convert

import (
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	pbGroup "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func Db2PbGroupInfo(m *relation.GroupModel, ownerUserID string, memberCount uint32) *sdkws.GroupInfo {
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

func Pb2DbGroupRequest(req *pbGroup.GroupApplicationResponseReq, handleUserID string) *relation.GroupRequestModel {
	return &relation.GroupRequestModel{
		UserID:       req.FromUserID,
		GroupID:      req.GroupID,
		HandleResult: req.HandleResult,
		HandledMsg:   req.HandledMsg,
		HandleUserID: handleUserID,
		HandledTime:  time.Now(),
	}
}

func Db2PbCMSGroup(m *relation.GroupModel, ownerUserID string, ownerUserName string, memberCount uint32) *pbGroup.CMSGroup {
	return &pbGroup.CMSGroup{
		GroupInfo:          Db2PbGroupInfo(m, ownerUserID, memberCount),
		GroupOwnerUserID:   ownerUserID,
		GroupOwnerUserName: ownerUserName,
	}
}

func Db2PbGroupMember(m *relation.GroupMemberModel) *sdkws.GroupMemberFullInfo {
	return &sdkws.GroupMemberFullInfo{
		GroupID:   m.GroupID,
		UserID:    m.UserID,
		RoleLevel: m.RoleLevel,
		JoinTime:  m.JoinTime.UnixMilli(),
		Nickname:  m.Nickname,
		FaceURL:   m.FaceURL,
		//AppMangerLevel: m.AppMangerLevel,
		JoinSource:     m.JoinSource,
		OperatorUserID: m.OperatorUserID,
		Ex:             m.Ex,
		MuteEndTime:    m.MuteEndTime.UnixMilli(),
		InviterUserID:  m.InviterUserID,
	}
}

func Db2PbGroupRequest(m *relation.GroupRequestModel, user *sdkws.PublicUserInfo, group *sdkws.GroupInfo) *sdkws.GroupRequest {
	return &sdkws.GroupRequest{
		UserInfo:      user,
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

func Db2PbGroupAbstractInfo(groupID string, groupMemberNumber uint32, groupMemberListHash uint64) *pbGroup.GroupAbstractInfo {
	return &pbGroup.GroupAbstractInfo{
		GroupID:             groupID,
		GroupMemberNumber:   groupMemberNumber,
		GroupMemberListHash: groupMemberListHash,
	}
}

func Pb2DBGroupInfo(m *sdkws.GroupInfo) *relation.GroupModel {
	return &relation.GroupModel{
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

func Pb2DbGroupMember(m *sdkws.UserInfo) *relation.GroupMemberModel {
	return &relation.GroupMemberModel{
		UserID:   m.UserID,
		Nickname: m.Nickname,
		FaceURL:  m.FaceURL,
		Ex:       m.Ex,
	}
}
