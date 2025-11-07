package group

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/tools/errs"
)

func (g *groupServer) GroupCreateCount(ctx context.Context, req *group.GroupCreateCountReq) (*group.GroupCreateCountResp, error) {
	if req.Start > req.End {
		return nil, errs.ErrArgs.WrapMsg("start > end: %d > %d", req.Start, req.End)
	}
	if err := authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	total, err := g.db.CountTotal(ctx, nil)
	if err != nil {
		return nil, err
	}
	start := time.UnixMilli(req.Start)
	before, err := g.db.CountTotal(ctx, &start)
	if err != nil {
		return nil, err
	}
	count, err := g.db.CountRangeEverydayTotal(ctx, start, time.UnixMilli(req.End))
	if err != nil {
		return nil, err
	}
	return &group.GroupCreateCountResp{Total: total, Before: before, Count: count}, nil
}
