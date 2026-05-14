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

// MsgBurnDeadline 单条消息的「阅后即焚截止时间」记录。
// 在 MarkMsgsAsRead / MarkConversationAsRead（单聊场景）时按
// (UserID, ConversationID, Seq) 维度写入。读取时间锁定后续不再覆盖。
//
// 当当前时间 > DeadlineMs 时，cron 会把该用户在该会话上的 min_seq
// 推进到 max(已过期 seq) + 1，从而让这些消息从对该用户的拉取结果中消失。
type MsgBurnDeadline struct {
	UserID         string `bson:"user_id"`
	ConversationID string `bson:"conversation_id"`
	Seq            int64  `bson:"seq"`
	// PeerID 单聊中的对端用户 ID。
	// cron 处理时可直接获取对端，无需额外查询 conversation 表。
	PeerID string `bson:"peer_id"`
	// DeadlineMs 截止时间戳（毫秒）；超过即可被 cron 收走推进 min_seq。
	DeadlineMs int64 `bson:"deadline_ms"`
	// CreateTime 写入时刻（毫秒）；用于排查/审计。
	CreateTime int64 `bson:"create_time"`
}
