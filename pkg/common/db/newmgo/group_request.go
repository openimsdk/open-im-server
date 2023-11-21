package newmgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/newmgo/mgotool"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewGroupRequestMgo(db *mongo.Database) (relation.GroupRequestModelInterface, error) {
	return &GroupRequestMgo{coll: db.Collection("group_request")}, nil
}

type GroupRequestMgo struct {
	coll *mongo.Collection
}

func (g *GroupRequestMgo) Create(ctx context.Context, groupRequests []*relation.GroupRequestModel) (err error) {
	return mgotool.InsertMany(ctx, g.coll, groupRequests)
}

func (g *GroupRequestMgo) Delete(ctx context.Context, groupID string, userID string) (err error) {
	return mgotool.DeleteOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID})
}

func (g *GroupRequestMgo) UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32) (err error) {
	return mgotool.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID}, bson.M{"$set": bson.M{"handle_msg": handledMsg, "handle_result": handleResult}}, true)
}

func (g *GroupRequestMgo) Take(ctx context.Context, groupID string, userID string) (groupRequest *relation.GroupRequestModel, err error) {
	return mgotool.FindOne[*relation.GroupRequestModel](ctx, g.coll, bson.M{"group_id": groupID, "user_id": userID})
}

func (g *GroupRequestMgo) FindGroupRequests(ctx context.Context, groupID string, userIDs []string) ([]*relation.GroupRequestModel, error) {
	return mgotool.Find[*relation.GroupRequestModel](ctx, g.coll, bson.M{"group_id": groupID, "user_id": bson.M{"$in": userIDs}})
}

func (g *GroupRequestMgo) Page(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, groups []*relation.GroupRequestModel, err error) {
	return mgotool.FindPage[*relation.GroupRequestModel](ctx, g.coll, bson.M{"user_id": userID}, pagination)
}

func (g *GroupRequestMgo) PageGroup(ctx context.Context, groupIDs []string, pagination pagination.Pagination) (total int64, groups []*relation.GroupRequestModel, err error) {
	return mgotool.FindPage[*relation.GroupRequestModel](ctx, g.coll, bson.M{"group_id": bson.M{"$in": groupIDs}}, pagination)
}
