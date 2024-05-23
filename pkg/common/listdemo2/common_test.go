package listdemo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func Result[V any](val V, err error) V {
	if err != nil {
		panic(err)
	}
	return val
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestName(t *testing.T) {
	cli := Result(mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.48:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))
	coll := cli.Database("openim_v3").Collection("friend_version")
	_ = coll
	//Result(coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
	//	{
	//		Keys: map[string]int{"user_id": 1},
	//	},
	//	{
	//		Keys: map[string]int{"friends.friend_user_id": 1},
	//	},
	//}))

	const num = 1
	lm := &LogModel{coll: coll}

	//start := time.Now()
	//eIds := make([]string, 0, num)
	//for i := 0; i < num; i++ {
	//	eIds = append(eIds, strconv.Itoa(1000+(i)))
	//}
	//lm.WriteLogBatch1(context.Background(), "100", eIds, false)
	//end := time.Now()
	//t.Log(end.Sub(start))       // 509.962208ms
	//t.Log(end.Sub(start) / num) // 511.496µs

	start := time.Now()
	wll, err := lm.FindChangeLog(context.Background(), "100", 3, 100)
	if err != nil {
		panic(err)
	}
	t.Log(time.Since(start))
	t.Log(wll)
}
