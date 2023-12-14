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

	"github.com/OpenIMSDK/protocol/sdkws"

	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

func UsersDB2Pb(users []*relationtb.UserModel) []*sdkws.UserInfo {
	result := make([]*sdkws.UserInfo, 0, len(users))
	for _, user := range users {
		userPb := &sdkws.UserInfo{
			UserID:           user.UserID,
			Nickname:         user.Nickname,
			FaceURL:          user.FaceURL,
			Ex:               user.Ex,
			CreateTime:       user.CreateTime.UnixMilli(),
			AppMangerLevel:   user.AppMangerLevel,
			GlobalRecvMsgOpt: user.GlobalRecvMsgOpt,
		}
		result = append(result, userPb)
	}
	return result
}

func UserPb2DB(user *sdkws.UserInfo) *relationtb.UserModel {
	return &relationtb.UserModel{
		UserID:           user.UserID,
		Nickname:         user.Nickname,
		FaceURL:          user.FaceURL,
		Ex:               user.Ex,
		CreateTime:       time.UnixMilli(user.CreateTime),
		AppMangerLevel:   user.AppMangerLevel,
		GlobalRecvMsgOpt: user.GlobalRecvMsgOpt,
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
		"app_manager_level":   user.AppMangerLevel,
		"global_recv_msg_opt": user.GlobalRecvMsgOpt,
	}
	for key, value := range fields {
		if v, ok := value.(string); ok && v != "" {
			val[key] = v
		} else if v, ok := value.(int32); ok && v != 0 {
			val[key] = v
		}
	}
	return val
}
