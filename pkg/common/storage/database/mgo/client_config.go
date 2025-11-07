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

	"github.com/openimsdk/tools/errs"
)

func NewClientConfig(db *mongo.Database) (database.ClientConfig, error) {
	coll := db.Collection("config")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "key", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &ClientConfig{
		coll: coll,
	}, nil
}

type ClientConfig struct {
	coll *mongo.Collection
}

func (x *ClientConfig) Set(ctx context.Context, userID string, config map[string]string) error {
	if len(config) == 0 {
		return nil
	}
	for key, value := range config {
		filter := bson.M{"key": key, "user_id": userID}
		update := bson.M{
			"value": value,
		}
		err := mongoutil.UpdateOne(ctx, x.coll, filter, bson.M{"$set": update}, false, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *ClientConfig) Get(ctx context.Context, userID string) (map[string]string, error) {
	cs, err := mongoutil.Find[*model.ClientConfig](ctx, x.coll, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	cm := make(map[string]string)
	for _, config := range cs {
		cm[config.Key] = config.Value
	}
	return cm, nil
}

func (x *ClientConfig) Del(ctx context.Context, userID string, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, x.coll, bson.M{"key": bson.M{"$in": keys}, "user_id": userID})
}

func (x *ClientConfig) GetPage(ctx context.Context, userID string, key string, pagination pagination.Pagination) (int64, []*model.ClientConfig, error) {
	filter := bson.M{}
	if userID != "" {
		filter["user_id"] = userID
	}
	if key != "" {
		filter["key"] = key
	}
	return mongoutil.FindPage[*model.ClientConfig](ctx, x.coll, filter, pagination)
}
