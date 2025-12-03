package redis

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newTestSeq() *seqConversationCacheRedis {
	mgocli, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://openIM:openIM123@127.0.0.1:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}
	model, err := mgo.NewSeqConversationMongo(mgocli.Database("openim_v3"))
	if err != nil {
		panic(err)
	}
	opt := &redis.Options{
		Addr:     "127.0.0.1:16379",
		Password: "openIM123",
		DB:       1,
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return NewSeqConversationCacheRedis(rdb, model).(*seqConversationCacheRedis)
}

func TestSeq(t *testing.T) {
	ts := newTestSeq()
	var (
		wg    sync.WaitGroup
		speed atomic.Int64
	)

	const count = 128
	wg.Add(count)
	for i := 0; i < count; i++ {
		index := i + 1
		go func() {
			defer wg.Done()
			var size int64 = 10
			cID := strconv.Itoa(index * 1)
			for i := 1; ; i++ {
				//first, err := ts.mgo.Malloc(context.Background(), cID, size) // mongo
				first, err := ts.Malloc(context.Background(), cID, size) // redis
				if err != nil {
					t.Logf("[%d-%d] %s %s", index, i, cID, err)
					return
				}
				speed.Add(size)
				_ = first
				//t.Logf("[%d] %d -> %d", i, first+1, first+size)
			}
		}()
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			value := speed.Swap(0)
			t.Logf("speed: %d/s", value)
		}
	}
}

func TestDel(t *testing.T) {
	ts := newTestSeq()
	for i := 1; i < 100; i++ {
		var size int64 = 100
		first, err := ts.Malloc(context.Background(), "100", size)
		if err != nil {
			t.Logf("[%d] %s", i, err)
			return
		}
		t.Logf("[%d] %d -> %d", i, first+1, first+size)
		time.Sleep(time.Second)
	}
}

func TestSeqMalloc(t *testing.T) {
	ts := newTestSeq()
	t.Log(ts.GetMaxSeq(context.Background(), "100"))
}

func TestMinSeq(t *testing.T) {
	ts := newTestSeq()
	t.Log(ts.GetMinSeq(context.Background(), "10000000"))
}

func TestMalloc(t *testing.T) {
	ts := newTestSeq()
	t.Log(ts.mallocTime(context.Background(), "10000000", 100))
}

func TestHMGET(t *testing.T) {
	ts := newTestSeq()
	res, err := ts.GetCacheMaxSeqWithTime(context.Background(), []string{"10000000", "123456"})
	if err != nil {
		panic(err)
	}
	t.Log(res)
}

func TestGetMaxSeqWithTime(t *testing.T) {
	ts := newTestSeq()
	t.Log(ts.GetMaxSeqWithTime(context.Background(), "10000000"))
}

func TestGetMaxSeqWithTime1(t *testing.T) {
	ts := newTestSeq()
	t.Log(ts.GetMaxSeqsWithTime(context.Background(), []string{"10000000", "12345", "111"}))
}
