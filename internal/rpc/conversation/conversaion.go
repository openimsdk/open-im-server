package conversation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/internal/common/check"
	"github.com/OpenIMSDK/Open-IM-Server/internal/common/notification"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	tableRelation "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

type conversationServer struct {
	groupChecker *check.GroupChecker
	controller.ConversationDatabase
	notify *notification.Check
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
	pbConversation.RegisterConversationServer(server, &conversationServer{
		groupChecker:         check.NewGroupChecker(client),
		ConversationDatabase: controller.NewConversationDatabase(relation.NewConversationGorm(db), cache.NewConversationRedis(rdb, cache.GetDefaultOpt()), tx.NewGorm(db)),
	})
	return nil
}

func (c *conversationServer) GetConversation(ctx context.Context, req *pbConversation.GetConversationReq) (*pbConversation.GetConversationResp, error) {
	resp := &pbConversation.GetConversationResp{Conversation: &pbConversation.Conversation{}}
	conversations, err := c.ConversationDatabase.FindConversations(ctx, req.OwnerUserID, []string{req.ConversationID})
	if err != nil {
		return nil, err
	}
	if len(conversations) > 0 {
		if err := utils.CopyStructFields(resp.Conversation, &conversations[0]); err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, nil
}

func (c *conversationServer) GetAllConversations(ctx context.Context, req *pbConversation.GetAllConversationsReq) (*pbConversation.GetAllConversationsResp, error) {
	resp := &pbConversation.GetAllConversationsResp{Conversations: []*pbConversation.Conversation{}}
	conversations, err := c.ConversationDatabase.GetUserAllConversation(ctx, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	if err := utils.CopyStructFields(&resp.Conversations, conversations); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *conversationServer) GetConversations(ctx context.Context, req *pbConversation.GetConversationsReq) (*pbConversation.GetConversationsResp, error) {
	resp := &pbConversation.GetConversationsResp{Conversations: []*pbConversation.Conversation{}}
	conversations, err := c.ConversationDatabase.FindConversations(ctx, req.OwnerUserID, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	if err := utils.CopyStructFields(&resp.Conversations, conversations); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *conversationServer) BatchSetConversations(ctx context.Context, req *pbConversation.BatchSetConversationsReq) (*pbConversation.BatchSetConversationsResp, error) {
	resp := &pbConversation.BatchSetConversationsResp{}
	var conversations []*tableRelation.ConversationModel
	if err := utils.CopyStructFields(&conversations, req.Conversations); err != nil {
		return nil, err
	}
	err := c.ConversationDatabase.SetUserConversations(ctx, req.OwnerUserID, conversations)
	if err != nil {
		return nil, err
	}
	c.notify.ConversationChangeNotification(ctx, req.OwnerUserID)
	return resp, nil
}

func (c *conversationServer) SetConversation(ctx context.Context, req *pbConversation.SetConversationReq) (*pbConversation.SetConversationResp, error) {
	panic("implement me")
}

func (c *conversationServer) SetRecvMsgOpt(ctx context.Context, req *pbConversation.SetRecvMsgOptReq) (*pbConversation.SetRecvMsgOptResp, error) {
	panic("implement me")
}

func (c *conversationServer) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) (*pbConversation.ModifyConversationFieldResp, error) {
	resp := &pbConversation.ModifyConversationFieldResp{}
	var err error
	isSyncConversation := true
	if req.Conversation.ConversationType == constant.GroupChatType {
		groupInfo, err := c.groupChecker.GetGroupInfo(ctx, req.Conversation.GroupID)
		if err != nil {
			return nil, err
		}
		if groupInfo.Status == constant.GroupStatusDismissed && req.FieldType != constant.FieldUnread {
			return nil, err
		}
	}
	var conversation tableRelation.ConversationModel
	if err := utils.CopyStructFields(&conversation, req.Conversation); err != nil {
		return nil, err
	}
	if req.FieldType == constant.FieldIsPrivateChat {
		err := c.ConversationDatabase.SyncPeerUserPrivateConversationTx(ctx, &conversation)
		if err != nil {
			return nil, err
		}
		c.notify.ConversationSetPrivateNotification(ctx, req.Conversation.OwnerUserID, req.Conversation.UserID, req.Conversation.IsPrivateChat)
		return resp, nil
	}
	//haveUserID, err := c.ConversationDatabase.GetUserIDExistConversation(ctx, req.UserIDList, req.Conversation.ConversationID)
	//if err != nil {
	//	return nil, err
	//}
	filedMap := make(map[string]interface{})
	switch req.FieldType {
	case constant.FieldRecvMsgOpt:
		filedMap["recv_msg_opt"] = req.Conversation.RecvMsgOpt
	case constant.FieldGroupAtType:
		filedMap["group_at_type"] = req.Conversation.GroupAtType
	case constant.FieldIsNotInGroup:
		filedMap["is_not_in_group"] = req.Conversation.IsNotInGroup
	case constant.FieldIsPinned:
		filedMap["is_pinned"] = req.Conversation.IsPinned
	case constant.FieldEx:
		filedMap["ex"] = req.Conversation.Ex
	case constant.FieldAttachedInfo:
		filedMap["attached_info"] = req.Conversation.AttachedInfo
	case constant.FieldUnread:
		isSyncConversation = false
		filedMap["update_unread_count_time"] = req.Conversation.UpdateUnreadCountTime
	case constant.FieldBurnDuration:
		filedMap["burn_duration"] = req.Conversation.BurnDuration
	}
	err = c.ConversationDatabase.SetUsersConversationFiledTx(ctx, req.UserIDList, &conversation, filedMap)
	if err != nil {
		return nil, err
	}

	if isSyncConversation {
		for _, v := range req.UserIDList {
			c.notify.ConversationChangeNotification(ctx, v)
		}
	} else {
		for _, v := range req.UserIDList {
			c.notify.ConversationUnreadChangeNotification(ctx, v, req.Conversation.ConversationID, req.Conversation.UpdateUnreadCountTime)
		}
	}
	return resp, nil
}

// 获取超级大群开启免打扰的用户ID
func (c *conversationServer) GetRecvMsgNotNotifyUserIDs(ctx context.Context, req *pbConversation.GetRecvMsgNotNotifyUserIDsReq) (*pbConversation.GetRecvMsgNotNotifyUserIDsResp, error) {
	userIDs, err := c.ConversationDatabase.FindRecvMsgNotNotifyUserIDs(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	resp := &pbConversation.GetRecvMsgNotNotifyUserIDsResp{}
	resp.UserIDs = userIDs
	return resp, nil
}
