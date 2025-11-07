package controller

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

type BlackDatabase interface {
	// Create add BlackList
	Create(ctx context.Context, blacks []*model.Black) (err error)
	// Delete delete BlackList
	Delete(ctx context.Context, blacks []*model.Black) (err error)
	// FindOwnerBlacks get BlackList list
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*model.Black, err error)
	FindBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*model.Black, err error)
	// CheckIn Check whether user2 is in the black list of user1 (inUser1Blacks==true) Check whether user1 is in the black list of user2 (inUser2Blacks==true)
	CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error)
}

type blackDatabase struct {
	black database.Black
	cache cache.BlackCache
}

func NewBlackDatabase(black database.Black, cache cache.BlackCache) BlackDatabase {
	return &blackDatabase{black, cache}
}

// Create Add Blacklist.
func (b *blackDatabase) Create(ctx context.Context, blacks []*model.Black) (err error) {
	if err := b.black.Create(ctx, blacks); err != nil {
		return err
	}
	return b.deleteBlackIDsCache(ctx, blacks)
}

// Delete Delete Blacklist.
func (b *blackDatabase) Delete(ctx context.Context, blacks []*model.Black) (err error) {
	if err := b.black.Delete(ctx, blacks); err != nil {
		return err
	}
	return b.deleteBlackIDsCache(ctx, blacks)
}

// FindOwnerBlacks Get Blacklist List.
func (b *blackDatabase) deleteBlackIDsCache(ctx context.Context, blacks []*model.Black) (err error) {
	cache := b.cache.CloneBlackCache()
	for _, black := range blacks {
		cache = cache.DelBlackIDs(ctx, black.OwnerUserID)
	}
	return cache.ChainExecDel(ctx)
}

// FindOwnerBlacks Get Blacklist List.
func (b *blackDatabase) FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*model.Black, err error) {
	return b.black.FindOwnerBlacks(ctx, ownerUserID, pagination)
}

// FindOwnerBlacks Get Blacklist List.
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
	return datautil.Contain(userID2, userID1BlackIDs...), datautil.Contain(userID1, userID2BlackIDs...), nil
}

// FindBlackIDs Get Blacklist List.
func (b *blackDatabase) FindBlackIDs(ctx context.Context, ownerUserID string) (blackIDs []string, err error) {
	return b.cache.GetBlackIDs(ctx, ownerUserID)
}

// FindBlackInfos Get Blacklist List.
func (b *blackDatabase) FindBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*model.Black, err error) {
	return b.black.FindOwnerBlackInfos(ctx, ownerUserID, userIDs)
}
