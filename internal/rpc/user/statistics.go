package user

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	pbuser "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"time"
)

func (s *userServer) UserRegisterCount(ctx context.Context, req *pbuser.UserRegisterCountReq) (*pbuser.UserRegisterCountResp, error) {
	if req.Start > req.End {
		return nil, errs.ErrArgs.Wrap("start > end")
	}
	total, err := s.CountTotal(ctx)
	if err != nil {
		return nil, err
	}
	count, err := s.CountRangeEverydayTotal(ctx, time.UnixMilli(req.Start), time.UnixMilli(req.End))
	if err != nil {
		return nil, err
	}
	return &pbuser.UserRegisterCountResp{Total: total, Count: count}, nil
}
