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

func NewUserMuteMongo(db *mongo.Database) (database.UserMute, error) {
	coll := db.Collection(database.UserMuteName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "owner_user_id", Value: 1},
			{Key: "muted_user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &UserMuteMgo{coll: coll}, nil
}

type UserMuteMgo struct {
	coll *mongo.Collection
}

func (u *UserMuteMgo) Upsert(ctx context.Context, mute *model.UserMute) error {
	if mute.CreateTime.IsZero() {
		mute.CreateTime = time.Now()
	}
	filter := bson.M{
		"owner_user_id": mute.OwnerUserID,
		"muted_user_id": mute.MutedUserID,
	}
	update := bson.M{
		"$set": bson.M{
			"mute_end_time": mute.MuteEndTime,
		},
		"$setOnInsert": bson.M{
			"owner_user_id": mute.OwnerUserID,
			"muted_user_id": mute.MutedUserID,
			"create_time":   mute.CreateTime,
		},
	}
	_, err := u.coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return errs.Wrap(err)
}

func (u *UserMuteMgo) Delete(ctx context.Context, ownerUserID, mutedUserID string) error {
	_, err := u.coll.DeleteOne(ctx, bson.M{
		"owner_user_id": ownerUserID,
		"muted_user_id": mutedUserID,
	})
	return errs.Wrap(err)
}

func (u *UserMuteMgo) IsMuted(ctx context.Context, ownerUserID, mutedUserID string) (bool, error) {
	now := time.Now().Unix()
	// mute_end_time == 0 means permanent; mute_end_time > now means still active
	filter := bson.M{
		"owner_user_id": ownerUserID,
		"muted_user_id": mutedUserID,
		"$or": bson.A{
			bson.M{"mute_end_time": 0},
			bson.M{"mute_end_time": bson.M{"$gt": now}},
		},
	}
	count, err := u.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, errs.Wrap(err)
	}
	return count > 0, nil
}

func (u *UserMuteMgo) Get(ctx context.Context, ownerUserID, mutedUserID string) (*model.UserMute, error) {
	var out model.UserMute
	err := u.coll.FindOne(ctx, bson.M{
		"owner_user_id": ownerUserID,
		"muted_user_id": mutedUserID,
	}).Decode(&out)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errs.Wrap(err)
	}
	return &out, nil
}
