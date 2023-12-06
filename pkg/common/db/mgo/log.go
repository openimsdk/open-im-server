package mgo

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

func NewLogMongo(db *mongo.Database) (relation.LogInterface, error) {
	coll := db.Collection("log")
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

func (l *LogMgo) Create(ctx context.Context, log []*relation.LogModel) error {
	return mgoutil.InsertMany(ctx, l.coll, log)
}

func (l *LogMgo) Search(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*relation.LogModel, error) {
	filter := bson.M{"create_time": bson.M{"$gte": start, "$lte": end}}
	if keyword != "" {
		filter["user_id"] = bson.M{"$regex": keyword}
	}
	return mgoutil.FindPage[*relation.LogModel](ctx, l.coll, filter, pagination, options.Find().SetSort(bson.M{"create_time": -1}))
}

func (l *LogMgo) Delete(ctx context.Context, logID []string, userID string) error {
	if userID == "" {
		return mgoutil.DeleteMany(ctx, l.coll, bson.M{"log_id": bson.M{"$in": logID}})
	}
	return mgoutil.DeleteMany(ctx, l.coll, bson.M{"log_id": bson.M{"$in": logID}, "user_id": userID})
}

func (l *LogMgo) Get(ctx context.Context, logIDs []string, userID string) ([]*relation.LogModel, error) {
	if userID == "" {
		return mgoutil.Find[*relation.LogModel](ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}})
	}
	return mgoutil.Find[*relation.LogModel](ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}, "user_id": userID})
}
