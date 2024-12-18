package tools

import (
	"fmt"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"os"
	"time"
)

func (c *cronServer) clearUserMsg() {
	now := time.Now()
	ctx := mcontext.SetOperationID(c.ctx, fmt.Sprintf("cron_%d_%d", os.Getpid(), now.UnixMilli()))
	log.ZDebug(ctx, "clear msg cron start", "now", now)

	conversations, err := c.conversationClient.GetConversationsNeedClearMsg(ctx, &pbconversation.GetConversationsNeedClearMsgReq{})
	if err != nil {
		log.ZError(ctx, "Get conversation need Destruct msgs failed.", err)
		return
	}

	_, err = c.msgClient.ClearMsg(ctx, &msg.ClearMsgReq{Conversations: conversations.Conversations})
	if err != nil {
		log.ZError(ctx, "Clear Msg failed.", err)
		return
	}

	log.ZDebug(ctx, "clear msg cron task completed", "cont", time.Since(now))
}
