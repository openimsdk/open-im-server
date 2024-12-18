package tools

import (
	"fmt"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"os"
	"time"
)

func (c *cronServer) clearS3() {
	start := time.Now()
	executeNum := 10
	// number of pagination. if need modify, need update value in third.DeleteOutdatedData
	pageShowNumber := 500
	deleteTime := start.Add(-time.Hour * 24 * time.Duration(c.config.CronTask.FileExpireTime))
	operationID := fmt.Sprintf("cron_%d_%d", os.Getpid(), deleteTime.UnixMilli())
	ctx := mcontext.SetOperationID(c.ctx, operationID)
	log.ZDebug(ctx, "deleteoutDatedData", "deletetime", deleteTime, "timestamp", deleteTime.UnixMilli())
	for i := 1; i <= executeNum; i++ {
		ctx := mcontext.SetOperationID(c.ctx, fmt.Sprintf("%s_%d", operationID, i))
		resp, err := c.thirdClient.DeleteOutdatedData(ctx, &third.DeleteOutdatedDataReq{ExpireTime: deleteTime.UnixMilli(), ObjectGroup: c.config.CronTask.DeleteObjectType, Count: int32(pageShowNumber)})
		if err != nil {
			log.ZError(ctx, "cron deleteoutDatedData failed", err, "deleteTime", deleteTime, "cont", time.Since(start))
			return
		}
		if resp.Count < int32(pageShowNumber) {
			break
		}
	}

	log.ZDebug(ctx, "cron deleteoutDatedData success", "deltime", deleteTime, "cont", time.Since(start))
}
