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
	"time"

	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/protocol/sdkws"
)

func UsersDB2Pb(users []*relationTb.UserModel) (result []*sdkws.UserInfo) {
	for _, user := range users {
		var userPb sdkws.UserInfo
		userPb.UserID = user.UserID
		userPb.Nickname = user.Nickname
		userPb.FaceURL = user.FaceURL
		userPb.Ex = user.Ex
		userPb.CreateTime = user.CreateTime.UnixMilli()
		userPb.AppMangerLevel = user.AppMangerLevel
		userPb.GlobalRecvMsgOpt = user.GlobalRecvMsgOpt
		result = append(result, &userPb)
	}
	return result
}

func UserPb2DB(user *sdkws.UserInfo) *relationTb.UserModel {
	var userDB relationTb.UserModel
	userDB.UserID = user.UserID
	userDB.Nickname = user.Nickname
	userDB.FaceURL = user.FaceURL
	userDB.Ex = user.Ex
	userDB.CreateTime = time.UnixMilli(user.CreateTime)
	userDB.AppMangerLevel = user.AppMangerLevel
	userDB.GlobalRecvMsgOpt = user.GlobalRecvMsgOpt
	return &userDB
}
