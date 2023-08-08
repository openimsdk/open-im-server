package tools

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/msgprocessor"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
)

func (c *MsgTool) ConvertTools() {
	ctx := mcontext.NewCtx("convert")
	conversationIDs, err := c.conversationDatabase.GetAllConversationIDs(ctx)
	if err != nil {
		log.ZError(ctx, "get all conversation ids failed", err)
		return
	}
	for _, conversationID := range conversationIDs {
		conversationIDs = append(conversationIDs, msgprocessor.GetNotificationConversationIDByConversationID(conversationID))
	}
	userIDs, err := c.userDatabase.GetAllUserID(ctx, 0, 0)
	if err != nil {
		log.ZError(ctx, "get all user ids failed", err)
		return
	}
	log.ZDebug(ctx, "all userIDs", "len userIDs", len(userIDs))
	for _, userID := range userIDs {
		conversationIDs = append(conversationIDs, msgprocessor.GetConversationIDBySessionType(constant.SingleChatType, userID, userID))
		conversationIDs = append(conversationIDs, msgprocessor.GetNotificationConversationID(constant.SingleChatType, userID, userID))
	}
	log.ZDebug(ctx, "all conversationIDs", "len userIDs", len(conversationIDs))
	c.msgDatabase.ConvertMsgsDocLen(ctx, conversationIDs)
}
