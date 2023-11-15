package newmgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewGroupRequestMgo(db *mongo.Database) (relation.GroupRequestModelInterface, error) {
	return &GroupRequestMgo{coll: db.Collection("group_request")}, nil
}

type GroupRequestMgo struct {
	coll *mongo.Collection
}

func (g *GroupRequestMgo) Create(ctx context.Context, groupRequests []*relation.GroupRequestModel) (err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupRequestMgo) Delete(ctx context.Context, groupID string, userID string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupRequestMgo) UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32) (err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupRequestMgo) Take(ctx context.Context, groupID string, userID string) (groupRequest *relation.GroupRequestModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupRequestMgo) FindGroupRequests(ctx context.Context, groupID string, userIDs []string) (int64, []*relation.GroupRequestModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupRequestMgo) Page(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, groups []*relation.GroupRequestModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupRequestMgo) PageGroup(ctx context.Context, groupIDs []string, pagination pagination.Pagination) (total int64, groups []*relation.GroupRequestModel, err error) {
	//TODO implement me
	panic("implement me")
}
