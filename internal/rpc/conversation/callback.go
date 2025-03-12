package conversation

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/utils/datautil"
)

func (c *conversationServer) webhookBeforeCreateSingleChatConversations(ctx context.Context, before *config.BeforeConfig, req *pbconversation.CreateSingleChatConversationsReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &callbackstruct.CallbackBeforeCreateSingleChatConversationsReq{
			CallbackCommand:  callbackstruct.CallbackBeforeCreateSingleChatConversationsCommand,
			RecvID:           req.RecvID,
			SendID:           req.SendID,
			ConversationID:   req.ConversationID,
			ConversationType: req.ConversationType,
		}

		resp := &callbackstruct.CallbackBeforeCreateSingleChatConversationsResp{}

		if err := c.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		datautil.NotNilReplace(&req.RecvID, resp.RecvID)
		datautil.NotNilReplace(&req.SendID, resp.SendID)
		datautil.NotNilReplace(&req.ConversationID, resp.ConversationID)
		datautil.NotNilReplace(&req.ConversationType, resp.ConversationType)
		return nil
	})
}

func (c *conversationServer) webhookAfterCreateSingleChatConversations(ctx context.Context, after *config.AfterConfig, req *pbconversation.CreateSingleChatConversationsReq) error {
	cbReq := &callbackstruct.CallbackAfterCreateSingleChatConversationsReq{
		CallbackCommand:  callbackstruct.CallbackAfterCreateSingleChatConversationsCommand,
		RecvID:           req.RecvID,
		SendID:           req.SendID,
		ConversationID:   req.ConversationID,
		ConversationType: req.ConversationType,
	}

	c.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterCreateSingleChatConversationsResp{}, after)
	return nil
}


func (c *conversationServer) webhookBeforeCreateGroupChatConversations(ctx context.Context, before *config.BeforeConfig, req *pbconversation.CreateGroupChatConversationsReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &callbackstruct.CallbackBeforeCreateGroupChatConversationsReq{
			CallbackCommand: callbackstruct.CallbackBeforeCreateGroupChatConversationsCommand,
			UserIDs:         req.UserIDs,
			GroupID:         req.GroupID,
		}

		resp := &callbackstruct.CallbackBeforeCreateGroupChatConversationsResp{}

		if err := c.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		datautil.NotNilReplace(&req.UserIDs, resp.UserIDs)
		datautil.NotNilReplace(&req.GroupID, resp.GroupID)
		return nil
	})
}

func (c *conversationServer) webhookAfterCreateGroupChatConversations(ctx context.Context, after *config.AfterConfig, req *pbconversation.CreateGroupChatConversationsReq) error {
	cbReq := &callbackstruct.CallbackAfterCreateGroupChatConversationsReq{
		CallbackCommand: callbackstruct.CallbackAfterCreateGroupChatConversationsCommand,
		UserIDs:         req.UserIDs,
		GroupID:         req.GroupID,
	}

	c.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterCreateGroupChatConversationsResp{}, after)
	return nil
}