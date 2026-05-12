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

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// ExpiredBurnGroup 表示某个 (UserID, ConversationID) 上需要被推进 min_seq 的 seq 集合。
type ExpiredBurnGroup struct {
	UserID         string
	ConversationID string
	// MaxSeq 当前批次中最大的过期 seq；推进 min_seq 时使用 MaxSeq + 1。
	MaxSeq int64
	// Seqs 当前批次实际涉及的所有过期 seq，便于精确删除已处理的 deadline 记录。
	Seqs []int64
}

// MsgBurnDeadline 持久化每条消息对每个用户的「阅后即焚截止时间」。
// 写入位置：msg 服务 MarkMsgsAsRead / MarkConversationAsRead（单聊）。
// 消费位置：conversation 服务 ClearBurnExpiredMsgs cron 入口。
type MsgBurnDeadline interface {
	// UpsertIfAbsent 仅在 (UserID, ConversationID, Seq) 不存在时插入；
	// 已存在则不覆盖，保证「首次阅读时刻」决定 deadline。
	UpsertIfAbsent(ctx context.Context, items []*model.MsgBurnDeadline) error

	// FindExpiredGroups 查询 deadline_ms <= nowMs 的记录，按 (UserID, ConversationID)
	// 聚合并返回每组的最大 seq 与所涉及的 seq 列表。limit 限制返回的 group 数量。
	FindExpiredGroups(ctx context.Context, nowMs int64, limit int) ([]*ExpiredBurnGroup, error)

	// DeleteByUserConversationSeqs 删除某 (UserID, ConversationID) 下指定 seq 列表的 deadline 记录。
	// 一般在成功推进 min_seq 后调用。
	DeleteByUserConversationSeqs(ctx context.Context, userID, conversationID string, seqs []int64) error
}
