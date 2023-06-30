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

package mysql

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"strconv"
	"time"
)

func ImportUserToSuperGroup() {
	for i := 18000000700; i <= 18000000800; i++ {
		user := db.User{
			UserID:           strconv.Itoa(i),
			Nickname:         strconv.Itoa(i),
			FaceURL:          "",
			Gender:           0,
			PhoneNumber:      strconv.Itoa(i),
			Birth:            time.Time{},
			Email:            "",
			Ex:               "",
			CreateTime:       time.Time{},
			AppMangerLevel:   0,
			GlobalRecvMsgOpt: 0,
		}
		err := im_mysql_model.UserRegister(user)
		if err != nil {
			log.NewError("", err.Error(), user)
			continue
		}

		groupMember := db.GroupMember{
			GroupID:        "3907826375",
			UserID:         strconv.Itoa(i),
			Nickname:       strconv.Itoa(i),
			FaceURL:        "",
			RoleLevel:      0,
			JoinTime:       time.Time{},
			JoinSource:     0,
			InviterUserID:  "openIMAdmin",
			OperatorUserID: "openIMAdmin",
			MuteEndTime:    time.Time{},
			Ex:             "",
		}

		err = im_mysql_model.InsertIntoGroupMember(groupMember)
		if err != nil {
			log.NewError("", err.Error(), user)
			continue
		}

		log.NewInfo("success", i)

	}

}
