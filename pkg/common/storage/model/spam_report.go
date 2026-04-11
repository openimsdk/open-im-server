// Copyright © 2024 OpenIM. All rights reserved.
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

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SpamReport status constants.
const (
	SpamReportStatusPending  int32 = 0 // 待处理
	SpamReportStatusHandled  int32 = 1 // 已处理
	SpamReportStatusIgnored  int32 = 2 // 已忽略
)

// SpamReport reason type constants.
const (
	SpamReasonTypeSpam    int32 = 1 // 垃圾消息
	SpamReasonTypePorn    int32 = 2 // 色情内容
	SpamReasonTypeIllegal int32 = 3 // 违法内容
	SpamReasonTypeOther   int32 = 4 // 其他
)

type SpamReport struct {
	ID             primitive.ObjectID `bson:"_id"`
	ReportID       string             `bson:"report_id"`
	ReporterUserID string             `bson:"reporter_user_id"`
	ReportedUserID string             `bson:"reported_user_id"`
	ConversationID string             `bson:"conversation_id"` // 举报具体消息时填写
	ClientMsgID    string             `bson:"client_msg_id"`   // 举报具体消息时填写
	Seq            int64              `bson:"seq"`
	ReasonType     int32              `bson:"reason_type"` // 1垃圾 2色情 3违法 4其他
	Reason         string             `bson:"reason"`
	Status         int32              `bson:"status"`          // 0待处理 1已处理 2已忽略
	CreateTime     time.Time          `bson:"create_time"`
	HandleTime     time.Time          `bson:"handle_time"`
	HandlerUserID  string             `bson:"handler_user_id"`
	Ex             string             `bson:"ex"`
}
