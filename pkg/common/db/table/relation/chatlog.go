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

package relation

import (
	"time"

	pbmsg "github.com/OpenIMSDK/protocol/msg"
)

const (
	ChatLogModelTableName = "chat_logs"
)

type ChatLogModel struct {
	ServerMsgID      string    `gorm:"column:server_msg_id;primary_key;type:char(64)"                                                                                                json:"serverMsgID"`
	ClientMsgID      string    `gorm:"column:client_msg_id;type:char(64)"                                                                                                            json:"clientMsgID"`
	SendID           string    `gorm:"column:send_id;type:char(64);index:send_id,priority:2"                                                                                         json:"sendID"`
	RecvID           string    `gorm:"column:recv_id;type:char(64);index:recv_id,priority:2"                                                                                         json:"recvID"`
	SenderPlatformID int32     `gorm:"column:sender_platform_id"                                                                                                                     json:"senderPlatformID"`
	SenderNickname   string    `gorm:"column:sender_nick_name;type:varchar(255)"                                                                                                     json:"senderNickname"`
	SenderFaceURL    string    `gorm:"column:sender_face_url;type:varchar(255);"                                                                                                     json:"senderFaceURL"`
	SessionType      int32     `gorm:"column:session_type;index:session_type,priority:2;index:session_type_alone"                                                                    json:"sessionType"`
	MsgFrom          int32     `gorm:"column:msg_from"                                                                                                                               json:"msgFrom"`
	ContentType      int32     `gorm:"column:content_type;index:content_type,priority:2;index:content_type_alone"                                                                    json:"contentType"`
	Content          string    `gorm:"column:content;type:varchar(3000)"                                                                                                             json:"content"`
	Status           int32     `gorm:"column:status"                                                                                                                                 json:"status"`
	SendTime         time.Time `gorm:"column:send_time;index:sendTime;index:content_type,priority:1;index:session_type,priority:1;index:recv_id,priority:1;index:send_id,priority:1" json:"sendTime"`
	CreateTime       time.Time `gorm:"column:create_time"                                                                                                                            json:"createTime"`
	Ex               string    `gorm:"column:ex;type:varchar(1024)"                                                                                                                  json:"ex"`
}

func (ChatLogModel) TableName() string {
	return ChatLogModelTableName
}

type ChatLogModelInterface interface {
	Create(msg *pbmsg.MsgDataToMQ) error
}
