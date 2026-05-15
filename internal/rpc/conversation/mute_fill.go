// Copyright © 2023 OpenIM. All rights reserved.

package conversation

import (
	"context"
	"time"

	pbconversation "github.com/openimsdk/protocol/conversation"
)

// computeIsMuted 根据会话模型中存储的 mute_duration 和 mute_end_time 计算当前是否处于静音状态：
//   - duration == 0 且 end == 0：未静音
//   - duration == -1 且 end == 0：永久静音
//   - end > 0 且 end > now：定时静音仍有效
//   - end > 0 且 end <= now：定时静音已过期，视为未静音
func computeIsMuted(muteDuration int32, muteEndTime int64) bool {
	if muteDuration == 0 && muteEndTime == 0 {
		return false
	}
	if muteDuration == -1 && muteEndTime == 0 {
		return true
	}
	return muteEndTime > time.Now().Unix()
}

// fillConversationUserMute 根据会话模型字段（已由 ConversationDB2Pb 通过 CopyStructFields 填入
// conv.MuteDuration / conv.MuteEndTime）计算并设置 conv.IsMuted，无需额外数据库查询。
func (c *conversationServer) fillConversationUserMute(_ context.Context, conv *pbconversation.Conversation) {
	if conv == nil {
		return
	}
	conv.IsMuted = computeIsMuted(conv.MuteDuration, conv.MuteEndTime)
}

func (c *conversationServer) fillConversationsUserMute(ctx context.Context, list []*pbconversation.Conversation) {
	for _, conv := range list {
		c.fillConversationUserMute(ctx, conv)
	}
}
