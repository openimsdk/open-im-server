package database

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type Group interface {
	Create(ctx context.Context, groups []*model.Group) (err error)
	UpdateMap(ctx context.Context, groupID string, args map[string]any) (err error)
	UpdateStatus(ctx context.Context, groupID string, status int32) (err error)
	Find(ctx context.Context, groupIDs []string) (groups []*model.Group, err error)
	Take(ctx context.Context, groupID string) (group *model.Group, err error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (total int64, groups []*model.Group, err error)
	// Get Group total quantity
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
	// Get Group total quantity every day
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)

	FindJoinSortGroupID(ctx context.Context, groupIDs []string) ([]string, error)

	SearchJoin(ctx context.Context, groupIDs []string, keyword string, pagination pagination.Pagination) (int64, []*model.Group, error)
}
