package redis

import (
	"context"
	"fmt"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func newTestOnline() *userOnline {
	opt := &redis.Options{
		Addr:     "172.16.8.48:16379",
		Password: "openIM123",
		DB:       0,
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return &userOnline{rdb: rdb, expire: time.Hour, channelName: "user_online"}
}

func TestOnline(t *testing.T) {
	ts := newTestOnline()
	var count atomic.Int64
	for i := 0; i < 64; i++ {
		go func(userID string) {
			var err error
			for i := 0; ; i++ {
				if i%2 == 0 {
					err = ts.SetUserOnline(context.Background(), userID, []int32{5, 6}, []int32{7, 8, 9})
				} else {
					err = ts.SetUserOnline(context.Background(), userID, []int32{1, 2, 3}, []int32{4, 5, 6})
				}
				if err != nil {
					panic(err)
				}
				count.Add(1)
			}
		}(strconv.Itoa(10000 + i))
	}

	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		t.Log(count.Swap(0))
	}
}

func TestGetOnline(t *testing.T) {
	ts := newTestOnline()
	ctx := context.Background()
	pIDs, err := ts.GetOnline(ctx, "10000")
	if err != nil {
		panic(err)
	}
	t.Log(pIDs)
}

func TestRecvOnline(t *testing.T) {
	ts := newTestOnline()
	ctx := context.Background()
	pubsub := ts.rdb.Subscribe(ctx, cachekey.OnlineChannel)

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatalf("Could not subscribe: %v", err)
	}

	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Printf("Received message from channel %s: %s\n", msg.Channel, msg.Payload)
	}
}
