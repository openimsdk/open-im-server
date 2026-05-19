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

// GroupMsgBurnRecord 记录群消息「阅后即焚」的阅读进度。
//
// 写入时机：某成员首次通过 MarkConversationAsRead 标记某 seq 为已读时创建记录，
// 后续成员已读时通过原子 $inc 累加 ReadCount。
//
// 删除时机：cron 发现 ReadCount >= MemberCount 且 BurnEndTime <= now 时触发删除，
// 同步推进所有成员的 min_seq 并发送会话变更通知。
type GroupMsgBurnRecord struct {
	// GroupID 群组 ID
	GroupID string `bson:"group_id"`
	// Seq 消息序列号
	Seq int64 `bson:"seq"`
	// ReadCount 已阅读该消息的成员数量（原子累加）
	ReadCount int32 `bson:"read_count"`
	// MemberCount 创建记录时的群成员总数；用于判断是否全员已读
	MemberCount int32 `bson:"member_count"`
	// BurnEndTime 消息焚毁截止时间戳（毫秒）；首次阅读时写入，不再覆盖
	BurnEndTime int64 `bson:"burn_end_time"`
	// CreateTime 记录创建时间戳（毫秒）
	CreateTime int64 `bson:"create_time"`
}
