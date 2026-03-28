package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewUserGlobalBlackMongo(db *mongo.Database) (database.UserGlobalBlack, error) {
	coll := db.Collection(database.UserGlobalBlackName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &UserGlobalBlackMgo{coll: coll}, nil
}

type UserGlobalBlackMgo struct {
	coll *mongo.Collection
}

func (u *UserGlobalBlackMgo) Add(ctx context.Context, blacks []*model.UserGlobalBlack) error {
	for _, b := range blacks {
		if b.CreateTime.IsZero() {
			b.CreateTime = time.Now()
		}
	}
	// 使用 upsert 避免重复插入报错
	for _, b := range blacks {
		filter := bson.M{"user_id": b.UserID}
		update := bson.M{
			"$set": bson.M{
				"nickname":    b.Nickname,
				"operator_id": b.OperatorID,
				"reason":      b.Reason,
			},
			"$setOnInsert": bson.M{
				"user_id":     b.UserID,
				"create_time": b.CreateTime,
			},
		}
		opts := options.Update().SetUpsert(true)
		if _, err := u.coll.UpdateOne(ctx, filter, update, opts); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func (u *UserGlobalBlackMgo) Remove(ctx context.Context, users []string) error {
	if len(users) == 0 {
		return nil
	}
	_, err := u.coll.DeleteMany(ctx, bson.M{"user_id": bson.M{"$in": users}})
	return errs.Wrap(err)
}

func (u *UserGlobalBlackMgo) Find(ctx context.Context, userIDs []string) ([]*model.UserGlobalBlack, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	return mongoutil.Find[*model.UserGlobalBlack](ctx, u.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (u *UserGlobalBlackMgo) IsBlocked(ctx context.Context, userID string) (bool, error) {
	count, err := u.coll.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.ZWarn(ctx, "IsBlocked failed", err, "collection", database.UserGlobalBlackName, "userID", userID, "count", count)
		return false, nil
	}

	return count > 0, nil
}

func (u *UserGlobalBlackMgo) Page(ctx context.Context, pagination pagination.Pagination) (int64, []*model.UserGlobalBlack, error) {
	return mongoutil.FindPage[*model.UserGlobalBlack](ctx, u.coll, bson.M{}, pagination)
}
