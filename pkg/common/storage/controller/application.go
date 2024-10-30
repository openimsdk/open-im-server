package controller

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationDatabase interface {
	LatestVersion(ctx context.Context, platform string) (*model.Application, error)
	AddVersion(ctx context.Context, val *model.Application) error
	UpdateVersion(ctx context.Context, id primitive.ObjectID, update map[string]any) error
	DeleteVersion(ctx context.Context, id []primitive.ObjectID) error
	PageVersion(ctx context.Context, platforms []string, page pagination.Pagination) (int64, []*model.Application, error)
}

func NewApplicationDatabase(db database.Application, cache cache.ApplicationCache) ApplicationDatabase {
	return &applicationDatabase{db: db, cache: cache}
}

type applicationDatabase struct {
	db    database.Application
	cache cache.ApplicationCache
}

func (a *applicationDatabase) LatestVersion(ctx context.Context, platform string) (*model.Application, error) {
	return a.cache.LatestVersion(ctx, platform)
}

func (a *applicationDatabase) AddVersion(ctx context.Context, val *model.Application) error {
	if err := a.db.AddVersion(ctx, val); err != nil {
		return err
	}
	return a.cache.DeleteCache(ctx, []string{val.Platform})
}

func (a *applicationDatabase) UpdateVersion(ctx context.Context, id primitive.ObjectID, update map[string]any) error {
	platforms, err := a.db.FindPlatform(ctx, []primitive.ObjectID{id})
	if err != nil {
		return err
	}
	if err := a.db.UpdateVersion(ctx, id, update); err != nil {
		return err
	}
	if p, ok := update["platform"]; ok {
		if val, ok := p.(string); ok {
			platforms = append(platforms, val)
		}
	}
	return a.cache.DeleteCache(ctx, platforms)
}

func (a *applicationDatabase) DeleteVersion(ctx context.Context, id []primitive.ObjectID) error {
	platforms, err := a.db.FindPlatform(ctx, id)
	if err != nil {
		return err
	}
	if err := a.db.DeleteVersion(ctx, id); err != nil {
		return err
	}
	return a.cache.DeleteCache(ctx, platforms)
}

func (a *applicationDatabase) PageVersion(ctx context.Context, platforms []string, page pagination.Pagination) (int64, []*model.Application, error) {
	return a.db.PageVersion(ctx, platforms, page)
}
