package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
)

type ClientConfig interface {
	Set(ctx context.Context, userID string, config map[string]string) error
	Get(ctx context.Context, userID string) (map[string]string, error)
	Del(ctx context.Context, userID string, keys []string) error
	GetPage(ctx context.Context, userID string, key string, pagination pagination.Pagination) (int64, []*model.ClientConfig, error)
}
