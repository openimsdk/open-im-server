package newmgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/newmgo/mgotool"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewFriendRequestMongo(db *mongo.Database) (relation.FriendRequestModelInterface, error) {
	return &FriendRequestMgo{
		coll: db.Collection("friend_request"),
	}, nil
}

type FriendRequestMgo struct {
	coll *mongo.Collection
}

func (f *FriendRequestMgo) FindToUserID(ctx context.Context, toUserID string, pagination pagination.Pagination) (total int64, friendRequests []*relation.FriendRequestModel, err error) {
	return mgotool.FindPage[*relation.FriendRequestModel](ctx, f.coll, bson.M{"to_user_id": toUserID}, pagination)
}

func (f *FriendRequestMgo) FindFromUserID(ctx context.Context, fromUserID string, pagination pagination.Pagination) (total int64, friendRequests []*relation.FriendRequestModel, err error) {
	return mgotool.FindPage[*relation.FriendRequestModel](ctx, f.coll, bson.M{"from_user_id": fromUserID}, pagination)
}

func (f *FriendRequestMgo) FindBothFriendRequests(ctx context.Context, fromUserID, toUserID string) (friends []*relation.FriendRequestModel, err error) {
	filter := bson.M{"$or": []bson.M{
		{"from_user_id": fromUserID, "to_user_id": toUserID},
		{"from_user_id": toUserID, "to_user_id": fromUserID},
	}}
	return mgotool.Find[*relation.FriendRequestModel](ctx, f.coll, filter)
}

func (f *FriendRequestMgo) Create(ctx context.Context, friendRequests []*relation.FriendRequestModel) error {
	return mgotool.InsertMany(ctx, f.coll, friendRequests)
}

func (f *FriendRequestMgo) Delete(ctx context.Context, fromUserID, toUserID string) (err error) {
	return mgotool.DeleteOne(ctx, f.coll, bson.M{"from_user_id": fromUserID, "to_user_id": toUserID})
}

func (f *FriendRequestMgo) UpdateByMap(ctx context.Context, formUserID, toUserID string, args map[string]any) (err error) {
	if len(args) == 0 {
		return nil
	}
	return mgotool.UpdateOne(ctx, f.coll, bson.M{"from_user_id": formUserID, "to_user_id": toUserID}, bson.M{"$set": args}, true)
}

func (f *FriendRequestMgo) Update(ctx context.Context, friendRequest *relation.FriendRequestModel) (err error) {
	return mgotool.UpdateOne(ctx, f.coll, bson.M{"_id": friendRequest.ID}, bson.M{"$set": friendRequest}, true)
}

func (f *FriendRequestMgo) Find(ctx context.Context, fromUserID, toUserID string) (friendRequest *relation.FriendRequestModel, err error) {
	return mgotool.FindOne[*relation.FriendRequestModel](ctx, f.coll, bson.M{"from_user_id": fromUserID, "to_user_id": toUserID})
}

func (f *FriendRequestMgo) Take(ctx context.Context, fromUserID, toUserID string) (friendRequest *relation.FriendRequestModel, err error) {
	return f.Find(ctx, fromUserID, toUserID)
}
