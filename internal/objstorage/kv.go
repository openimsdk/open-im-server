package objstorage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type KV interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	IsNotFound(err error) bool
}

func NewKV() KV {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "",
		Username: "",
		Password: "",
	})
	return &redisImpl{
		rdb: rdb,
	}
}

type redisImpl struct {
	rdb *redis.Client
}

func (r *redisImpl) Del(ctx context.Context, key string) error {
	log.Println("redis del", key)
	return r.rdb.Del(ctx, key).Err()
}

func (r *redisImpl) Get(ctx context.Context, key string) (string, error) {
	log.Println("redis get", key)
	return r.rdb.Get(ctx, key).Result()
}

func (r *redisImpl) Set(ctx context.Context, key string, val string, expiration time.Duration) error {
	log.Println("redis set", key, val, expiration.String())
	return r.rdb.Set(ctx, key, val, expiration).Err()
}

func (r *redisImpl) IsNotFound(err error) bool {
	return err == redis.Nil
}
