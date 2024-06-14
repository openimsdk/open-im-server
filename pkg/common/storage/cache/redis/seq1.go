package redis

import (
	"context"
	"fmt"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type RedisHash struct {
	NextSeq int64
	LastSeq int64
}

func NewTestSeq() *SeqMalloc {
	mgocli, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://openIM:openIM123@172.16.8.48:37017/openim_v3?maxPoolSize=100").SetConnectTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}
	model, err := mgo.NewSeqMongo(mgocli.Database("openim_v3"))
	if err != nil {
		panic(err)
	}
	opt := &redis.Options{
		Addr:     "172.16.8.48:16379",
		Password: "openIM123",
		DB:       1,
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return &SeqMalloc{
		rdb: rdb,
		mgo: model,
		//lockTime: time.Second * 30,
		lockTime: time.Second * 60 * 60 * 24 * 1,
		dataTime: time.Second * 60 * 60 * 24 * 7,
	}
}

type SeqMalloc struct {
	rdb      redis.UniversalClient
	mgo      database.Seq
	lockTime time.Duration
	dataTime time.Duration
}

func (s *SeqMalloc) getSeqMallocKey(conversationID string) string {
	return cachekey.GetMallocSeqKey(conversationID)
}

func (s *SeqMalloc) setSeq(ctx context.Context, key string, owner int64, currSeq int64, lastSeq int64) (int64, error) {
	if lastSeq < currSeq {
		return 0, errs.New("lastSeq must be greater than currSeq")
	}
	// 0： 成功
	// 1： 成功 锁过期，但未被其他人锁
	// 2： 已经被锁，但是锁的不是自己
	script := `
local key = KEYS[1]
local lockValue = ARGV[1]
local dataSecond = ARGV[2]
local curr_seq = tonumber(ARGV[3])
local last_seq = tonumber(ARGV[4])
if redis.call("EXISTS", key) == 0 then
	redis.call("HSET", key, "CURR", curr_seq, "LAST", last_seq)
	redis.call("EXPIRE", key, dataSecond)
	return 1
end
if redis.call("HGET", key, "LOCK") ~= lockValue then
	return 2
end
redis.call("HDEL", key, "LOCK")
redis.call("HSET", key, "CURR", curr_seq, "LAST", last_seq)
redis.call("EXPIRE", key, dataSecond)
return 0
`
	result, err := s.rdb.Eval(ctx, script, []string{key}, owner, int64(s.dataTime/time.Second), currSeq, lastSeq).Int64()
	if err != nil {
		return 0, errs.Wrap(err)
	}
	return result, nil
}

// malloc size=0为获取当前seq size>0为分配seq
func (s *SeqMalloc) malloc(ctx context.Context, key string, size int64) ([]int64, error) {
	// 0： 成功
	// 1： 需要获取，并加锁
	// 2： 已经被锁
	// 3： 超过最大值，并加锁
	script := `
local key = KEYS[1]
local size = tonumber(ARGV[1])
local lockSecond = ARGV[2]
local dataSecond = ARGV[3]
local result = {}
if redis.call("EXISTS", key) == 0 then
	local lockValue = math.random(0, 999999999)
	redis.call("HSET", key, "LOCK", lockValue)
	redis.call("EXPIRE", key, lockSecond)
	table.insert(result, 1)
	table.insert(result, lockValue)
	return result
end
if redis.call("HEXISTS", key, "LOCK") == 1 then
	table.insert(result, 2)
	return result
end
local curr_seq = tonumber(redis.call("HGET", key, "CURR"))
local last_seq = tonumber(redis.call("HGET", key, "LAST"))
if size == 0 then
	table.insert(result, 0)
	table.insert(result, curr_seq)
	table.insert(result, last_seq)
	return result
end
local max_seq = curr_seq + size
if max_seq > last_seq then
	local lockValue = math.random(0, 999999999)
	redis.call("HSET", key, "LOCK", lockValue)
	redis.call("HSET", key, "CURR", last_seq)
	redis.call("EXPIRE", key, lockSecond)
	table.insert(result, 3)
	table.insert(result, curr_seq)
	table.insert(result, last_seq)
	table.insert(result, lockValue)
	return result
end
redis.call("HSET", key, "CURR", max_seq)
table.insert(result, 0)
table.insert(result, curr_seq)
table.insert(result, last_seq)
return result
`
	result, err := s.rdb.Eval(ctx, script, []string{key}, size, int64(s.lockTime/time.Second), int64(s.dataTime/time.Second)).Int64Slice()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return result, nil
}

func (s *SeqMalloc) wait(ctx context.Context) error {
	timer := time.NewTimer(time.Second / 4)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *SeqMalloc) setSeqRetry(ctx context.Context, key string, owner int64, currSeq int64, lastSeq int64) {
	for i := 0; i < 10; i++ {
		state, err := s.setSeq(ctx, key, owner, currSeq, lastSeq)
		if err != nil {
			log.ZError(ctx, "set seq cache failed", err, "key", key, "owner", owner, "currSeq", currSeq, "lastSeq", lastSeq, "count", i+1)
			if err := s.wait(ctx); err != nil {
				return
			}
			continue
		}
		switch state {
		case 0: // ideal state
		case 1:
			log.ZWarn(ctx, "set seq cache lock not found", nil, "key", key, "owner", owner, "currSeq", currSeq, "lastSeq", lastSeq)
		case 2:
			log.ZWarn(ctx, "set seq cache lock to be held by someone else", nil, "key", key, "owner", owner, "currSeq", currSeq, "lastSeq", lastSeq)
		default:
			log.ZError(ctx, "set seq cache lock unknown state", nil, "key", key, "owner", owner, "currSeq", currSeq, "lastSeq", lastSeq)
		}
		return
	}
	log.ZError(ctx, "set seq cache retrying still failed", nil, "key", key, "owner", owner, "currSeq", currSeq, "lastSeq", lastSeq)
}

func (s *SeqMalloc) getMallocSize(conversationID string, size int64) int64 {
	if size == 0 {
		return 0
	}
	var basicSize int64
	if msgprocessor.IsGroupConversationID(conversationID) {
		basicSize = 200
	} else {
		basicSize = 50
	}
	basicSize += size
	return basicSize
}

func (s *SeqMalloc) Malloc(ctx context.Context, conversationID string, size int64) (int64, error) {
	if size < 0 {
		return 0, errs.New("size must be greater than 0")
	}
	key := s.getSeqMallocKey(conversationID)
	for i := 0; i < 10; i++ {
		states, err := s.malloc(ctx, key, size)
		if err != nil {
			return 0, err
		}
		switch states[0] {
		case 0: // success
			return states[1], nil
		case 1: // not found
			mallocSize := s.getMallocSize(conversationID, size)
			seq, err := s.mgo.Malloc(ctx, conversationID, mallocSize)
			if err != nil {
				return 0, err
			}
			s.setSeqRetry(ctx, key, states[1], seq+size, seq+mallocSize)
			return seq, nil
		case 2: // locked
			if err := s.wait(ctx); err != nil {
				return 0, err
			}
			continue
		case 3: // exceeded cache max value
			currSeq := states[1]
			lastSeq := states[2]
			mallocSize := s.getMallocSize(conversationID, size)
			seq, err := s.mgo.Malloc(ctx, conversationID, mallocSize)
			if err != nil {
				return 0, err
			}
			if lastSeq == seq {
				s.setSeqRetry(ctx, key, states[3], currSeq+size, seq+mallocSize)
				return currSeq, nil
			} else {
				log.ZWarn(ctx, "malloc seq not equal cache last seq", nil, "conversationID", conversationID, "currSeq", currSeq, "lastSeq", lastSeq, "mallocSeq", seq)
				s.setSeqRetry(ctx, key, states[3], seq+size, seq+mallocSize)
				return seq, nil
			}
		default:
			log.ZError(ctx, "malloc seq unknown state", nil, "state", states[0], "conversationID", conversationID, "size", size)
			return 0, errs.New(fmt.Sprintf("unknown state: %d", states[0]))
		}
	}
	log.ZError(ctx, "malloc seq retrying still failed", nil, "conversationID", conversationID, "size", size)
	return 0, errs.New("malloc seq waiting for lock timeout", "conversationID", conversationID, "size", size)
}

func (s *SeqMalloc) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	return s.Malloc(ctx, conversationID, 0)
}
