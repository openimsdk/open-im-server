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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

func NewUserMongo(db *mongo.Database) (database.User, error) {
	coll := db.Collection(database.UserName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &UserMgo{coll: coll}, nil
}

type UserMgo struct {
	coll *mongo.Collection
}

func (u *UserMgo) Create(ctx context.Context, users []*model.User) error {
	return mongoutil.InsertMany(ctx, u.coll, users)
}

func (u *UserMgo) UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error) {
	if len(args) == 0 {
		return nil
	}
	return mongoutil.UpdateOne(ctx, u.coll, bson.M{"user_id": userID}, bson.M{"$set": args}, true)
}

func (u *UserMgo) Find(ctx context.Context, userIDs []string) (users []*model.User, err error) {
	return mongoutil.Find[*model.User](ctx, u.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (u *UserMgo) Take(ctx context.Context, userID string) (user *model.User, err error) {
	return mongoutil.FindOne[*model.User](ctx, u.coll, bson.M{"user_id": userID})
}

func (u *UserMgo) TakeNotification(ctx context.Context, level int64) (user []*model.User, err error) {
	return mongoutil.Find[*model.User](ctx, u.coll, bson.M{"app_manger_level": level})
}

func (u *UserMgo) TakeByNickname(ctx context.Context, nickname string) (user []*model.User, err error) {
	return mongoutil.Find[*model.User](ctx, u.coll, bson.M{"nickname": nickname})
}

func (u *UserMgo) Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*model.User, err error) {
	return mongoutil.FindPage[*model.User](ctx, u.coll, bson.M{}, pagination)
}

func (u *UserMgo) PageFindUser(ctx context.Context, level1 int64, level2 int64, pagination pagination.Pagination) (count int64, users []*model.User, err error) {
	query := bson.M{
		"$or": []bson.M{
			{"app_manger_level": level1},
			{"app_manger_level": level2},
		},
	}

	return mongoutil.FindPage[*model.User](ctx, u.coll, query, pagination)
}

func (u *UserMgo) PageFindUserWithKeyword(
	ctx context.Context,
	level1 int64,
	level2 int64,
	userID string,
	nickName string,
	pagination pagination.Pagination,
) (count int64, users []*model.User, err error) {
	// Initialize the base query with level conditions
	query := bson.M{
		"$and": []bson.M{
			{"app_manger_level": bson.M{"$in": []int64{level1, level2}}},
		},
	}

	// Add userID and userName conditions to the query if they are provided
	if userID != "" || nickName != "" {
		userConditions := []bson.M{}
		if userID != "" {
			// Use regex for userID
			regexPattern := primitive.Regex{Pattern: userID, Options: "i"} // 'i' for case-insensitive matching
			userConditions = append(userConditions, bson.M{"user_id": regexPattern})
		}
		if nickName != "" {
			// Use regex for userName
			regexPattern := primitive.Regex{Pattern: nickName, Options: "i"} // 'i' for case-insensitive matching
			userConditions = append(userConditions, bson.M{"nickname": regexPattern})
		}
		query["$and"] = append(query["$and"].([]bson.M), bson.M{"$or": userConditions})
	}

	// Perform the paginated search
	return mongoutil.FindPage[*model.User](ctx, u.coll, query, pagination)
}

func (u *UserMgo) GetAllUserID(ctx context.Context, pagination pagination.Pagination) (int64, []string, error) {
	return mongoutil.FindPage[string](ctx, u.coll, bson.M{}, pagination, options.Find().SetProjection(bson.M{"_id": 0, "user_id": 1}))
}

func (u *UserMgo) Exist(ctx context.Context, userID string) (exist bool, err error) {
	return mongoutil.Exist(ctx, u.coll, bson.M{"user_id": userID})
}

func (u *UserMgo) GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error) {
	return mongoutil.FindOne[int](ctx, u.coll, bson.M{"user_id": userID}, options.FindOne().SetProjection(bson.M{"_id": 0, "global_recv_msg_opt": 1}))
}

func (u *UserMgo) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	if before == nil {
		return mongoutil.Count(ctx, u.coll, bson.M{})
	}
	return mongoutil.Count(ctx, u.coll, bson.M{"create_time": bson.M{"$lt": before}})
}

func (u *UserMgo) AddUserCommand(ctx context.Context, userID string, Type int32, UUID string, value string, ex string) error {
	collection := u.coll.Database().Collection("userCommands")

	// Create a new document instead of updating an existing one
	doc := bson.M{
		"userID":     userID,
		"type":       Type,
		"uuid":       UUID,
		"createTime": time.Now().Unix(), // assuming you want the creation time in Unix timestamp
		"value":      value,
		"ex":         ex,
	}

	_, err := collection.InsertOne(ctx, doc)
	return errs.Wrap(err)
}

