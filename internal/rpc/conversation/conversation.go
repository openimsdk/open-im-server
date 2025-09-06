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
	"sort"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"

	"google.golang.org/grpc"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	dbModel "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/constant"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

type conversationServer struct {
	pbconversation.UnimplementedConversationServer
	conversationDatabase controller.ConversationDatabase

	conversationNotificationSender *ConversationNotificationSender
	config                         *Config

	webhookClient *webhook.Client
	userClient    *rpcli.UserClient
	msgClient     *rpcli.MsgClient
	groupClient   *rpcli.GroupClient
}

type Config struct {
	RpcConfig          config.Conversation
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	NotificationConfig config.Notification
	Share              config.Share
	WebhooksConfig     config.Webhooks
	LocalCacheConfig   config.LocalCache
	Discovery          config.Discovery
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server grpc.ServiceRegistrar) error {
	dbb := dbbuild.NewBuilder(&config.MongodbConfig, &config.RedisConfig)
	mgocli, err := dbb.Mongo(ctx)
	if err != nil {
		return err
	}
	rdb, err := dbb.Redis(ctx)
	if err != nil {
		return err
	}
	conversationDB, err := mgo.NewConversationMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	userConn, err := client.GetConn(ctx, config.Discovery.RpcService.User)
	if err != nil {
		return err
	}
	groupConn, err := client.GetConn(ctx, config.Discovery.RpcService.Group)
	if err != nil {
		return err
	}
	msgConn, err := client.GetConn(ctx, config.Discovery.RpcService.Msg)
	if err != nil {
		return err
	}

	msgClient := rpcli.NewMsgClient(msgConn)

	cs := conversationServer{
		config:        config,
		webhookClient: webhook.NewWebhookClient(config.WebhooksConfig.URL),
		userClient:    rpcli.NewUserClient(userConn),
		groupClient:   rpcli.NewGroupClient(groupConn),
		msgClient:     msgClient,
	}

	cs.conversationNotificationSender = NewConversationNotificationSender(&config.NotificationConfig, msgClient)
	cs.conversationDatabase = controller.NewConversationDatabase(
		conversationDB,
		redis.NewConversationRedis(rdb, &config.LocalCacheConfig, conversationDB),
		mgocli.GetTx())

	localcache.InitLocalCache(&config.LocalCacheConfig)
	pbconversation.RegisterConversationServer(server, &cs)
	return nil
}

func (c *conversationServer) GetConversation(ctx context.Context, req *pbconversation.GetConversationReq) (*pbconversation.GetConversationResp, error) {
	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	conversations, err := c.conversationDatabase.FindConversations(ctx, req.OwnerUserID, []string{req.ConversationID})
	if err != nil {
		return nil, err
	}
	if len(conversations) < 1 {
		return nil, errs.ErrRecordNotFound.WrapMsg("conversation not found")
	}
	resp := &pbconversation.GetConversationResp{Conversation: &pbconversation.Conversation{}}
	resp.Conversation = convert.ConversationDB2Pb(conversations[0])
	return resp, nil
}

