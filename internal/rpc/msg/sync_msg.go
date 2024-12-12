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
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/open-im-server/v3/pkg/util/conversationutil"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/timeutil"
)

func (m *msgServer) PullMessageBySeqs(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	resp := &sdkws.PullMessageBySeqsResp{}
	resp.Msgs = make(map[string]*sdkws.PullMsgs)
	resp.NotificationMsgs = make(map[string]*sdkws.PullMsgs)
	for _, seq := range req.SeqRanges {
		if !msgprocessor.IsNotification(seq.ConversationID) {
			conversation, err := m.ConversationLocalCache.GetConversation(ctx, req.UserID, seq.ConversationID)
			if err != nil {
				log.ZError(ctx, "GetConversation error", err, "conversationID", seq.ConversationID)
				continue
			}
			minSeq, maxSeq, msgs, err := m.MsgDatabase.GetMsgBySeqsRange(ctx, req.UserID, seq.ConversationID,
				seq.Begin, seq.End, seq.Num, conversation.MaxSeq)
			if err != nil {
				log.ZWarn(ctx, "GetMsgBySeqsRange error", err, "conversationID", seq.ConversationID, "seq", seq)
				continue
			}
			var isEnd bool
			switch req.Order {
			case sdkws.PullOrder_PullOrderAsc:
				isEnd = maxSeq <= seq.End
			case sdkws.PullOrder_PullOrderDesc:
				isEnd = seq.Begin <= minSeq
			}
			if len(msgs) == 0 {
				log.ZWarn(ctx, "not have msgs", nil, "conversationID", seq.ConversationID, "seq", seq)
				continue
			}
			resp.Msgs[seq.ConversationID] = &sdkws.PullMsgs{Msgs: msgs, IsEnd: isEnd}
		} else {
			var seqs []int64
			for i := seq.Begin; i <= seq.End; i++ {
				seqs = append(seqs, i)
			}
			minSeq, maxSeq, notificationMsgs, err := m.MsgDatabase.GetMsgBySeqs(ctx, req.UserID, seq.ConversationID, seqs)
			if err != nil {
				log.ZWarn(ctx, "GetMsgBySeqs error", err, "conversationID", seq.ConversationID, "seq", seq)

				continue
			}
			var isEnd bool
			switch req.Order {
			case sdkws.PullOrder_PullOrderAsc:
				isEnd = maxSeq <= seq.End
			case sdkws.PullOrder_PullOrderDesc:
				isEnd = seq.Begin <= minSeq
			}
			if len(notificationMsgs) == 0 {
				log.ZWarn(ctx, "not have notificationMsgs", nil, "conversationID", seq.ConversationID, "seq", seq)

				continue
			}
			resp.NotificationMsgs[seq.ConversationID] = &sdkws.PullMsgs{Msgs: notificationMsgs, IsEnd: isEnd}
		}
	}
	return resp, nil
}

func (m *msgServer) GetSeqMessage(ctx context.Context, req *msg.GetSeqMessageReq) (*msg.GetSeqMessageResp, error) {
	resp := &msg.GetSeqMessageResp{
		Msgs:             make(map[string]*sdkws.PullMsgs),
		NotificationMsgs: make(map[string]*sdkws.PullMsgs),
	}
	for _, conv := range req.Conversations {
		isEnd, endSeq, msgs, err := m.MsgDatabase.GetMessagesBySeqWithBounds(ctx, req.UserID, conv.ConversationID, conv.Seqs, req.GetOrder())
		if err != nil {
			return nil, err
		}
		var pullMsgs *sdkws.PullMsgs
		pullMsgs.IsEnd = isEnd
		pullMsgs.EndSeq = endSeq
		if ok := false; conversationutil.IsNotificationConversationID(conv.ConversationID) {
			pullMsgs, ok = resp.NotificationMsgs[conv.ConversationID]
			if !ok {
				pullMsgs = &sdkws.PullMsgs{}
				resp.NotificationMsgs[conv.ConversationID] = pullMsgs
			}
		} else {
			pullMsgs, ok = resp.Msgs[conv.ConversationID]
			if !ok {
				pullMsgs = &sdkws.PullMsgs{}
				resp.Msgs[conv.ConversationID] = pullMsgs
			}
		}
		pullMsgs.Msgs = append(pullMsgs.Msgs, msgs...)
		pullMsgs.IsEnd = isEnd
		pullMsgs.EndSeq = endSeq
	}
	return resp, nil
}

