package mgo

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

func TestName(t *testing.T) {
	cli := Result(mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.48:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))
	tmp, err := NewSeqConversationMongo(cli.Database("openim_v3"))
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		var size int64 = 100
		firstSeq, err := tmp.Malloc(context.Background(), "1", size)
		if err != nil {
			t.Log(err)
			return
		}
		t.Logf("%d -> %d", firstSeq, firstSeq+size-1)
	}
}
