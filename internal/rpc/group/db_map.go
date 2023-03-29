package group

import (
	"time"

	pbGroup "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func UpdateGroupInfoMap(group *sdkws.GroupInfoForSet) map[string]any {
	m := make(map[string]any)
	if group.GroupName != "" {
		m["name"] = group.GroupName
	}
	if group.Notification != "" {
		m["Notification"] = group.Notification
	}
	if group.Introduction != "" {
		m["introduction"] = group.Introduction
	}
	if group.FaceURL != "" {
		m["face_url"] = group.FaceURL
	}
	if group.NeedVerification != nil {
		m["need_verification"] = group.NeedVerification.Value
	}
	if group.LookMemberInfo != nil {
		m["look_member_info"] = group.LookMemberInfo.Value
	}
	if group.ApplyMemberFriend != nil {
		m["apply_member_friend"] = group.ApplyMemberFriend.Value
	}
	return m
}

func UpdateGroupStatusMap(status int) map[string]any {
	return map[string]any{
		"status": status,
	}
}

func UpdateGroupMemberMutedTimeMap(t time.Time) map[string]any {
	return map[string]any{
		"mute_end_time": t,
	}
}

func UpdateGroupMemberMap(req *pbGroup.SetGroupMemberInfo) map[string]any {
	m := make(map[string]any)
	if req.Nickname != nil {
		m["nickname"] = req.Nickname.Value
	}
	if req.FaceURL != nil {
		m["user_group_face_url"] = req.FaceURL.Value
	}
	if req.RoleLevel != nil {
		m["role_level"] = req.RoleLevel.Value
	}
	if req.Ex != nil {
		m["ex"] = req.Ex.Value
	}
	return m
}
