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
	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewGroupMember(db *mongo.Database) (database.GroupMember, error) {
	coll := db.Collection(database.GroupMemberName)
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
	member, err := NewVersionLog(db.Collection(database.GroupMemberVersionName))
	if err != nil {
		return nil, err
	}
	join, err := NewVersionLog(db.Collection(database.GroupJoinVersionName))
	if err != nil {
		return nil, err
	}
	return &GroupMemberMgo{coll: coll, member: member, join: join}, nil
}

type GroupMemberMgo struct {
	coll   *mongo.Collection
	member database.VersionLog
	join   database.VersionLog
}

func (g *GroupMemberMgo) memberSort() any {
	return bson.D{{Key: "role_level", Value: -1}, {Key: "create_time", Value: 1}}
}

func (g *GroupMemberMgo) Create(ctx context.Context, groupMembers []*model.GroupMember) (err error) {
	return mongoutil.IncrVersion(func() error {
		return mongoutil.InsertMany(ctx, g.coll, groupMembers)
	}, func() error {
		gms := make(map[string][]string)
		for _, member := range groupMembers {
			gms[member.GroupID] = append(gms[member.GroupID], member.UserID)
		}
		for groupID, userIDs := range gms {
			if err := g.member.IncrVersion(ctx, groupID, userIDs, model.VersionStateInsert); err != nil {
				return err
			}
		}
		return nil
	}, func() error {
		gms := make(map[string][]string)
		for _, member := range groupMembers {
			gms[member.UserID] = append(gms[member.UserID], member.GroupID)
		}
		for userID, groupIDs := range gms {
			if err := g.join.IncrVersion(ctx, userID, groupIDs, model.VersionStateInsert); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *GroupMemberMgo) Delete(ctx context.Context, groupID string, userIDs []string) (err error) {
	filter := bson.M{"group_id": groupID}
	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{"$in": userIDs}
	}
	return mongoutil.IncrVersion(func() error {
		return mongoutil.DeleteMany(ctx, g.coll, filter)
	}, func() error {
		if len(userIDs) == 0 {
			return g.member.Delete(ctx, groupID)
		} else {
			return g.member.IncrVersion(ctx, groupID, userIDs, model.VersionStateDelete)
		}
	}, func() error {
		for _, userID := range userIDs {
			if err := g.join.IncrVersion(ctx, userID, []string{groupID}, model.VersionStateDelete); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *GroupMemberMgo) UpdateRoleLevel(ctx context.Context, groupID string, userID string, roleLevel int32) error {
	return mongoutil.IncrVersion(func() error {
		return mongoutil.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID},
			bson.M{"$set": bson.M{"role_level": roleLevel}}, true)
	}, func() error {
		return g.member.IncrVersion(ctx, groupID, []string{model.VersionSortChangeID, userID}, model.VersionStateUpdate)
	})
}
func (g *GroupMemberMgo) UpdateUserRoleLevels(ctx context.Context, groupID string, firstUserID string, firstUserRoleLevel int32, secondUserID string, secondUserRoleLevel int32) error {
	return mongoutil.IncrVersion(func() error {
		if err := mongoutil.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": firstUserID},
			bson.M{"$set": bson.M{"role_level": firstUserRoleLevel}}, true); err != nil {
			return err
		}
		if err := mongoutil.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": secondUserID},
			bson.M{"$set": bson.M{"role_level": secondUserRoleLevel}}, true); err != nil {
			return err
		}
		return nil
	}, func() error {
		return g.member.IncrVersion(ctx, groupID, []string{model.VersionSortChangeID, firstUserID, secondUserID}, model.VersionStateUpdate)
	})
}

func (g *GroupMemberMgo) Update(ctx context.Context, groupID string, userID string, data map[string]any) (err error) {
	if len(data) == 0 {
		return nil
	}
	return mongoutil.IncrVersion(func() error {
		return mongoutil.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID}, bson.M{"$set": data}, true)
	}, func() error {
		var userIDs []string
		if g.IsUpdateRoleLevel(data) {
			userIDs = []string{model.VersionSortChangeID, userID}
		} else {
			userIDs = []string{userID}
		}
		return g.member.IncrVersion(ctx, groupID, userIDs, model.VersionStateUpdate)
	})
}

func (g *GroupMemberMgo) FindMemberUserID(ctx context.Context, groupID string) (userIDs []string, err error) {
	return mongoutil.Find[string](ctx, g.coll, bson.M{"group_id": groupID}, options.Find().SetProjection(bson.M{"_id": 0, "user_id": 1}).SetSort(g.memberSort()))
}

func (g *GroupMemberMgo) Find(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupMember, error) {
	filter := bson.M{"group_id": groupID}
	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{"$in": userIDs}
	}
	return mongoutil.Find[*model.GroupMember](ctx, g.coll, filter)
}

func (g *GroupMemberMgo) FindInGroup(ctx context.Context, userID string, groupIDs []string) ([]*model.GroupMember, error) {
	filter := bson.M{"user_id": userID}
	if len(groupIDs) > 0 {
		filter["group_id"] = bson.M{"$in": groupIDs}
	}
	return mongoutil.Find[*model.GroupMember](ctx, g.coll, filter)
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

func (g *GroupMemberMgo) SearchMember(ctx context.Context, keyword string, groupID string, pagination pagination.Pagination) (int64, []*model.GroupMember, error) {
	filter := bson.M{"group_id": groupID, "nickname": bson.M{"$regex": keyword}}
	return mongoutil.FindPage[*model.GroupMember](ctx, g.coll, filter, pagination, options.Find().SetSort(g.memberSort()))
}

func (g *GroupMemberMgo) FindUserJoinedGroupID(ctx context.Context, userID string) (groupIDs []string, err error) {
	return mongoutil.Find[string](ctx, g.coll, bson.M{"user_id": userID}, options.Find().SetProjection(bson.M{"_id": 0, "group_id": 1}).SetSort(g.memberSort()))
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

func (g *GroupMemberMgo) JoinGroupIncrVersion(ctx context.Context, userID string, groupIDs []string, state int32) error {
	return g.join.IncrVersion(ctx, userID, groupIDs, state)
}

func (g *GroupMemberMgo) MemberGroupIncrVersion(ctx context.Context, groupID string, userIDs []string, state int32) error {
	return g.member.IncrVersion(ctx, groupID, userIDs, state)
}

func (g *GroupMemberMgo) FindMemberIncrVersion(ctx context.Context, groupID string, version uint, limit int) (*model.VersionLog, error) {
	log.ZDebug(ctx, "find member incr version", "groupID", groupID, "version", version)
	return g.member.FindChangeLog(ctx, groupID, version, limit)
}

func (g *GroupMemberMgo) BatchFindMemberIncrVersion(ctx context.Context, groupIDs []string, versions []uint, limits []int) ([]*model.VersionLog, error) {
	log.ZDebug(ctx, "Batch find member incr version", "groupIDs", groupIDs, "versions", versions)
	return g.member.BatchFindChangeLog(ctx, groupIDs, versions, limits)
}

func (g *GroupMemberMgo) FindJoinIncrVersion(ctx context.Context, userID string, version uint, limit int) (*model.VersionLog, error) {
	log.ZDebug(ctx, "find join incr version", "userID", userID, "version", version)
	return g.join.FindChangeLog(ctx, userID, version, limit)
}
