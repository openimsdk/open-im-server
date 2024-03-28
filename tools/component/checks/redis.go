package checks

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/log"
)

type RedisCheck struct {
	Redis *config.Redis
}

func CheckRedis(ctx context.Context, config *RedisCheck) error {
	redisConfig := &redisutil.Config{
		Address:  config.Redis.Address,
		Username: config.Redis.Username,
		Password: config.Redis.Password,
	}

	log.CInfo(ctx, "Checking Redis connection", "Address", redisConfig.Address)

	err := redisutil.CheckRedis(ctx, redisConfig)
	if err != nil {
		log.CInfo(ctx, "Redis connection failed", "error", err)
		return err
	}

	log.CInfo(ctx, "Redis connection established successfully")
	return nil
}
