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

func UpdateGroupInfoExMap(ctx context.Context, group *pbgroup.SetGroupInfoExReq) (map[string]any, error) {
	m := make(map[string]any)

	if group.GroupName != nil {
		if group.GroupName.Value != "" {
			m["group_name"] = group.GroupName.Value
		} else {
			return nil, errs.ErrArgs.WrapMsg("group name is empty")
		}
	}
	if group.Notification != nil {
		m["notification"] = group.Notification.Value
		m["notification_update_time"] = time.Now()
		m["notification_user_id"] = mcontext.GetOpUserID(ctx)
	}
	if group.Introduction != nil {
		m["introduction"] = group.Introduction.Value
	}
	if group.FaceURL != nil {
		m["face_url"] = group.FaceURL.Value
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

	return m, nil
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
