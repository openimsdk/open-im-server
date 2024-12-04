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

package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

func NewBlackMongo(db *mongo.Database) (database.Black, error) {
	coll := db.Collection(database.BlackName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "owner_user_id", Value: 1},
			{Key: "block_user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &BlackMgo{coll: coll}, nil
}

type BlackMgo struct {
	coll *mongo.Collection
}

func (b *BlackMgo) blackFilter(ownerUserID, blockUserID string) bson.M {
	return bson.M{
		"owner_user_id": ownerUserID,
		"block_user_id": blockUserID,
	}
}

func (b *BlackMgo) blacksFilter(blacks []*model.Black) bson.M {
	if len(blacks) == 0 {
		return nil
	}
	or := make(bson.A, 0, len(blacks))
	for _, black := range blacks {
		or = append(or, b.blackFilter(black.OwnerUserID, black.BlockUserID))
	}
	return bson.M{"$or": or}
}

func (b *BlackMgo) Create(ctx context.Context, blacks []*model.Black) (err error) {
	return mongoutil.InsertMany(ctx, b.coll, blacks)
}

func (b *BlackMgo) Delete(ctx context.Context, blacks []*model.Black) (err error) {
	if len(blacks) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, b.coll, b.blacksFilter(blacks))
}

func (b *BlackMgo) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]any) (err error) {
	if len(args) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, b.coll, b.blackFilter(ownerUserID, blockUserID), bson.M{"$set": args}, false)
}

func (b *BlackMgo) Find(ctx context.Context, blacks []*model.Black) (blackList []*model.Black, err error) {
	return mongoutil.Find[*model.Black](ctx, b.coll, b.blacksFilter(blacks))
}

func (b *BlackMgo) Take(ctx context.Context, ownerUserID, blockUserID string) (black *model.Black, err error) {
	return mongoutil.FindOne[*model.Black](ctx, b.coll, b.blackFilter(ownerUserID, blockUserID))
}

func (b *BlackMgo) FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*model.Black, err error) {
	return mongoutil.FindPage[*model.Black](ctx, b.coll, bson.M{"owner_user_id": ownerUserID}, pagination)
}

func (b *BlackMgo) FindOwnerBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*model.Black, err error) {
	if len(userIDs) == 0 {
		return mongoutil.Find[*model.Black](ctx, b.coll, bson.M{"owner_user_id": ownerUserID})
	}
	return mongoutil.Find[*model.Black](ctx, b.coll, bson.M{"owner_user_id": ownerUserID, "block_user_id": bson.M{"$in": userIDs}})
}

func (b *BlackMgo) FindBlackUserIDs(ctx context.Context, ownerUserID string) (blackUserIDs []string, err error) {
	return mongoutil.Find[string](ctx, b.coll, bson.M{"owner_user_id": ownerUserID}, options.Find().SetProjection(bson.M{"_id": 0, "block_user_id": 1}))
}
