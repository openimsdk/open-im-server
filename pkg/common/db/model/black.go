package model

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/mysql"
	"context"
)

type BlackModel struct {
	db    *mysql.Black
	cache *cache.GroupCache
}

func (b *BlackModel) Create(ctx context.Context, blacks []*mysql.Black) (err error) {
	return b.db.Create(ctx, blacks)
}

func (b *BlackModel) Delete(ctx context.Context, blacks []*mysql.Black) (err error) {
	return b.db.Delete(ctx, blacks)
}

func (b *BlackModel) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error) {
	return b.db.UpdateByMap(ctx, ownerUserID, blockUserID, args)
}

func (b *BlackModel) Update(ctx context.Context, blacks []*mysql.Black) (err error) {
	return b.db.Update(ctx, blacks)
}

func (b *BlackModel) Find(ctx context.Context, blacks []*mysql.Black) (blackList []*mysql.Black, err error) {
	return b.db.Find(ctx, blacks)
}

func (b *BlackModel) Take(ctx context.Context, blackID string) (black *mysql.Black, err error) {
	return b.db.Take(ctx, blackID)
}

func (b *BlackModel) FindByOwnerUserID(ctx context.Context, ownerUserID string) (blackList []*mysql.Black, err error) {
	return b.db.FindByOwnerUserID(ctx, ownerUserID)
}
