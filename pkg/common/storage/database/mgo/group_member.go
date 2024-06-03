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

	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewGroupMember(db *mongo.Database) (database.GroupMember, error) {
	coll := db.Collection("group_member")
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
	return &GroupMemberMgo{coll: coll}, nil
}

type GroupMemberMgo struct {
	coll *mongo.Collection
}

func (g *GroupMemberMgo) Create(ctx context.Context, groupMembers []*model.GroupMember) (err error) {
	return mongoutil.InsertMany(ctx, g.coll, groupMembers)
}

func (g *GroupMemberMgo) Delete(ctx context.Context, groupID string, userIDs []string) (err error) {
	filter := bson.M{"group_id": groupID}
	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{"$in": userIDs}
	}
	return mongoutil.DeleteMany(ctx, g.coll, filter)
}

func (g *GroupMemberMgo) UpdateRoleLevel(ctx context.Context, groupID string, userID string, roleLevel int32) error {
	return g.Update(ctx, groupID, userID, bson.M{"role_level": roleLevel})
}

func (g *GroupMemberMgo) Update(ctx context.Context, groupID string, userID string, data map[string]any) (err error) {
	return mongoutil.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID}, bson.M{"$set": data}, true)
}

func (g *GroupMemberMgo) Find(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) (groupMembers []*model.GroupMember, err error) {
	// TODO implement me
	panic("implement me")
}

func (g *GroupMemberMgo) FindMemberUserID(ctx context.Context, groupID string) (userIDs []string, err error) {
	return mongoutil.Find[string](ctx, g.coll, bson.M{"group_id": groupID}, options.Find().SetProjection(bson.M{"_id": 0, "user_id": 1}))
}

func (g *GroupMemberMgo) Take(ctx context.Context, groupID string, userID string) (groupMember *model.GroupMember, err error) {
	return mongoutil.FindOne[*model.GroupMember](ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID})
}

func (g *GroupMemberMgo) TakeOwner(ctx context.Context, groupID string) (groupMember *model.GroupMember, err error) {
	return mongoutil.FindOne[*model.GroupMember](ctx, g.coll, bson.M{"group_id": groupID, "role_level": constant.GroupOwner})
}

func (g *GroupMemberMgo) FindRoleLevelUserIDs(ctx context.Context, groupID string, roleLevel int32) ([]string, error) {
	return mongoutil.Find[string](ctx, g.coll, bson.M{"group_id": groupID, "role_level": roleLevel}, options.Find().SetProjection(bson.M{"_id": 0, "user_id": 1}))
}

func (g *GroupMemberMgo) SearchMember(ctx context.Context, keyword string, groupID string, pagination pagination.Pagination) (total int64, groupList []*model.GroupMember, err error) {
	filter := bson.M{"group_id": groupID, "nickname": bson.M{"$regex": keyword}}
	return mongoutil.FindPage[*model.GroupMember](ctx, g.coll, filter, pagination)
}

func (g *GroupMemberMgo) FindUserJoinedGroupID(ctx context.Context, userID string) (groupIDs []string, err error) {
	return mongoutil.Find[string](ctx, g.coll, bson.M{"user_id": userID}, options.Find().SetProjection(bson.M{"_id": 0, "group_id": 1}))
}

func (g *GroupMemberMgo) TakeGroupMemberNum(ctx context.Context, groupID string) (count int64, err error) {
	return mongoutil.Count(ctx, g.coll, bson.M{"group_id": groupID})
}

func (g *GroupMemberMgo) FindUserManagedGroupID(ctx context.Context, userID string) (groupIDs []string, err error) {
	filter := bson.M{
		"user_id": userID,
		"role_level": bson.M{
			"$in": []int{constant.GroupOwner, constant.GroupAdmin},
		},
	}
	return mongoutil.Find[string](ctx, g.coll, filter, options.Find().SetProjection(bson.M{"_id": 0, "group_id": 1}))
}

func (g *GroupMemberMgo) IsUpdateRoleLevel(data map[string]any) bool {
	if len(data) == 0 {
		return false
	}
	_, ok := data["role_level"]
	return ok
}
