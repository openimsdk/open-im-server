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

package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"fmt"
)

func GetChatLog(chatLog *db.ChatLog, pageNumber, showNumber int32, contentTypeList []int32) (int64, []db.ChatLog, error) {
	mdb := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs")
	if chatLog.SendTime.Unix() > 0 {
		mdb = mdb.Where("send_time > ? and send_time < ?", chatLog.SendTime, chatLog.SendTime.AddDate(0, 0, 1))
	}
	if chatLog.Content != "" {
		mdb = mdb.Where(" content like ? ", fmt.Sprintf("%%%s%%", chatLog.Content))
	}
	if chatLog.SessionType == 1 {
		mdb = mdb.Where("session_type = ?", chatLog.SessionType)
	} else if chatLog.SessionType == 2 {
		mdb = mdb.Where("session_type in (?)", []int{constant.GroupChatType, constant.SuperGroupChatType})
	}
	if chatLog.ContentType != 0 {
		mdb = mdb.Where("content_type = ?", chatLog.ContentType)
	}
	if chatLog.SendID != "" {
		mdb = mdb.Where("send_id = ?", chatLog.SendID)
	}
	if chatLog.RecvID != "" {
		mdb = mdb.Where("recv_id = ?", chatLog.RecvID)
	}
	if len(contentTypeList) > 0 {
		mdb = mdb.Where("content_type in (?)", contentTypeList)
	}
	var count int64
	if err := mdb.Count(&count).Error; err != nil {
		return 0, nil, err
	}
	var chatLogs []db.ChatLog
	mdb = mdb.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1)))
	if err := mdb.Find(&chatLogs).Error; err != nil {
		return 0, nil, err
	}
	return count, chatLogs, nil
}
