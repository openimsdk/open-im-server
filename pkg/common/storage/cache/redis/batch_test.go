package redis

import (
	"context"
	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"testing"
	"time"
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
	userMgo, err := mgo.NewUserMongo(mgocli.GetDB())
	if err != nil {
		panic(err)
	}
	rock := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())
	//var keys []string
	//for i := 1; i <= 10; i++ {
	//	keys = append(keys, fmt.Sprintf("test%d", i))
	//}
	//res, err := cli.FetchBatch2(ctx, keys, time.Hour, func(idx []int) (map[int]string, error) {
	//	t.Log("FetchBatch2=>", idx)
	//	time.Sleep(time.Second * 1)
	//	res := make(map[int]string)
	//	for _, i := range idx {
	//		res[i] = fmt.Sprintf("hello_%d", i)
	//	}
	//	t.Log("FetchBatch2=>", res)
	//	return res, nil
	//})
	//if err != nil {
	//	t.Log(err)
	//	return
	//}
	//t.Log(res)

	userIDs := []string{"1814217053", "2110910952", "1234567890"}

	res, err := batchGetCache2(ctx, rock, time.Hour, userIDs, func(id string) string {
		return "TEST_USER:" + id
	}, func(v *model.User) string {
		return v.UserID
	}, func(ctx context.Context, ids []string) ([]*model.User, error) {
		t.Log("find mongo", ids)
		return userMgo.Find(ctx, ids)
	})
	if err != nil {
		panic(err)
	}
	t.Log("==>", res)
}
