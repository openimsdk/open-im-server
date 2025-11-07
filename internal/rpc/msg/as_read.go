package msg

import (
	"context"
	"errors"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
)

func (m *msgServer) GetConversationsHasReadAndMaxSeq(ctx context.Context, req *msg.GetConversationsHasReadAndMaxSeqReq) (*msg.GetConversationsHasReadAndMaxSeqResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
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
	if req.ReturnPinned {
		pinnedConversationIDs, err := m.ConversationLocalCache.GetPinnedConversationIDs(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		resp.PinnedConversationIDs = pinnedConversationIDs
	}
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
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
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
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
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
	m.sendMarkAsReadNotification(ctx, req.ConversationID, conversation.ConversationType, req.UserID,
		m.conversationAndGetRecvID(conversation, req.UserID), req.Seqs, hasReadSeq)
	return &msg.MarkMsgsAsReadResp{}, nil
}

func (m *msgServer) MarkConversationAsRead(ctx context.Context, req *msg.MarkConversationAsReadReq) (*msg.MarkConversationAsReadResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
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

func (m *msgServer) sendMarkAsReadNotification(ctx context.Context, conversationID string, sessionType int32, sendID, recvID string, seqs []int64, hasReadSeq int64) {
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
		HasReadSeq:       hasReadSeq,
	}
	m.notificationSender.NotificationWithSessionType(ctx, sendID, recvID, constant.HasReadReceipt, sessionType, tips)
}
