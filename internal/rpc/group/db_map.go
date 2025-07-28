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
	"strings"
	"time"

	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
)

func UpdateGroupInfoMap(ctx context.Context, group *sdkws.GroupInfoForSet) map[string]any {
	m := make(map[string]any)
	if group.GroupName != "" {
		m["group_name"] = group.GroupName
	}
	if group.Notification != "" {
		m["notification"] = group.Notification
		m["notification_update_time"] = time.Now()
		m["notification_user_id"] = mcontext.GetOpUserID(ctx)
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
	if group.Ex != nil {
		m["ex"] = group.Ex.Value
	}
	return m
}

func UpdateGroupInfoExMap(ctx context.Context, group *pbgroup.SetGroupInfoExReq) (m map[string]any, normalFlag, groupNameFlag, notificationFlag bool, err error) {
	m = make(map[string]any)

	if group.GroupName != nil {
		if strings.TrimSpace(group.GroupName.Value) != "" {
			m["group_name"] = group.GroupName.Value
			groupNameFlag = true
		} else {
			return nil, normalFlag, notificationFlag, groupNameFlag, errs.ErrArgs.WrapMsg("group name is empty")
		}
	}

	if group.Notification != nil {
		notificationFlag = true
		group.Notification.Value = strings.TrimSpace(group.Notification.Value) // if Notification only contains spaces, set it to empty string

		m["notification"] = group.Notification.Value
		m["notification_user_id"] = mcontext.GetOpUserID(ctx)
		m["notification_update_time"] = time.Now()
	}
	if group.Introduction != nil {
		m["introduction"] = group.Introduction.Value
		normalFlag = true
	}
	if group.FaceURL != nil {
		m["face_url"] = group.FaceURL.Value
		normalFlag = true
	}
	if group.NeedVerification != nil {
		m["need_verification"] = group.NeedVerification.Value
		normalFlag = true
	}
	if group.LookMemberInfo != nil {
		m["look_member_info"] = group.LookMemberInfo.Value
		normalFlag = true
	}
	if group.ApplyMemberFriend != nil {
		m["apply_member_friend"] = group.ApplyMemberFriend.Value
		normalFlag = true
	}
	if group.Ex != nil {
		m["ex"] = group.Ex.Value
		normalFlag = true
	}

	return m, normalFlag, groupNameFlag, notificationFlag, nil
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

func UpdateGroupMemberMap(req *pbgroup.SetGroupMemberInfo) map[string]any {
	m := make(map[string]any)
	if req.Nickname != nil {
		m["nickname"] = req.Nickname.Value
	}
	if req.FaceURL != nil {
		m["face_url"] = req.FaceURL.Value
	}
	if req.RoleLevel != nil {
		m["role_level"] = req.RoleLevel.Value
	}
	if req.Ex != nil {
		m["ex"] = req.Ex.Value
	}
	return m
}
