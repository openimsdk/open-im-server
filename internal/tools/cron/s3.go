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

//	var req *third.DeleteOutdatedDataReq
//	count1, err := ExtractField(ctx, c.thirdClient.DeleteOutdatedData, req, (*third.DeleteOutdatedDataResp).GetCount)
//
//	c.thirdClient.DeleteOutdatedData(ctx, &third.DeleteOutdatedDataReq{})
//	msggateway.GetUsersOnlineStatusCaller.Invoke(ctx, &msggateway.GetUsersOnlineStatusReq{})
//
//	var cli ThirdClient
//
//	c111, err := cli.DeleteOutdatedData(ctx, 100)
//
//	cli.ThirdClient.DeleteOutdatedData(ctx, &third.DeleteOutdatedDataReq{})
//
//	cli.AuthSign(ctx, &third.AuthSignReq{})
//
//	cli.SetAppBadge()
//
//}
//
//func extractField[A, B, C any](ctx context.Context, fn func(ctx context.Context, req *A, opts ...grpc.CallOption) (*B, error), req *A, get func(*B) C) (C, error) {
//	resp, err := fn(ctx, req)
//	if err != nil {
//		var c C
//		return c, err
//	}
//	return get(resp), nil
//}
//
//func ignore(_ any, err error) error {
//	return err
//}
//
//type ThirdClient struct {
//	third.ThirdClient
//}
//
//func (c *ThirdClient) DeleteOutdatedData(ctx context.Context, expireTime int64) (int32, error) {
//	return extractField(ctx, c.ThirdClient.DeleteOutdatedData, &third.DeleteOutdatedDataReq{ExpireTime: expireTime}, (*third.DeleteOutdatedDataResp).GetCount)
//}
//
//func (c *ThirdClient) DeleteOutdatedData1(ctx context.Context, expireTime int64) error {
//	return ignore(c.ThirdClient.DeleteOutdatedData(ctx, &third.DeleteOutdatedDataReq{ExpireTime: expireTime}))
//}
