package mgo

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/tools/up35/pkg/internal/rtc/mongo/table"
)

func NewSignal(db *mongo.Database) (table.SignalInterface, error) {
	coll := db.Collection("signal")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "sid", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "inviter_user_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "initiate_time", Value: -1},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &signal{coll: coll}, nil
}

type signal struct {
	coll *mongo.Collection
}

func (x *signal) Find(ctx context.Context, sids []string) ([]*table.SignalModel, error) {
	return mgoutil.Find[*table.SignalModel](ctx, x.coll, bson.M{"sid": bson.M{"$in": sids}})
}

func (x *signal) CreateSignal(ctx context.Context, signalModel *table.SignalModel) error {
	return mgoutil.InsertMany(ctx, x.coll, []*table.SignalModel{signalModel})
}

func (x *signal) Update(ctx context.Context, sid string, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, x.coll, bson.M{"sid": sid}, bson.M{"$set": update}, false)
}

func (x *signal) UpdateSignalFileURL(ctx context.Context, sID, fileURL string) error {
	return x.Update(ctx, sID, map[string]any{"file_url": fileURL})
}

func (x *signal) UpdateSignalEndTime(ctx context.Context, sID string, endTime time.Time) error {
	return x.Update(ctx, sID, map[string]any{"end_time": endTime})
}

func (x *signal) Delete(ctx context.Context, sids []string) error {
	if len(sids) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, x.coll, bson.M{"sid": bson.M{"$in": sids}})
}

func (x *signal) PageSignal(ctx context.Context, sesstionType int32, sendID string, startTime, endTime time.Time, pagination pagination.Pagination) (int64, []*table.SignalModel, error) {
	var and []bson.M
	if !startTime.IsZero() {
		and = append(and, bson.M{"initiate_time": bson.M{"$gte": startTime}})
	}
	if !endTime.IsZero() {
		and = append(and, bson.M{"initiate_time": bson.M{"$lte": endTime}})
	}
	if sesstionType != 0 {
		and = append(and, bson.M{"sesstion_type": sesstionType})
	}
	if sendID != "" {
		and = append(and, bson.M{"inviter_user_id": sendID})
	}
	return mgoutil.FindPage[*table.SignalModel](ctx, x.coll, bson.M{"$and": and}, pagination, options.Find().SetSort(bson.M{"initiate_time": -1}))
}
