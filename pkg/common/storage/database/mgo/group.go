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
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"time"

	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewGroupMongo(db *mongo.Database) (database.Group, error) {
	coll := db.Collection(database.GroupName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "group_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &GroupMgo{coll: coll}, nil
}

type GroupMgo struct {
	coll *mongo.Collection
}

func (g *GroupMgo) sortGroup() any {
	return bson.D{{"group_name", 1}, {"create_time", 1}}
}

func (g *GroupMgo) Create(ctx context.Context, groups []*model.Group) (err error) {
	return mongoutil.InsertMany(ctx, g.coll, groups)
}

func (g *GroupMgo) UpdateStatus(ctx context.Context, groupID string, status int32) (err error) {
	return g.UpdateMap(ctx, groupID, map[string]any{"status": status})
}

func (g *GroupMgo) UpdateMap(ctx context.Context, groupID string, args map[string]any) (err error) {
	if len(args) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID}, bson.M{"$set": args}, true)
}

func (g *GroupMgo) Find(ctx context.Context, groupIDs []string) (groups []*model.Group, err error) {
	return mongoutil.Find[*model.Group](ctx, g.coll, bson.M{"group_id": bson.M{"$in": groupIDs}})
}

func (g *GroupMgo) Take(ctx context.Context, groupID string) (group *model.Group, err error) {
	return mongoutil.FindOne[*model.Group](ctx, g.coll, bson.M{"group_id": groupID})
}

func (g *GroupMgo) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (total int64, groups []*model.Group, err error) {
	// Define the sorting options
	opts := options.Find().SetSort(bson.D{{Key: "create_time", Value: -1}})

	// Perform the search with pagination and sorting
	return mongoutil.FindPage[*model.Group](ctx, g.coll, bson.M{
		"group_name": bson.M{"$regex": keyword},
		"status":     bson.M{"$ne": constant.GroupStatusDismissed},
	}, pagination, opts)
}

func (g *GroupMgo) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	if before == nil {
		return mongoutil.Count(ctx, g.coll, bson.M{})
	}
	return mongoutil.Count(ctx, g.coll, bson.M{"create_time": bson.M{"$lt": before}})
}

func (g *GroupMgo) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"create_time": bson.M{
					"$gte": start,
					"$lt":  end,
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$create_time",
					},
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}
	type Item struct {
		Date  string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	items, err := mongoutil.Aggregate[Item](ctx, g.coll, pipeline)
	if err != nil {
		return nil, err
	}
	res := make(map[string]int64, len(items))
	for _, item := range items {
		res[item.Date] = item.Count
	}
	return res, nil
}

func (g *GroupMgo) FindJoinSortGroupID(ctx context.Context, groupIDs []string) ([]string, error) {
	if len(groupIDs) < 2 {
		return groupIDs, nil
	}
	filter := bson.M{
		"group_id": bson.M{"$in": groupIDs},
		"status":   bson.M{"$ne": constant.GroupStatusDismissed},
	}
	opt := options.Find().SetSort(g.sortGroup()).SetProjection(bson.M{"_id": 0, "group_id": 1})
	return mongoutil.Find[string](ctx, g.coll, filter, opt)
}

func (g *GroupMgo) SearchJoin(ctx context.Context, groupIDs []string, keyword string, pagination pagination.Pagination) (int64, []*model.Group, error) {
	if len(groupIDs) == 0 {
		return 0, nil, nil
	}
	filter := bson.M{
		"group_id": bson.M{"$in": groupIDs},
		"status":   bson.M{"$ne": constant.GroupStatusDismissed},
	}
	if keyword != "" {
		filter["group_name"] = bson.M{"$regex": keyword}
	}
	// Define the sorting options
	opts := options.Find().SetSort(g.sortGroup())
	// Perform the search with pagination and sorting
	return mongoutil.FindPage[*model.Group](ctx, g.coll, filter, pagination, opts)
}
