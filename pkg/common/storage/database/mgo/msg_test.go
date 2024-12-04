package mgo

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/db/mongoutil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

func TestName1(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	cli := Result(mongo.Connect(ctx, options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.48:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))

	v := &MsgMgo{
		coll: cli.Database("openim_v3").Collection("msg3"),
	}

	req := &msg.SearchMessageReq{
		//RecvID: "3187706596",
		//SendID:      "7009965934",
		ContentType: 101,
		//SendTime:    "2024-05-06",
		//SessionType: 3,
		Pagination: &sdkws.RequestPagination{
			PageNumber: 1,
			ShowNumber: 10,
		},
	}
	total, res, err := v.SearchMessage(ctx, req)
	if err != nil {
		panic(err)
	}

	for i, re := range res {
		t.Logf("%d => %d | %+v", i+1, re.Msg.Seq, re.Msg.Content)
	}

	t.Log(total)
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
