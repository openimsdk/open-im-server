package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewUserOfflineRecordMongo(db *mongo.Database) (database.UserOfflineRecord, error) {
	coll := db.Collection(database.UserOfflineRecordName)
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "delete_user_deadline", Value: 1}},
		},
	}
	if _, err := coll.Indexes().CreateMany(context.Background(), indexes); err != nil {
		return nil, errs.Wrap(err)
	}
	return &userOfflineRecordMgo{coll: coll}, nil
}

type userOfflineRecordMgo struct {
	coll *mongo.Collection
}

// Upsert 写入用户的离线记录；若记录已存在则不覆盖（$setOnInsert），
// 保留最早一次的全离线时刻作为计时起点。
// deadline = offlineTime + delete_account_interval，供范围查询快速定位过期账号。
func (u *userOfflineRecordMgo) Upsert(ctx context.Context, userID string, offlineTime, deadline time.Time) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$setOnInsert": bson.M{
			"user_id":              userID,
			"offline_time":         offlineTime,
			"delete_user_deadline": deadline,
		},
	}
	opt := options.Update().SetUpsert(true)
	_, err := u.coll.UpdateOne(ctx, filter, update, opt)
	return errs.Wrap(err)
}

// RefreshOfflineTime 将离线记录的 offline_time 与 delete_user_deadline 同时覆盖写为新值（$set），
// 仅更新已存在的记录；用户在线时（无记录）不做任何操作。
// 适用场景：用户修改 delete_account_interval，让倒计时从设置时刻重新起算。
func (u *userOfflineRecordMgo) RefreshOfflineTime(ctx context.Context, userID string, newOfflineTime, newDeadline time.Time) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{
		"offline_time":         newOfflineTime,
		"delete_user_deadline": newDeadline,
	}}
	_, err := u.coll.UpdateOne(ctx, filter, update)
	return errs.Wrap(err)
}

// Delete 删除用户的离线记录（用户重新上线时调用，停止计时）。
func (u *userOfflineRecordMgo) Delete(ctx context.Context, userID string) error {
	_, err := u.coll.DeleteOne(ctx, bson.M{"user_id": userID})
	return errs.Wrap(err)
}

// FindExpiredUsers 返回 delete_user_deadline <= now 的用户。
// 通过 $lookup 联表 user 集合获取完整 *model.User，$unwind 同时起到过滤孤儿记录的作用
// （若 user 文档已不存在，$unwind 会将其丢弃，避免对无效账号重复触发删除）。
func (u *userOfflineRecordMgo) FindExpiredUsers(ctx context.Context, now time.Time, limit int) ([]*model.User, error) {
	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"delete_user_deadline": bson.M{"$lte": now},
		}},
		bson.M{"$limit": limit},
		bson.M{"$lookup": bson.M{
			"from":         database.UserName,
			"localField":   "user_id",
			"foreignField": "user_id",
			"as":           "u",
		}},
		bson.M{"$unwind": "$u"},
		bson.M{"$replaceRoot": bson.M{"newRoot": "$u"}},
	}
	return mongoutil.Aggregate[*model.User](ctx, u.coll, pipeline)
}
