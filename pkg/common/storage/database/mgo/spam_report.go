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
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewSpamReportMongo(db *mongo.Database) (database.SpamReport, error) {
	coll := db.Collection(database.SpamReportName)
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "report_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "reporter_user_id", Value: 1},
				{Key: "create_time", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "reported_user_id", Value: 1},
				{Key: "create_time", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "create_time", Value: -1},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &SpamReportMgo{coll: coll}, nil
}

type SpamReportMgo struct {
	coll *mongo.Collection
}

func (s *SpamReportMgo) Create(ctx context.Context, report *model.SpamReport) error {
	return mongoutil.InsertOne(ctx, s.coll, report)
}

func (s *SpamReportMgo) Find(ctx context.Context, status int32, reportedUserID, reporterUserID string,
	start, end time.Time, pagination pagination.Pagination) (int64, []*model.SpamReport, error) {
	filter := bson.M{}
	if status >= 0 {
		filter["status"] = status
	}
	if reportedUserID != "" {
		filter["reported_user_id"] = reportedUserID
	}
	if reporterUserID != "" {
		filter["reporter_user_id"] = reporterUserID
	}
	if !start.IsZero() || !end.IsZero() {
		timeFilter := bson.M{}
		if !start.IsZero() {
			timeFilter["$gte"] = start
		}
		if !end.IsZero() {
			timeFilter["$lte"] = end
		}
		filter["create_time"] = timeFilter
	}
	return mongoutil.FindPage[*model.SpamReport](ctx, s.coll, filter, pagination,
		options.Find().SetSort(bson.D{{Key: "create_time", Value: -1}}))
}

func (s *SpamReportMgo) UpdateStatus(ctx context.Context, reportID string, status int32, handlerUserID string, handleTime time.Time) error {
	return mongoutil.UpdateOne(ctx, s.coll,
		bson.M{"report_id": reportID},
		bson.M{"$set": bson.M{
			"status":          status,
			"handler_user_id": handlerUserID,
			"handle_time":     handleTime,
		}},
		false,
	)
}

func (s *SpamReportMgo) Get(ctx context.Context, reportID string) (*model.SpamReport, error) {
	return mongoutil.FindOne[*model.SpamReport](ctx, s.coll, bson.M{"report_id": reportID})
}
