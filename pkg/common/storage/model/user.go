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

package model

import (
	"time"
)

// PhoneVisibility 手机号可见性枚举。
// 0=所有人可见, 1=仅好友可见, 2=完全隐藏
const (
	PhoneVisibilityPublic  int32 = 0
	PhoneVisibilityFriends int32 = 1
	PhoneVisibilityHidden  int32 = 2
)

// CallAcceptSetting 音视频通话接受权限枚举。
// 0=所有人可发起, 1=仅好友可发起, 2=不接受任何通话
const (
	CallAcceptSettingPublic  int32 = 0
	CallAcceptSettingFriends int32 = 1
	CallAcceptSettingNobody  int32 = 2
)

// MsgReceiveSetting 会话消息接收权限枚举。
// 0=所有人可发送, 1=仅好友可发送, 2=所有人不可发送
const (
	MsgReceiveSettingPublic  int32 = 0
	MsgReceiveSettingFriends int32 = 1
	MsgReceiveSettingNobody  int32 = 2
)

type User struct {
	UserID            string    `bson:"user_id"`
	Nickname          string    `bson:"nickname"`
	FaceURL           string    `bson:"face_url"`
	Ex                string    `bson:"ex"`
	AppMangerLevel    int32     `bson:"app_manger_level"`
	GlobalRecvMsgOpt  int32     `bson:"global_recv_msg_opt"`
	CreateTime        time.Time `bson:"create_time"`
	FirstName         string    `bson:"first_name"`
	LastName          string    `bson:"last_name"`
	FullName          string    `bson:"full_name"`
	Phone             string    `bson:"phone"`
	AreaCode          string    `bson:"area_code"`
	PhoneVisibility   int32     `bson:"phone_visibility"`
	CallAcceptSetting int32     `bson:"call_accept_setting"`
	MsgReceiveSetting int32     `bson:"msg_receive_setting"`
	// CallRingtoneURL 用户自定义来电铃声 URL；对方来电时播放此铃声
	CallRingtoneURL string `bson:"call_ringtone_url"`
}

func (u *User) GetNickname() string {
	return u.Nickname
}

func (u *User) GetFaceURL() string {
	return u.FaceURL
}

func (u *User) GetUserID() string {
	return u.UserID
}

func (u *User) GetEx() string {
	return u.Ex
}
