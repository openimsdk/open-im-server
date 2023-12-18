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

package conversion

import (
	v2 "github.com/openimsdk/open-im-server/v3/tools/data-conversion/chat/v2"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/chat/v3/admin"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/chat/v3/chat"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/utils"
)

// ########## chat ##########

func Account(v v2.Account) (chat.Account, bool) {
	utils.InitTime(&v.CreateTime, &v.ChangeTime)
	return chat.Account{
		UserID:         v.UserID,
		Password:       v.Password,
		CreateTime:     v.CreateTime,
		ChangeTime:     v.ChangeTime,
		OperatorUserID: v.OperatorUserID,
	}, true
}

func Attribute(v v2.Attribute) (chat.Attribute, bool) {
	utils.InitTime(&v.CreateTime, &v.ChangeTime, &v.BirthTime)
	return chat.Attribute{
		UserID:           v.UserID,
		Account:          v.Account,
		PhoneNumber:      v.PhoneNumber,
		AreaCode:         v.AreaCode,
		Email:            v.Email,
		Nickname:         v.Nickname,
		FaceURL:          v.FaceURL,
		Gender:           v.Gender,
		CreateTime:       v.CreateTime,
		ChangeTime:       v.ChangeTime,
		BirthTime:        v.BirthTime,
		Level:            v.Level,
		AllowVibration:   v.AllowVibration,
		AllowBeep:        v.AllowBeep,
		AllowAddFriend:   v.AllowAddFriend,
		GlobalRecvMsgOpt: 0,
	}, true
}

func Register(v v2.Register) (chat.Register, bool) {
	utils.InitTime(&v.CreateTime)
	return chat.Register{
		UserID:      v.UserID,
		DeviceID:    v.DeviceID,
		IP:          v.IP,
		Platform:    v.Platform,
		AccountType: v.AccountType,
		Mode:        v.Mode,
		CreateTime:  v.CreateTime,
	}, true
}

func UserLoginRecord(v v2.UserLoginRecord) (chat.UserLoginRecord, bool) {
	utils.InitTime(&v.LoginTime)
	return chat.UserLoginRecord{
		UserID:    v.UserID,
		LoginTime: v.LoginTime,
		IP:        v.IP,
		DeviceID:  v.DeviceID,
		Platform:  v.Platform,
	}, true
}

// ########## admin ##########

func Admin(v v2.Admin) (admin.Admin, bool) {
	utils.InitTime(&v.CreateTime)
	return admin.Admin{
		Account:    v.Account,
		Password:   v.Password,
		FaceURL:    v.FaceURL,
		Nickname:   v.Nickname,
		UserID:     v.UserID,
		Level:      v.Level,
		CreateTime: v.CreateTime,
	}, true
}

func Applet(v v2.Applet) (admin.Applet, bool) {
	utils.InitTime(&v.CreateTime)
	return admin.Applet{
		ID:         v.ID,
		Name:       v.Name,
		AppID:      v.AppID,
		Icon:       v.Icon,
		URL:        v.URL,
		MD5:        v.MD5,
		Size:       v.Size,
		Version:    v.Version,
		Priority:   v.Priority,
		Status:     v.Status,
		CreateTime: v.CreateTime,
	}, true
}

func ForbiddenAccount(v v2.ForbiddenAccount) (admin.ForbiddenAccount, bool) {
	utils.InitTime(&v.CreateTime)
	return admin.ForbiddenAccount{
		UserID:         v.UserID,
		Reason:         v.Reason,
		OperatorUserID: v.OperatorUserID,
		CreateTime:     v.CreateTime,
	}, true
}

func InvitationRegister(v v2.InvitationRegister) (admin.InvitationRegister, bool) {
	utils.InitTime(&v.CreateTime)
	return admin.InvitationRegister{
		InvitationCode: v.InvitationCode,
		UsedByUserID:   v.UsedByUserID,
		CreateTime:     v.CreateTime,
	}, true
}

func IPForbidden(v v2.IPForbidden) (admin.IPForbidden, bool) {
	utils.InitTime(&v.CreateTime)
	return admin.IPForbidden{
		IP:            v.IP,
		LimitRegister: v.LimitRegister > 0,
		LimitLogin:    v.LimitLogin > 0,
		CreateTime:    v.CreateTime,
	}, true
}

func LimitUserLoginIP(v v2.LimitUserLoginIP) (admin.LimitUserLoginIP, bool) {
	utils.InitTime(&v.CreateTime)
	return admin.LimitUserLoginIP{
		UserID:     v.UserID,
		IP:         v.IP,
		CreateTime: v.CreateTime,
	}, true
}

func RegisterAddFriend(v v2.RegisterAddFriend) (admin.RegisterAddFriend, bool) {
	utils.InitTime(&v.CreateTime)
	return admin.RegisterAddFriend{
		UserID:     v.UserID,
		CreateTime: v.CreateTime,
	}, true
}

func RegisterAddGroup(v v2.RegisterAddGroup) (admin.RegisterAddGroup, bool) {
	utils.InitTime(&v.CreateTime)
	return admin.RegisterAddGroup{
		GroupID:    v.GroupID,
		CreateTime: v.CreateTime,
	}, true
}
