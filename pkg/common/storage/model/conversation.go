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

type Conversation struct {
	OwnerUserID           string    `bson:"owner_user_id"`
	ConversationID        string    `bson:"conversation_id"`
	ConversationType      int32     `bson:"conversation_type"`
	UserID                string    `bson:"user_id"`
	GroupID               string    `bson:"group_id"`
	RecvMsgOpt            int32     `bson:"recv_msg_opt"`
	IsPinned              bool      `bson:"is_pinned"`
	IsPrivateChat         bool      `bson:"is_private_chat"`
	BurnDuration          int32     `bson:"burn_duration"`
	GroupAtType           int32     `bson:"group_at_type"`
	AttachedInfo          string    `bson:"attached_info"`
	Ex                    string    `bson:"ex"`
	MaxSeq                int64     `bson:"max_seq"`
	MinSeq                int64     `bson:"min_seq"`
	CreateTime            time.Time `bson:"create_time"`
	IsMsgDestruct         bool      `bson:"is_msg_destruct"`
	MsgDestructTime       int64     `bson:"msg_destruct_time"`
	LatestMsgDestructTime time.Time `bson:"latest_msg_destruct_time"`
}