// Deprecated: Use `GetConversations` instead.
func (c *conversationServer) GetSortedConversationList(ctx context.Context, req *pbconversation.GetSortedConversationListReq) (resp *pbconversation.GetSortedConversationListResp, err error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	var conversationIDs []string
	if len(req.ConversationIDs) == 0 {
		conversationIDs, err = c.conversationDatabase.GetConversationIDs(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		conversationIDs = req.ConversationIDs
	}

	conversations, err := c.conversationDatabase.FindConversations(ctx, req.UserID, conversationIDs)
	if err != nil {
		return nil, err
	}
	if len(conversations) == 0 {
		return nil, errs.ErrRecordNotFound.Wrap()
	}
	maxSeqs, err := c.msgClient.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}

	chatLogs, err := c.msgClient.GetMsgByConversationIDs(ctx, conversationIDs, maxSeqs)
	if err != nil {
		return nil, err
	}

	conversationMsg, err := c.getConversationInfo(ctx, chatLogs, req.UserID)
	if err != nil {
		return nil, err
	}

	hasReadSeqs, err := c.msgClient.GetHasReadSeqs(ctx, conversationIDs, req.UserID)
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
		var time int64
		if _, ok := conversationMsg[conversationID]; ok {
			time = conversationMsg[conversationID].MsgInfo.LatestMsgRecvTime
		} else {
			conversationMsg[conversationID] = &pbconversation.ConversationElem{
				ConversationID: conversationID,
				IsPinned:       v.IsPinned,
				MsgInfo:        nil,
			}
			time = v.CreateTime.UnixMilli()
		}

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

	c.conversationSort(conversation_isPinTime, resp, conversation_unreadCount, conversationMsg)
	c.conversationSort(conversation_notPinTime, resp, conversation_unreadCount, conversationMsg)

	resp.ConversationElems = datautil.Paginate(resp.ConversationElems, int(req.Pagination.GetPageNumber()), int(req.Pagination.GetShowNumber()))
	return resp, nil
}

