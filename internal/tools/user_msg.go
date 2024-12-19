package tools

import (
	"fmt"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"os"
	"time"
)

func (c *cronServer) clearUserMsg() {
	now := time.Now()
	operationID := fmt.Sprintf("cron_user_msg_%d_%d", os.Getpid(), now.UnixMilli())
	ctx := mcontext.SetOperationID(c.ctx, operationID)
	log.ZDebug(ctx, "clear user msg cron start")
	const (
		deleteCount = 200
		deleteLimit = 100
	)
	var count int
	for i := 1; i <= deleteCount; i++ {
		resp, err := c.conversationClient.ClearUserConversationMsg(ctx, &pbconversation.ClearUserConversationMsgReq{Timestamp: now.UnixMilli(), Limit: deleteLimit})
		if err != nil {
			log.ZError(ctx, "ClearUserConversationMsg failed.", err)
			return
		}
		count += int(resp.Count)
		if resp.Count < deleteLimit {
			break
		}
	}
	log.ZDebug(ctx, "clear user msg cron task completed", "cont", time.Since(now), "count", count)
}
