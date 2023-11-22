package mgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo/mtool"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewLogMongo(db *mongo.Database) (relation.LogInterface, error) {
	lm := &LogMgo{
		coll: db.Collection("log"),
	}
	return lm, nil
}

type LogMgo struct {
	coll *mongo.Collection
}

func (l *LogMgo) Create(ctx context.Context, log []*relation.LogModel) error {
	return mtool.InsertMany(ctx, l.coll, log)
}

func (l *LogMgo) Search(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*relation.LogModel, error) {
	filter := bson.M{"create_time": bson.M{"$gte": start, "$lte": end}}
	if keyword != "" {
		filter["user_id"] = bson.M{"$regex": keyword}
	}
	return mtool.FindPage[*relation.LogModel](ctx, l.coll, filter, pagination, options.Find().SetSort(bson.M{"create_time": -1}))
}

func (l *LogMgo) Delete(ctx context.Context, logID []string, userID string) error {
	if userID == "" {
		return mtool.DeleteMany(ctx, l.coll, bson.M{"log_id": bson.M{"$in": logID}})
	}
	return mtool.DeleteMany(ctx, l.coll, bson.M{"log_id": bson.M{"$in": logID}, "user_id": userID})
}

func (l *LogMgo) Get(ctx context.Context, logIDs []string, userID string) ([]*relation.LogModel, error) {
	if userID == "" {
		return mtool.Find[*relation.LogModel](ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}})
	}
	return mtool.Find[*relation.LogModel](ctx, l.coll, bson.M{"log_id": bson.M{"$in": logIDs}, "user_id": userID})
}
