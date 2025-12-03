package mgo

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestName1111(t *testing.T) {
	coll := Mongodb().Collection("temp")
	pipeline := mongo.Pipeline{
		{
			{"$set", bson.M{
				"value": bson.M{
					"$toString": bson.M{
						"$add": bson.A{
							bson.M{"$toInt": "$value"},
							1,
						},
					},
				},
			}},
		},
	}

	opt := options.FindOneAndUpdate().SetUpsert(false).SetReturnDocument(options.After)
	res, err := mongoutil.FindOneAndUpdate[model.Cache](context.Background(), coll, bson.M{"key": "123456"}, pipeline, opt)
	if err != nil {
		panic(err)
	}
	t.Log(res)
}

func TestName33333(t *testing.T) {
	c, err := NewCacheMgo(Mongodb())
	if err != nil {
		panic(err)
	}
	if err := c.Set(context.Background(), "123456", "123456", time.Hour); err != nil {
		panic(err)
	}

	if err := c.Set(context.Background(), "123666", "123666", time.Hour); err != nil {
		panic(err)
	}

	res1, err := c.Get(context.Background(), []string{"123456"})
	if err != nil {
		panic(err)
	}
	t.Log(res1)

	res2, err := c.Prefix(context.Background(), "123")
	if err != nil {
		panic(err)
	}
	t.Log(res2)
}

func TestName1111aa(t *testing.T) {

	c, err := NewCacheMgo(Mongodb())
	if err != nil {
		panic(err)
	}
	var count int

	key := "123456"

	doFunc := func() {
		value, err := c.Lock(context.Background(), key, time.Second*30)
		if err != nil {
			t.Log("Lock error", err)
			return
		}
		tmp := count
		tmp++
		count = tmp
		t.Log("count", tmp)
		if err := c.Unlock(context.Background(), key, value); err != nil {
			t.Log("Unlock error", err)
			return
		}
	}

	if _, err := c.Lock(context.Background(), key, time.Second*10); err != nil {
		t.Log(err)
		return
	}

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				doFunc()
			}
		}()
	}

	wg.Wait()

}

func TestName111111a(t *testing.T) {
	arr := strings.SplitN("1:testkakskdask:1111", ":", 2)
	t.Log(arr)
}
