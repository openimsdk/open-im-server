package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/db/tx"
)

type ClientConfigDatabase interface {
	SetUserConfig(ctx context.Context, userID string, config map[string]string) error
	GetUserConfig(ctx context.Context, userID string) (map[string]string, error)
	DelUserConfig(ctx context.Context, userID string, keys []string) error
	GetUserConfigPage(ctx context.Context, userID string, key string, pagination pagination.Pagination) (int64, []*model.ClientConfig, error)
}

func NewClientConfigDatabase(db database.ClientConfig, cache cache.ClientConfigCache, tx tx.Tx) ClientConfigDatabase {
	return &clientConfigDatabase{
		tx:    tx,
		db:    db,
		cache: cache,
	}
}

type clientConfigDatabase struct {
	tx    tx.Tx
	db    database.ClientConfig
	cache cache.ClientConfigCache
}

func (x *clientConfigDatabase) SetUserConfig(ctx context.Context, userID string, config map[string]string) error {
	return x.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := x.db.Set(ctx, userID, config); err != nil {
			return err
		}
		return x.cache.DeleteUserCache(ctx, []string{userID})
	})
}

func (x *clientConfigDatabase) GetUserConfig(ctx context.Context, userID string) (map[string]string, error) {
	return x.cache.GetUserConfig(ctx, userID)
}

func (x *clientConfigDatabase) DelUserConfig(ctx context.Context, userID string, keys []string) error {
	return x.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := x.db.Del(ctx, userID, keys); err != nil {
			return err
		}
		return x.cache.DeleteUserCache(ctx, []string{userID})
	})
}

func (x *clientConfigDatabase) GetUserConfigPage(ctx context.Context, userID string, key string, pagination pagination.Pagination) (int64, []*model.ClientConfig, error) {
	return x.db.GetPage(ctx, userID, key, pagination)
}
