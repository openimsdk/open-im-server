package group

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
)

func (s *groupServer) GroupCreateCount(
	ctx context.Context,
	req *group.GroupCreateCountReq,
) (*group.GroupCreateCountResp, error) {
	if req.Start > req.End {
		return nil, errs.ErrArgs.Wrap("start > end")
	}
	total, err := s.GroupDatabase.CountTotal(ctx, nil)
	if err != nil {
		return nil, err
	}
	start := time.UnixMilli(req.Start)
	before, err := s.GroupDatabase.CountTotal(ctx, &start)
	if err != nil {
		return nil, err
	}
	count, err := s.GroupDatabase.CountRangeEverydayTotal(ctx, start, time.UnixMilli(req.End))
	if err != nil {
		return nil, err
	}
	return &group.GroupCreateCountResp{Total: total, Before: before, Count: count}, nil
}
