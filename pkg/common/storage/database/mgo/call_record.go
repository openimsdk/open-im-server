// Copyright © 2024 OpenIM. All rights reserved.
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
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCallRecordMongo(db *mongo.Database) (database.CallRecordDatabase, error) {
	coll := db.Collection(database.CallRecordName)
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "sid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "inviter_user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "invitee_user_id_list", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "create_time", Value: -1}},
		},
	})
	if err != nil {
		return nil, err
	}
	return &callRecordMgo{coll: coll}, nil
}

type callRecordMgo struct {
	coll *mongo.Collection
}

func (c *callRecordMgo) CreateCallRecord(ctx context.Context, record *model.CallRecord) error {
	return mongoutil.InsertMany(ctx, c.coll, []*model.CallRecord{record})
}

func (c *callRecordMgo) SearchCallRecords(ctx context.Context, userID string, status int32, startTime, endTime int64, keyword string, pg pagination.Pagination) (int64, []*model.CallRecord, error) {
	filter := bson.M{}
	if userID != "" {
		filter["$or"] = bson.A{
			bson.M{"inviter_user_id": userID},
			bson.M{"invitee_user_id_list": userID},
		}
	}
	if status != 0 {
		filter["status"] = status
	}
	if startTime > 0 || endTime > 0 {
		timeFilter := bson.M{}
		if startTime > 0 {
			timeFilter["$gte"] = startTime
		}
		if endTime > 0 {
			timeFilter["$lte"] = endTime
		}
		filter["create_time"] = timeFilter
	}
	if keyword != "" {
		filter["inviter_user_nickname"] = bson.M{"$regex": keyword, "$options": "i"}
	}
	return mongoutil.FindPage[*model.CallRecord](ctx, c.coll, filter, pg, options.Find().SetSort(bson.M{"create_time": -1}))
}

func (c *callRecordMgo) DeleteCallRecords(ctx context.Context, sids []string) error {
	return mongoutil.DeleteMany(ctx, c.coll, bson.M{"sid": bson.M{"$in": sids}})
}
