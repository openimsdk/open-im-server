package listdemo

import (
	"context"
	"fmt"
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
	coll := cli.Database("openim_v3").Collection("demo")
	_ = coll
	//Result(coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
	//	{
	//		Keys: map[string]int{"user_id": 1},
	//	},
	//	{
	//		Keys: map[string]int{"friends.friend_user_id": 1},
	//	},
	//}))

	wl := WriteLog{
		DID: "100",
		Logs: []LogElem{
			{
				EID:        "1000",
				Deleted:    false,
				Version:    1,
				UpdateTime: time.Now(),
			},
			{
				EID:        "2000",
				Deleted:    false,
				Version:    1,
				UpdateTime: time.Now(),
			},
		},
		Version:       2,
		DeleteVersion: 0,
		LastUpdate:    time.Now(),
	}

	fmt.Println(Result(coll.InsertOne(context.Background(), wl)))

}
