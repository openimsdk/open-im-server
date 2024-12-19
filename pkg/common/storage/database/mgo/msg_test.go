package mgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestName1(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	cli := Result(mongo.Connect(ctx, options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.66:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))

	//v := &MsgMgo{
	//	coll: cli.Database("openim_v3").Collection("msg3"),
	//}
	//
	//req := &msg.SearchMessageReq{
	//	//RecvID: "3187706596",
	//	//SendID:      "7009965934",
	//	ContentType: 101,
	//	//SendTime:    "2024-05-06",
	//	//SessionType: 3,
	//	Pagination: &sdkws.RequestPagination{
	//		PageNumber: 1,
	//		ShowNumber: 10,
	//	},
	//}
	//total, res, err := v.SearchMessage(ctx, req)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for i, re := range res {
	//	t.Logf("%d => %d | %+v", i+1, re.Msg.Seq, re.Msg.Content)
	//}
	//
	//t.Log(total)

	msg, err := NewMsgMongo(cli.Database("openim_v3"))
	if err != nil {
		panic(err)
	}
	res, err := msg.GetBeforeMsg(ctx, time.Now().UnixMilli(), []string{"1:0"}, 1000)
	if err != nil {
		panic(err)
	}
	t.Log(len(res))
}

func TestName10(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	cli := Result(mongo.Connect(ctx, options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.48:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))

	v := &MsgMgo{
		coll: cli.Database("openim_v3").Collection("msg3"),
	}
	opt := options.Find().SetLimit(1000)

	res, err := mongoutil.Find[model.MsgDocModel](ctx, v.coll, bson.M{}, opt)
	if err != nil {
		panic(err)
	}
	ctx = context.Background()
	for i := 0; i < 100000; i++ {
		for j := range res {
			res[j].DocID = strconv.FormatUint(rand.Uint64(), 10) + ":0"
		}
		if err := mongoutil.InsertMany(ctx, v.coll, res); err != nil {
			panic(err)
		}
		t.Log("====>", time.Now(), i)
	}

}

func TestName3(t *testing.T) {
	t.Log(uint64(math.MaxUint64))
	t.Log(int64(math.MaxInt64))

	t.Log(int64(math.MinInt64))
}

func TestName4(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	cli := Result(mongo.Connect(ctx, options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.66:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))

	msg, err := NewMsgMongo(cli.Database("openim_v3"))
	if err != nil {
		panic(err)
	}
	res, err := msg.GetBeforeMsg(ctx, time.Now().UnixMilli(), []string{"1:0"}, 1000)
	if err != nil {
		panic(err)
	}
	t.Log(len(res))
}
