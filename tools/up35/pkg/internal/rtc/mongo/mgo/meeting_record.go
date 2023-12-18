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

	"github.com/OpenIMSDK/tools/mgoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/openimsdk/open-im-server/v3/tools/up35/pkg/internal/rtc/mongo/table"
)

func NewMeetingRecord(db *mongo.Database) (table.MeetingRecordInterface, error) {
	coll := db.Collection("meeting_record")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "room_id", Value: 1},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &meetingRecord{coll: coll}, nil
}

type meetingRecord struct {
	coll *mongo.Collection
}

func (x *meetingRecord) CreateMeetingVideoRecord(ctx context.Context, meetingVideoRecord *table.MeetingVideoRecord) error {
	return mgoutil.InsertMany(ctx, x.coll, []*table.MeetingVideoRecord{meetingVideoRecord})
}
