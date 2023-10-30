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

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msg"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"
)

func (m *msgServer) PullMessageBySeqs(
	ctx context.Context,
	req *sdkws.PullMessageBySeqsReq,
) (*sdkws.PullMessageBySeqsResp, error) {
	resp := &sdkws.PullMessageBySeqsResp{}
	resp.Msgs = make(map[string]*sdkws.PullMsgs)
	resp.NotificationMsgs = make(map[string]*sdkws.PullMsgs)
	for _, seq := range req.SeqRanges {
		if !msgprocessor.IsNotification(seq.ConversationID) {
			conversation, err := m.Conversation.GetConversation(ctx, req.UserID, seq.ConversationID)
			if err != nil {
				log.ZError(ctx, "GetConversation error", err, "conversationID", seq.ConversationID)
				continue
			}
			minSeq, maxSeq, msgs, err := m.MsgDatabase.GetMsgBySeqsRange(
				ctx,
				req.UserID,
				seq.ConversationID,
				seq.Begin,
				seq.End,
				seq.Num,
				conversation.MaxSeq,
			)
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

func (m *msgServer) GetMaxSeq(ctx context.Context, req *sdkws.GetMaxSeqReq) (*sdkws.GetMaxSeqResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	conversationIDs, err := m.ConversationLocalCache.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, conversationID := range conversationIDs {
		conversationIDs = append(conversationIDs, utils.GetNotificationConversationIDByConversationID(conversationID))
	}
	conversationIDs = append(conversationIDs, utils.GetSelfNotificationConversationID(req.UserID))
	log.ZDebug(ctx, "GetMaxSeq", "conversationIDs", conversationIDs)
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		log.ZWarn(ctx, "GetMaxSeqs error", err, "conversationIDs", conversationIDs, "maxSeqs", maxSeqs)
		return nil, err
	}
	resp := new(sdkws.GetMaxSeqResp)
	resp.MaxSeqs = maxSeqs
	return resp, nil
}

func (m *msgServer) SearchMessage(ctx context.Context, req *msg.SearchMessageReq) (resp *msg.SearchMessageResp, err error) {
	var chatLogs []*sdkws.MsgData
	var total int32
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
		if chatLog.SenderNickname == "" {
			sendIDs = append(sendIDs, chatLog.SendID)
		}
		switch chatLog.SessionType {
		case constant.SingleChatType:
			recvIDs = append(recvIDs, chatLog.RecvID)
		case constant.GroupChatType, constant.SuperGroupChatType:
			groupIDs = append(groupIDs, chatLog.GroupID)
		}
	}
	if len(sendIDs) != 0 {
		sendInfos, err := m.User.GetUsersInfo(ctx, sendIDs)
		if err != nil {
			return nil, err
		}
		for _, sendInfo := range sendInfos {
			sendMap[sendInfo.UserID] = sendInfo.Nickname
		}
	}
	if len(recvIDs) != 0 {
		recvInfos, err := m.User.GetUsersInfo(ctx, recvIDs)
		if err != nil {
			return nil, err
		}
		for _, recvInfo := range recvInfos {
			recvMap[recvInfo.UserID] = recvInfo.Nickname
		}
	}
	if len(groupIDs) != 0 {
		groupInfos, err := m.Group.GetGroupInfos(ctx, groupIDs, true)
		if err != nil {
			return nil, err
		}
		for _, groupInfo := range groupInfos {
			groupMap[groupInfo.GroupID] = groupInfo
		}
	}
	for _, chatLog := range chatLogs {
		pbchatLog := &msg.ChatLog{}
		utils.CopyStructFields(pbchatLog, chatLog)
		pbchatLog.SendTime = chatLog.SendTime
		pbchatLog.CreateTime = chatLog.CreateTime
		if chatLog.SenderNickname == "" {
			pbchatLog.SenderNickname = sendMap[chatLog.SendID]
		}
		switch chatLog.SessionType {
		case constant.SingleChatType:
			pbchatLog.RecvNickname = recvMap[chatLog.RecvID]

		case constant.GroupChatType, constant.SuperGroupChatType:
			pbchatLog.SenderFaceURL = groupMap[chatLog.GroupID].FaceURL
			pbchatLog.GroupMemberCount = groupMap[chatLog.GroupID].MemberCount
			pbchatLog.RecvID = groupMap[chatLog.GroupID].GroupID
			pbchatLog.GroupName = groupMap[chatLog.GroupID].GroupName
			pbchatLog.GroupOwner = groupMap[chatLog.GroupID].OwnerUserID
			pbchatLog.GroupType = groupMap[chatLog.GroupID].GroupType
		}
		resp.ChatLogs = append(resp.ChatLogs, pbchatLog)
	}
	resp.ChatLogsNum = total
	return resp, nil
}

func (m *msgServer) GetServerTime(ctx context.Context, _ *msg.GetServerTimeReq) (*msg.GetServerTimeResp, error) {
	return &msg.GetServerTimeResp{ServerTime: utils.GetCurrentTimestampByMill()}, nil
}
