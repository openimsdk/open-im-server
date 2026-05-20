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

package database

import "context"

const GroupMsgBurnRecordName = "group_msg_burn_record"

// ExpiredGroupBurn 表示一个群中某批到期消息的聚合结果。
type ExpiredGroupBurn struct {
	// GroupID 群组 ID
	GroupID string
	// MaxSeq 该批中最大的 seq；推进 min_seq 时使用 MaxSeq + 1
	MaxSeq int64
	// Seqs 该批所有到期的 seq 列表，用于精确删除已处理记录
	Seqs []int64
}

// GroupMsgBurnRecord 持久化群消息「阅后即焚」的阅读计数与截止时间。
//
// 写入：msg 服务 MarkConversationAsRead 群聊分支。
// 消费：conversation 服务 ClearGroupBurnExpiredMsgs cron 入口。
type GroupMsgBurnRecord interface {
	// UpsertOnRead 批量原子更新阅读记录：
	//   - 若 (group_id, seq) 不存在：插入 {member_count, burn_end_time, create_time, send_id, read_count=1}；send_id 来自 seqSenderID[seq]，可为空。
	//   - 若已存在：仅对 read_count 执行 $inc，不覆盖首次写入的 burn_end_time、send_id
	UpsertOnRead(ctx context.Context, groupID string, seqs []int64, seqSenderID map[int64]string, memberCount int32, burnEndTimeMs int64) error

	// FindExpired 查询满足以下条件的记录并按 group_id 聚合：
	//   burn_end_time <= nowMs AND read_count >= member_count
	// limit 限制返回的 group 数量。
	FindExpired(ctx context.Context, nowMs int64, limit int) ([]*ExpiredGroupBurn, error)

	// DeleteByGroupSeqs 删除指定群下一批 seq 的记录，在成功推进 min_seq 后调用。
	DeleteByGroupSeqs(ctx context.Context, groupID string, seqs []int64) error
}
