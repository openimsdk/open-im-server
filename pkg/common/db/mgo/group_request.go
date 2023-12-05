package mgo

import (
	"context"

	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

func NewGroupRequestMgo(db *mongo.Database) (relation.GroupRequestModelInterface, error) {
	coll := db.Collection("group_request")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "group_id", Value: 1},
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &GroupRequestMgo{coll: coll}, nil
}

type GroupRequestMgo struct {
	coll *mongo.Collection
}

func (g *GroupRequestMgo) Create(ctx context.Context, groupRequests []*relation.GroupRequestModel) (err error) {
	return mgoutil.InsertMany(ctx, g.coll, groupRequests)
}

func (g *GroupRequestMgo) Delete(ctx context.Context, groupID string, userID string) (err error) {
	return mgoutil.DeleteOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID})
}

func (g *GroupRequestMgo) UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32) (err error) {
	return mgoutil.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID}, bson.M{"$set": bson.M{"handle_msg": handledMsg, "handle_result": handleResult}}, true)
}

func (g *GroupRequestMgo) Take(ctx context.Context, groupID string, userID string) (groupRequest *relation.GroupRequestModel, err error) {
	return mgoutil.FindOne[*relation.GroupRequestModel](ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID})
}

func (g *GroupRequestMgo) FindGroupRequests(ctx context.Context, groupID string, userIDs []string) ([]*relation.GroupRequestModel, error) {
	return mgoutil.Find[*relation.GroupRequestModel](ctx, g.coll, bson.M{"group_id": groupID, "user_id": bson.M{"$in": userIDs}})
}

func (g *GroupRequestMgo) Page(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, groups []*relation.GroupRequestModel, err error) {
	return mgoutil.FindPage[*relation.GroupRequestModel](ctx, g.coll, bson.M{"user_id": userID}, pagination)
}

func (g *GroupRequestMgo) PageGroup(ctx context.Context, groupIDs []string, pagination pagination.Pagination) (total int64, groups []*relation.GroupRequestModel, err error) {
	return mgoutil.FindPage[*relation.GroupRequestModel](ctx, g.coll, bson.M{"group_id": bson.M{"$in": groupIDs}}, pagination)
}
