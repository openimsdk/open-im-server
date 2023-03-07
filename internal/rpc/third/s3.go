package third

import (
	"OpenIM/pkg/proto/third"
	"context"
	"time"
)

func (t *thirdServer) ApplyPut(ctx context.Context, req *third.ApplyPutReq) (*third.ApplyPutResp, error) {
	return t.s3dataBase.ApplyPut(ctx, req)
}

func (t *thirdServer) GetPut(ctx context.Context, req *third.GetPutReq) (*third.GetPutResp, error) {
	return t.s3dataBase.GetPut(ctx, req)
}

func (t *thirdServer) ConfirmPut(ctx context.Context, req *third.ConfirmPutReq) (*third.ConfirmPutResp, error) {
	return t.s3dataBase.ConfirmPut(ctx, req)
}

func (t *thirdServer) CleanObject(ctx context.Context, now time.Time) {
	t.s3dataBase.CleanExpirationObject(ctx, now)
}
