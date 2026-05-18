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

package msg

import (
	"context"
	"errors"
	"time"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
)

func (m *msgServer) GetConversationsHasReadAndMaxSeq(ctx context.Context, req *msg.GetConversationsHasReadAndMaxSeqReq) (*msg.GetConversationsHasReadAndMaxSeqResp, error) {
	var conversationIDs []string
	if len(req.ConversationIDs) == 0 {
		var err error
		conversationIDs, err = m.ConversationLocalCache.GetConversationIDs(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		conversationIDs = req.ConversationIDs
	}

	hasReadSeqs, err := m.MsgDatabase.GetHasReadSeqs(ctx, req.UserID, conversationIDs)
	if err != nil {
		return nil, err
	}

	conversations, err := m.ConversationLocalCache.GetConversations(ctx, req.UserID, conversationIDs)
	if err != nil {
		return nil, err
	}

	conversationMaxSeqMap := make(map[string]int64)
	for _, conversation := range conversations {
		if conversation.MaxSeq != 0 {
			conversationMaxSeqMap[conversation.ConversationID] = conversation.MaxSeq
		}
	}
	maxSeqs, err := m.MsgDatabase.GetMaxSeqsWithTime(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}
	resp := &msg.GetConversationsHasReadAndMaxSeqResp{Seqs: make(map[string]*msg.Seqs)}
	for conversationID, maxSeq := range maxSeqs {
		resp.Seqs[conversationID] = &msg.Seqs{
			HasReadSeq: hasReadSeqs[conversationID],
			MaxSeq:     maxSeq.Seq,
			MaxSeqTime: maxSeq.Time,
		}
		if v, ok := conversationMaxSeqMap[conversationID]; ok {
			resp.Seqs[conversationID].MaxSeq = v
		}
	}
	return resp, nil
}

func (m *msgServer) SetConversationHasReadSeq(ctx context.Context, req *msg.SetConversationHasReadSeqReq) (*msg.SetConversationHasReadSeqResp, error) {
	maxSeq, err := m.MsgDatabase.GetMaxSeq(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	if req.HasReadSeq > maxSeq {
		return nil, errs.ErrArgs.WrapMsg("hasReadSeq must not be bigger than maxSeq")
	}
	if err := m.MsgDatabase.SetHasReadSeq(ctx, req.UserID, req.ConversationID, req.HasReadSeq); err != nil {
		return nil, err
	}
	m.sendMarkAsReadNotification(ctx, req.ConversationID, constant.SingleChatType, req.UserID, req.UserID, nil, req.HasReadSeq)
	return &msg.SetConversationHasReadSeqResp{}, nil
}

func (m *msgServer) MarkMsgsAsRead(ctx context.Context, req *msg.MarkMsgsAsReadReq) (*msg.MarkMsgsAsReadResp, error) {
	if len(req.Seqs) < 1 {
		return nil, errs.ErrArgs.WrapMsg("seqs must not be empty")
	}
	maxSeq, err := m.MsgDatabase.GetMaxSeq(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	hasReadSeq := req.Seqs[len(req.Seqs)-1]
	if hasReadSeq > maxSeq {
		return nil, errs.ErrArgs.WrapMsg("hasReadSeq must not be bigger than maxSeq")
	}
	conversation, err := m.ConversationLocalCache.GetConversation(ctx, req.UserID, req.ConversationID)
	if err != nil {
		return nil, err
	}
	if err := m.MsgDatabase.MarkSingleChatMsgsAsRead(ctx, req.UserID, req.ConversationID, req.Seqs); err != nil {
		return nil, err
	}
	currentHasReadSeq, err := m.MsgDatabase.GetHasReadSeq(ctx, req.UserID, req.ConversationID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if hasReadSeq > currentHasReadSeq {
		err = m.MsgDatabase.SetHasReadSeq(ctx, req.UserID, req.ConversationID, hasReadSeq)
		if err != nil {
			return nil, err
		}
	}

	reqCallback := &cbapi.CallbackSingleMsgReadReq{
		ConversationID: conversation.ConversationID,
		UserID:         req.UserID,
		Seqs:           req.Seqs,
		ContentType:    conversation.ConversationType,
	}
	m.webhookAfterSingleMsgRead(ctx, &m.config.WebhooksConfig.AfterSingleMsgRead, reqCallback)
	m.recordBurnDeadlines(ctx, conversation, req.UserID, req.Seqs)
	m.sendMarkAsReadNotification(ctx, req.ConversationID, conversation.ConversationType, req.UserID,
		m.conversationAndGetRecvID(conversation, req.UserID), req.Seqs, hasReadSeq)
	return &msg.MarkMsgsAsReadResp{}, nil
}

func (m *msgServer) MarkConversationAsRead(ctx context.Context, req *msg.MarkConversationAsReadReq) (*msg.MarkConversationAsReadResp, error) {
	conversation, err := m.ConversationLocalCache.GetConversation(ctx, req.UserID, req.ConversationID)
	if err != nil {
		return nil, err
	}
	hasReadSeq, err := m.MsgDatabase.GetHasReadSeq(ctx, req.UserID, req.ConversationID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	var seqs []int64

	log.ZDebug(ctx, "MarkConversationAsRead", "hasReadSeq", hasReadSeq, "req.HasReadSeq", req.HasReadSeq)
	if conversation.ConversationType == constant.SingleChatType {
		for i := hasReadSeq + 1; i <= req.HasReadSeq; i++ {
			seqs = append(seqs, i)
		}
		// avoid client missed call MarkConversationMessageAsRead by order
		for _, val := range req.Seqs {
			if !datautil.Contain(val, seqs...) {
				seqs = append(seqs, val)
			}
		}
		if len(seqs) > 0 {
			log.ZDebug(ctx, "MarkConversationAsRead", "seqs", seqs, "conversationID", req.ConversationID)
			if err = m.MsgDatabase.MarkSingleChatMsgsAsRead(ctx, req.UserID, req.ConversationID, seqs); err != nil {
				return nil, err
			}
		}
		if req.HasReadSeq > hasReadSeq {
			err = m.MsgDatabase.SetHasReadSeq(ctx, req.UserID, req.ConversationID, req.HasReadSeq)
			if err != nil {
				return nil, err
			}
			hasReadSeq = req.HasReadSeq
		}
		m.recordBurnDeadlines(ctx, conversation, req.UserID, seqs)
		m.sendMarkAsReadNotification(ctx, req.ConversationID, conversation.ConversationType, req.UserID,
			m.conversationAndGetRecvID(conversation, req.UserID), seqs, hasReadSeq)
	} else if conversation.ConversationType == constant.ReadGroupChatType ||
		conversation.ConversationType == constant.NotificationChatType {
		if req.HasReadSeq > hasReadSeq {
			err = m.MsgDatabase.SetHasReadSeq(ctx, req.UserID, req.ConversationID, req.HasReadSeq)
			if err != nil {
				return nil, err
			}
			hasReadSeq = req.HasReadSeq
		}
		m.sendMarkAsReadNotification(ctx, req.ConversationID, constant.SingleChatType, req.UserID,
			req.UserID, seqs, hasReadSeq)
	}

	if conversation.ConversationType == constant.SingleChatType {
		reqCall := &cbapi.CallbackSingleMsgReadReq{
			ConversationID: conversation.ConversationID,
			UserID:         conversation.OwnerUserID,
			Seqs:           req.Seqs,
			ContentType:    conversation.ConversationType,
		}
		m.webhookAfterSingleMsgRead(ctx, &m.config.WebhooksConfig.AfterSingleMsgRead, reqCall)
	} else if conversation.ConversationType == constant.ReadGroupChatType {
		reqCall := &cbapi.CallbackGroupMsgReadReq{
			SendID:       conversation.OwnerUserID,
			ReceiveID:    req.UserID,
			UnreadMsgNum: req.HasReadSeq,
			ContentType:  int64(conversation.ConversationType),
		}
		m.webhookAfterGroupMsgRead(ctx, &m.config.WebhooksConfig.AfterGroupMsgRead, reqCall)
	}
	return &msg.MarkConversationAsReadResp{}, nil
}

// recordBurnDeadlines 在「单聊」场景下，为本次已读的每条消息同时给接收者和发送者
// 各记录一份「阅后即焚」截止时间。cron 到期后会分别删除双方该消息，双方都看不到。
//
// 销毁时长优先级：
//  1. 会话级 BurnDuration（通过 /conversation/set_burn 设置）；
//  2. 对端（发送者）全局 MsgBurnDuration。
//
// 设计要点：
//  1. 仅单聊。
//  2. $setOnInsert 确保同一 (UserID, ConversationID, Seq) 已存在时不覆盖，
//     以「首次阅读时刻」为 deadline 基准，多端重复 MarkAsRead 不会往后推。
//  3. 失败仅记录日志，不影响已读主流程。
func (m *msgServer) recordBurnDeadlines(ctx context.Context, conv *conversation.Conversation, readerUserID string, seqs []int64) {
	if len(seqs) == 0 {
		return
	}
	if conv.ConversationType != constant.SingleChatType {
		return
	}
	peerID := m.conversationAndGetRecvID(conv, readerUserID)
	if peerID == "" || peerID == readerUserID {
		return
	}

	// 优先使用会话级 BurnDuration（双方协商后保存到会话），否则回退到发送者全局设置。
	burnSeconds := conv.BurnDuration
	if burnSeconds <= 0 {
		peerInfo, err := m.UserLocalCache.GetUserInfo(ctx, peerID)
		if err != nil {
			log.ZWarn(ctx, "recordBurnDeadlines GetUserInfo failed", err, "peerID", peerID)
			return
		}
		if peerInfo == nil || peerInfo.MsgBurnDuration <= 0 {
			return
		}
		burnSeconds = peerInfo.MsgBurnDuration
	}

	now := time.Now().UnixMilli()
	deadline := now + int64(burnSeconds)*1000
	// 每条消息同时为接收者和发送者各写一条 deadline，双方消息同步焚毁。
	items := make([]*model.MsgBurnDeadline, 0, len(seqs)*2)
	for _, seq := range seqs {
		items = append(items,
			&model.MsgBurnDeadline{
				UserID:         readerUserID,
				ConversationID: conv.ConversationID,
				Seq:            seq,
				PeerID:         peerID,
				DeadlineMs:     deadline,
				CreateTime:     now,
			},
			&model.MsgBurnDeadline{
				UserID:         peerID,
				ConversationID: conv.ConversationID,
				Seq:            seq,
				PeerID:         readerUserID,
				DeadlineMs:     deadline,
				CreateTime:     now,
			},
		)
	}
	if err := m.msgBurnDeadlineDB.UpsertIfAbsent(ctx, items); err != nil {
		log.ZError(ctx, "recordBurnDeadlines UpsertIfAbsent failed", err,
			"readerUserID", readerUserID, "peerID", peerID,
			"conversationID", conv.ConversationID, "seqs", seqs)
	}
}

func (m *msgServer) sendMarkAsReadNotification(ctx context.Context, conversationID string, sessionType int32, sendID, recvID string, seqs []int64, hasReadSeq int64) {
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
		HasReadSeq:       hasReadSeq,
	}
	m.notificationSender.NotificationWithSessionType(ctx, sendID, recvID, constant.HasReadReceipt, sessionType, tips)

}
