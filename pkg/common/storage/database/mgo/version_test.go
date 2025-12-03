package mgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestName(t *testing.T) {
	cli := Result(mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.48:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))
	coll := cli.Database("openim_v3").Collection("version_test")
	tmp, err := NewVersionLog(coll)
	if err != nil {
		panic(err)
	}
	vl := tmp.(*VersionLogMgo)
	res, err := vl.incrVersionResult(context.Background(), "100", []string{"1000", "1001", "1003"}, model.VersionStateInsert)
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("%+v", res)
}
