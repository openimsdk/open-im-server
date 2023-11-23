package mgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo/mtool"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewFriendRequestMongo(db *mongo.Database) (relation.FriendRequestModelInterface, error) {
	coll := db.Collection("friend_request")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"from_user_id": 1, "to_user_id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &FriendRequestMgo{coll: coll}, nil
}

type FriendRequestMgo struct {
	coll *mongo.Collection
}

func (f *FriendRequestMgo) FindToUserID(ctx context.Context, toUserID string, pagination pagination.Pagination) (total int64, friendRequests []*relation.FriendRequestModel, err error) {
	return mtool.FindPage[*relation.FriendRequestModel](ctx, f.coll, bson.M{"to_user_id": toUserID}, pagination)
}

func (f *FriendRequestMgo) FindFromUserID(ctx context.Context, fromUserID string, pagination pagination.Pagination) (total int64, friendRequests []*relation.FriendRequestModel, err error) {
	return mtool.FindPage[*relation.FriendRequestModel](ctx, f.coll, bson.M{"from_user_id": fromUserID}, pagination)
}

func (f *FriendRequestMgo) FindBothFriendRequests(ctx context.Context, fromUserID, toUserID string) (friends []*relation.FriendRequestModel, err error) {
	filter := bson.M{"$or": []bson.M{
		{"from_user_id": fromUserID, "to_user_id": toUserID},
		{"from_user_id": toUserID, "to_user_id": fromUserID},
	}}
	return mtool.Find[*relation.FriendRequestModel](ctx, f.coll, filter)
}

func (f *FriendRequestMgo) Create(ctx context.Context, friendRequests []*relation.FriendRequestModel) error {
	return mtool.InsertMany(ctx, f.coll, friendRequests)
}

func (f *FriendRequestMgo) Delete(ctx context.Context, fromUserID, toUserID string) (err error) {
	return mtool.DeleteOne(ctx, f.coll, bson.M{"from_user_id": fromUserID, "to_user_id": toUserID})
}

func (f *FriendRequestMgo) UpdateByMap(ctx context.Context, formUserID, toUserID string, args map[string]any) (err error) {
	if len(args) == 0 {
		return nil
	}
	return mtool.UpdateOne(ctx, f.coll, bson.M{"from_user_id": formUserID, "to_user_id": toUserID}, bson.M{"$set": args}, true)
}

func (f *FriendRequestMgo) Update(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	updater := bson.M{}
	if friendRequest.HandleResult != 0 {
		updater["handle_result"] = friendRequest.HandleResult
	}
	if friendRequest.ReqMsg != "" {
		updater["req_msg"] = friendRequest.ReqMsg
	}
	if friendRequest.HandlerUserID != "" {
		updater["handler_user_id"] = friendRequest.HandlerUserID
	}
	if friendRequest.HandleMsg != "" {
		updater["handle_msg"] = friendRequest.HandleMsg
	}
	if !friendRequest.HandleTime.IsZero() {
		updater["handle_time"] = friendRequest.HandleTime
	}
	if friendRequest.Ex != "" {
		updater["ex"] = friendRequest.Ex
	}
	if len(updater) == 0 {
		return nil
	}
	filter := bson.M{"from_user_id": friendRequest.FromUserID, "to_user_id": friendRequest.ToUserID}
	return mtool.UpdateOne(ctx, f.coll, filter, bson.M{"$set": updater}, true)
}

func (f *FriendRequestMgo) Find(ctx context.Context, fromUserID, toUserID string) (friendRequest *relation.FriendRequestModel, err error) {
	return mtool.FindOne[*relation.FriendRequestModel](ctx, f.coll, bson.M{"from_user_id": fromUserID, "to_user_id": toUserID})
}

func (f *FriendRequestMgo) Take(ctx context.Context, fromUserID, toUserID string) (friendRequest *relation.FriendRequestModel, err error) {
	return f.Find(ctx, fromUserID, toUserID)
}
