// Package msg provides message handling functionalities for the messaging server.
package msg

import (
	"context"

	"github.com/redis/go-redis/v9"

	// Importing necessary components from OpenIMSDK.
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/conversation"
	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
)

// GetConversationsHasReadAndMaxSeq retrieves the read status and maximum sequence number of specified conversations.
func (m *msgServer) GetConversationsHasReadAndMaxSeq(ctx context.Context, req *msg.GetConversationsHasReadAndMaxSeqReq) (resp *msg.GetConversationsHasReadAndMaxSeqResp, err error) {
	var conversationIDs []string

	// If no conversation IDs are provided, retrieve them from the local cache.
	if len(req.ConversationIDs) == 0 {
		conversationIDs, err = m.ConversationLocalCache.GetConversationIDs(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		conversationIDs = req.ConversationIDs
	}

	// Retrieve the read sequence numbers for the conversations.
	hasReadSeqs, err := m.MsgDatabase.GetHasReadSeqs(ctx, req.UserID, conversationIDs)
	if err != nil {
		return nil, err
	}

	// Retrieve conversation details.
	conversations, err := m.Conversation.GetConversations(ctx, req.UserID, conversationIDs)
	if err != nil {
		return nil, err
	}

	// Prepare a map to store the maximum sequence numbers.
	conversationMaxSeqMap := make(map[string]int64)
	for _, conversation := range conversations {
		if conversation.MaxSeq != 0 {
			conversationMaxSeqMap[conversation.ConversationID] = conversation.MaxSeq
		}
	}

	// Retrieve the maximum sequence numbers for the conversations.
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}

	// Prepare the response with the sequence information.
	resp = &msg.GetConversationsHasReadAndMaxSeqResp{Seqs: make(map[string]*msg.Seqs)}
	for conversationID, maxSeq := range maxSeqs {
		resp.Seqs[conversationID] = &msg.Seqs{
			HasReadSeq: hasReadSeqs[conversationID],
			MaxSeq:     maxSeq,
		}

		// Override the maximum sequence number if available in the map.
		if v, ok := conversationMaxSeqMap[conversationID]; ok {
			resp.Seqs[conversationID].MaxSeq = v
		}
	}
	return resp, nil
}

// SetConversationHasReadSeq updates the read sequence number for a specific conversation.
func (m *msgServer) SetConversationHasReadSeq(
	ctx context.Context,
	req *msg.SetConversationHasReadSeqReq,
) (resp *msg.SetConversationHasReadSeqResp, err error) {
	// Retrieve the maximum sequence number for the conversation.
	maxSeq, err := m.MsgDatabase.GetMaxSeq(ctx, req.ConversationID)
	if err != nil {
		return
	}

	// Validate the provided read sequence number.
	if req.HasReadSeq > maxSeq {
		return nil, errs.ErrArgs.Wrap("hasReadSeq must not be bigger than maxSeq")
	}

	// Update the read sequence number in the database.
	if err := m.MsgDatabase.SetHasReadSeq(ctx, req.UserID, req.ConversationID, req.HasReadSeq); err != nil {
		return nil, err
	}

	// Send a notification for the read status update.
	if err = m.sendMarkAsReadNotification(ctx, req.ConversationID, constant.SingleChatType, req.UserID,
		req.UserID, nil, req.HasReadSeq); err != nil {
		return
	}
	return &msg.SetConversationHasReadSeqResp{}, nil
}

// MarkMsgsAsRead marks specific messages in a conversation as read.
func (m *msgServer) MarkMsgsAsRead(ctx context.Context, req *msg.MarkMsgsAsReadReq) (resp *msg.MarkMsgsAsReadResp, err error) {
	// Ensure that the sequence numbers are provided.
	// Ensure that the sequence numbers are provided.
	if len(req.Seqs) < 1 {
		return nil, errs.ErrArgs.Wrap("seqs must not be empty")
	}

	// Retrieve the maximum sequence number for the conversation.
	maxSeq, err := m.MsgDatabase.GetMaxSeq(ctx, req.ConversationID)
	if err != nil {
		return
	}

	// Determine the highest sequence number from the request.
	hasReadSeq := req.Seqs[len(req.Seqs)-1]
	if hasReadSeq > maxSeq {
		return nil, errs.ErrArgs.Wrap("hasReadSeq must not be bigger than maxSeq")
	}

	// Retrieve conversation details.
	conversation, err := m.Conversation.GetConversation(ctx, req.UserID, req.ConversationID)
	if err != nil {
		return
	}

	// Mark the specified messages as read in the database.
	if err = m.MsgDatabase.MarkSingleChatMsgsAsRead(ctx, req.UserID, req.ConversationID, req.Seqs); err != nil {
		return
	}

	// Get the current read sequence number.
	currentHasReadSeq, err := m.MsgDatabase.GetHasReadSeq(ctx, req.UserID, req.ConversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return
	}

	// Update the read sequence number if the new value is greater.
	if hasReadSeq > currentHasReadSeq {
		err = m.MsgDatabase.SetHasReadSeq(ctx, req.UserID, req.ConversationID, hasReadSeq)
		if err != nil {
			return
		}
	}

	// Send a notification to indicate that messages have been marked as read.
	if err = m.sendMarkAsReadNotification(ctx, req.ConversationID, conversation.ConversationType, req.UserID,
		m.conversationAndGetRecvID(conversation, req.UserID), req.Seqs, hasReadSeq); err != nil {
		return
	}
	return &msg.MarkMsgsAsReadResp{}, nil
}

