package cron

import (
	"fmt"
	"os"
	"time"

	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

func (c *cronServer) clearS3() {
	start := time.Now()
	deleteTime := start.Add(-time.Hour * 24 * time.Duration(c.config.CronTask.FileExpireTime))
	operationID := fmt.Sprintf("cron_s3_%d_%d", os.Getpid(), deleteTime.UnixMilli())
	ctx := mcontext.SetOperationID(c.ctx, operationID)
	log.ZDebug(ctx, "deleteoutDatedData", "deletetime", deleteTime, "timestamp", deleteTime.UnixMilli())
	const (
		deleteCount = 10000
		deleteLimit = 100
	)

	var count int
	for i := 1; i <= deleteCount; i++ {
		resp, err := c.thirdClient.DeleteOutdatedData(ctx, &third.DeleteOutdatedDataReq{ExpireTime: deleteTime.UnixMilli(), ObjectGroup: c.config.CronTask.DeleteObjectType, Limit: deleteLimit})
		if err != nil {
			log.ZError(ctx, "cron deleteoutDatedData failed", err)
			return
		}
		count += int(resp.Count)
		if resp.Count < deleteLimit {
			break
		}
	}
	log.ZDebug(ctx, "cron deleteoutDatedData success", "deltime", deleteTime, "cont", time.Since(start), "count", count)
}
