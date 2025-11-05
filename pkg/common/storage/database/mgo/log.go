package mgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"time"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewLogMongo(db *mongo.Database) (database.Log, error) {
	coll := db.Collection(database.LogName)
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "log_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "create_time", Value: -1},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &LogMgo{coll: coll}, nil
}

type LogMgo struct {
	coll *mongo.Collection
}

func (l *LogMgo) Create(ctx context.Context, log []*model.Log) error {
	return mongoutil.InsertMany(ctx, l.coll, log)
}

func (l *LogMgo) Search(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*model.Log, error) {
	filter := bson.M{"create_time": bson.M{"$gte": start, "$lte": end}}
	if keyword != "" {
		filter["user_id"] = bson.M{"$regex": keyword}
	}
	return mongoutil.FindPage[*model.Log](ctx, l.coll, filter, pagination, options.Find().SetSort(bson.M{"create_time": -1}))
}

func (l *LogMgo) Delete(ctx context.Context, logID []string, userID string) error {
	if userID == "" {
		return mongoutil.DeleteMany(ctx, l.coll, bson.M{"log_id": bson.M{"$in": logID}})
	}
	return mongoutil.DeleteMany(ctx, l.coll, bson.M{"log_id": bson.M{"$in": logID}, "user_id": userID})
}

func (l *LogMgo) Get(ctx context.Context, logIDs []string, userID string) ([]*model.Log, error) {
	if userID == "" {
		return mongoutil.Find[*model.Log](ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}})
	}
	return mongoutil.Find[*model.Log](ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}, "user_id": userID})
}
