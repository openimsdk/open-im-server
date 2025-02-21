package dbbuild

import (
	"context"
	"sync"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/redis/go-redis/v9"
)

const (
	standaloneMongo = "mongo"
	standaloneRedis = "redis"
)

var globalStandalone = &standalone{}

type standaloneConn[C any] struct {
	Conn C
	Err  error
}

func (x *standaloneConn[C]) result() (C, error) {
	return x.Conn, x.Err
}

type standalone struct {
	lock  sync.Mutex
	mongo *config.Mongo
	redis *config.Redis
	conn  map[string]any
}

func (x *standalone) setConfig(mongoConf *config.Mongo, redisConf *config.Redis) {
	x.lock.Lock()
	defer x.lock.Unlock()
	x.mongo = mongoConf
	x.redis = redisConf
}

func (x *standalone) Mongo(ctx context.Context) (*mongoutil.Client, error) {
	x.lock.Lock()
	defer x.lock.Unlock()
	if x.conn == nil {
		x.conn = make(map[string]any)
	}
	v, ok := x.conn[standaloneMongo]
	if !ok {
		var val standaloneConn[*mongoutil.Client]
		val.Conn, val.Err = mongoutil.NewMongoDB(ctx, x.mongo.Build())
		v = &val
		x.conn[standaloneMongo] = v
	}
	return v.(*standaloneConn[*mongoutil.Client]).result()
}

func (x *standalone) Redis(ctx context.Context) (redis.UniversalClient, error) {
	x.lock.Lock()
	defer x.lock.Unlock()
	if x.redis.Disable {
		return nil, nil
	}
	if x.conn == nil {
		x.conn = make(map[string]any)
	}
	v, ok := x.conn[standaloneRedis]
	if !ok {
		var val standaloneConn[redis.UniversalClient]
		val.Conn, val.Err = redisutil.NewRedisClient(ctx, x.redis.Build())
		v = &val
		x.conn[standaloneRedis] = v
	}
	return v.(*standaloneConn[redis.UniversalClient]).result()
}
