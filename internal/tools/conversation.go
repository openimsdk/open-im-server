package tools

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func (c *MsgTool) ConversationsDestructMsgs() {
	log.ZInfo(context.Background(), "start msg destruct cron task")
	ctx := mcontext.NewCtx(utils.GetSelfFuncName())
	conversations, err := c.conversationDatabase.GetConversationIDsNeedDestruct(ctx)
	if err != nil {
		log.ZError(ctx, "get conversation id need destruct failed", err)
		return
	}
	log.ZDebug(context.Background(), "nums conversations need destruct", "nums", len(conversations))
	for _, conversation := range conversations {
		log.ZDebug(ctx, "UserMsgsDestruct", "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID, "msgDestructTime", conversation.MsgDestructTime, "lastMsgDestructTime", conversation.LatestMsgDestructTime)
		seqs, err := c.msgDatabase.UserMsgsDestruct(ctx, conversation.OwnerUserID, conversation.ConversationID, conversation.MsgDestructTime, conversation.LatestMsgDestructTime)
		if err != nil {
			log.ZError(ctx, "user msg destruct failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
			continue
		}
		if err := c.conversationDatabase.UpdateUsersConversationFiled(ctx, []string{conversation.OwnerUserID}, conversation.ConversationID, map[string]interface{}{"latest_msg_destruct_time": time.Now()}); err != nil {
			log.ZError(ctx, "updateUsersConversationFiled failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
			continue
		}
		if len(seqs) > 0 {
			if err := c.msgNotificationSender.UserDeleteMsgsNotification(ctx, conversation.OwnerUserID, conversation.ConversationID, seqs); err != nil {
				log.ZError(ctx, "userDeleteMsgsNotification failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
			}
		}
	}
}
