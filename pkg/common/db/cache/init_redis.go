package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw/specialerror"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/redis/go-redis/v9"
)

func NewRedis() (redis.UniversalClient, error) {
	if len(config.Config.Redis.Address) == 0 {
		return nil, errors.New("redis address is empty")
	}
	specialerror.AddReplace(redis.Nil, errs.ErrRecordNotFound)
	var rdb redis.UniversalClient
	if len(config.Config.Redis.Address) > 1 {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.Config.Redis.Address,
			Username: config.Config.Redis.Username,
			Password: config.Config.Redis.Password, // no password set
			PoolSize: 50,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     config.Config.Redis.Address[0],
			Username: config.Config.Redis.Username,
			Password: config.Config.Redis.Password, // no password set
			DB:       0,                            // use default DB
			PoolSize: 100,                          // 连接池大小
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis ping %w", err)
	}
	return rdb, nil
}
