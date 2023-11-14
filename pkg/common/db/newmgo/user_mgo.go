package newmgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/newmgo/mgotool"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewUserMongo(db *mongo.Database) relation.UserModelInterface {
	return &UserMgo{
		coll: db.Collection("user"),
	}
}

type UserMgo struct {
	coll *mongo.Collection
}

func (u *UserMgo) Create(ctx context.Context, users []*relation.UserModel) error {
	return mgotool.InsertMany(ctx, u.coll, users)
}

func (u *UserMgo) UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error) {
	if len(args) == 0 {
		return nil
	}
	return mgotool.UpdateOne(ctx, u.coll, bson.M{"user_id": userID}, bson.M{"$set": args}, true)
}

func (u *UserMgo) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	return mgotool.Find[*relation.UserModel](ctx, u.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (u *UserMgo) Take(ctx context.Context, userID string) (user *relation.UserModel, err error) {
	return mgotool.FindOne[*relation.UserModel](ctx, u.coll, bson.M{"user_id": userID})
}

func (u *UserMgo) Page(ctx context.Context, pagination mgotool.Pagination) (count int64, users []*relation.UserModel, err error) {
	return mgotool.FindPage[*relation.UserModel](ctx, u.coll, bson.M{}, pagination)
}

func (u *UserMgo) GetAllUserID(ctx context.Context, pagination mgotool.Pagination) (int64, []string, error) {
	return mgotool.FindPage[string](ctx, u.coll, bson.M{}, pagination, options.Find().SetProjection(bson.M{"user_id": 1}))
}

func (u *UserMgo) Exist(ctx context.Context, userID string) (exist bool, err error) {
	return mgotool.Exist(ctx, u.coll, bson.M{"user_id": userID})
}

func (u *UserMgo) GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error) {
	return mgotool.FindOne[int](ctx, u.coll, bson.M{"user_id": userID}, options.FindOne().SetProjection(bson.M{"global_recv_msg_opt": 1}))
}

func (u *UserMgo) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	return mgotool.Count(ctx, u.coll, bson.M{"create_time": bson.M{"$lt": before}})
}

func (u *UserMgo) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	//type Temp struct {
	//	CreateTime time.Time `bson:"create_time"`
	//	Number     int64     `bson:"number"`
	//}
	//mgotool.Find(ctx, u.coll, bson.M{"create_time": bson.M{"$gte": start, "$lt": end}}, options.Find().SetProjection(bson.M{"create_time": 1}))
	panic("implement me")
	return nil, nil
}