func (m *msgServer) GetMaxSeq(ctx context.Context, req *sdkws.GetMaxSeqReq) (*sdkws.GetMaxSeqResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	conversationIDs, err := m.ConversationLocalCache.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, conversationID := range conversationIDs {
		conversationIDs = append(conversationIDs, conversationutil.GetNotificationConversationIDByConversationID(conversationID))
	}
	conversationIDs = append(conversationIDs, conversationutil.GetSelfNotificationConversationID(req.UserID))
	log.ZDebug(ctx, "GetMaxSeq", "conversationIDs", conversationIDs)
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		log.ZWarn(ctx, "GetMaxSeqs error", err, "conversationIDs", conversationIDs, "maxSeqs", maxSeqs)
		return nil, err
	}
	// avoid pulling messages from sessions with a large number of max seq values of 0
	for conversationID, seq := range maxSeqs {
		if seq == 0 {
			delete(maxSeqs, conversationID)
		}
	}
	resp := new(sdkws.GetMaxSeqResp)
	resp.MaxSeqs = maxSeqs
	return resp, nil
}

func (m *msgServer) SearchMessage(ctx context.Context, req *msg.SearchMessageReq) (resp *msg.SearchMessageResp, err error) {
	// var chatLogs []*sdkws.MsgData
	var chatLogs []*msg.SearchedMsgData
	var total int64
	resp = &msg.SearchMessageResp{}
	if total, chatLogs, err = m.MsgDatabase.SearchMessage(ctx, req); err != nil {
		return nil, err
	}

	var (
		sendIDs  []string
		recvIDs  []string
		groupIDs []string
		sendMap  = make(map[string]string)
		recvMap  = make(map[string]string)
		groupMap = make(map[string]*sdkws.GroupInfo)
	)

	for _, chatLog := range chatLogs {
		if chatLog.MsgData.SenderNickname == "" {
			sendIDs = append(sendIDs, chatLog.MsgData.SendID)
		}
		switch chatLog.MsgData.SessionType {
		case constant.SingleChatType, constant.NotificationChatType:
			recvIDs = append(recvIDs, chatLog.MsgData.RecvID)
		case constant.WriteGroupChatType, constant.ReadGroupChatType:
			groupIDs = append(groupIDs, chatLog.MsgData.GroupID)
		}
	}

	// Retrieve sender and receiver information
	if len(sendIDs) != 0 {
		sendInfos, err := m.UserLocalCache.GetUsersInfo(ctx, sendIDs)
		if err != nil {
			return nil, err
		}
		for _, sendInfo := range sendInfos {
			sendMap[sendInfo.UserID] = sendInfo.Nickname
		}
	}

	if len(recvIDs) != 0 {
		recvInfos, err := m.UserLocalCache.GetUsersInfo(ctx, recvIDs)
		if err != nil {
			return nil, err
		}
		for _, recvInfo := range recvInfos {
			recvMap[recvInfo.UserID] = recvInfo.Nickname
		}
	}

	// Retrieve group information including member counts
	if len(groupIDs) != 0 {
		groupInfos, err := m.GroupLocalCache.GetGroupInfos(ctx, groupIDs)
		if err != nil {
			return nil, err
		}
		for _, groupInfo := range groupInfos {
			groupMap[groupInfo.GroupID] = groupInfo
			// Get actual member count
			memberIDs, err := m.GroupLocalCache.GetGroupMemberIDs(ctx, groupInfo.GroupID)
			if err == nil {
				groupInfo.MemberCount = uint32(len(memberIDs)) // Update the member count with actual number
			}
		}
	}

	// Construct response with updated information
	for _, chatLog := range chatLogs {
		pbchatLog := &msg.ChatLog{}
		datautil.CopyStructFields(pbchatLog, chatLog.MsgData)
		pbchatLog.SendTime = chatLog.MsgData.SendTime
		pbchatLog.CreateTime = chatLog.MsgData.CreateTime
		if chatLog.MsgData.SenderNickname == "" {
			pbchatLog.SenderNickname = sendMap[chatLog.MsgData.SendID]
		}
		switch chatLog.MsgData.SessionType {
		case constant.SingleChatType, constant.NotificationChatType:
			pbchatLog.RecvNickname = recvMap[chatLog.MsgData.RecvID]
		case constant.ReadGroupChatType:
			groupInfo := groupMap[chatLog.MsgData.GroupID]
			pbchatLog.SenderFaceURL = groupInfo.FaceURL
			pbchatLog.GroupMemberCount = groupInfo.MemberCount // Reflects actual member count
			pbchatLog.RecvID = groupInfo.GroupID
			pbchatLog.GroupName = groupInfo.GroupName
			pbchatLog.GroupOwner = groupInfo.OwnerUserID
			pbchatLog.GroupType = groupInfo.GroupType
		}
		searchChatLog := &msg.SearchChatLog{ChatLog: pbchatLog, IsRevoked: chatLog.IsRevoked}

		resp.ChatLogs = append(resp.ChatLogs, searchChatLog)
	}
	resp.ChatLogsNum = int32(total)
	return resp, nil
}

func (m *msgServer) GetServerTime(ctx context.Context, _ *msg.GetServerTimeReq) (*msg.GetServerTimeResp, error) {
	return &msg.GetServerTimeResp{ServerTime: timeutil.GetCurrentTimestampByMill()}, nil
}
