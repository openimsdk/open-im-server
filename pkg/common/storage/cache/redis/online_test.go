package redis

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/db/redisutil"
	"testing"
	"time"
)

/*
address: [ 172.16.8.48:7001, 172.16.8.48:7002, 172.16.8.48:7003, 172.16.8.48:7004, 172.16.8.48:7005, 172.16.8.48:7006 ]
username:
password: passwd123
clusterMode: true
db: 0
maxRetry: 10
*/
func TestName111111(t *testing.T) {
	conf := config.Redis{
		Address: []string{
			"172.16.8.124:7001",
			"172.16.8.124:7002",
			"172.16.8.124:7003",
			"172.16.8.124:7004",
			"172.16.8.124:7005",
			"172.16.8.124:7006",
		},
		RedisMode: "cluster",
		Password:    "passwd123",
		//Address:  []string{"localhost:16379"},
		//Password: "openIM123",
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()
	rdb, err := redisutil.NewRedisClient(ctx, conf.Build())
	if err != nil {
		panic(err)
	}
	online := NewUserOnline(rdb)

	userID := "a123456"
	t.Log(online.GetOnline(ctx, userID))
	t.Log(online.SetUserOnline(ctx, userID, []int32{1, 2, 3, 4}, nil))
	t.Log(online.GetOnline(ctx, userID))

}

func TestName111(t *testing.T) {

}
