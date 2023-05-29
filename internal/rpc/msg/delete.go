package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

func (m *msgServer) getMinSeqs(maxSeqs map[string]int64) map[string]int64 {
	minSeqs := make(map[string]int64)
	for k, v := range maxSeqs {
		minSeqs[k] = v + 1
	}
	return minSeqs
}

func (m *msgServer) ClearConversationsMsg(ctx context.Context, req *msg.ClearConversationsMsgReq) (*msg.ClearConversationsMsgResp, error) {
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	if err := m.MsgDatabase.SetUserConversationsMinSeqs(ctx, req.UserID, m.getMinSeqs(maxSeqs)); err != nil {
		return nil, err
	}
	m.conversationClearSync(ctx, req.DeleteSyncOpt, req.UserID, req.ConversationIDs)
	return &msg.ClearConversationsMsgResp{}, nil
}

func (m *msgServer) UserClearAllMsg(ctx context.Context, req *msg.UserClearAllMsgReq) (*msg.UserClearAllMsgResp, error) {
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	conversationIDs, err := m.ConversationLocalCache.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "GetMaxSeq", "conversationIDs", conversationIDs)
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}
	if err := m.MsgDatabase.SetUserConversationsMinSeqs(ctx, req.UserID, m.getMinSeqs(maxSeqs)); err != nil {
		return nil, err
	}
	m.conversationClearSync(ctx, req.DeleteSyncOpt, req.UserID, conversationIDs)
	return &msg.UserClearAllMsgResp{}, nil
}

func (m *msgServer) DeleteMsgs(ctx context.Context, req *msg.DeleteMsgsReq) (*msg.DeleteMsgsResp, error) {
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	if err := m.MsgDatabase.DeleteUserMsgsBySeqs(ctx, req.UserID, req.ConversationID, req.Seqs); err != nil {
		return nil, err
	}
	m.DeleteMsgsNotification(ctx, req.ConversationID, req.UserID, req.Seqs, req.DeleteSyncOpt)
	return &msg.DeleteMsgsResp{}, nil
}

func (m *msgServer) DeleteMsgPhysicalBySeq(ctx context.Context, req *msg.DeleteMsgPhysicalBySeqReq) (*msg.DeleteMsgPhysicalBySeqResp, error) {
	err := m.MsgDatabase.DeleteMsgsPhysicalBySeqs(ctx, req.ConversationID, req.Seqs)
	if err != nil {
		return nil, err
	}
	return &msg.DeleteMsgPhysicalBySeqResp{}, nil
}

func (m *msgServer) DeleteMsgPhysical(ctx context.Context, req *msg.DeleteMsgPhysicalReq) (*msg.DeleteMsgPhysicalResp, error) {
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	for _, conversationID := range req.ConversationIDs {
		if err := m.MsgDatabase.DeleteConversationMsgsAndSetMinSeq(ctx, conversationID, req.RemainTime); err != nil {
			log.ZWarn(ctx, "DeleteConversationMsgsAndSetMinSeq error", err, "conversationID", conversationID, "err", err)
		}
	}
	return &msg.DeleteMsgPhysicalResp{}, nil
}

func (m *msgServer) conversationClearSync(ctx context.Context, opt *msg.DeleteSyncOpt, userID string, conversationIDs []string) {
	if opt == nil {
		return
	}
	for _, conversationID := range conversationIDs {
		conversation, err := m.Conversation.GetConversation(ctx, userID, conversationID)
		if err != nil {
			log.ZWarn(ctx, "GetConversation error", err, "conversationID", conversationID, "userID", userID)
			continue
		}
		if conversation.ConversationType == constant.SingleChatType || conversation.ConversationType == constant.NotificationChatType {

		} else if conversation.ConversationType == constant.SuperGroupChatType {

		}
	}

	if opt.IsSyncSelf {

	} else if opt.IsSyncOther {
	}

}

func (m *msgServer) DeleteMsgsNotification(ctx context.Context, conversationID, userID string, seqs []int64, opt *msg.DeleteSyncOpt) error {
	return nil
}
