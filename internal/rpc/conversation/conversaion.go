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

package conversation

import (
	"context"
	"errors"
	"sort"

	"github.com/OpenIMSDK/protocol/sdkws"

	"github.com/OpenIMSDK/tools/tx"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/protocol/constant"
	pbconversation "github.com/OpenIMSDK/protocol/conversation"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	tablerelation "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient/notification"
)

type conversationServer struct {
	msgRpcClient                   *rpcclient.MessageRpcClient
	user                           *rpcclient.UserRpcClient
	groupRpcClient                 *rpcclient.GroupRpcClient
	conversationDatabase           controller.ConversationDatabase
	conversationNotificationSender *notification.ConversationNotificationSender
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}
	conversationDB, err := mgo.NewConversationMongo(mongo.GetDatabase())
	if err != nil {
		return err
	}
	groupRpcClient := rpcclient.NewGroupRpcClient(client)
	msgRpcClient := rpcclient.NewMessageRpcClient(client)
	userRpcClient := rpcclient.NewUserRpcClient(client)
	pbconversation.RegisterConversationServer(server, &conversationServer{
		msgRpcClient:                   &msgRpcClient,
		user:                           &userRpcClient,
		conversationNotificationSender: notification.NewConversationNotificationSender(&msgRpcClient),
		groupRpcClient:                 &groupRpcClient,
		conversationDatabase:           controller.NewConversationDatabase(conversationDB, cache.NewConversationRedis(rdb, cache.GetDefaultOpt(), conversationDB), tx.NewMongo(mongo.GetClient())),
	})
	return nil
}

func (c *conversationServer) GetConversation(ctx context.Context, req *pbconversation.GetConversationReq) (*pbconversation.GetConversationResp, error) {
	conversations, err := c.conversationDatabase.FindConversations(ctx, req.OwnerUserID, []string{req.ConversationID})
	if err != nil {
		return nil, err
	}
	if len(conversations) < 1 {
		return nil, errs.ErrRecordNotFound.Wrap("conversation not found")
	}
	resp := &pbconversation.GetConversationResp{Conversation: &pbconversation.Conversation{}}
	resp.Conversation = convert.ConversationDB2Pb(conversations[0])
	return resp, nil
}

func (m *conversationServer) GetSortedConversationList(ctx context.Context, req *pbconversation.GetSortedConversationListReq) (resp *pbconversation.GetSortedConversationListResp, err error) {
	log.ZDebug(ctx, "GetSortedConversationList", "seqs", req, "userID", req.UserID)
	var conversationIDs []string
	if len(req.ConversationIDs) == 0 {
		conversationIDs, err = m.conversationDatabase.GetConversationIDs(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		conversationIDs = req.ConversationIDs
	}

	conversations, err := m.conversationDatabase.FindConversations(ctx, req.UserID, conversationIDs)
	if err != nil {
		return nil, err
	}
	if len(conversations) == 0 {
		return nil, errs.ErrRecordNotFound.Wrap()
	}

	maxSeqs, err := m.msgRpcClient.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}

	chatLogs, err := m.msgRpcClient.GetMsgByConversationIDs(ctx, conversationIDs, maxSeqs)
	if err != nil {
		return nil, err
	}

	conversationMsg, err := m.getConversationInfo(ctx, chatLogs, req.UserID)
	if err != nil {
		return nil, err
	}

	hasReadSeqs, err := m.msgRpcClient.GetHasReadSeqs(ctx, req.UserID, conversationIDs)
	if err != nil {
		return nil, err
	}

	var unreadTotal int64
	conversation_unreadCount := make(map[string]int64)
	for conversationID, maxSeq := range maxSeqs {
		unreadCount := maxSeq - hasReadSeqs[conversationID]
		conversation_unreadCount[conversationID] = unreadCount
		unreadTotal += unreadCount
	}

	conversation_isPinTime := make(map[int64]string)
	conversation_notPinTime := make(map[int64]string)
	for _, v := range conversations {
		conversationID := v.ConversationID
		time := conversationMsg[conversationID].MsgInfo.LatestMsgRecvTime
		conversationMsg[conversationID].RecvMsgOpt = v.RecvMsgOpt
		if v.IsPinned {
			conversationMsg[conversationID].IsPinned = v.IsPinned
			conversation_isPinTime[time] = conversationID
			continue
		}
		conversation_notPinTime[time] = conversationID
	}
	resp = &pbconversation.GetSortedConversationListResp{
		ConversationTotal: int64(len(chatLogs)),
		ConversationElems: []*pbconversation.ConversationElem{},
		UnreadTotal:       unreadTotal,
	}

	m.conversationSort(conversation_isPinTime, resp, conversation_unreadCount, conversationMsg)
	m.conversationSort(conversation_notPinTime, resp, conversation_unreadCount, conversationMsg)

	resp.ConversationElems = utils.Paginate(resp.ConversationElems, int(req.Pagination.GetPageNumber()), int(req.Pagination.GetShowNumber()))
	return resp, nil
}

