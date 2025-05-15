package mgo

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Result[V any](val V, err error) V {
	if err != nil {
		panic(err)
	}
	return val
}

func Mongodb() *mongo.Database {
	return Result(
		mongo.Connect(context.Background(),
			options.Client().
				ApplyURI("mongodb://openIM:openIM123@172.16.8.135:37017/openim_v3?maxPoolSize=100").
				SetConnectTimeout(5*time.Second)),
	).Database("openim_v3")
}

func TestUserSeq(t *testing.T) {
	uSeq := Result(NewSeqUserMongo(Mongodb())).(*seqUserMongo)
	t.Log(uSeq.SetUserMinSeq(context.Background(), "1000", "2000", 4))
}

func TestConversationSeq(t *testing.T) {
	cSeq := Result(NewSeqConversationMongo(Mongodb())).(*seqConversationMongo)
	t.Log(cSeq.SetMaxSeq(context.Background(), "2000", 10))
	t.Log(cSeq.Malloc(context.Background(), "2000", 10))
	t.Log(cSeq.GetMaxSeq(context.Background(), "2000"))
}

func TestUserGetUserReadSeqs(t *testing.T) {
	uSeq := Result(NewSeqUserMongo(Mongodb())).(*seqUserMongo)
	t.Log(uSeq.GetUserReadSeqs(context.Background(), "2110910952", []string{"sg_345762580", "2000", "3000"}))
}
