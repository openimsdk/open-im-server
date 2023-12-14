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

func NewSignalInvitation(db *mongo.Database) (table.SignalInvitationInterface, error) {
	coll := db.Collection("signal_invitation")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "sid", Value: 1},
				{Key: "user_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
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
	return &signalInvitation{coll: coll}, nil
}

type signalInvitation struct {
	coll *mongo.Collection
}

func (x *signalInvitation) Find(ctx context.Context, sid string) ([]*table.SignalInvitationModel, error) {
	return mgoutil.Find[*table.SignalInvitationModel](ctx, x.coll, bson.M{"sid": sid})
}

func (x *signalInvitation) CreateSignalInvitation(ctx context.Context, sid string, inviteeUserIDs []string) error {
	now := time.Now()
	return mgoutil.InsertMany(ctx, x.coll, utils.Slice(inviteeUserIDs, func(userID string) *table.SignalInvitationModel {
		return &table.SignalInvitationModel{
			UserID:       userID,
			SID:          sid,
			InitiateTime: now,
			HandleTime:   time.Unix(0, 0),
		}
	}))
}

func (x *signalInvitation) HandleSignalInvitation(ctx context.Context, sID, InviteeUserID string, status int32) error {
	return mgoutil.UpdateOne(ctx, x.coll, bson.M{"sid": sID, "user_id": InviteeUserID}, bson.M{"$set": bson.M{"status": status, "handle_time": time.Now()}}, true)
}

func (x *signalInvitation) PageSID(ctx context.Context, recvID string, startTime, endTime time.Time, pagination pagination.Pagination) (int64, []string, error) {
	var and []bson.M
	and = append(and, bson.M{"user_id": recvID})
	if !startTime.IsZero() {
		and = append(and, bson.M{"initiate_time": bson.M{"$gte": startTime}})
	}
	if !endTime.IsZero() {
		and = append(and, bson.M{"initiate_time": bson.M{"$lte": endTime}})
	}
	return mgoutil.FindPage[string](ctx, x.coll, bson.M{"$and": and}, pagination, options.Find().SetProjection(bson.M{"_id": 0, "sid": 1}).SetSort(bson.M{"initiate_time": -1}))
}

func (x *signalInvitation) Delete(ctx context.Context, sids []string) error {
	if len(sids) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, x.coll, bson.M{"sid": bson.M{"$in": sids}})
}
