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
	"strings"
	"time"

	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/protocol/sdkws"
)

func BuildFullName(firstName, lastName string) string {
	if firstName == "" {
		return lastName
	}
	if lastName == "" {
		return firstName
	}
	return strings.TrimSpace(firstName + " " + lastName)
}

// MemberDisplayNickname 非好友场景下的群成员展示名：优先 firstName+lastName，否则 nickname。
func MemberDisplayNickname(u *sdkws.UserInfo) string {
	if u == nil {
		return ""
	}
	if name := BuildFullName(u.FirstName, u.LastName); name != "" {
		return name
	}
	return u.Nickname
}

func UserDB2Pb(user *relationtb.User) *sdkws.UserInfo {
	return &sdkws.UserInfo{
		UserID:            user.UserID,
		Nickname:          user.Nickname,
		FaceURL:           user.FaceURL,
		Ex:                user.Ex,
		CreateTime:        user.CreateTime.UnixMilli(),
		AppMangerLevel:    user.AppMangerLevel,
		GlobalRecvMsgOpt:  user.GlobalRecvMsgOpt,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Phone:             user.Phone,
		AreaCode:          user.AreaCode,
		PhoneVisibility:   user.PhoneVisibility,
		CallAcceptSetting:  user.CallAcceptSetting,
		MsgReceiveSetting:  user.MsgReceiveSetting,
		GroupInviteSetting: user.GroupInviteSetting,
		CallRingtoneURL:    user.CallRingtoneURL,
		MsgBurnDuration:    user.MsgBurnDuration,
	}
}

func UsersDB2Pb(users []*relationtb.User) []*sdkws.UserInfo {
	return datautil.Slice(users, UserDB2Pb)
}

func UserPb2DB(user *sdkws.UserInfo) *relationtb.User {
	fullName := BuildFullName(user.FirstName, user.LastName)
	return &relationtb.User{
		UserID:          user.UserID,
		Nickname:        user.Nickname,
		FaceURL:         user.FaceURL,
		Ex:              user.Ex,
		CreateTime:      time.UnixMilli(user.CreateTime),
		AppMangerLevel:  user.AppMangerLevel,
		GlobalRecvMsgOpt: user.GlobalRecvMsgOpt,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		FullName:        fullName,
		AreaCode:        user.AreaCode,
		CallRingtoneURL: user.CallRingtoneURL,
	}
}

func UserPb2DBMap(user *sdkws.UserInfo) map[string]any {
	if user == nil {
		return nil
	}
	val := make(map[string]any)
	fields := map[string]any{
		"nickname":            user.Nickname,
		"face_url":            user.FaceURL,
		"ex":                  user.Ex,
		"first_name":          user.FirstName,
		"last_name":           user.LastName,
		"area_code":           user.AreaCode,
		"app_manager_level":   user.AppMangerLevel,
		"global_recv_msg_opt": user.GlobalRecvMsgOpt,
		"call_ringtone_url":   user.CallRingtoneURL,
		"msg_burn_duration":   user.MsgBurnDuration,
	}
	for key, value := range fields {
		if v, ok := value.(string); ok && v != "" {
			val[key] = v
		} else if v, ok := value.(int32); ok && v != 0 {
			val[key] = v
		}
	}
	if user.FirstName != "" || user.LastName != "" {
		fullName := BuildFullName(user.FirstName, user.LastName)
		val["full_name"] = fullName
	}
	return val
}
func UserPb2DBMapEx(user *sdkws.UserInfoWithEx) map[string]any {
	if user == nil {
		return nil
	}
	val := make(map[string]any)

	// Map fields from UserInfoWithEx to val
	if user.Nickname != nil {
		val["nickname"] = user.Nickname.Value
	}
	if user.FaceURL != nil {
		val["face_url"] = user.FaceURL.Value
	}
	if user.Ex != nil {
		val["ex"] = user.Ex.Value
	}
	if user.FirstName != nil {
		val["first_name"] = user.FirstName.Value
	}
	if user.LastName != nil {
		val["last_name"] = user.LastName.Value
	}
	if user.FirstName != nil || user.LastName != nil {
		firstName := ""
		lastName := ""
		if user.FirstName != nil {
			firstName = user.FirstName.Value
		}
		if user.LastName != nil {
			lastName = user.LastName.Value
		}
		val["full_name"] = BuildFullName(firstName, lastName)
	}
	if user.GlobalRecvMsgOpt != nil {
		val["global_recv_msg_opt"] = user.GlobalRecvMsgOpt.Value
	}
	if user.Phone != nil {
		val["phone"] = user.Phone.Value
	}
	if user.AreaCode != nil {
		val["area_code"] = user.AreaCode.Value
	}
	if user.PhoneVisibility != nil {
		val["phone_visibility"] = user.PhoneVisibility.Value
	}
	if user.CallAcceptSetting != nil {
		val["call_accept_setting"] = user.CallAcceptSetting.Value
	}
	if user.MsgReceiveSetting != nil {
		val["msg_receive_setting"] = user.MsgReceiveSetting.Value
	}
	if user.GroupInviteSetting != nil {
		val["group_invite_setting"] = user.GroupInviteSetting.Value
	}
	if user.CallRingtoneURL != nil {
		val["call_ringtone_url"] = user.CallRingtoneURL.Value
	}
	if user.MsgBurnDuration != nil {
		val["msg_burn_duration"] = user.MsgBurnDuration.Value
	}
	return val
}
