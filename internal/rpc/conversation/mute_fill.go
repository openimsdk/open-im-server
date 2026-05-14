// Copyright © 2023 OpenIM. All rights reserved.

package conversation

import (
	"context"
	"math"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/constant"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/log"
)

// int64MuteDurationToProto 将 user_mute 的秒数配置写入 Conversation.muteDuration（int32）；正数过大时截断。
func int64MuteDurationToProto(d int64) int32 {
	if d > int64(math.MaxInt32) {
		return math.MaxInt32
	}
	if d < int64(math.MinInt32) {
		return math.MinInt32
	}
	return int32(d)
}

// conversationMuteFromRecord 与 relation.GetMute 判定一致：未记录/已过期则未静音；永久为 duration=-1 且 end=0。
func conversationMuteFromRecord(rec *model.UserMute, nowUnix int64) (isMuted bool, muteDuration int32, muteEndTime int64) {
	if rec == nil {
		return false, 0, 0
	}
	if rec.MuteEndTime != 0 && rec.MuteEndTime <= nowUnix {
		return false, 0, 0
	}
	d := rec.MuteDuration
	if d == 0 && rec.MuteEndTime == 0 {
		d = -1
	}
	md := int64MuteDurationToProto(d)
	me := rec.MuteEndTime
	isMuted = (md == -1) || (me > nowUnix)
	return isMuted, md, me
}

func (c *conversationServer) fillConversationUserMute(ctx context.Context, conv *pbconversation.Conversation) {
	if c == nil || c.userMuteDB == nil || conv == nil {
		return
	}
	if conv.ConversationType != constant.SingleChatType || conv.UserID == "" {
		return
	}
	rec, err := c.userMuteDB.Get(ctx, conv.OwnerUserID, conv.UserID)
	if err != nil {
		log.ZWarn(ctx, "fillConversationUserMute Get", err, "owner", conv.OwnerUserID, "peer", conv.UserID)
		return
	}
	now := time.Now().Unix()
	isMuted, dur, end := conversationMuteFromRecord(rec, now)
	conv.IsMuted = isMuted
	conv.MuteDuration = dur
	conv.MuteEndTime = end
}

func (c *conversationServer) fillConversationsUserMute(ctx context.Context, list []*pbconversation.Conversation) {
	if len(list) == 0 {
		return
	}
	for _, conv := range list {
		c.fillConversationUserMute(ctx, conv)
	}
}

func (c *conversationServer) fillConversationElemUserMute(
	ctx context.Context,
	db controller.UserMuteDatabase,
	ownerUserID string,
	elem *pbconversation.ConversationElem,
	conversationType int32,
	peerUserID string,
) {
	if db == nil || elem == nil || ownerUserID == "" {
		return
	}
	if conversationType != constant.SingleChatType || peerUserID == "" {
		return
	}
	rec, err := db.Get(ctx, ownerUserID, peerUserID)
	if err != nil {
		log.ZWarn(ctx, "fillConversationElemUserMute Get", err, "owner", ownerUserID, "peer", peerUserID)
		return
	}
	now := time.Now().Unix()
	isMuted, dur, end := conversationMuteFromRecord(rec, now)
	elem.IsMuted = isMuted
	elem.MuteDuration = dur
	elem.MuteEndTime = end
}
