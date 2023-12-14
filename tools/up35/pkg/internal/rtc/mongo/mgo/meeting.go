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

func NewMeeting(db *mongo.Database) (table.MeetingInterface, error) {
	coll := db.Collection("meeting")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "room_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "host_user_id", Value: 1},
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
	return &meeting{coll: coll}, nil
}

type meeting struct {
	coll *mongo.Collection
}

func (x *meeting) Find(ctx context.Context, roomIDs []string) ([]*table.MeetingInfo, error) {
	return mgoutil.Find[*table.MeetingInfo](ctx, x.coll, bson.M{"room_id": bson.M{"$in": roomIDs}})
}

func (x *meeting) CreateMeetingInfo(ctx context.Context, meetingInfo *table.MeetingInfo) error {
	return mgoutil.InsertMany(ctx, x.coll, []*table.MeetingInfo{meetingInfo})
}

func (x *meeting) UpdateMeetingInfo(ctx context.Context, roomID string, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, x.coll, bson.M{"room_id": roomID}, bson.M{"$set": update}, false)
}

func (x *meeting) GetUnCompleteMeetingIDList(ctx context.Context, roomIDs []string) ([]string, error) {
	if len(roomIDs) == 0 {
		return nil, nil
	}
	return mgoutil.Find[string](ctx, x.coll, bson.M{"room_id": bson.M{"$in": roomIDs}, "status": 0}, options.Find().SetProjection(bson.M{"_id": 0, "room_id": 1}))
}

func (x *meeting) Delete(ctx context.Context, roomIDs []string) error {
	return mgoutil.DeleteMany(ctx, x.coll, bson.M{"room_id": bson.M{"$in": roomIDs}})
}

func (x *meeting) GetMeetingRecords(ctx context.Context, hostUserID string, startTime, endTime time.Time, pagination pagination.Pagination) (int64, []*table.MeetingInfo, error) {
	var and []bson.M
	if hostUserID != "" {
		and = append(and, bson.M{"host_user_id": hostUserID})
	}
	if !startTime.IsZero() {
		and = append(and, bson.M{"create_time": bson.M{"$gte": startTime}})
	}
	if !endTime.IsZero() {
		and = append(and, bson.M{"create_time": bson.M{"$lte": endTime}})
	}
	filter := bson.M{}
	if len(and) > 0 {
		filter["$and"] = and
	}
	return mgoutil.FindPage[*table.MeetingInfo](ctx, x.coll, filter, pagination, options.Find().SetSort(bson.M{"create_time": -1}))
}
