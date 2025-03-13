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
			CallbackCommand:       callbackstruct.CallbackBeforeCreateSingleChatConversationsCommand,
			OwnerUserID:           req.OwnerUserID,
			ConversationID:        req.ConversationID,
			ConversationType:      req.ConversationType,
			UserID:                req.UserID,
			GroupID:               req.GroupID,
			RecvMsgOpt:            req.RecvMsgOpt,
			IsPinned:              req.IsPinned,
			IsPrivateChat:         req.IsPrivateChat,
			BurnDuration:          req.BurnDuration,
			GroupAtType:           req.GroupAtType,
			AttachedInfo:          req.AttachedInfo,
			Ex:                    req.Ex,
			MaxSeq:                req.MaxSeq,
			MinSeq:                req.MinSeq,
			CreateTime:            req.CreateTime,
			IsMsgDestruct:         req.IsMsgDestruct,
			MsgDestructTime:       req.MsgDestructTime,
			LatestMsgDestructTime: req.LatestMsgDestructTime,
		}

		resp := &callbackstruct.CallbackBeforeCreateSingleChatConversationsResp{}

		if err := c.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		datautil.NotNilReplace(&req.OwnerUserID, resp.OwnerUserID)
		datautil.NotNilReplace(&req.ConversationID, resp.ConversationID)
		datautil.NotNilReplace(&req.ConversationType, resp.ConversationType)
		datautil.NotNilReplace(&req.UserID, resp.UserID)
		datautil.NotNilReplace(&req.GroupID, resp.GroupID)
		datautil.NotNilReplace(&req.RecvMsgOpt, resp.RecvMsgOpt)
		datautil.NotNilReplace(&req.IsPinned, resp.IsPinned)
		datautil.NotNilReplace(&req.IsPrivateChat, resp.IsPrivateChat)
		datautil.NotNilReplace(&req.BurnDuration, resp.BurnDuration)
		datautil.NotNilReplace(&req.GroupAtType, resp.GroupAtType)
		datautil.NotNilReplace(&req.AttachedInfo, resp.AttachedInfo)
		datautil.NotNilReplace(&req.Ex, resp.Ex)
		datautil.NotNilReplace(&req.MaxSeq, resp.MaxSeq)
		datautil.NotNilReplace(&req.MinSeq, resp.MinSeq)
		datautil.NotNilReplace(&req.CreateTime, resp.CreateTime)
		datautil.NotNilReplace(&req.IsMsgDestruct, resp.IsMsgDestruct)
		datautil.NotNilReplace(&req.MsgDestructTime, resp.MsgDestructTime)
		datautil.NotNilReplace(&req.LatestMsgDestructTime, resp.LatestMsgDestructTime)

		return nil
	})
}

func (c *conversationServer) webhookAfterCreateSingleChatConversations(ctx context.Context, after *config.AfterConfig, req *dbModel.Conversation) error {
	cbReq := &callbackstruct.CallbackAfterCreateSingleChatConversationsReq{
		CallbackCommand:       callbackstruct.CallbackAfterCreateSingleChatConversationsCommand,
		OwnerUserID:           req.OwnerUserID,
		ConversationID:        req.ConversationID,
		ConversationType:      req.ConversationType,
		UserID:                req.UserID,
		GroupID:               req.GroupID,
		RecvMsgOpt:            req.RecvMsgOpt,
		IsPinned:              req.IsPinned,
		IsPrivateChat:         req.IsPrivateChat,
		BurnDuration:          req.BurnDuration,
		GroupAtType:           req.GroupAtType,
		AttachedInfo:          req.AttachedInfo,
		Ex:                    req.Ex,
		MaxSeq:                req.MaxSeq,
		MinSeq:                req.MinSeq,
		CreateTime:            req.CreateTime,
		IsMsgDestruct:         req.IsMsgDestruct,
		MsgDestructTime:       req.MsgDestructTime,
		LatestMsgDestructTime: req.LatestMsgDestructTime,
	}

	c.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterCreateSingleChatConversationsResp{}, after)
	return nil
}

