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

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestName(t *testing.T) {
	cli := Result(mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.48:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second)))
	conv := Result(NewConversationMongo(cli.Database("openim_v3")))
	num, err := conv.GetAllConversationIDsNumber(context.Background())
	if err != nil {
		panic(err)
	}
	t.Log(num)
	ids := Result(conv.GetAllConversationIDs(context.Background()))
	t.Log(ids)
}
