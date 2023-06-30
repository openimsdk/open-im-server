package third

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
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

func (t *thirdServer) GetUrl(ctx context.Context, req *third.GetUrlReq) (*third.GetUrlResp, error) {
	if req.Expires <= 0 {
		if err := tokenverify.CheckAdmin(ctx); err != nil {
			return nil, err
		}
	}
	return t.s3dataBase.GetUrl(ctx, req)
}

func (t *thirdServer) GetHashInfo(ctx context.Context, req *third.GetHashInfoReq) (*third.GetHashInfoResp, error) {
	return t.s3dataBase.GetHashInfo(ctx, req)
}

func (t *thirdServer) CleanObject(ctx context.Context, now time.Time) {
	t.s3dataBase.CleanExpirationObject(ctx, now)
}
