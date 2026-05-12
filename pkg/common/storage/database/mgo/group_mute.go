package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewGroupMuteMongo(db *mongo.Database) (database.GroupMute, error) {
	coll := db.Collection(database.GroupMuteName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "owner_user_id", Value: 1},
			{Key: "group_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &GroupMuteMgo{coll: coll}, nil
}

type GroupMuteMgo struct {
	coll *mongo.Collection
}

func (g *GroupMuteMgo) Upsert(ctx context.Context, mute *model.GroupMute) error {
	if mute.CreateTime.IsZero() {
		mute.CreateTime = time.Now()
	}
	filter := bson.M{
		"owner_user_id": mute.OwnerUserID,
		"group_id":      mute.GroupID,
	}
	update := bson.M{
		"$set": bson.M{
			"mute_end_time":   mute.MuteEndTime,
			"mute_duration": mute.MuteDuration,
		},
		"$setOnInsert": bson.M{
			"owner_user_id": mute.OwnerUserID,
			"group_id":      mute.GroupID,
			"create_time":   mute.CreateTime,
		},
	}
	_, err := g.coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return errs.Wrap(err)
}

func (g *GroupMuteMgo) Delete(ctx context.Context, ownerUserID, groupID string) error {
	_, err := g.coll.DeleteOne(ctx, bson.M{
		"owner_user_id": ownerUserID,
		"group_id":      groupID,
	})
	return errs.Wrap(err)
}

func (g *GroupMuteMgo) ListActiveMutedUserIDs(ctx context.Context, groupID string, candidateUserIDs []string) ([]string, error) {
	if len(candidateUserIDs) == 0 {
		return nil, nil
	}
	now := time.Now().Unix()
	filter := bson.M{
		"group_id":      groupID,
		"owner_user_id": bson.M{"$in": candidateUserIDs},
		"$or": bson.A{
			bson.M{"mute_end_time": 0},
			bson.M{"mute_end_time": bson.M{"$gt": now}},
		},
	}
	cur, err := g.coll.Find(ctx, filter, options.Find().SetProjection(bson.M{"owner_user_id": 1, "_id": 0}))
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer cur.Close(ctx)
	var out []string
	for cur.Next(ctx) {
		var doc struct {
			OwnerUserID string `bson:"owner_user_id"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, errs.Wrap(err)
		}
		out = append(out, doc.OwnerUserID)
	}
	return out, cur.Err()
}

func (g *GroupMuteMgo) Get(ctx context.Context, ownerUserID, groupID string) (*model.GroupMute, error) {
	var out model.GroupMute
	err := g.coll.FindOne(ctx, bson.M{
		"owner_user_id": ownerUserID,
		"group_id":      groupID,
	}).Decode(&out)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errs.Wrap(err)
	}
	return &out, nil
}
