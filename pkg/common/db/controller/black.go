// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"

	"github.com/OpenIMSDK/tools/pagination"

	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

type BlackDatabase interface {
	// Create add BlackList
	Create(ctx context.Context, blacks []*relation.BlackModel) (err error)
	// Delete delete BlackList
	Delete(ctx context.Context, blacks []*relation.BlackModel) (err error)
	// FindOwnerBlacks get BlackList list
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*relation.BlackModel, err error)
	FindBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*relation.BlackModel, err error)
	// CheckIn Check whether user2 is in the black list of user1 (inUser1Blacks==true) Check whether user1 is in the black list of user2 (inUser2Blacks==true)
	CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error)
}

type blackDatabase struct {
	black relation.BlackModelInterface
	cache cache.BlackCache
}

func NewBlackDatabase(black relation.BlackModelInterface, cache cache.BlackCache) BlackDatabase {
	return &blackDatabase{black, cache}
}

// Create Add Blacklist.
func (b *blackDatabase) Create(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	if err := b.black.Create(ctx, blacks); err != nil {
		return err
	}
	return b.deleteBlackIDsCache(ctx, blacks)
}

// Delete Delete Blacklist.
func (b *blackDatabase) Delete(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	if err := b.black.Delete(ctx, blacks); err != nil {
		return err
	}
	return b.deleteBlackIDsCache(ctx, blacks)
}

// FindOwnerBlacks Get Blacklist List.
func (b *blackDatabase) deleteBlackIDsCache(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	cache := b.cache.NewCache()
	for _, black := range blacks {
		cache = cache.DelBlackIDs(ctx, black.OwnerUserID)
	}
	return cache.ExecDel(ctx)
}

// FindOwnerBlacks Get Blacklist List.
func (b *blackDatabase) FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*relation.BlackModel, err error) {
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
	return utils.IsContain(userID2, userID1BlackIDs), utils.IsContain(userID1, userID2BlackIDs), nil
}

// FindBlackIDs Get Blacklist List.
func (b *blackDatabase) FindBlackIDs(ctx context.Context, ownerUserID string) (blackIDs []string, err error) {
	return b.cache.GetBlackIDs(ctx, ownerUserID)
}

// FindBlackInfos Get Blacklist List.
func (b *blackDatabase) FindBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*relation.BlackModel, err error) {
	return b.black.FindOwnerBlackInfos(ctx, ownerUserID, userIDs)
}
