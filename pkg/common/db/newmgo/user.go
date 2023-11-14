package newmgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/newmgo/mgotool"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserModel struct {
	UserID           string    `bson:"user_id"`
	Nickname         string    `bson:"nickname"`
	FaceURL          string    `bson:"face_url"`
	Ex               string    `bson:"ex"`
	AppMangerLevel   int32     `bson:"app_manger_level"`
	GlobalRecvMsgOpt int32     `bson:"global_recv_msg_opt"`
	CreateTime       time.Time `bson:"create_time"`
}

type UserModelInterface interface {
	Create(ctx context.Context, users []*UserModel) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]any) (err error)
	// 获取指定用户信息  不存在，也不返回错误
	Find(ctx context.Context, userIDs []string) (users []*UserModel, err error)
	// 获取某个用户信息  不存在，则返回错误
	Take(ctx context.Context, userID string) (user *UserModel, err error)
	// 获取用户信息 不存在，不返回错误
	Page(ctx context.Context, pageNumber, showNumber int32) (users []*UserModel, count int64, err error)
	GetAllUserID(ctx context.Context, pageNumber, showNumber int32) (userIDs []string, err error)
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
	// 获取用户总数
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
	// 获取范围内用户增量
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
}

type UserMgo struct {
	coll *mongo.Collection
}

func (u *UserMgo) Create(ctx context.Context, users []*UserModel) error {
	return mgotool.InsertMany(ctx, u.coll, users)
}

func (u *UserMgo) UpdateOneByMap(ctx context.Context, userID string, args map[string]any) error {
	if len(args) == 0 {
		return nil
	}
	return mgotool.UpdateOne(ctx, u.coll, bson.M{"user_id": userID}, bson.M{"$set": args}, true)
}

func (u *UserMgo) Find(ctx context.Context, userIDs []string) (users []*UserModel, err error) {
	return mgotool.Find[*UserModel](ctx, u.coll, bson.M{"user_id": bson.M{"$in": userIDs}})
}

func (u *UserMgo) Take(ctx context.Context, userID string) (user *UserModel, err error) {
	return mgotool.FindOne[*UserModel](ctx, u.coll, bson.M{"user_id": userID})
}

func (u *UserMgo) Page(ctx context.Context, pagination mgotool.Pagination) (count int64, users []*UserModel, err error) {
	return mgotool.FindPage[*UserModel](ctx, u.coll, bson.M{}, pagination)
}

func (u *UserMgo) GetAllUserID(ctx context.Context, pagination mgotool.Pagination) (int64, []string, error) {
	return mgotool.FindPage[string](ctx, u.coll, bson.M{}, pagination, options.Find().SetProjection(bson.M{"user_id": 1}))
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
