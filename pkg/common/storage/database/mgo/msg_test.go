package mgo

import (
	"context"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestName1(t *testing.T) {
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
	cli := Result(mongo.Connect(ctx, options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.135:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))

	msg, err := NewMsgMongo(cli.Database("openim_v3"))
	if err != nil {
		panic(err)
	}
	ts := time.Now().Add(-time.Hour * 24 * 5).UnixMilli()
	t.Log(ts)
	res, err := msg.GetLastMessageSeqByTime(ctx, "sg_1523453548", ts)
	if err != nil {
		panic(err)
	}
	t.Log(res)
}

func TestName5(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	cli := Result(mongo.Connect(ctx, options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.135:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))

	tmp, err := NewMsgMongo(cli.Database("openim_v3"))
	if err != nil {
		panic(err)
	}
	msg := tmp.(*MsgMgo)
	ts := time.Now().Add(-time.Hour * 24 * 5).UnixMilli()
	t.Log(ts)
	var seqs []int64
	for i := 1; i < 256; i++ {
		seqs = append(seqs, int64(i))
	}
	res, err := msg.FindSeqs(ctx, "si_4924054191_9511766539", seqs)
	if err != nil {
		panic(err)
	}
	t.Log(res)
}

func TestSearchMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	cli := Result(mongo.Connect(ctx, options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.135:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))

	msgMongo, err := NewMsgMongo(cli.Database("openim_v3"))
	if err != nil {
		panic(err)
	}
	ts := time.Now().Add(-time.Hour * 24 * 5).UnixMilli()
	t.Log(ts)
	req := &msg.SearchMessageReq{
		//SendID: "yjz",
		//RecvID: "aibot",
		Pagination: &sdkws.RequestPagination{
			PageNumber: 1,
			ShowNumber: 20,
		},
	}
	count, resp, err := msgMongo.SearchMessage(ctx, req)
	if err != nil {
		panic(err)
	}
	t.Log(resp, count)
}