// MarkConversationAsRead marks an entire conversation as read.
func (m *msgServer) MarkConversationAsRead(ctx context.Context, req *msg.MarkConversationAsReadReq) (resp *msg.MarkConversationAsReadResp, err error) {
	// Retrieve conversation details.
	conversation, err := m.Conversation.GetConversation(ctx, req.UserID, req.ConversationID)
	if err != nil {
		return nil, err
	}

	// Get the current read sequence number.
	hasReadSeq, err := m.MsgDatabase.GetHasReadSeq(ctx, req.UserID, req.ConversationID)
	if err != nil && errs.Unwrap(err) != redis.Nil {
		return nil, err
	}

	// Generate the sequence numbers to be marked as read.
	seqs := generateSeqs(hasReadSeq, req)

	// Update the read status if there are new sequences to mark or if the hasReadSeq is greater.
	if len(seqs) > 0 || req.HasReadSeq > hasReadSeq {
		err = m.updateReadStatus(ctx, req, conversation, seqs, hasReadSeq)
		if err != nil {
			return nil, err
		}
	}
	return &msg.MarkConversationAsReadResp{}, nil
}

// generateSeqs creates a slice of sequence numbers that are greater than the provided hasReadSeq.
func generateSeqs(hasReadSeq int64, req *msg.MarkConversationAsReadReq) []int64 {
	var seqs []int64
	for _, val := range req.Seqs {
		if val > hasReadSeq && !utils2.Contain(val, seqs...) {
			seqs = append(seqs, val)
		}
	}
	return seqs
}

// updateReadStatus updates the read status for messages in a conversation.
func (m *msgServer) updateReadStatus(ctx context.Context, req *msg.MarkConversationAsReadReq, conversation *conversation.Conversation, seqs []int64, hasReadSeq int64) error {
	// Special handling for single chat type conversations.
	if conversation.ConversationType == constant.SingleChatType && len(seqs) > 0 {
		log.ZDebug(ctx, "MarkConversationAsRead", "seqs", seqs, "conversationID", req.ConversationID)
		if err := m.MsgDatabase.MarkSingleChatMsgsAsRead(ctx, req.UserID, req.ConversationID, seqs); err != nil {
			return err
		}
	}

	// Update the hasReadSeq if the new value is greater.
	if req.HasReadSeq > hasReadSeq {
		if err := m.MsgDatabase.SetHasReadSeq(ctx, req.UserID, req.ConversationID, req.HasReadSeq); err != nil {
			return err
		}
	}

	// Determine the receiver ID for the read receipt.
	recvID := m.conversationAndGetRecvID(conversation, req.UserID)
	// Adjust the receiver ID for specific conversation types.
	if conversation.ConversationType == constant.SuperGroupChatType || conversation.ConversationType == constant.NotificationChatType {
		recvID = req.UserID
	}

	// Send a notification to indicate the read status update.
	return m.sendMarkAsReadNotification(ctx, req.ConversationID, conversation.ConversationType, req.UserID, recvID, seqs, req.HasReadSeq)
}

// sendMarkAsReadNotification sends a notification about the read status update.
func (m *msgServer) sendMarkAsReadNotification(
	ctx context.Context,
	conversationID string,
	sessionType int32,
	sendID, recvID string,
	seqs []int64,
	hasReadSeq int64,
) error {
	// Construct the read receipt notification.
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
		HasReadSeq:       hasReadSeq,
	}
	// Send the notification with session type information.
	err := m.notificationSender.NotificationWithSesstionType(ctx, sendID, recvID, constant.HasReadReceipt, sessionType, tips)
	if err != nil {
		log.ZWarn(ctx, "send has read Receipt err", err)
	}
	return nil
}
