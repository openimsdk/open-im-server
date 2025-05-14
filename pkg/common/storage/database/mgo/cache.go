package mgo

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewCacheMgo(db *mongo.Database) (*CacheMgo, error) {
	coll := db.Collection(database.CacheName)
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "key", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "expire_at", Value: 1},
			},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &CacheMgo{coll: coll}, nil
}

type CacheMgo struct {
	coll *mongo.Collection
}

func (x *CacheMgo) findToMap(res []model.Cache, now time.Time) map[string]string {
	kv := make(map[string]string)
	for _, re := range res {
		if re.ExpireAt != nil && re.ExpireAt.Before(now) {
			continue
		}
		kv[re.Key] = re.Value
	}
	return kv

}

func (x *CacheMgo) Get(ctx context.Context, key []string) (map[string]string, error) {
	if len(key) == 0 {
		return nil, nil
	}
	now := time.Now()
	res, err := mongoutil.Find[model.Cache](ctx, x.coll, bson.M{
		"key": bson.M{"$in": key},
		"$or": []bson.M{
			{"expire_at": bson.M{"$gt": now}},
			{"expire_at": nil},
		},
	})
	if err != nil {
		return nil, err
	}
	return x.findToMap(res, now), nil
}

func (x *CacheMgo) Prefix(ctx context.Context, prefix string) (map[string]string, error) {
	now := time.Now()
	res, err := mongoutil.Find[model.Cache](ctx, x.coll, bson.M{
		"key": bson.M{"$regex": "^" + prefix},
		"$or": []bson.M{
			{"expire_at": bson.M{"$gt": now}},
			{"expire_at": nil},
		},
	})
	if err != nil {
		return nil, err
	}
	return x.findToMap(res, now), nil
}

func (x *CacheMgo) Set(ctx context.Context, key string, value string, expireAt time.Duration) error {
	cv := &model.Cache{
		Key:   key,
		Value: value,
	}
	if expireAt > 0 {
		now := time.Now().Add(expireAt)
		cv.ExpireAt = &now
	}
	opt := options.Update().SetUpsert(true)
	return mongoutil.UpdateOne(ctx, x.coll, bson.M{"key": key}, bson.M{"$set": cv}, false, opt)
}

func (x *CacheMgo) Incr(ctx context.Context, key string, value int) (int, error) {
	pipeline := mongo.Pipeline{
		{
			{"$set", bson.M{
				"value": bson.M{
					"$toString": bson.M{
						"$add": bson.A{
							bson.M{"$toInt": "$value"},
							value,
						},
					},
				},
			}},
		},
	}
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res, err := mongoutil.FindOneAndUpdate[model.Cache](ctx, x.coll, bson.M{"key": key}, pipeline, opt)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res.Value)
}

func (x *CacheMgo) Del(ctx context.Context, key []string) error {
	if len(key) == 0 {
		return nil
	}
	_, err := x.coll.DeleteMany(ctx, bson.M{"key": bson.M{"$in": key}})
	return errs.Wrap(err)
}

func (x *CacheMgo) lockKey(key string) string {
	return "LOCK_" + key
}

func (x *CacheMgo) Lock(ctx context.Context, key string, duration time.Duration) (string, error) {
	tmp, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	if duration <= 0 || duration > time.Minute*10 {
		duration = time.Minute * 10
	}
	cv := &model.Cache{
		Key:      x.lockKey(key),
		Value:    tmp.String(),
		ExpireAt: nil,
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	wait := func() error {
		timeout := time.NewTimer(time.Millisecond * 100)
		defer timeout.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout.C:
			return nil
		}
	}
	for {
		if err := mongoutil.DeleteOne(ctx, x.coll, bson.M{"key": key, "expire_at": bson.M{"$lt": time.Now()}}); err != nil {
			return "", err
		}
		expireAt := time.Now().Add(duration)
		cv.ExpireAt = &expireAt
		if err := mongoutil.InsertMany[*model.Cache](ctx, x.coll, []*model.Cache{cv}); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				if err := wait(); err != nil {
					return "", err
				}
				continue
			}
			return "", err
		}
		return cv.Value, nil
	}
}

func (x *CacheMgo) Unlock(ctx context.Context, key string, value string) error {
	return mongoutil.DeleteOne(ctx, x.coll, bson.M{"key": x.lockKey(key), "value": value})
}
