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

func NewUserMongo(db *mongo.Database) (relation.UserModelInterface, error) {
	coll := db.Collection("user")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"user_id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}
	return &UserMgo{coll: coll}, nil
}

type UserMgo struct {
	coll *mongo.Collection
}

func (u *UserMgo) Create(ctx context.Context, users []*relation.UserModel) error {
	return mtool.InsertMany(ctx, u.coll, users)
}

func (u *UserMgo) UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error) {
	if len(args) == 0 {
		return nil
	}
	return mtool.UpdateOne(ctx, u.coll, bson.M{"user_id": userID}, bson.M{"$set": args}, true)
}

func (u *UserMgo) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	return mtool.Find[*relation.UserModel](ctx, u.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (u *UserMgo) Take(ctx context.Context, userID string) (user *relation.UserModel, err error) {
	return mtool.FindOne[*relation.UserModel](ctx, u.coll, bson.M{"user_id": userID})
}

func (u *UserMgo) Page(ctx context.Context, pagination pagination.Pagination) (count int64, users []*relation.UserModel, err error) {
	return mtool.FindPage[*relation.UserModel](ctx, u.coll, bson.M{}, pagination)
}

func (u *UserMgo) GetAllUserID(ctx context.Context, pagination pagination.Pagination) (int64, []string, error) {
	return mtool.FindPage[string](ctx, u.coll, bson.M{}, pagination, options.Find().SetProjection(bson.M{"user_id": 1}))
}

func (u *UserMgo) Exist(ctx context.Context, userID string) (exist bool, err error) {
	return mtool.Exist(ctx, u.coll, bson.M{"user_id": userID})
}

func (u *UserMgo) GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error) {
	return mtool.FindOne[int](ctx, u.coll, bson.M{"user_id": userID}, options.FindOne().SetProjection(bson.M{"global_recv_msg_opt": 1}))
}

func (u *UserMgo) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	if before == nil {
		return mtool.Count(ctx, u.coll, bson.M{})
	}
	return mtool.Count(ctx, u.coll, bson.M{"create_time": bson.M{"$lt": before}})
}

func (u *UserMgo) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
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
	items, err := mtool.Aggregate[Item](ctx, u.coll, pipeline)
	if err != nil {
		return nil, err
	}
	res := make(map[string]int64, len(items))
	for _, item := range items {
		res[item.Date] = item.Count
	}
	return res, nil
}
