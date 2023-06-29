package conversation

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	tableRelation "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient/notification"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

type conversationServer struct {
	groupRpcClient                 *rpcclient.GroupRpcClient
	conversationDatabase           controller.ConversationDatabase
	conversationNotificationSender *notification.ConversationNotificationSender
	msgRpcClient                   *rpcclient.MessageRpcClient
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&tableRelation.ConversationModel{}); err != nil {
		return err
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	conversationDB := relation.NewConversationGorm(db)
	groupRpcClient := rpcclient.NewGroupRpcClient(client)
	msgRpcClient := rpcclient.NewMessageRpcClient(client)
	pbConversation.RegisterConversationServer(server, &conversationServer{
		conversationNotificationSender: notification.NewConversationNotificationSender(client),
		groupRpcClient:                 &groupRpcClient,
		msgRpcClient:                   &msgRpcClient,
		conversationDatabase:           controller.NewConversationDatabase(conversationDB, cache.NewConversationRedis(rdb, cache.GetDefaultOpt(), conversationDB), tx.NewGorm(db)),
	})
	return nil
}

func (c *conversationServer) GetConversation(ctx context.Context, req *pbConversation.GetConversationReq) (*pbConversation.GetConversationResp, error) {
	conversations, err := c.conversationDatabase.FindConversations(ctx, req.OwnerUserID, []string{req.ConversationID})
	if err != nil {
		return nil, err
	}
	if len(conversations) < 1 {
		return nil, errs.ErrRecordNotFound.Wrap("conversation not found")
	}
	resp := &pbConversation.GetConversationResp{Conversation: &pbConversation.Conversation{}}
	resp.Conversation = convert.ConversationDB2Pb(conversations[0])
	return resp, nil
}

