package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
	"time"
)

func newTestOnline() *userOnline {
	opt := &redis.Options{
		Addr:     "172.16.8.48:16379",
		Password: "openIM123",
		DB:       1,
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return &userOnline{rdb: rdb, expire: time.Hour, channelName: "user_online"}
}

func TestOnline(t *testing.T) {
	ts := newTestOnline()

	//err := ts.SetUserOnline(context.Background(), "1000", []int32{1, 2, 3}, []int32{4, 5, 6})
	err := ts.SetUserOnline(context.Background(), "1000", nil, []int32{1, 2, 3})

	t.Log(err)

}

/*

local function tableToString(tbl, separator)
	local result = {}
    for _, v in ipairs(tbl) do
        table.insert(result, tostring(v))
    end
    return table.concat(result, separator)
end

local myTable = {"one", "two", "three"}
local result = tableToString(myTable, ":")

print(result)

*/

func TestRecvOnline(t *testing.T) {
	ts := newTestOnline()
	ctx := context.Background()
	pubsub := ts.rdb.Subscribe(ctx, "user_online")

	// 等待订阅确认
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatalf("Could not subscribe: %v", err)
	}

	// 创建一个通道来接收消息
	ch := pubsub.Channel()

	// 处理接收到的消息
	for msg := range ch {
		fmt.Printf("Received message from channel %s: %s\n", msg.Channel, msg.Payload)
	}
}
