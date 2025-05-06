package msg

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
)

func (m *msgServer) getAndClearUnread(ctx context.Context, userID string, conversationIDs []string, clearCount bool) (map[string]int64, error) {
	if err := authverify.CheckAccessV3(ctx, userID, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	readSeqs, err := m.MsgDatabase.GetHasReadSeqs(ctx, userID, conversationIDs)
	if err != nil {
		return nil, err
	}
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}
	conversations, err := m.ConversationLocalCache.GetConversations(ctx, userID, conversationIDs)
	if err != nil {
		return nil, err
	}
	unreadCount := make(map[string]int64)
	for _, conversation := range conversations {
		maxSeq := conversation.MaxSeq
		if maxSeq == 0 {
			maxSeq = maxSeqs[conversation.ConversationID]
		}
		count := maxSeq - readSeqs[conversation.ConversationID]
		if count < 0 {
			count = 0
		}
		unreadCount[conversation.ConversationID] = count
		if count > 0 && clearCount {
			if err := m.MsgDatabase.SetHasReadSeq(ctx, userID, conversation.ConversationID, maxSeq); err != nil {
				return nil, err
			}
			m.sendMarkAsReadNotification(ctx, conversation.ConversationID, constant.SingleChatType, userID, userID, nil, maxSeq)
		}
	}
	return unreadCount, nil
}

func (m *msgServer) GetConversationsUnreadCount(ctx context.Context, req *msg.GetConversationsUnreadCountReq) (*msg.GetConversationsUnreadCountResp, error) {
	res, err := m.getAndClearUnread(ctx, req.UserID, req.ConversationIDs, false)
	if err != nil {
		return nil, err
	}
	return &msg.GetConversationsUnreadCountResp{UnreadCount: res}, nil
}

func (m *msgServer) ClearConversationsUnreadCount(ctx context.Context, req *msg.ClearConversationsUnreadCountReq) (*msg.ClearConversationsUnreadCountResp, error) {
	res, err := m.getAndClearUnread(ctx, req.UserID, req.ConversationIDs, true)
	if err != nil {
		return nil, err
	}
	return &msg.ClearConversationsUnreadCountResp{UnreadCount: res}, nil
}