func (c *conversationServer) GetAllConversations(ctx context.Context, req *pbconversation.GetAllConversationsReq) (*pbconversation.GetAllConversationsResp, error) {
	conversations, err := c.conversationDatabase.GetUserAllConversation(ctx, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	resp := &pbconversation.GetAllConversationsResp{Conversations: []*pbconversation.Conversation{}}
	resp.Conversations = convert.ConversationsDB2Pb(conversations)
	return resp, nil
}

func (c *conversationServer) GetConversations(ctx context.Context, req *pbconversation.GetConversationsReq) (*pbconversation.GetConversationsResp, error) {
	conversations, err := c.conversationDatabase.FindConversations(ctx, req.OwnerUserID, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	resp := &pbconversation.GetConversationsResp{Conversations: []*pbconversation.Conversation{}}
	resp.Conversations = convert.ConversationsDB2Pb(conversations)
	return resp, nil
}

func (c *conversationServer) SetConversation(ctx context.Context, req *pbconversation.SetConversationReq) (*pbconversation.SetConversationResp, error) {
	var conversation tablerelation.ConversationModel
	if err := utils.CopyStructFields(&conversation, req.Conversation); err != nil {
		return nil, err
	}
	err := c.conversationDatabase.SetUserConversations(ctx, req.Conversation.OwnerUserID, []*tablerelation.ConversationModel{&conversation})
	if err != nil {
		return nil, err
	}
	_ = c.conversationNotificationSender.ConversationChangeNotification(ctx, req.Conversation.OwnerUserID, []string{req.Conversation.ConversationID})
	resp := &pbconversation.SetConversationResp{}
	return resp, nil
}

// nolint
func (c *conversationServer) SetConversations(ctx context.Context,
	req *pbconversation.SetConversationsReq,
) (*pbconversation.SetConversationsResp, error) {
	if req.Conversation == nil {
		return nil, errs.ErrArgs.Wrap("conversation must not be nil")
	}
	if req.Conversation.ConversationType == constant.GroupChatType {
		groupInfo, err := c.groupRpcClient.GetGroupInfo(ctx, req.Conversation.GroupID)
		if err != nil {
			return nil, err
		}
		if groupInfo.Status == constant.GroupStatusDismissed {
			return nil, errs.ErrDismissedAlready.Wrap("group dismissed")
		}
	}
	var unequal int
	var conv tablerelation.ConversationModel
	if len(req.UserIDs) == 1 {
		cs, err := c.conversationDatabase.FindConversations(ctx, req.UserIDs[0], []string{req.Conversation.ConversationID})
		if err != nil {
			return nil, err
		}
		if len(cs) == 0 {
			return nil, errs.ErrRecordNotFound.Wrap("conversation not found")
		}
		conv = *cs[0]
	}
	var conversation tablerelation.ConversationModel
	conversation.ConversationID = req.Conversation.ConversationID
	conversation.ConversationType = req.Conversation.ConversationType
	conversation.UserID = req.Conversation.UserID
	conversation.GroupID = req.Conversation.GroupID
	m := make(map[string]any)
	if req.Conversation.RecvMsgOpt != nil {
		m["recv_msg_opt"] = req.Conversation.RecvMsgOpt.Value
		if req.Conversation.RecvMsgOpt.Value != conv.RecvMsgOpt {
			unequal++
		}
	}
	if req.Conversation.AttachedInfo != nil {
		m["attached_info"] = req.Conversation.AttachedInfo.Value
		if req.Conversation.AttachedInfo.Value != conv.AttachedInfo {
			unequal++
		}
	}
	if req.Conversation.Ex != nil {
		m["ex"] = req.Conversation.Ex.Value
		if req.Conversation.Ex.Value != conv.Ex {
			unequal++
		}
	}
	if req.Conversation.IsPinned != nil {
		m["is_pinned"] = req.Conversation.IsPinned.Value
		if req.Conversation.IsPinned.Value != conv.IsPinned {
			unequal++
		}
	}
	if req.Conversation.GroupAtType != nil {
		m["group_at_type"] = req.Conversation.GroupAtType.Value
		if req.Conversation.GroupAtType.Value != conv.GroupAtType {
			unequal++
		}
	}
	if req.Conversation.MsgDestructTime != nil {
		m["msg_destruct_time"] = req.Conversation.MsgDestructTime.Value
		if req.Conversation.MsgDestructTime.Value != conv.MsgDestructTime {
			unequal++
		}
	}
	if req.Conversation.IsMsgDestruct != nil {
		m["is_msg_destruct"] = req.Conversation.IsMsgDestruct.Value
		if req.Conversation.IsMsgDestruct.Value != conv.IsMsgDestruct {
			unequal++
		}
	}
	if req.Conversation.IsPrivateChat != nil && req.Conversation.ConversationType != constant.SuperGroupChatType {
		var conversations []*tablerelation.ConversationModel
		for _, ownerUserID := range req.UserIDs {
			conversation2 := conversation
			conversation2.OwnerUserID = ownerUserID
			conversation2.IsPrivateChat = req.Conversation.IsPrivateChat.Value
			conversations = append(conversations, &conversation2)
		}
		if err := c.conversationDatabase.SyncPeerUserPrivateConversationTx(ctx, conversations); err != nil {
			return nil, err
		}
		for _, userID := range req.UserIDs {
			err := c.conversationNotificationSender.ConversationSetPrivateNotification(ctx, userID, req.Conversation.UserID,
				req.Conversation.IsPrivateChat.Value, req.Conversation.ConversationID)
			if err != nil {
				log.ZWarn(ctx, "send conversation set private notification failed", err,
					"userID", userID, "conversationID", req.Conversation.ConversationID)

				continue
			}
		}
	}
	if req.Conversation.BurnDuration != nil {
		m["burn_duration"] = req.Conversation.BurnDuration.Value
		if req.Conversation.BurnDuration.Value != conv.BurnDuration {
			unequal++
		}
	}
	if err := c.conversationDatabase.SetUsersConversationFiledTx(ctx, req.UserIDs, &conversation, m); err != nil {
		return nil, err
	}
	if unequal > 0 {
		for _, v := range req.UserIDs {
			c.conversationNotificationSender.ConversationChangeNotification(ctx, v, []string{req.Conversation.ConversationID})
		}
	}
	return &pbconversation.SetConversationsResp{}, nil
}

// 获取超级大群开启免打扰的用户ID.
func (c *conversationServer) GetRecvMsgNotNotifyUserIDs(ctx context.Context, req *pbconversation.GetRecvMsgNotNotifyUserIDsReq) (*pbconversation.GetRecvMsgNotNotifyUserIDsResp, error) {
	//userIDs, err := c.conversationDatabase.FindRecvMsgNotNotifyUserIDs(ctx, req.GroupID)
	//if err != nil {
	//	return nil, err
	//}
	//return &pbconversation.GetRecvMsgNotNotifyUserIDsResp{UserIDs: userIDs}, nil
	return nil, errors.New("deprecated")
}

// create conversation without notification for msg redis transfer.
func (c *conversationServer) CreateSingleChatConversations(ctx context.Context,
	req *pbconversation.CreateSingleChatConversationsReq,
) (*pbconversation.CreateSingleChatConversationsResp, error) {
	switch req.ConversationType {
	case constant.SingleChatType:
		var conversation tablerelation.ConversationModel
		conversation.ConversationID = req.ConversationID
		conversation.ConversationType = req.ConversationType
		conversation.OwnerUserID = req.SendID
		conversation.UserID = req.RecvID
		err := c.conversationDatabase.CreateConversation(ctx, []*tablerelation.ConversationModel{&conversation})
		if err != nil {
			log.ZWarn(ctx, "create conversation failed", err, "conversation", conversation)
		}

		conversation2 := conversation
		conversation2.OwnerUserID = req.RecvID
		conversation2.UserID = req.SendID
		err = c.conversationDatabase.CreateConversation(ctx, []*tablerelation.ConversationModel{&conversation2})
		if err != nil {
			log.ZWarn(ctx, "create conversation failed", err, "conversation2", conversation)
		}
	case constant.NotificationChatType:
		var conversation tablerelation.ConversationModel
		conversation.ConversationID = req.ConversationID
		conversation.ConversationType = req.ConversationType
		conversation.OwnerUserID = req.RecvID
		conversation.UserID = req.SendID
		err := c.conversationDatabase.CreateConversation(ctx, []*tablerelation.ConversationModel{&conversation})
		if err != nil {
			log.ZWarn(ctx, "create conversation failed", err, "conversation2", conversation)
		}
	}

	return &pbconversation.CreateSingleChatConversationsResp{}, nil
}

func (c *conversationServer) CreateGroupChatConversations(ctx context.Context, req *pbconversation.CreateGroupChatConversationsReq) (*pbconversation.CreateGroupChatConversationsResp, error) {
	err := c.conversationDatabase.CreateGroupChatConversation(ctx, req.GroupID, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &pbconversation.CreateGroupChatConversationsResp{}, nil
}

func (c *conversationServer) SetConversationMaxSeq(ctx context.Context, req *pbconversation.SetConversationMaxSeqReq) (*pbconversation.SetConversationMaxSeqResp, error) {
	if err := c.conversationDatabase.UpdateUsersConversationFiled(ctx, req.OwnerUserID, req.ConversationID,
		map[string]any{"max_seq": req.MaxSeq}); err != nil {
		return nil, err
	}
	return &pbconversation.SetConversationMaxSeqResp{}, nil
}

func (c *conversationServer) GetConversationIDs(ctx context.Context, req *pbconversation.GetConversationIDsReq) (*pbconversation.GetConversationIDsResp, error) {
	conversationIDs, err := c.conversationDatabase.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetConversationIDsResp{ConversationIDs: conversationIDs}, nil
}

func (c *conversationServer) GetUserConversationIDsHash(ctx context.Context, req *pbconversation.GetUserConversationIDsHashReq) (*pbconversation.GetUserConversationIDsHashResp, error) {
	hash, err := c.conversationDatabase.GetUserConversationIDsHash(ctx, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetUserConversationIDsHashResp{Hash: hash}, nil
}

func (c *conversationServer) GetConversationsByConversationID(
	ctx context.Context,
	req *pbconversation.GetConversationsByConversationIDReq,
) (*pbconversation.GetConversationsByConversationIDResp, error) {
	conversations, err := c.conversationDatabase.GetConversationsByConversationID(ctx, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetConversationsByConversationIDResp{Conversations: convert.ConversationsDB2Pb(conversations)}, nil
}

func (c *conversationServer) GetConversationOfflinePushUserIDs(
	ctx context.Context,
	req *pbconversation.GetConversationOfflinePushUserIDsReq,
) (*pbconversation.GetConversationOfflinePushUserIDsResp, error) {
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.Wrap("conversationID is empty")
	}
	if len(req.UserIDs) == 0 {
		return &pbconversation.GetConversationOfflinePushUserIDsResp{}, nil
	}
	userIDs, err := c.conversationDatabase.GetConversationNotReceiveMessageUserIDs(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	if len(userIDs) == 0 {
		return &pbconversation.GetConversationOfflinePushUserIDsResp{UserIDs: req.UserIDs}, nil
	}
	userIDSet := make(map[string]struct{})
	for _, userID := range req.UserIDs {
		userIDSet[userID] = struct{}{}
	}
	for _, userID := range userIDs {
		delete(userIDSet, userID)
	}
	return &pbconversation.GetConversationOfflinePushUserIDsResp{UserIDs: utils.Keys(userIDSet)}, nil
}

func (c *conversationServer) conversationSort(
	conversations map[int64]string,
	resp *pbconversation.GetSortedConversationListResp,
	conversation_unreadCount map[string]int64,
	conversationMsg map[string]*pbconversation.ConversationElem,
) {
	keys := []int64{}
	for key := range conversations {
		keys = append(keys, key)
	}

	sort.Slice(keys[:], func(i, j int) bool {
		return keys[i] > keys[j]
	})
	index := 0

	cons := make([]*pbconversation.ConversationElem, len(conversations))
	for _, v := range keys {
		conversationID := conversations[v]
		conversationElem := conversationMsg[conversationID]
		conversationElem.UnreadCount = conversation_unreadCount[conversationID]
		cons[index] = conversationElem
		index++
	}
	resp.ConversationElems = append(resp.ConversationElems, cons...)
}

func (c *conversationServer) getConversationInfo(
	ctx context.Context,
	chatLogs map[string]*sdkws.MsgData,
	userID string) (map[string]*pbconversation.ConversationElem, error) {
	var (
		sendIDs         []string
		groupIDs        []string
		sendMap         = make(map[string]*sdkws.UserInfo)
		groupMap        = make(map[string]*sdkws.GroupInfo)
		conversationMsg = make(map[string]*pbconversation.ConversationElem)
	)
	for _, chatLog := range chatLogs {
		switch chatLog.SessionType {
		case constant.SingleChatType:
			if chatLog.SendID == userID {
				sendIDs = append(sendIDs, chatLog.RecvID)
			}
			sendIDs = append(sendIDs, chatLog.SendID)
		case constant.GroupChatType, constant.SuperGroupChatType:
			groupIDs = append(groupIDs, chatLog.GroupID)
			sendIDs = append(sendIDs, chatLog.SendID)
		}
	}
	if len(sendIDs) != 0 {
		sendInfos, err := c.user.GetUsersInfo(ctx, sendIDs)
		if err != nil {
			return nil, err
		}
		for _, sendInfo := range sendInfos {
			sendMap[sendInfo.UserID] = sendInfo
		}
	}
	if len(groupIDs) != 0 {
		groupInfos, err := c.groupRpcClient.GetGroupInfos(ctx, groupIDs, false)
		if err != nil {
			return nil, err
		}
		for _, groupInfo := range groupInfos {
			groupMap[groupInfo.GroupID] = groupInfo
		}
	}
	for conversationID, chatLog := range chatLogs {
		pbchatLog := &pbconversation.ConversationElem{}
		msgInfo := &pbconversation.MsgInfo{}
		if err := utils.CopyStructFields(msgInfo, chatLog); err != nil {
			return nil, err
		}
		switch chatLog.SessionType {
		case constant.SingleChatType:
			if chatLog.SendID == userID {
				msgInfo.FaceURL = sendMap[chatLog.RecvID].FaceURL
				msgInfo.SenderName = sendMap[chatLog.RecvID].Nickname
				break
			}
			msgInfo.FaceURL = sendMap[chatLog.SendID].FaceURL
			msgInfo.SenderName = sendMap[chatLog.SendID].Nickname
		case constant.GroupChatType, constant.SuperGroupChatType:
			msgInfo.GroupName = groupMap[chatLog.GroupID].GroupName
			msgInfo.GroupFaceURL = groupMap[chatLog.GroupID].FaceURL
			msgInfo.GroupMemberCount = groupMap[chatLog.GroupID].MemberCount
			msgInfo.GroupID = chatLog.GroupID
			msgInfo.GroupType = groupMap[chatLog.GroupID].GroupType
			msgInfo.SenderName = sendMap[chatLog.SendID].Nickname
		}
		pbchatLog.ConversationID = conversationID
		msgInfo.LatestMsgRecvTime = chatLog.SendTime
		pbchatLog.MsgInfo = msgInfo
		conversationMsg[conversationID] = pbchatLog
	}
	return conversationMsg, nil
}
