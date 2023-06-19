package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw/specialerror"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/redis/go-redis/v9"
)

func NewRedis() (redis.UniversalClient, error) {
	specialerror.AddReplace(redis.Nil, errs.ErrRecordNotFound)
	var rdb redis.UniversalClient
	if config.Config.Redis.EnableCluster {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.Config.Redis.DBAddress,
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			PoolSize: 50,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     config.Config.Redis.DBAddress[0],
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			DB:       0,                              // use default DB
			PoolSize: 100,                            // 连接池大小
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
