package group

import (
	"Open_IM/pkg/common/db/table/relation"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"time"
)

func DbToPbGroupInfo(m *relation.GroupModel, ownerUserID string, memberCount uint32) *open_im_sdk.GroupInfo {
	return &open_im_sdk.GroupInfo{
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

func PbToDbGroupRequest(req *pbGroup.GroupApplicationResponseReq, handleUserID string) *relation.GroupRequestModel {
	return &relation.GroupRequestModel{
		UserID:       req.FromUserID,
		GroupID:      req.GroupID,
		HandleResult: req.HandleResult,
		HandledMsg:   req.HandledMsg,
		HandleUserID: handleUserID,
		HandledTime:  time.Now(),
	}
}

func DbToPbCMSGroup(m *relation.GroupModel, ownerUserID string, ownerUserName string, memberCount uint32) *pbGroup.CMSGroup {
	return &pbGroup.CMSGroup{
		GroupInfo:          DbToPbGroupInfo(m, ownerUserID, memberCount),
		GroupOwnerUserID:   ownerUserID,
		GroupOwnerUserName: ownerUserName,
	}
}

func DbToPbGroupMembersCMSResp(m *relation.GroupMemberModel) *open_im_sdk.GroupMemberFullInfo {
	return &open_im_sdk.GroupMemberFullInfo{
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

func DbToPbGroupRequest(m *relation.GroupRequestModel, user *open_im_sdk.PublicUserInfo, group *open_im_sdk.GroupInfo) *open_im_sdk.GroupRequest {
	return &open_im_sdk.GroupRequest{
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

func DbToPbGroupAbstractInfo(groupID string, groupMemberNumber int32, groupMemberListHash uint64) *pbGroup.GroupAbstractInfo {
	return &pbGroup.GroupAbstractInfo{
		GroupID:             groupID,
		GroupMemberNumber:   groupMemberNumber,
		GroupMemberListHash: groupMemberListHash,
	}
}
