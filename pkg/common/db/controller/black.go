package controller

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type BlackDatabase interface {
	// Create 增加黑名单
	Create(ctx context.Context, blacks []*relation.BlackModel) (err error)
	// Delete 删除黑名单
	Delete(ctx context.Context, blacks []*relation.BlackModel) (err error)
	// FindOwnerBlacks 获取黑名单列表
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*relation.BlackModel, total int64, err error)
	FindBlackIDs(ctx context.Context, ownerUserID string) (blackIDs []string, err error)
	// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
	CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error)
}

type blackDatabase struct {
	black relation.BlackModelInterface
	cache cache.BlackCache
}

func NewBlackDatabase(black relation.BlackModelInterface, cache cache.BlackCache) BlackDatabase {
	return &blackDatabase{black, cache}
}

// Create 增加黑名单
func (b *blackDatabase) Create(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	if err := b.black.Create(ctx, blacks); err != nil {
		return err
	}
	return b.deleteBlackIDsCache(ctx, blacks)
}

// Delete 删除黑名单
func (b *blackDatabase) Delete(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	if err := b.black.Delete(ctx, blacks); err != nil {
		return err
	}
	return b.deleteBlackIDsCache(ctx, blacks)
}

func (b *blackDatabase) deleteBlackIDsCache(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	cache := b.cache.NewCache()
	for _, black := range blacks {
		cache = cache.DelBlackIDs(ctx, black.OwnerUserID)
	}
	return cache.ExecDel(ctx)
}

// FindOwnerBlacks 获取黑名单列表
func (b *blackDatabase) FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*relation.BlackModel, total int64, err error) {
	return b.black.FindOwnerBlacks(ctx, ownerUserID, pageNumber, showNumber)
}

// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
func (b *blackDatabase) CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error) {
	userID1BlackIDs, err := b.cache.GetBlackIDs(ctx, userID1)
	if err != nil {
		return
	}
	userID2BlackIDs, err := b.cache.GetBlackIDs(ctx, userID2)
	if err != nil {
		return
	}
	log.ZDebug(ctx, "blackIDs", "user1BlackIDs", userID1BlackIDs, "user2BlackIDs", userID2BlackIDs)
	return utils.IsContain(userID2, userID1BlackIDs), utils.IsContain(userID1, userID2BlackIDs), nil
}

func (b *blackDatabase) FindBlackIDs(ctx context.Context, ownerUserID string) (blackIDs []string, err error) {
	return b.cache.GetBlackIDs(ctx, ownerUserID)
}
