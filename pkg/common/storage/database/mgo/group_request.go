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
	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

func NewGroupRequestMgo(db *mongo.Database) (database.GroupRequest, error) {
	coll := db.Collection(database.GroupRequestName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "group_id", Value: 1},
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &GroupRequestMgo{coll: coll}, nil
}

type GroupRequestMgo struct {
	coll *mongo.Collection
}

func (g *GroupRequestMgo) Create(ctx context.Context, groupRequests []*model.GroupRequest) (err error) {
	return mongoutil.InsertMany(ctx, g.coll, groupRequests)
}

func (g *GroupRequestMgo) Delete(ctx context.Context, groupID string, userID string) (err error) {
	return mongoutil.DeleteOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID})
}

func (g *GroupRequestMgo) UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32) (err error) {
	return mongoutil.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID}, bson.M{"$set": bson.M{"handle_msg": handledMsg, "handle_result": handleResult}}, true)
}

func (g *GroupRequestMgo) Take(ctx context.Context, groupID string, userID string) (groupRequest *model.GroupRequest, err error) {
	return mongoutil.FindOne[*model.GroupRequest](ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID})
}

func (g *GroupRequestMgo) FindGroupRequests(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupRequest, error) {
	return mongoutil.Find[*model.GroupRequest](ctx, g.coll, bson.M{"group_id": groupID, "user_id": bson.M{"$in": userIDs}})
}

func (g *GroupRequestMgo) Page(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, groups []*model.GroupRequest, err error) {
	return mongoutil.FindPage[*model.GroupRequest](ctx, g.coll, bson.M{"user_id": userID}, pagination)
}

func (g *GroupRequestMgo) PageGroup(ctx context.Context, groupIDs []string, pagination pagination.Pagination) (total int64, groups []*model.GroupRequest, err error) {
	return mongoutil.FindPage[*model.GroupRequest](ctx, g.coll, bson.M{"group_id": bson.M{"$in": groupIDs}}, pagination)
}