func (c *conversationServer) GetAllConversations(ctx context.Context, req *pbConversation.GetAllConversationsReq) (*pbConversation.GetAllConversationsResp, error) {
	conversations, err := c.conversationDatabase.GetUserAllConversation(ctx, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	resp := &pbConversation.GetAllConversationsResp{Conversations: []*pbConversation.Conversation{}}
	resp.Conversations = convert.ConversationsDB2Pb(conversations)
	return resp, nil
}

func (c *conversationServer) GetConversations(ctx context.Context, req *pbConversation.GetConversationsReq) (*pbConversation.GetConversationsResp, error) {
	conversations, err := c.conversationDatabase.FindConversations(ctx, req.OwnerUserID, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	resp := &pbConversation.GetConversationsResp{Conversations: []*pbConversation.Conversation{}}
	resp.Conversations = convert.ConversationsDB2Pb(conversations)
	return resp, nil
}

func (c *conversationServer) BatchSetConversations(ctx context.Context, req *pbConversation.BatchSetConversationsReq) (*pbConversation.BatchSetConversationsResp, error) {
	conversations := convert.ConversationsPb2DB(req.Conversations)
	err := c.conversationDatabase.SetUserConversations(ctx, req.OwnerUserID, conversations)
	if err != nil {
		return nil, err
	}
	_ = c.conversationNotificationSender.ConversationChangeNotification(ctx, req.OwnerUserID)
	return &pbConversation.BatchSetConversationsResp{}, nil
}

func (c *conversationServer) SetConversation(ctx context.Context, req *pbConversation.SetConversationReq) (*pbConversation.SetConversationResp, error) {
	var conversation tableRelation.ConversationModel
	if err := utils.CopyStructFields(&conversation, req.Conversation); err != nil {
		return nil, err
	}
	err := c.conversationDatabase.SetUserConversations(ctx, req.Conversation.OwnerUserID, []*tableRelation.ConversationModel{&conversation})
	if err != nil {
		return nil, err
	}
	_ = c.conversationNotificationSender.ConversationChangeNotification(ctx, req.Conversation.OwnerUserID)
	resp := &pbConversation.SetConversationResp{}
	return resp, nil
}

func (c *conversationServer) SetRecvMsgOpt(ctx context.Context, req *pbConversation.SetRecvMsgOptReq) (*pbConversation.SetRecvMsgOptResp, error) {
	if err := c.conversationDatabase.SetUsersConversationFiledTx(ctx, []string{req.OwnerUserID}, &tableRelation.ConversationModel{OwnerUserID: req.OwnerUserID, ConversationID: req.ConversationID, RecvMsgOpt: req.RecvMsgOpt}, map[string]interface{}{"recv_msg_opt": req.RecvMsgOpt}); err != nil {
		return nil, err
	}
	_ = c.conversationNotificationSender.ConversationChangeNotification(ctx, req.OwnerUserID)
	return &pbConversation.SetRecvMsgOptResp{}, nil
}

// deprecated
func (c *conversationServer) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) (*pbConversation.ModifyConversationFieldResp, error) {
	resp := &pbConversation.ModifyConversationFieldResp{}
	var err error
	isSyncConversation := true
	if req.Conversation.ConversationType == constant.GroupChatType {
		groupInfo, err := c.groupRpcClient.GetGroupInfo(ctx, req.Conversation.GroupID)
		if err != nil {
			return nil, err
		}
		if groupInfo.Status == constant.GroupStatusDismissed && req.FieldType != constant.FieldUnread {
			return nil, err
		}
	}
	conversation := convert.ConversationPb2DB(req.Conversation)
	if req.FieldType == constant.FieldIsPrivateChat {
		err := c.conversationDatabase.SyncPeerUserPrivateConversationTx(ctx, []*tableRelation.ConversationModel{conversation})
		if err != nil {
			return nil, err
		}
		c.conversationNotificationSender.ConversationSetPrivateNotification(ctx, req.Conversation.OwnerUserID, req.Conversation.UserID, req.Conversation.IsPrivateChat)
		return resp, nil
	}
	filedMap := make(map[string]interface{})
	switch req.FieldType {
	case constant.FieldRecvMsgOpt:
		filedMap["recv_msg_opt"] = req.Conversation.RecvMsgOpt
	case constant.FieldGroupAtType:
		filedMap["group_at_type"] = req.Conversation.GroupAtType
	case constant.FieldIsPinned:
		filedMap["is_pinned"] = req.Conversation.IsPinned
	case constant.FieldEx:
		filedMap["ex"] = req.Conversation.Ex
	case constant.FieldAttachedInfo:
		filedMap["attached_info"] = req.Conversation.AttachedInfo
	case constant.FieldUnread:
		isSyncConversation = false
		filedMap["update_unread_count_time"] = req.Conversation.UpdateUnreadCountTime
		filedMap["has_read_seq"] = req.Conversation.HasReadSeq
	case constant.FieldBurnDuration:
		filedMap["burn_duration"] = req.Conversation.BurnDuration
	}
	err = c.conversationDatabase.SetUsersConversationFiledTx(ctx, req.UserIDList, conversation, filedMap)
	if err != nil {
		return nil, err
	}

	if isSyncConversation {
		for _, v := range req.UserIDList {
			c.conversationNotificationSender.ConversationChangeNotification(ctx, v)
		}
	} else {
		for _, v := range req.UserIDList {
			c.conversationNotificationSender.ConversationUnreadChangeNotification(ctx, v, req.Conversation.ConversationID, req.Conversation.UpdateUnreadCountTime, req.Conversation.HasReadSeq)
		}
	}
	return resp, nil
}

func (c *conversationServer) SetConversations(ctx context.Context, req *pbConversation.SetConversationsReq) (*pbConversation.SetConversationsResp, error) {
	if req.Conversation == nil {
		return nil, errs.ErrArgs.Wrap("conversation must not be nil")
	}
	isSyncConversation := true
	if req.Conversation.ConversationType == constant.GroupChatType {
		groupInfo, err := c.groupRpcClient.GetGroupInfo(ctx, req.Conversation.GroupID)
		if err != nil {
			return nil, err
		}
		if groupInfo.Status == constant.GroupStatusDismissed {
			return nil, err
		}
	}
	var conversation tableRelation.ConversationModel
	conversation.ConversationID = req.Conversation.ConversationID
	conversation.ConversationType = req.Conversation.ConversationType
	conversation.UserID = req.Conversation.UserID
	conversation.GroupID = req.Conversation.GroupID
	m := make(map[string]interface{})
	if req.Conversation.RecvMsgOpt != nil {
		m["recv_msg_opt"] = req.Conversation.RecvMsgOpt.Value
	}
	if req.Conversation.DraftTextTime != nil {
		m["draft_text_time"] = req.Conversation.DraftTextTime.Value
	}
	if req.Conversation.AttachedInfo != nil {
		m["attached_info"] = req.Conversation.AttachedInfo.Value
	}
	if req.Conversation.Ex != nil {
		m["ex"] = req.Conversation.Ex.Value
	}
	if req.Conversation.IsPinned != nil {
		m["is_pinned"] = req.Conversation.IsPinned.Value
	}
	if req.Conversation.GroupAtType != nil {
		m["group_at_type"] = req.Conversation.GroupAtType.Value
	}
	if req.Conversation.IsPrivateChat != nil {
		var conversations []*tableRelation.ConversationModel
		for _, ownerUserID := range req.UserIDs {
			conversation2 := conversation
			conversation.OwnerUserID = ownerUserID
			conversation.IsPrivateChat = req.Conversation.IsPrivateChat.Value
			conversations = append(conversations, &conversation2)
		}
		if err := c.conversationDatabase.SyncPeerUserPrivateConversationTx(ctx, conversations); err != nil {
			return nil, err
		}
		for _, ownerUserID := range req.UserIDs {
			c.conversationNotificationSender.ConversationSetPrivateNotification(ctx, ownerUserID, req.Conversation.UserID, req.Conversation.IsPrivateChat.Value)
		}
	}
	if req.Conversation.BurnDuration != nil {
		m["burn_duration"] = req.Conversation.BurnDuration.Value
	}
	if req.Conversation.HasReadSeq != nil && req.Conversation.UpdateUnreadCountTime != nil {
		isSyncConversation = false
		m["has_read_seq"] = req.Conversation.HasReadSeq.Value
		m["update_unread_count_time"] = req.Conversation.UpdateUnreadCountTime.Value
	}
	err := c.conversationDatabase.SetUsersConversationFiledTx(ctx, req.UserIDs, &conversation, m)
	if err != nil {
		return nil, err
	}

	if isSyncConversation {
		for _, v := range req.UserIDs {
			c.conversationNotificationSender.ConversationChangeNotification(ctx, v)
		}
	} else {
		for _, v := range req.UserIDs {
			c.conversationNotificationSender.ConversationUnreadChangeNotification(ctx, v, req.Conversation.ConversationID, req.Conversation.UpdateUnreadCountTime.Value, req.Conversation.HasReadSeq.Value)
		}
	}
	return &pbConversation.SetConversationsResp{}, nil
}

// 获取超级大群开启免打扰的用户ID
func (c *conversationServer) GetRecvMsgNotNotifyUserIDs(ctx context.Context, req *pbConversation.GetRecvMsgNotNotifyUserIDsReq) (*pbConversation.GetRecvMsgNotNotifyUserIDsResp, error) {
	userIDs, err := c.conversationDatabase.FindRecvMsgNotNotifyUserIDs(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	return &pbConversation.GetRecvMsgNotNotifyUserIDsResp{UserIDs: userIDs}, nil
}

// create conversation without notification for msg redis transfer
func (c *conversationServer) CreateSingleChatConversations(ctx context.Context, req *pbConversation.CreateSingleChatConversationsReq) (*pbConversation.CreateSingleChatConversationsResp, error) {
	var conversation tableRelation.ConversationModel
	conversation.ConversationID = utils.GetConversationIDBySessionType(constant.SingleChatType, req.RecvID, req.SendID)
	conversation.ConversationType = constant.SingleChatType
	conversation.OwnerUserID = req.SendID
	conversation.UserID = req.RecvID
	err := c.conversationDatabase.CreateConversation(ctx, []*tableRelation.ConversationModel{&conversation})
	if err != nil {
		log.ZWarn(ctx, "create conversation failed", err, "conversation", conversation)
	}

	conversation2 := conversation
	conversation2.OwnerUserID = req.RecvID
	conversation2.UserID = req.SendID
	err = c.conversationDatabase.CreateConversation(ctx, []*tableRelation.ConversationModel{&conversation2})
	if err != nil {
		log.ZWarn(ctx, "create conversation failed", err, "conversation2", conversation)
	}
	return &pbConversation.CreateSingleChatConversationsResp{}, nil
}

func (c *conversationServer) CreateGroupChatConversations(ctx context.Context, req *pbConversation.CreateGroupChatConversationsReq) (*pbConversation.CreateGroupChatConversationsResp, error) {
	err := c.conversationDatabase.CreateGroupChatConversation(ctx, req.GroupID, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &pbConversation.CreateGroupChatConversationsResp{}, nil
}

func (c *conversationServer) SetConversationMaxSeq(ctx context.Context, req *pbConversation.SetConversationMaxSeqReq) (*pbConversation.SetConversationMaxSeqResp, error) {
	if err := c.conversationDatabase.UpdateUsersConversationFiled(ctx, req.OwnerUserID, req.ConversationID,
		map[string]interface{}{"max_seq": req.MaxSeq}); err != nil {
		return nil, err
	}
	return &pbConversation.SetConversationMaxSeqResp{}, nil
}

func (c *conversationServer) GetConversationIDs(ctx context.Context, req *pbConversation.GetConversationIDsReq) (*pbConversation.GetConversationIDsResp, error) {
	conversationIDs, err := c.conversationDatabase.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pbConversation.GetConversationIDsResp{ConversationIDs: conversationIDs}, nil
}

func (c *conversationServer) GetUserConversationIDsHash(ctx context.Context, req *pbConversation.GetUserConversationIDsHashReq) (*pbConversation.GetUserConversationIDsHashResp, error) {
	hash, err := c.conversationDatabase.GetUserConversationIDsHash(ctx, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	return &pbConversation.GetUserConversationIDsHashResp{Hash: hash}, nil
}

func (c *conversationServer) GetConversationsByConversationID(ctx context.Context, req *pbConversation.GetConversationsByConversationIDReq) (*pbConversation.GetConversationsByConversationIDResp, error) {
	conversations, err := c.conversationDatabase.GetConversationsByConversationID(ctx, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	return &pbConversation.GetConversationsByConversationIDResp{Conversations: convert.ConversationsDB2Pb(conversations)}, nil
}
