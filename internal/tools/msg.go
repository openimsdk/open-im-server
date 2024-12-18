package tools

import (
	"fmt"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"os"
	"time"
)

func (c *cronServer) deleteMsg() {
	now := time.Now()
	deltime := now.Add(-time.Hour * 24 * time.Duration(c.config.CronTask.RetainChatRecords))
	ctx := mcontext.SetOperationID(c.ctx, fmt.Sprintf("cron_%d_%d", os.Getpid(), deltime.UnixMilli()))
	log.ZDebug(ctx, "Destruct chat records", "deltime", deltime, "timestamp", deltime.UnixMilli())

	if _, err := c.msgClient.DestructMsgs(ctx, &msg.DestructMsgsReq{Timestamp: deltime.UnixMilli()}); err != nil {
		log.ZError(ctx, "cron destruct chat records failed", err, "deltime", deltime, "cont", time.Since(now))
		return
	}
	log.ZDebug(ctx, "cron destruct chat records success", "deltime", deltime, "cont", time.Since(now))
}
