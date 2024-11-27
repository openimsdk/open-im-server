package redis

import (
	"context"
	"testing"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
)

func TestName(t *testing.T) {
	//var rocks rockscache.Client
	//rdb := getRocksCacheRedisClient(&rocks)
	//t.Log(rdb == nil)

	ctx := context.Background()
	rdb, err := redisutil.NewRedisClient(ctx, (&config.Redis{
		Address:  []string{"172.16.8.48:16379"},
		Password: "openIM123",
		DB:       3,
	}).Build())
	if err != nil {
		panic(err)
	}
	mgocli, err := mongoutil.NewMongoDB(ctx, (&config.Mongo{
		Address:     []string{"172.16.8.48:37017"},
		Database:    "openim_v3",
		Username:    "openIM",
		Password:    "openIM123",
		MaxPoolSize: 100,
		MaxRetry:    1,
	}).Build())
	if err != nil {
		panic(err)
	}
	//userMgo, err := mgo.NewUserMongo(mgocli.GetDB())
	//if err != nil {
	//	panic(err)
	//}
	//rock := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())
	mgoSeqUser, err := mgo.NewSeqUserMongo(mgocli.GetDB())
	if err != nil {
		panic(err)
	}
	seqUser := NewSeqUserCacheRedis(rdb, mgoSeqUser)

	res, err := seqUser.GetUserReadSeqs(ctx, "2110910952", []string{"sg_2920732023", "sg_345762580"})
	if err != nil {
		panic(err)
	}

	t.Log(res)

}