func (c *conversationServer) GetAllConversations(ctx context.Context, req *pbconversation.GetAllConversationsReq) (*pbconversation.GetAllConversationsResp, error) {
	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	conversations, err := c.conversationDatabase.GetUserAllConversation(ctx, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	resp := &pbconversation.GetAllConversationsResp{Conversations: []*pbconversation.Conversation{}}
	resp.Conversations = convert.ConversationsDB2Pb(conversations)
	return resp, nil
}

func (c *conversationServer) GetConversations(ctx context.Context, req *pbconversation.GetConversationsReq) (*pbconversation.GetConversationsResp, error) {
	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	conversations, err := c.getConversations(ctx, req.OwnerUserID, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetConversationsResp{
		Conversations: conversations,
	}, nil
}

func (c *conversationServer) getConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*pbconversation.Conversation, error) {
	conversations, err := c.conversationDatabase.FindConversations(ctx, ownerUserID, conversationIDs)
	if err != nil {
		return nil, err
	}
	resp := &pbconversation.GetConversationsResp{Conversations: []*pbconversation.Conversation{}}
	resp.Conversations = convert.ConversationsDB2Pb(conversations)
	return convert.ConversationsDB2Pb(conversations), nil
}

// Deprecated
func (c *conversationServer) SetConversation(ctx context.Context, req *pbconversation.SetConversationReq) (*pbconversation.SetConversationResp, error) {
	if err := authverify.CheckAccess(ctx, req.GetConversation().GetUserID()); err != nil {
		return nil, err
	}
	var conversation dbModel.Conversation
	conversation.CreateTime = time.Now()

	if err := datautil.CopyStructFields(&conversation, req.Conversation); err != nil {
		return nil, err
	}
	err := c.conversationDatabase.SetUserConversations(ctx, req.Conversation.OwnerUserID, []*dbModel.Conversation{&conversation})
	if err != nil {
		return nil, err
	}
	c.conversationNotificationSender.ConversationChangeNotification(ctx, req.Conversation.OwnerUserID, []string{req.Conversation.ConversationID})
	resp := &pbconversation.SetConversationResp{}
	return resp, nil
}

func (c *conversationServer) SetConversations(ctx context.Context, req *pbconversation.SetConversationsReq) (*pbconversation.SetConversationsResp, error) {
	for _, userID := range req.UserIDs {
		if err := authverify.CheckAccess(ctx, userID); err != nil {
			return nil, err
		}
	}
	if req.Conversation.ConversationType == constant.WriteGroupChatType {
		groupInfo, err := c.groupClient.GetGroupInfo(ctx, req.Conversation.GroupID)
		if err != nil {
			return nil, err
		}
		if groupInfo == nil {
			return nil, servererrs.ErrGroupIDNotFound.WrapMsg(req.Conversation.GroupID)
		}
		if groupInfo.Status == constant.GroupStatusDismissed {
			return nil, servererrs.ErrDismissedAlready.WrapMsg("group dismissed")
		}
	}

	conversationMap := make(map[string]*dbModel.Conversation)
	var needUpdateUsersList []string

	for _, userID := range req.UserIDs {
		conversationList, err := c.conversationDatabase.FindConversations(ctx, userID, []string{req.Conversation.ConversationID})
		if err != nil {
			return nil, err
		}
		if len(conversationList) != 0 {
			conversationMap[userID] = conversationList[0]
		} else {
			needUpdateUsersList = append(needUpdateUsersList, userID)
		}
	}

	var conversation dbModel.Conversation
	conversation.ConversationID = req.Conversation.ConversationID
	conversation.ConversationType = req.Conversation.ConversationType
	conversation.UserID = req.Conversation.UserID
	conversation.GroupID = req.Conversation.GroupID
	conversation.CreateTime = time.Now()

	m, conversation, err := UpdateConversationsMap(ctx, req)
	if err != nil {
		return nil, err
	}

	for userID := range conversationMap {
		unequal := UserUpdateCheckMap(ctx, userID, req.Conversation, conversationMap[userID])

		if unequal {
			needUpdateUsersList = append(needUpdateUsersList, userID)
		}
	}

	if len(m) != 0 && len(needUpdateUsersList) != 0 {
		if err := c.conversationDatabase.SetUsersConversationFieldTx(ctx, needUpdateUsersList, &conversation, m); err != nil {
			return nil, err
		}

		for _, userID := range needUpdateUsersList {
			c.conversationNotificationSender.ConversationChangeNotification(ctx, userID, []string{req.Conversation.ConversationID})
		}
	}

	if req.Conversation.IsPrivateChat != nil && req.Conversation.ConversationType != constant.ReadGroupChatType {
		var conversations []*dbModel.Conversation
		for _, ownerUserID := range req.UserIDs {
			transConversation := conversation
			transConversation.OwnerUserID = ownerUserID
			transConversation.IsPrivateChat = req.Conversation.IsPrivateChat.Value
			conversations = append(conversations, &transConversation)
		}

		if err := c.conversationDatabase.SyncPeerUserPrivateConversationTx(ctx, conversations); err != nil {
			return nil, err
		}

		for _, userID := range req.UserIDs {
			c.conversationNotificationSender.ConversationSetPrivateNotification(ctx, userID, req.Conversation.UserID,
				req.Conversation.IsPrivateChat.Value, req.Conversation.ConversationID)
		}
	}

	return &pbconversation.SetConversationsResp{}, nil
}

func (c *conversationServer) UpdateConversationsByUser(ctx context.Context, req *pbconversation.UpdateConversationsByUserReq) (*pbconversation.UpdateConversationsByUserResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	m := make(map[string]any)
	if req.Ex != nil {
		m["ex"] = req.Ex.Value
	}
	if len(m) > 0 {
		if err := c.conversationDatabase.UpdateUserConversations(ctx, req.UserID, m); err != nil {
			return nil, err
		}
	}
	return &pbconversation.UpdateConversationsByUserResp{}, nil
}

// create conversation without notification for msg redis transfer.
func (c *conversationServer) CreateSingleChatConversations(ctx context.Context, req *pbconversation.CreateSingleChatConversationsReq) (*pbconversation.CreateSingleChatConversationsResp, error) {
	var conversation dbModel.Conversation
	conversation.CreateTime = time.Now()

	switch req.ConversationType {
	case constant.SingleChatType:
		// sendUser create
		conversation.ConversationID = req.ConversationID
		conversation.ConversationType = req.ConversationType
		conversation.OwnerUserID = req.SendID
		conversation.UserID = req.RecvID

		if err := c.webhookBeforeCreateSingleChatConversations(ctx, &c.config.WebhooksConfig.BeforeCreateSingleChatConversations, &conversation); err != nil && err != servererrs.ErrCallbackContinue {
			return nil, err
		}

		err := c.conversationDatabase.CreateConversation(ctx, []*dbModel.Conversation{&conversation})
		if err != nil {
			log.ZWarn(ctx, "create conversation failed", err, "conversation", conversation)
		}

		c.webhookAfterCreateSingleChatConversations(ctx, &c.config.WebhooksConfig.AfterCreateSingleChatConversations, &conversation)

		// recvUser create
		conversation2 := conversation
		conversation2.OwnerUserID = req.RecvID
		conversation2.UserID = req.SendID

		if err := c.webhookBeforeCreateSingleChatConversations(ctx, &c.config.WebhooksConfig.BeforeCreateSingleChatConversations, &conversation); err != nil && err != servererrs.ErrCallbackContinue {
			return nil, err
		}

		err = c.conversationDatabase.CreateConversation(ctx, []*dbModel.Conversation{&conversation2})
		if err != nil {
			log.ZWarn(ctx, "create conversation failed", err, "conversation2", conversation)
		}

		c.webhookAfterCreateSingleChatConversations(ctx, &c.config.WebhooksConfig.AfterCreateSingleChatConversations, &conversation2)
	case constant.NotificationChatType:
		conversation.ConversationID = req.ConversationID
		conversation.ConversationType = req.ConversationType
		conversation.OwnerUserID = req.RecvID
		conversation.UserID = req.SendID

		if err := c.webhookBeforeCreateSingleChatConversations(ctx, &c.config.WebhooksConfig.BeforeCreateSingleChatConversations, &conversation); err != nil && err != servererrs.ErrCallbackContinue {
			return nil, err
		}

		err := c.conversationDatabase.CreateConversation(ctx, []*dbModel.Conversation{&conversation})
		if err != nil {
			log.ZWarn(ctx, "create conversation failed", err, "conversation2", conversation)
		}

		c.webhookAfterCreateSingleChatConversations(ctx, &c.config.WebhooksConfig.AfterCreateSingleChatConversations, &conversation)
	}

	return &pbconversation.CreateSingleChatConversationsResp{}, nil
}

func (c *conversationServer) CreateGroupChatConversations(ctx context.Context, req *pbconversation.CreateGroupChatConversationsReq) (*pbconversation.CreateGroupChatConversationsResp, error) {
	var conversation dbModel.Conversation

	conversation.ConversationID = msgprocessor.GetConversationIDBySessionType(constant.ReadGroupChatType, req.GroupID)
	conversation.GroupID = req.GroupID
	conversation.ConversationType = constant.ReadGroupChatType
	conversation.CreateTime = time.Now()

	if err := c.webhookBeforeCreateGroupChatConversations(ctx, &c.config.WebhooksConfig.BeforeCreateGroupChatConversations, &conversation); err != nil {
		return nil, err
	}

	err := c.conversationDatabase.CreateGroupChatConversation(ctx, req.GroupID, req.UserIDs, &conversation)
	if err != nil {
		return nil, err
	}
	if err := c.msgClient.SetUserConversationMaxSeq(ctx, conversation.ConversationID, req.UserIDs, 0); err != nil {
		return nil, err
	}

	c.webhookAfterCreateGroupChatConversations(ctx, &c.config.WebhooksConfig.AfterCreateGroupChatConversations, &conversation)
	return &pbconversation.CreateGroupChatConversationsResp{}, nil
}

func (c *conversationServer) SetConversationMaxSeq(ctx context.Context, req *pbconversation.SetConversationMaxSeqReq) (*pbconversation.SetConversationMaxSeqResp, error) {
	if err := c.msgClient.SetUserConversationMaxSeq(ctx, req.ConversationID, req.OwnerUserID, req.MaxSeq); err != nil {
		return nil, err
	}
	if err := c.conversationDatabase.UpdateUsersConversationField(ctx, req.OwnerUserID, req.ConversationID,
		map[string]any{"max_seq": req.MaxSeq}); err != nil {
		return nil, err
	}
	for _, userID := range req.OwnerUserID {
		c.conversationNotificationSender.ConversationChangeNotification(ctx, userID, []string{req.ConversationID})
	}
	return &pbconversation.SetConversationMaxSeqResp{}, nil
}

func (c *conversationServer) SetConversationMinSeq(ctx context.Context, req *pbconversation.SetConversationMinSeqReq) (*pbconversation.SetConversationMinSeqResp, error) {
	if err := c.msgClient.SetUserConversationMin(ctx, req.ConversationID, req.OwnerUserID, req.MinSeq); err != nil {
		return nil, err
	}
	if err := c.conversationDatabase.UpdateUsersConversationField(ctx, req.OwnerUserID, req.ConversationID,
		map[string]any{"min_seq": req.MinSeq}); err != nil {
		return nil, err
	}
	for _, userID := range req.OwnerUserID {
		c.conversationNotificationSender.ConversationChangeNotification(ctx, userID, []string{req.ConversationID})
	}
	return &pbconversation.SetConversationMinSeqResp{}, nil
}

func (c *conversationServer) GetConversationIDs(ctx context.Context, req *pbconversation.GetConversationIDsReq) (*pbconversation.GetConversationIDsResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	conversationIDs, err := c.conversationDatabase.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetConversationIDsResp{ConversationIDs: conversationIDs}, nil
}

func (c *conversationServer) GetUserConversationIDsHash(ctx context.Context, req *pbconversation.GetUserConversationIDsHashReq) (*pbconversation.GetUserConversationIDsHashResp, error) {
	if err := authverify.CheckAccess(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	hash, err := c.conversationDatabase.GetUserConversationIDsHash(ctx, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetUserConversationIDsHashResp{Hash: hash}, nil
}

func (c *conversationServer) GetConversationsByConversationID(ctx context.Context, req *pbconversation.GetConversationsByConversationIDReq) (*pbconversation.GetConversationsByConversationIDResp, error) {
	conversations, err := c.conversationDatabase.GetConversationsByConversationID(ctx, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetConversationsByConversationIDResp{Conversations: convert.ConversationsDB2Pb(conversations)}, nil
}

func (c *conversationServer) GetConversationOfflinePushUserIDs(ctx context.Context, req *pbconversation.GetConversationOfflinePushUserIDsReq) (*pbconversation.GetConversationOfflinePushUserIDsResp, error) {
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversationID is empty")
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
	return &pbconversation.GetConversationOfflinePushUserIDsResp{UserIDs: datautil.Keys(userIDSet)}, nil
}

func (c *conversationServer) conversationSort(conversations map[int64]string, resp *pbconversation.GetSortedConversationListResp, conversation_unreadCount map[string]int64, conversationMsg map[string]*pbconversation.ConversationElem) {
	keys := []int64{}
	for key := range conversations {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
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

func (c *conversationServer) getConversationInfo(ctx context.Context, chatLogs map[string]*sdkws.MsgData, userID string) (map[string]*pbconversation.ConversationElem, error) {
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
		case constant.WriteGroupChatType, constant.ReadGroupChatType:
			groupIDs = append(groupIDs, chatLog.GroupID)
			sendIDs = append(sendIDs, chatLog.SendID)
		}
	}
	if len(sendIDs) != 0 {
		sendInfos, err := c.userClient.GetUsersInfo(ctx, sendIDs)
		if err != nil {
			return nil, err
		}
		for _, sendInfo := range sendInfos {
			sendMap[sendInfo.UserID] = sendInfo
		}
	}
	if len(groupIDs) != 0 {
		groupInfos, err := c.groupClient.GetGroupsInfo(ctx, groupIDs)
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
		if err := datautil.CopyStructFields(msgInfo, chatLog); err != nil {
			return nil, err
		}
		switch chatLog.SessionType {
		case constant.SingleChatType:
			if chatLog.SendID == userID {
				if recv, ok := sendMap[chatLog.RecvID]; ok {
					msgInfo.FaceURL = recv.FaceURL
					msgInfo.SenderName = recv.Nickname
				}
				break
			}
			if send, ok := sendMap[chatLog.SendID]; ok {
				msgInfo.FaceURL = send.FaceURL
				msgInfo.SenderName = send.Nickname
			}
		case constant.WriteGroupChatType, constant.ReadGroupChatType:
			msgInfo.GroupID = chatLog.GroupID
			if group, ok := groupMap[chatLog.GroupID]; ok {
				msgInfo.GroupName = group.GroupName
				msgInfo.GroupFaceURL = group.FaceURL
				msgInfo.GroupMemberCount = group.MemberCount
				msgInfo.GroupType = group.GroupType
			}
			if send, ok := sendMap[chatLog.SendID]; ok {
				msgInfo.SenderName = send.Nickname
			}
		}
		pbchatLog.ConversationID = conversationID
		msgInfo.LatestMsgRecvTime = chatLog.SendTime
		pbchatLog.MsgInfo = msgInfo
		conversationMsg[conversationID] = pbchatLog
	}
	return conversationMsg, nil
}

func (c *conversationServer) GetConversationNotReceiveMessageUserIDs(ctx context.Context, req *pbconversation.GetConversationNotReceiveMessageUserIDsReq) (*pbconversation.GetConversationNotReceiveMessageUserIDsResp, error) {
	userIDs, err := c.conversationDatabase.GetConversationNotReceiveMessageUserIDs(ctx, req.ConversationID)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetConversationNotReceiveMessageUserIDsResp{UserIDs: userIDs}, nil
}

func (c *conversationServer) UpdateConversation(ctx context.Context, req *pbconversation.UpdateConversationReq) (*pbconversation.UpdateConversationResp, error) {
	for _, userID := range req.UserIDs {
		if err := authverify.CheckAccess(ctx, userID); err != nil {
			return nil, err
		}
	}
	m := make(map[string]any)
	if req.RecvMsgOpt != nil {
		m["recv_msg_opt"] = req.RecvMsgOpt.Value
	}
	if req.AttachedInfo != nil {
		m["attached_info"] = req.AttachedInfo.Value
	}
	if req.Ex != nil {
		m["ex"] = req.Ex.Value
	}
	if req.IsPinned != nil {
		m["is_pinned"] = req.IsPinned.Value
	}
	if req.GroupAtType != nil {
		m["group_at_type"] = req.GroupAtType.Value
	}
	if req.MsgDestructTime != nil {
		m["msg_destruct_time"] = req.MsgDestructTime.Value
	}
	if req.IsMsgDestruct != nil {
		m["is_msg_destruct"] = req.IsMsgDestruct.Value
	}
	if req.BurnDuration != nil {
		m["burn_duration"] = req.BurnDuration.Value
	}
	if req.IsPrivateChat != nil {
		m["is_private_chat"] = req.IsPrivateChat.Value
	}
	if req.MinSeq != nil {
		m["min_seq"] = req.MinSeq.Value
	}
	if req.MaxSeq != nil {
		m["max_seq"] = req.MaxSeq.Value
	}
	if req.LatestMsgDestructTime != nil {
		m["latest_msg_destruct_time"] = time.UnixMilli(req.LatestMsgDestructTime.Value)
	}
	if len(m) > 0 {
		if err := c.conversationDatabase.UpdateUsersConversationField(ctx, req.UserIDs, req.ConversationID, m); err != nil {
			return nil, err
		}
	}
	return &pbconversation.UpdateConversationResp{}, nil
}

func (c *conversationServer) GetOwnerConversation(ctx context.Context, req *pbconversation.GetOwnerConversationReq) (*pbconversation.GetOwnerConversationResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	total, conversations, err := c.conversationDatabase.GetOwnerConversation(ctx, req.UserID, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetOwnerConversationResp{
		Total:         total,
		Conversations: convert.ConversationsDB2Pb(conversations),
	}, nil
}

func (c *conversationServer) GetConversationsNeedClearMsg(ctx context.Context, _ *pbconversation.GetConversationsNeedClearMsgReq) (*pbconversation.GetConversationsNeedClearMsgResp, error) {
	num, err := c.conversationDatabase.GetAllConversationIDsNumber(ctx)
	if err != nil {
		log.ZError(ctx, "GetAllConversationIDsNumber failed", err)
		return nil, err
	}
	const batchNum = 100

	if num == 0 {
		return nil, errs.New("Need Destruct Msg is nil").Wrap()
	}

	maxPage := (num + batchNum - 1) / batchNum

	temp := make([]*dbModel.Conversation, 0, maxPage*batchNum)

	for pageNumber := 0; pageNumber < int(maxPage); pageNumber++ {
		pagination := &sdkws.RequestPagination{
			PageNumber: int32(pageNumber),
			ShowNumber: batchNum,
		}

		conversationIDs, err := c.conversationDatabase.PageConversationIDs(ctx, pagination)
		if err != nil {
			log.ZError(ctx, "PageConversationIDs failed", err, "pageNumber", pageNumber)
			continue
		}

		// log.ZDebug(ctx, "PageConversationIDs success", "pageNumber", pageNumber, "conversationIDsNum", len(conversationIDs), "conversationIDs", conversationIDs)
		if len(conversationIDs) == 0 {
			continue
		}

		conversations, err := c.conversationDatabase.GetConversationsByConversationID(ctx, conversationIDs)
		if err != nil {
			log.ZError(ctx, "GetConversationsByConversationID failed", err, "conversationIDs", conversationIDs)
			continue
		}

		for _, conversation := range conversations {
			if conversation.IsMsgDestruct && conversation.MsgDestructTime != 0 && ((time.Now().UnixMilli() > (conversation.MsgDestructTime + conversation.LatestMsgDestructTime.UnixMilli() + 8*60*60)) || // 8*60*60 is UTC+8
				conversation.LatestMsgDestructTime.IsZero()) {
				temp = append(temp, conversation)
			}
		}
	}

	return &pbconversation.GetConversationsNeedClearMsgResp{Conversations: convert.ConversationsDB2Pb(temp)}, nil
}

func (c *conversationServer) GetNotNotifyConversationIDs(ctx context.Context, req *pbconversation.GetNotNotifyConversationIDsReq) (*pbconversation.GetNotNotifyConversationIDsResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	conversationIDs, err := c.conversationDatabase.GetNotNotifyConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetNotNotifyConversationIDsResp{ConversationIDs: conversationIDs}, nil
}

func (c *conversationServer) GetPinnedConversationIDs(ctx context.Context, req *pbconversation.GetPinnedConversationIDsReq) (*pbconversation.GetPinnedConversationIDsResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	conversationIDs, err := c.conversationDatabase.GetPinnedConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pbconversation.GetPinnedConversationIDsResp{ConversationIDs: conversationIDs}, nil
}

func (c *conversationServer) ClearUserConversationMsg(ctx context.Context, req *pbconversation.ClearUserConversationMsgReq) (*pbconversation.ClearUserConversationMsgResp, error) {
	conversations, err := c.conversationDatabase.FindRandConversation(ctx, req.Timestamp, int(req.Limit))
	if err != nil {
		return nil, err
	}
	latestMsgDestructTime := time.UnixMilli(req.Timestamp)
	for i, conversation := range conversations {
		if conversation.IsMsgDestruct == false || conversation.MsgDestructTime == 0 {
			continue
		}
		seq, err := c.msgClient.GetLastMessageSeqByTime(ctx, conversation.ConversationID, req.Timestamp-(conversation.MsgDestructTime*1000))
		if err != nil {
			return nil, err
		}
		if seq <= 0 {
			log.ZDebug(ctx, "ClearUserConversationMsg GetLastMessageSeqByTime seq <= 0", "index", i, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID, "msgDestructTime", conversation.MsgDestructTime, "seq", seq)
			if err := c.setConversationMinSeqAndLatestMsgDestructTime(ctx, conversation.ConversationID, conversation.OwnerUserID, -1, latestMsgDestructTime); err != nil {
				return nil, err
			}
			continue
		}
		seq++
		if err := c.setConversationMinSeqAndLatestMsgDestructTime(ctx, conversation.ConversationID, conversation.OwnerUserID, seq, latestMsgDestructTime); err != nil {
			return nil, err
		}
		log.ZDebug(ctx, "ClearUserConversationMsg set min seq", "index", i, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID, "seq", seq, "msgDestructTime", conversation.MsgDestructTime)
	}
	return &pbconversation.ClearUserConversationMsgResp{Count: int32(len(conversations))}, nil
}

func (c *conversationServer) setConversationMinSeqAndLatestMsgDestructTime(ctx context.Context, conversationID string, ownerUserID string, minSeq int64, latestMsgDestructTime time.Time) error {
	update := map[string]any{
		"latest_msg_destruct_time": latestMsgDestructTime,
	}
	if minSeq >= 0 {
		if err := c.msgClient.SetUserConversationMin(ctx, conversationID, []string{ownerUserID}, minSeq); err != nil {
			return err
		}
		update["min_seq"] = minSeq
	}

	if err := c.conversationDatabase.UpdateUsersConversationField(ctx, []string{ownerUserID}, conversationID, update); err != nil {
		return err
	}
	c.conversationNotificationSender.ConversationChangeNotification(ctx, ownerUserID, []string{conversationID})
	return nil
}

func (c *conversationServer) GetConversationReadCursors(ctx context.Context, req *pbconversation.GetConversationReadCursorsReq) (*pbconversation.GetConversationReadCursorsResp, error) {
	if err := req.Check(); err != nil {
		return nil, err
	}

	conversations, err := c.conversationDatabase.GetConversationsByConversationID(ctx, req.ConversationIDs)
	if err != nil {
		return nil, err
	}

	conversationMap := make(map[string]*dbModel.Conversation)
	for _, conv := range conversations {
		conversationMap[conv.ConversationID] = conv
	}

	result := make(map[string]*pbconversation.ConversationReadCursorList)
	for _, conversationID := range req.ConversationIDs {
		conversation, exists := conversationMap[conversationID]
		if !exists {
			continue
		}

		var userIDs []string
		if conversation.ConversationType == constant.ReadGroupChatType {
			if conversation.GroupID == "" {
				log.ZWarn(ctx, "groupID is empty for group conversation", nil, "conversationID", conversationID)
				result[conversationID] = &pbconversation.ConversationReadCursorList{Cursors: []*pbconversation.ConversationReadCursor{}}
				continue
			}
			userIDs, err = c.groupClient.GetGroupMemberUserIDs(ctx, conversation.GroupID)
			if err != nil {
				log.ZWarn(ctx, "GetGroupMemberUserIDs failed", err, "conversationID", conversationID, "groupID", conversation.GroupID)
				result[conversationID] = &pbconversation.ConversationReadCursorList{Cursors: []*pbconversation.ConversationReadCursor{}}
				continue
			}
		} else {
			continue
		}

		var cursors []*pbconversation.ConversationReadCursor
		for _, userID := range userIDs {
			hasReadSeqs, err := c.msgClient.GetHasReadSeqs(ctx, []string{conversationID}, userID)
			if err != nil {
				log.ZWarn(ctx, "GetHasReadSeqs failed", err, "userID", userID, "conversationID", conversationID)
				continue
			}

			readSeq := hasReadSeqs[conversationID]
			cursors = append(cursors, &pbconversation.ConversationReadCursor{
				UserID:     userID,
				MaxReadSeq: readSeq,
			})
		}

		result[conversationID] = &pbconversation.ConversationReadCursorList{Cursors: cursors}
	}

	return &pbconversation.GetConversationReadCursorsResp{
		Cursors: result,
	}, nil
}
