package controller

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	"context"
	"errors"
	"gorm.io/gorm"
)

type BlackModel struct {
	db    *relation.Black
	cache *cache.GroupCache
}

func (b *BlackModel) Create(ctx context.Context, blacks []*relation.Black) (err error) {
	return b.db.Create(ctx, blacks)
}

func (b *BlackModel) Delete(ctx context.Context, blacks []*relation.Black) (err error) {
	return b.db.Delete(ctx, blacks)
}

func (b *BlackModel) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error) {
	return b.db.UpdateByMap(ctx, ownerUserID, blockUserID, args)
}

func (b *BlackModel) Update(ctx context.Context, blacks []*relation.Black) (err error) {
	return b.db.Update(ctx, blacks)
}

func (b *BlackModel) Find(ctx context.Context, blacks []*relation.Black) (blackList []*relation.Black, err error) {
	return b.db.Find(ctx, blacks)
}

func (b *BlackModel) Take(ctx context.Context, ownerUserID, blockUserID string) (black *relation.Black, err error) {
	return b.db.Take(ctx, ownerUserID, blockUserID)
}

func (b *BlackModel) FindByOwnerUserID(ctx context.Context, ownerUserID string) (blackList []*relation.Black, err error) {
	return b.db.FindByOwnerUserID(ctx, ownerUserID)
}

func (b *BlackModel) IsExist(ctx context.Context, ownerUserID, blockUserID string) (bool, error) {
	if _, err := b.Take(ctx, ownerUserID, blockUserID); err == nil {
		return true, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else {
		return false, err
	}
}
