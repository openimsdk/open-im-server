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

func NewGroupMongo(db *mongo.Database) (relation.GroupModelInterface, error) {
	coll := db.Collection("group")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"group_id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &GroupMgo{coll: coll}, nil
}

type GroupMgo struct {
	coll *mongo.Collection
}

func (g *GroupMgo) Create(ctx context.Context, groups []*relation.GroupModel) (err error) {
	return mtool.InsertMany(ctx, g.coll, groups)
}

func (g *GroupMgo) UpdateState(ctx context.Context, groupID string, state int32) (err error) {
	return g.UpdateMap(ctx, groupID, map[string]any{"state": state})
}

func (g *GroupMgo) UpdateMap(ctx context.Context, groupID string, args map[string]any) (err error) {
	if len(args) == 0 {
		return nil
	}
	return mtool.UpdateOne(ctx, g.coll, bson.M{"group_id": groupID}, bson.M{"$set": args}, true)
}

func (g *GroupMgo) Find(ctx context.Context, groupIDs []string) (groups []*relation.GroupModel, err error) {
	return mtool.Find[*relation.GroupModel](ctx, g.coll, bson.M{"group_id": bson.M{"$in": groupIDs}})
}

func (g *GroupMgo) Take(ctx context.Context, groupID string) (group *relation.GroupModel, err error) {
	return mtool.FindOne[*relation.GroupModel](ctx, g.coll, bson.M{"group_id": groupID})
}

func (g *GroupMgo) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (total int64, groups []*relation.GroupModel, err error) {
	return mtool.FindPage[*relation.GroupModel](ctx, g.coll, bson.M{"group_name": bson.M{"$regex": keyword}}, pagination)
}

func (g *GroupMgo) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	if before == nil {
		return mtool.Count(ctx, g.coll, bson.M{})
	}
	return mtool.Count(ctx, g.coll, bson.M{"create_time": bson.M{"$lt": before}})
}

func (g *GroupMgo) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"create_time": bson.M{
					"$gte": start,
					"$lt":  end,
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$create_time",
					},
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
	}
	type Item struct {
		Date  string `bson:"_id"`
		Count int64  `bson:"count"`
	}
	items, err := mtool.Aggregate[Item](ctx, g.coll, pipeline)
	if err != nil {
		return nil, err
	}
	res := make(map[string]int64, len(items))
	for _, item := range items {
		res[item.Date] = item.Count
	}
	return res, nil
}
