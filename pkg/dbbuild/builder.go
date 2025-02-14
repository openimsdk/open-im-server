package dbbuild

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/redis/go-redis/v9"
)

type Builder interface {
	Mongo(ctx context.Context) (*mongoutil.Client, error)
	Redis(ctx context.Context) (redis.UniversalClient, error)
}

func NewBuilder(mongoConf *config.Mongo, redisConf *config.Redis) Builder {
	if config.Standalone() {
		globalStandalone.setConfig(mongoConf, redisConf)
		return globalStandalone
	}
	return &microservices{
		mongo: mongoConf,
		redis: redisConf,
	}
}
