package pkg

import (
	"time"

	mongoModel "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	mysqlModel "github.com/openimsdk/open-im-server/v3/tools/data-conversion/openim/mysql/v3"
	mongoModelRtc "github.com/openimsdk/open-im-server/v3/tools/up35/pkg/internal/rtc/mongo/table"
	mysqlModelRtc "github.com/openimsdk/open-im-server/v3/tools/up35/pkg/internal/rtc/mysql"
)

type convert struct{}

func (convert) User(v mysqlModel.UserModel) mongoModel.UserModel {
	return mongoModel.UserModel{
		UserID:           v.UserID,
		Nickname:         v.Nickname,
		FaceURL:          v.FaceURL,
		Ex:               v.Ex,
		AppMangerLevel:   v.AppMangerLevel,
		GlobalRecvMsgOpt: v.GlobalRecvMsgOpt,
		CreateTime:       v.CreateTime,
	}
}

func (convert) Friend(v mysqlModel.FriendModel) mongoModel.FriendModel {
	return mongoModel.FriendModel{
		OwnerUserID:    v.OwnerUserID,
		FriendUserID:   v.FriendUserID,
		Remark:         v.Remark,
		CreateTime:     v.CreateTime,
		AddSource:      v.AddSource,
		OperatorUserID: v.OperatorUserID,
		Ex:             v.Ex,
	}
}

func (convert) FriendRequest(v mysqlModel.FriendRequestModel) mongoModel.FriendRequestModel {
	return mongoModel.FriendRequestModel{
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

func (convert) Black(v mysqlModel.BlackModel) mongoModel.BlackModel {
	return mongoModel.BlackModel{
		OwnerUserID:    v.OwnerUserID,
		BlockUserID:    v.BlockUserID,
		CreateTime:     v.CreateTime,
		AddSource:      v.AddSource,
		OperatorUserID: v.OperatorUserID,
		Ex:             v.Ex,
	}
}

func (convert) Group(v mysqlModel.GroupModel) mongoModel.GroupModel {
	return mongoModel.GroupModel{
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

func (convert) GroupMember(v mysqlModel.GroupMemberModel) mongoModel.GroupMemberModel {
	return mongoModel.GroupMemberModel{
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

func (convert) GroupRequest(v mysqlModel.GroupRequestModel) mongoModel.GroupRequestModel {
	return mongoModel.GroupRequestModel{
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

func (convert) Conversation(v mysqlModel.ConversationModel) mongoModel.ConversationModel {
	return mongoModel.ConversationModel{
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

func (convert) Object(engine string) func(v mysqlModel.ObjectModel) mongoModel.ObjectModel {
	return func(v mysqlModel.ObjectModel) mongoModel.ObjectModel {
		return mongoModel.ObjectModel{
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

func (convert) Log(v mysqlModel.Log) mongoModel.LogModel {
	return mongoModel.LogModel{
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

func (convert) SignalModel(v mysqlModelRtc.SignalModel) mongoModelRtc.SignalModel {
	return mongoModelRtc.SignalModel{
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

func (convert) SignalInvitationModel(v mysqlModelRtc.SignalInvitationModel) mongoModelRtc.SignalInvitationModel {
	return mongoModelRtc.SignalInvitationModel{
		SID:          v.SID,
		UserID:       v.UserID,
		Status:       v.Status,
		InitiateTime: v.InitiateTime,
		HandleTime:   v.HandleTime,
	}
}

func (convert) Meeting(v mysqlModelRtc.MeetingInfo) mongoModelRtc.MeetingInfo {
	return mongoModelRtc.MeetingInfo{
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

func (convert) MeetingInvitationInfo(v mysqlModelRtc.MeetingInvitationInfo) mongoModelRtc.MeetingInvitationInfo {
	return mongoModelRtc.MeetingInvitationInfo{
		RoomID:     v.RoomID,
		UserID:     v.UserID,
		CreateTime: v.CreateTime,
	}
}

func (convert) MeetingVideoRecord(v mysqlModelRtc.MeetingVideoRecord) mongoModelRtc.MeetingVideoRecord {
	return mongoModelRtc.MeetingVideoRecord{
		RoomID:     v.RoomID,
		FileURL:    v.FileURL,
		CreateTime: v.CreateTime,
	}
}