func (u *UserMgo) DeleteUserCommand(ctx context.Context, userID string, Type int32, UUID string) error {
	collection := u.coll.Database().Collection("userCommands")

	filter := bson.M{"userID": userID, "type": Type, "uuid": UUID}

	result, err := collection.DeleteOne(ctx, filter)
	// when err is not nil, result might be nil
	if err != nil {
		return errs.Wrap(err)
	}
	if result.DeletedCount == 0 {
		// No records found to update
		return errs.Wrap(errs.ErrRecordNotFound)
	}
	return errs.Wrap(err)
}
func (u *UserMgo) UpdateUserCommand(ctx context.Context, userID string, Type int32, UUID string, val map[string]any) error {
	if len(val) == 0 {
		return nil
	}

	collection := u.coll.Database().Collection("userCommands")

	filter := bson.M{"userID": userID, "type": Type, "uuid": UUID}
	update := bson.M{"$set": val}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errs.Wrap(err)
	}

	if result.MatchedCount == 0 {
		// No records found to update
		return errs.Wrap(errs.ErrRecordNotFound)
	}

	return nil
}

func (u *UserMgo) GetUserCommand(ctx context.Context, userID string, Type int32) ([]*user.CommandInfoResp, error) {
	collection := u.coll.Database().Collection("userCommands")
	filter := bson.M{"userID": userID, "type": Type}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Initialize commands as a slice of pointers
	commands := []*user.CommandInfoResp{}

	for cursor.Next(ctx) {
		var document struct {
			Type       int32  `bson:"type"`
			UUID       string `bson:"uuid"`
			Value      string `bson:"value"`
			CreateTime int64  `bson:"createTime"`
			Ex         string `bson:"ex"`
		}

		if err := cursor.Decode(&document); err != nil {
			return nil, err
		}

		commandInfo := &user.CommandInfoResp{
			Type:       document.Type,
			Uuid:       document.UUID,
			Value:      document.Value,
			CreateTime: document.CreateTime,
			Ex:         document.Ex,
		}

		commands = append(commands, commandInfo)
	}

	if err := cursor.Err(); err != nil {
		return nil, errs.Wrap(err)
	}

	return commands, nil
}
func (u *UserMgo) GetAllUserCommand(ctx context.Context, userID string) ([]*user.AllCommandInfoResp, error) {
	collection := u.coll.Database().Collection("userCommands")
	filter := bson.M{"userID": userID}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer cursor.Close(ctx)

	// Initialize commands as a slice of pointers
	commands := []*user.AllCommandInfoResp{}

	for cursor.Next(ctx) {
		var document struct {
			Type       int32  `bson:"type"`
			UUID       string `bson:"uuid"`
			Value      string `bson:"value"`
			CreateTime int64  `bson:"createTime"`
			Ex         string `bson:"ex"`
		}

		if err := cursor.Decode(&document); err != nil {
			return nil, errs.Wrap(err)
		}

		commandInfo := &user.AllCommandInfoResp{
			Type:       document.Type,
			Uuid:       document.UUID,
			Value:      document.Value,
			CreateTime: document.CreateTime,
			Ex:         document.Ex,
		}

		commands = append(commands, commandInfo)
	}

	if err := cursor.Err(); err != nil {
		return nil, errs.Wrap(err)
	}
	return commands, nil
}
func (u *UserMgo) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
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
	items, err := mongoutil.Aggregate[Item](ctx, u.coll, pipeline)
	if err != nil {
		return nil, err
	}
	res := make(map[string]int64, len(items))
	for _, item := range items {
		res[item.Date] = item.Count
	}
	return res, nil
}

func (u *UserMgo) SortQuery(ctx context.Context, userIDName map[string]string, asc bool) ([]*model.User, error) {
	if len(userIDName) == 0 {
		return nil, nil
	}
	userIDs := make([]string, 0, len(userIDName))
	attached := make(map[string]string)
	for userID, name := range userIDName {
		userIDs = append(userIDs, userID)
		if name == "" {
			continue
		}
		attached[userID] = name
	}
	var sortValue int
	if asc {
		sortValue = 1
	} else {
		sortValue = -1
	}
	if len(attached) == 0 {
		filter := bson.M{"user_id": bson.M{"$in": userIDs}}
		opt := options.Find().SetSort(bson.M{"nickname": sortValue})
		return mongoutil.Find[*model.User](ctx, u.coll, filter, opt)
	}
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id": bson.M{"$in": userIDs},
			},
		},
		{
			"$addFields": bson.M{
				"_query_sort_name": bson.M{
					"$arrayElemAt": []any{
						bson.M{
							"$filter": bson.M{
								"input": bson.M{
									"$objectToArray": attached,
								},
								"as": "item",
								"cond": bson.M{
									"$eq": []any{"$$item.k", "$user_id"},
								},
							},
						},
						0,
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"_query_sort_name": bson.M{
					"$ifNull": []any{"$_query_sort_name.v", "$nickname"},
				},
			},
		},
		{
			"$sort": bson.M{
				"_query_sort_name": sortValue,
			},
		},
	}
	return mongoutil.Aggregate[*model.User](ctx, u.coll, pipeline)
}
