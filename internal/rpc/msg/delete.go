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

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/timeutil"
)

func (m *msgServer) getMinSeqs(maxSeqs map[string]int64) map[string]int64 {
	minSeqs := make(map[string]int64)
	for k, v := range maxSeqs {
		minSeqs[k] = v + 1
	}
	return minSeqs
}

func (m *msgServer) validateDeleteSyncOpt(opt *msg.DeleteSyncOpt) (isSyncSelf, isSyncOther bool) {
	if opt == nil {
		return
	}
	return opt.IsSyncSelf, opt.IsSyncOther
}

func (m *msgServer) ClearConversationsMsg(ctx context.Context, req *msg.ClearConversationsMsgReq) (*msg.ClearConversationsMsgResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	if err := m.clearConversation(ctx, req.ConversationIDs, req.UserID, req.DeleteSyncOpt); err != nil {
		return nil, err
	}
	return &msg.ClearConversationsMsgResp{}, nil
}

func (m *msgServer) UserClearAllMsg(ctx context.Context, req *msg.UserClearAllMsgReq) (*msg.UserClearAllMsgResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	conversationIDs, err := m.ConversationLocalCache.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if err := m.clearConversation(ctx, conversationIDs, req.UserID, req.DeleteSyncOpt); err != nil {
		return nil, err
	}
	return &msg.UserClearAllMsgResp{}, nil
}

func (m *msgServer) DeleteMsgs(ctx context.Context, req *msg.DeleteMsgsReq) (*msg.DeleteMsgsResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	isSyncSelf, isSyncOther := m.validateDeleteSyncOpt(req.DeleteSyncOpt)
	if isSyncOther {
		if err := m.MsgDatabase.DeleteMsgsPhysicalBySeqs(ctx, req.ConversationID, req.Seqs); err != nil {
			return nil, err
		}
		conversations, err := m.Conversation.GetConversationsByConversationID(ctx, []string{req.ConversationID})
		if err != nil {
			return nil, err
		}
		tips := &sdkws.DeleteMsgsTips{UserID: req.UserID, ConversationID: req.ConversationID, Seqs: req.Seqs}
		m.notificationSender.NotificationWithSessionType(
			ctx,
			req.UserID,
			m.conversationAndGetRecvID(conversations[0], req.UserID),
			constant.DeleteMsgsNotification,
			conversations[0].ConversationType,
			tips,
		)
	} else {
		if err := m.MsgDatabase.DeleteUserMsgsBySeqs(ctx, req.UserID, req.ConversationID, req.Seqs); err != nil {
			return nil, err
		}
		if isSyncSelf {
			tips := &sdkws.DeleteMsgsTips{UserID: req.UserID, ConversationID: req.ConversationID, Seqs: req.Seqs}
			m.notificationSender.NotificationWithSessionType(ctx, req.UserID, req.UserID, constant.DeleteMsgsNotification, constant.SingleChatType, tips)
		}
	}
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
	if err := authverify.CheckAdmin(ctx, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	remainTime := timeutil.GetCurrentTimestampBySecond() - req.Timestamp
	for _, conversationID := range req.ConversationIDs {
		if err := m.MsgDatabase.DeleteConversationMsgsAndSetMinSeq(ctx, conversationID, remainTime); err != nil {
			log.ZWarn(ctx, "DeleteConversationMsgsAndSetMinSeq error", err, "conversationID", conversationID, "err", err)
		}
	}
	return &msg.DeleteMsgPhysicalResp{}, nil
}

func (m *msgServer) clearConversation(ctx context.Context, conversationIDs []string, userID string, deleteSyncOpt *msg.DeleteSyncOpt) error {
	conversations, err := m.Conversation.GetConversationsByConversationID(ctx, conversationIDs)
	if err != nil {
		return err
	}
	var existConversations []*conversation.Conversation
	var existConversationIDs []string
	for _, conversation := range conversations {
		existConversations = append(existConversations, conversation)
		existConversationIDs = append(existConversationIDs, conversation.ConversationID)
	}
	log.ZDebug(ctx, "ClearConversationsMsg", "existConversationIDs", existConversationIDs)
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, existConversationIDs)
	if err != nil {
		return err
	}
	isSyncSelf, isSyncOther := m.validateDeleteSyncOpt(deleteSyncOpt)
	if !isSyncOther {
		setSeqs := m.getMinSeqs(maxSeqs)
		if err := m.MsgDatabase.SetUserConversationsMinSeqs(ctx, userID, setSeqs); err != nil {
			return err
		}
		ownerUserIDs := []string{userID}
		for conversationID, seq := range setSeqs {
			if err := m.Conversation.SetConversationMinSeq(ctx, ownerUserIDs, conversationID, seq); err != nil {
				return err
			}
		}
		// notification 2 self
		if isSyncSelf {
			tips := &sdkws.ClearConversationTips{UserID: userID, ConversationIDs: existConversationIDs}
			m.notificationSender.NotificationWithSessionType(ctx, userID, userID, constant.ClearConversationNotification, constant.SingleChatType, tips)
		}
	} else {
		if err := m.MsgDatabase.SetMinSeqs(ctx, m.getMinSeqs(maxSeqs)); err != nil {
			return err
		}
		for _, conversation := range existConversations {
			tips := &sdkws.ClearConversationTips{UserID: userID, ConversationIDs: []string{conversation.ConversationID}}
			m.notificationSender.NotificationWithSessionType(ctx, userID, m.conversationAndGetRecvID(conversation, userID), constant.ClearConversationNotification, conversation.ConversationType, tips)
		}
	}
	if err := m.MsgDatabase.UserSetHasReadSeqs(ctx, userID, maxSeqs); err != nil {
		return err
	}
	return nil
}
