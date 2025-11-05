package database

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
	"time"
)

type Log interface {
	Create(ctx context.Context, log []*model.Log) error
	Search(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*model.Log, error)
	Delete(ctx context.Context, logID []string, userID string) error
	Get(ctx context.Context, logIDs []string, userID string) ([]*model.Log, error)
}
