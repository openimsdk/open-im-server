package mgo

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/OpenIMSDK/tools/pagination"
	"github.com/OpenIMSDK/tools/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/tools/up35/pkg/internal/rtc/mongo/table"
)

func NewMeetingInvitation(db *mongo.Database) (table.MeetingInvitationInterface, error) {
	coll := db.Collection("meeting_invitation")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "room_id", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
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
	return &meetingInvitation{coll: coll}, nil
}

type meetingInvitation struct {
	coll *mongo.Collection
}

func (x *meetingInvitation) FindUserIDs(ctx context.Context, roomID string) ([]string, error) {
	return mgoutil.Find[string](ctx, x.coll, bson.M{"room_id": roomID}, options.Find().SetProjection(bson.M{"_id": 0, "user_id": 1}))
}

func (x *meetingInvitation) CreateMeetingInvitationInfo(ctx context.Context, roomID string, inviteeUserIDs []string) error {
	now := time.Now()
	return mgoutil.InsertMany(ctx, x.coll, utils.Slice(inviteeUserIDs, func(userID string) *table.MeetingInvitationInfo {
		return &table.MeetingInvitationInfo{
			RoomID:     roomID,
			UserID:     userID,
			CreateTime: now,
		}
	}))
}

func (x *meetingInvitation) GetUserInvitedMeetingIDs(ctx context.Context, userID string) (meetingIDs []string, err error) {
	fiveDaysAgo := time.Now().AddDate(0, 0, -5)
	return mgoutil.Find[string](
		ctx,
		x.coll,
		bson.M{"user_id": userID, "create_time": bson.M{"$gte": fiveDaysAgo}},
		options.Find().SetSort(bson.M{"create_time": -1}).SetProjection(bson.M{"_id": 0, "room_id": 1}),
	)
}

func (x *meetingInvitation) Delete(ctx context.Context, roomIDs []string) error {
	return mgoutil.DeleteMany(ctx, x.coll, bson.M{"room_id": bson.M{"$in": roomIDs}})
}

func (x *meetingInvitation) GetMeetingRecords(ctx context.Context, joinedUserID string, startTime, endTime time.Time, pagination pagination.Pagination) (int64, []string, error) {
	var and []bson.M
	and = append(and, bson.M{"user_id": joinedUserID})
	if !startTime.IsZero() {
		and = append(and, bson.M{"create_time": bson.M{"$gte": startTime}})
	}
	if !endTime.IsZero() {
		and = append(and, bson.M{"create_time": bson.M{"$lte": endTime}})
	}
	opt := options.Find().SetSort(bson.M{"create_time": -1}).SetProjection(bson.M{"_id": 0, "room_id": 1})
	return mgoutil.FindPage[string](ctx, x.coll, bson.M{"$and": and}, pagination, opt)
}
