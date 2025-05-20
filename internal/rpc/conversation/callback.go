package conversation

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	dbModel "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/tools/utils/datautil"
)

func (c *conversationServer) webhookBeforeCreateSingleChatConversations(ctx context.Context, before *config.BeforeConfig, req *dbModel.Conversation) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &callbackstruct.CallbackBeforeCreateSingleChatConversationsReq{
			CallbackCommand:  callbackstruct.CallbackBeforeCreateSingleChatConversationsCommand,
			OwnerUserID:      req.OwnerUserID,
			ConversationID:   req.ConversationID,
			ConversationType: req.ConversationType,
			UserID:           req.UserID,
			RecvMsgOpt:       req.RecvMsgOpt,
			IsPinned:         req.IsPinned,
			IsPrivateChat:    req.IsPrivateChat,
			BurnDuration:     req.BurnDuration,
			GroupAtType:      req.GroupAtType,
			AttachedInfo:     req.AttachedInfo,
			Ex:               req.Ex,
		}

		resp := &callbackstruct.CallbackBeforeCreateSingleChatConversationsResp{}

		if err := c.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		datautil.NotNilReplace(&req.RecvMsgOpt, resp.RecvMsgOpt)
		datautil.NotNilReplace(&req.IsPinned, resp.IsPinned)
		datautil.NotNilReplace(&req.IsPrivateChat, resp.IsPrivateChat)
		datautil.NotNilReplace(&req.BurnDuration, resp.BurnDuration)
		datautil.NotNilReplace(&req.GroupAtType, resp.GroupAtType)
		datautil.NotNilReplace(&req.AttachedInfo, resp.AttachedInfo)
		datautil.NotNilReplace(&req.Ex, resp.Ex)
		return nil
	})
}

func (c *conversationServer) webhookAfterCreateSingleChatConversations(ctx context.Context, after *config.AfterConfig, req *dbModel.Conversation) error {
	cbReq := &callbackstruct.CallbackAfterCreateSingleChatConversationsReq{
		CallbackCommand:  callbackstruct.CallbackAfterCreateSingleChatConversationsCommand,
		OwnerUserID:      req.OwnerUserID,
		ConversationID:   req.ConversationID,
		ConversationType: req.ConversationType,
		UserID:           req.UserID,
		RecvMsgOpt:       req.RecvMsgOpt,
		IsPinned:         req.IsPinned,
		IsPrivateChat:    req.IsPrivateChat,
		BurnDuration:     req.BurnDuration,
		GroupAtType:      req.GroupAtType,
		AttachedInfo:     req.AttachedInfo,
		Ex:               req.Ex,
	}

	c.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterCreateSingleChatConversationsResp{}, after)
	return nil
}

func (c *conversationServer) webhookBeforeCreateGroupChatConversations(ctx context.Context, before *config.BeforeConfig, req *dbModel.Conversation) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &callbackstruct.CallbackBeforeCreateGroupChatConversationsReq{
			CallbackCommand:  callbackstruct.CallbackBeforeCreateGroupChatConversationsCommand,
			ConversationID:   req.ConversationID,
			ConversationType: req.ConversationType,
			GroupID:          req.GroupID,
			RecvMsgOpt:       req.RecvMsgOpt,
			IsPinned:         req.IsPinned,
			IsPrivateChat:    req.IsPrivateChat,
			BurnDuration:     req.BurnDuration,
			GroupAtType:      req.GroupAtType,
			AttachedInfo:     req.AttachedInfo,
			Ex:               req.Ex,
		}

		resp := &callbackstruct.CallbackBeforeCreateGroupChatConversationsResp{}

		if err := c.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		datautil.NotNilReplace(&req.RecvMsgOpt, resp.RecvMsgOpt)
		datautil.NotNilReplace(&req.IsPinned, resp.IsPinned)
		datautil.NotNilReplace(&req.IsPrivateChat, resp.IsPrivateChat)
		datautil.NotNilReplace(&req.BurnDuration, resp.BurnDuration)
		datautil.NotNilReplace(&req.GroupAtType, resp.GroupAtType)
		datautil.NotNilReplace(&req.AttachedInfo, resp.AttachedInfo)
		datautil.NotNilReplace(&req.Ex, resp.Ex)
		return nil
	})
}

func (c *conversationServer) webhookAfterCreateGroupChatConversations(ctx context.Context, after *config.AfterConfig, req *dbModel.Conversation) error {
	cbReq := &callbackstruct.CallbackAfterCreateGroupChatConversationsReq{
		CallbackCommand:  callbackstruct.CallbackAfterCreateGroupChatConversationsCommand,
		ConversationID:   req.ConversationID,
		ConversationType: req.ConversationType,
		GroupID:          req.GroupID,
		RecvMsgOpt:       req.RecvMsgOpt,
		IsPinned:         req.IsPinned,
		IsPrivateChat:    req.IsPrivateChat,
		BurnDuration:     req.BurnDuration,
		GroupAtType:      req.GroupAtType,
		AttachedInfo:     req.AttachedInfo,
		Ex:               req.Ex,
	}

	c.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterCreateGroupChatConversationsResp{}, after)
	return nil
}