func (c *conversationServer) webhookBeforeCreateGroupChatConversations(ctx context.Context, before *config.BeforeConfig, req *dbModel.Conversation) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &callbackstruct.CallbackBeforeCreateGroupChatConversationsReq{
			CallbackCommand:       callbackstruct.CallbackBeforeCreateGroupChatConversationsCommand,
			ConversationID:        req.ConversationID,
			ConversationType:      req.ConversationType,
			GroupID:               req.GroupID,
			RecvMsgOpt:            req.RecvMsgOpt,
			IsPinned:              req.IsPinned,
			IsPrivateChat:         req.IsPrivateChat,
			BurnDuration:          req.BurnDuration,
			GroupAtType:           req.GroupAtType,
			AttachedInfo:          req.AttachedInfo,
			Ex:                    req.Ex,
			MaxSeq:                req.MaxSeq,
			MinSeq:                req.MinSeq,
			CreateTime:            req.CreateTime,
			IsMsgDestruct:         req.IsMsgDestruct,
			MsgDestructTime:       req.MsgDestructTime,
			LatestMsgDestructTime: req.LatestMsgDestructTime,
		}

		resp := &callbackstruct.CallbackBeforeCreateGroupChatConversationsResp{}

		if err := c.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}

		datautil.NotNilReplace(&req.ConversationID, resp.ConversationID)
		datautil.NotNilReplace(&req.ConversationType, resp.ConversationType)
		datautil.NotNilReplace(&req.GroupID, resp.GroupID)
		datautil.NotNilReplace(&req.RecvMsgOpt, resp.RecvMsgOpt)
		datautil.NotNilReplace(&req.IsPinned, resp.IsPinned)
		datautil.NotNilReplace(&req.IsPrivateChat, resp.IsPrivateChat)
		datautil.NotNilReplace(&req.BurnDuration, resp.BurnDuration)
		datautil.NotNilReplace(&req.GroupAtType, resp.GroupAtType)
		datautil.NotNilReplace(&req.AttachedInfo, resp.AttachedInfo)
		datautil.NotNilReplace(&req.Ex, resp.Ex)
		datautil.NotNilReplace(&req.MaxSeq, resp.MaxSeq)
		datautil.NotNilReplace(&req.MinSeq, resp.MinSeq)
		datautil.NotNilReplace(&req.CreateTime, resp.CreateTime)
		datautil.NotNilReplace(&req.IsMsgDestruct, resp.IsMsgDestruct)
		datautil.NotNilReplace(&req.MsgDestructTime, resp.MsgDestructTime)
		datautil.NotNilReplace(&req.LatestMsgDestructTime, resp.LatestMsgDestructTime)
		return nil
	})
}

func (c *conversationServer) webhookAfterCreateGroupChatConversations(ctx context.Context, after *config.AfterConfig, req *dbModel.Conversation) error {
	cbReq := &callbackstruct.CallbackAfterCreateGroupChatConversationsReq{
		CallbackCommand:       callbackstruct.CallbackAfterCreateGroupChatConversationsCommand,
		ConversationID:        req.ConversationID,
		ConversationType:      req.ConversationType,
		GroupID:               req.GroupID,
		RecvMsgOpt:            req.RecvMsgOpt,
		IsPinned:              req.IsPinned,
		IsPrivateChat:         req.IsPrivateChat,
		BurnDuration:          req.BurnDuration,
		GroupAtType:           req.GroupAtType,
		AttachedInfo:          req.AttachedInfo,
		Ex:                    req.Ex,
		MaxSeq:                req.MaxSeq,
		MinSeq:                req.MinSeq,
		CreateTime:            req.CreateTime,
		IsMsgDestruct:         req.IsMsgDestruct,
		MsgDestructTime:       req.MsgDestructTime,
		LatestMsgDestructTime: req.LatestMsgDestructTime,
	}

	c.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &callbackstruct.CallbackAfterCreateGroupChatConversationsResp{}, after)
	return nil
}
