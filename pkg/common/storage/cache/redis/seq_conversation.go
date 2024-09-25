package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func NewSeqConversationCacheRedis(rdb redis.UniversalClient, mgo database.SeqConversation) cache.SeqConversationCache {
	return &seqConversationCacheRedis{
		rdb:              rdb,
		mgo:              mgo,
		lockTime:         time.Second * 3,
		dataTime:         time.Hour * 24 * 365,
		minSeqExpireTime: time.Hour,
		rocks:            rockscache.NewClient(rdb, *GetRocksCacheOptions()),
	}
}

type seqConversationCacheRedis struct {
	rdb              redis.UniversalClient
	mgo              database.SeqConversation
	rocks            *rockscache.Client
	lockTime         time.Duration
	dataTime         time.Duration
	minSeqExpireTime time.Duration
}

func (s *seqConversationCacheRedis) getMinSeqKey(conversationID string) string {
	return cachekey.GetMallocMinSeqKey(conversationID)
}

func (s *seqConversationCacheRedis) SetMinSeq(ctx context.Context, conversationID string, seq int64) error {
	return s.SetMinSeqs(ctx, map[string]int64{conversationID: seq})
}

func (s *seqConversationCacheRedis) GetMinSeq(ctx context.Context, conversationID string) (int64, error) {
	return getCache(ctx, s.rocks, s.getMinSeqKey(conversationID), s.minSeqExpireTime, func(ctx context.Context) (int64, error) {
		return s.mgo.GetMinSeq(ctx, conversationID)
	})
}

func (s *seqConversationCacheRedis) getSingleMaxSeq(ctx context.Context, conversationID string) (map[string]int64, error) {
	seq, err := s.GetMaxSeq(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	return map[string]int64{conversationID: seq}, nil
}

func (s *seqConversationCacheRedis) getSingleMaxSeqWithTime(ctx context.Context, conversationID string) (map[string]database.SeqTime, error) {
	seq, err := s.GetMaxSeqWithTime(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	return map[string]database.SeqTime{conversationID: seq}, nil
}

func (s *seqConversationCacheRedis) batchGetMaxSeq(ctx context.Context, keys []string, keyConversationID map[string]string, seqs map[string]int64) error {
	result := make([]*redis.StringCmd, len(keys))
	pipe := s.rdb.Pipeline()
	for i, key := range keys {
		result[i] = pipe.HGet(ctx, key, "CURR")
	}
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return errs.Wrap(err)
	}
	var notFoundKey []string
	for i, r := range result {
		req, err := r.Int64()
		if err == nil {
			seqs[keyConversationID[keys[i]]] = req
		} else if errors.Is(err, redis.Nil) {
			notFoundKey = append(notFoundKey, keys[i])
		} else {
			return errs.Wrap(err)
		}
	}
	for _, key := range notFoundKey {
		conversationID := keyConversationID[key]
		seq, err := s.GetMaxSeq(ctx, conversationID)
		if err != nil {
			return err
		}
		seqs[conversationID] = seq
	}
	return nil
}

func (s *seqConversationCacheRedis) batchGetMaxSeqWithTime(ctx context.Context, keys []string, keyConversationID map[string]string, seqs map[string]database.SeqTime) error {
	result := make([]*redis.SliceCmd, len(keys))
	pipe := s.rdb.Pipeline()
	for i, key := range keys {
		result[i] = pipe.HMGet(ctx, key, "CURR", "TIME")
	}
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return errs.Wrap(err)
	}
	var notFoundKey []string
	for i, r := range result {
		val, err := r.Result()
		if len(val) != 2 {
			return errs.WrapMsg(err, "batchGetMaxSeqWithTime invalid result", "key", keys[i], "res", val)
		}
		if val[0] == nil {
			notFoundKey = append(notFoundKey, keys[i])
			continue
		}
		seq, err := s.parseInt64(val[0])
		if err != nil {
			return err
		}
		mill, err := s.parseInt64(val[1])
		if err != nil {
			return err
		}
		seqs[keyConversationID[keys[i]]] = database.SeqTime{Seq: seq, Time: mill}
	}
	for _, key := range notFoundKey {
		conversationID := keyConversationID[key]
		seq, err := s.GetMaxSeqWithTime(ctx, conversationID)
		if err != nil {
			return err
		}
		seqs[conversationID] = seq
	}
	return nil
}

func (s *seqConversationCacheRedis) GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	switch len(conversationIDs) {
	case 0:
		return map[string]int64{}, nil
	case 1:
		return s.getSingleMaxSeq(ctx, conversationIDs[0])
	}
	keys := make([]string, 0, len(conversationIDs))
	keyConversationID := make(map[string]string, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		key := s.getSeqMallocKey(conversationID)
		if _, ok := keyConversationID[key]; ok {
			continue
		}
		keys = append(keys, key)
		keyConversationID[key] = conversationID
	}
	if len(keys) == 1 {
		return s.getSingleMaxSeq(ctx, conversationIDs[0])
	}
	slotKeys, err := groupKeysBySlot(ctx, s.rdb, keys)
	if err != nil {
		return nil, err
	}
	seqs := make(map[string]int64, len(conversationIDs))
	for _, keys := range slotKeys {
		if err := s.batchGetMaxSeq(ctx, keys, keyConversationID, seqs); err != nil {
			return nil, err
		}
	}
	return seqs, nil
}

func (s *seqConversationCacheRedis) GetMaxSeqsWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error) {
	switch len(conversationIDs) {
	case 0:
		return map[string]database.SeqTime{}, nil
	case 1:
		return s.getSingleMaxSeqWithTime(ctx, conversationIDs[0])
	}
	keys := make([]string, 0, len(conversationIDs))
	keyConversationID := make(map[string]string, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		key := s.getSeqMallocKey(conversationID)
		if _, ok := keyConversationID[key]; ok {
			continue
		}
		keys = append(keys, key)
		keyConversationID[key] = conversationID
	}
	if len(keys) == 1 {
		return s.getSingleMaxSeqWithTime(ctx, conversationIDs[0])
	}
	slotKeys, err := groupKeysBySlot(ctx, s.rdb, keys)
	if err != nil {
		return nil, err
	}
	seqs := make(map[string]database.SeqTime, len(conversationIDs))
	for _, keys := range slotKeys {
		if err := s.batchGetMaxSeqWithTime(ctx, keys, keyConversationID, seqs); err != nil {
			return nil, err
		}
	}
	return seqs, nil
}

func (s *seqConversationCacheRedis) getSeqMallocKey(conversationID string) string {
	return cachekey.GetMallocSeqKey(conversationID)
}

func (s *seqConversationCacheRedis) setSeq(ctx context.Context, key string, owner int64, currSeq int64, lastSeq int64, mill int64) (int64, error) {
	if lastSeq < currSeq {
		return 0, errs.New("lastSeq must be greater than currSeq")
	}
	// 0: success
	// 1: success the lock has expired, but has not been locked by anyone else
	// 2: already locked, but not by yourself
	script := `
local key = KEYS[1]
local lockValue = ARGV[1]
local dataSecond = ARGV[2]
local curr_seq = tonumber(ARGV[3])
local last_seq = tonumber(ARGV[4])
local mallocTime = ARGV[5]
if redis.call("EXISTS", key) == 0 then
	redis.call("HSET", key, "CURR", curr_seq, "LAST", last_seq, "TIME", mallocTime)
	redis.call("EXPIRE", key, dataSecond)
	return 1
end
if redis.call("HGET", key, "LOCK") ~= lockValue then
	return 2
end
redis.call("HDEL", key, "LOCK")
redis.call("HSET", key, "CURR", curr_seq, "LAST", last_seq, "TIME", mallocTime)
redis.call("EXPIRE", key, dataSecond)
return 0
`
	result, err := s.rdb.Eval(ctx, script, []string{key}, owner, int64(s.dataTime/time.Second), currSeq, lastSeq, mill).Int64()
	if err != nil {
		return 0, errs.Wrap(err)
	}
	return result, nil
}

// malloc size=0 is to get the current seq size>0 is to allocate seq
func (s *seqConversationCacheRedis) malloc(ctx context.Context, key string, size int64) ([]int64, error) {
	// 0: success
	// 1: need to obtain and lock
	// 2: already locked
	// 3: exceeded the maximum value and locked
	script := `
local key = KEYS[1]
local size = tonumber(ARGV[1])
local lockSecond = ARGV[2]
local dataSecond = ARGV[3]
local mallocTime = ARGV[4]
local result = {}
if redis.call("EXISTS", key) == 0 then
	local lockValue = math.random(0, 999999999)
	redis.call("HSET", key, "LOCK", lockValue)
	redis.call("EXPIRE", key, lockSecond)
	table.insert(result, 1)
	table.insert(result, lockValue)
	table.insert(result, mallocTime)
	return result
end
if redis.call("HEXISTS", key, "LOCK") == 1 then
	table.insert(result, 2)
	return result
end
local curr_seq = tonumber(redis.call("HGET", key, "CURR"))
local last_seq = tonumber(redis.call("HGET", key, "LAST"))
if size == 0 then
	redis.call("EXPIRE", key, dataSecond)
	table.insert(result, 0)
	table.insert(result, curr_seq)
	table.insert(result, last_seq)
	local setTime = redis.call("HGET", key, "TIME")
	if setTime then
		table.insert(result, setTime)	
	else
		table.insert(result, 0)
	end
	return result
end
local max_seq = curr_seq + size
if max_seq > last_seq then
	local lockValue = math.random(0, 999999999)
	redis.call("HSET", key, "LOCK", lockValue)
	redis.call("HSET", key, "CURR", last_seq)
	redis.call("HSET", key, "TIME", mallocTime)
	redis.call("EXPIRE", key, lockSecond)
	table.insert(result, 3)
	table.insert(result, curr_seq)
	table.insert(result, last_seq)
	table.insert(result, lockValue)
	table.insert(result, mallocTime)
	return result
end
redis.call("HSET", key, "CURR", max_seq)
redis.call("HSET", key, "TIME", ARGV[4])
redis.call("EXPIRE", key, dataSecond)
table.insert(result, 0)
table.insert(result, curr_seq)
table.insert(result, last_seq)
table.insert(result, mallocTime)
return result
`
	result, err := s.rdb.Eval(ctx, script, []string{key}, size, int64(s.lockTime/time.Second), int64(s.dataTime/time.Second), time.Now().UnixMilli()).Int64Slice()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return result, nil
}

func (s *seqConversationCacheRedis) wait(ctx context.Context) error {
	timer := time.NewTimer(time.Second / 4)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *seqConversationCacheRedis) setSeqRetry(ctx context.Context, key string, owner int64, currSeq int64, lastSeq int64, mill int64) {
	for i := 0; i < 10; i++ {
		state, err := s.setSeq(ctx, key, owner, currSeq, lastSeq, mill)
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

func (s *seqConversationCacheRedis) getMallocSize(conversationID string, size int64) int64 {
	if size == 0 {
		return 0
	}
	var basicSize int64
	if msgprocessor.IsGroupConversationID(conversationID) {
		basicSize = 100
	} else {
		basicSize = 50
	}
	basicSize += size
	return basicSize
}

func (s *seqConversationCacheRedis) Malloc(ctx context.Context, conversationID string, size int64) (int64, error) {
	seq, _, err := s.mallocTime(ctx, conversationID, size)
	return seq, err
}

func (s *seqConversationCacheRedis) mallocTime(ctx context.Context, conversationID string, size int64) (int64, int64, error) {
	if size < 0 {
		return 0, 0, errs.New("size must be greater than 0")
	}
	key := s.getSeqMallocKey(conversationID)
	for i := 0; i < 10; i++ {
		states, err := s.malloc(ctx, key, size)
		if err != nil {
			return 0, 0, err
		}
		switch states[0] {
		case 0: // success
			return states[1], states[3], nil
		case 1: // not found
			mallocSize := s.getMallocSize(conversationID, size)
			seq, err := s.mgo.Malloc(ctx, conversationID, mallocSize)
			if err != nil {
				return 0, 0, err
			}
			s.setSeqRetry(ctx, key, states[1], seq+size, seq+mallocSize, states[2])
			return seq, 0, nil
		case 2: // locked
			if err := s.wait(ctx); err != nil {
				return 0, 0, err
			}
			continue
		case 3: // exceeded cache max value
			currSeq := states[1]
			lastSeq := states[2]
			mill := states[4]
			mallocSize := s.getMallocSize(conversationID, size)
			seq, err := s.mgo.Malloc(ctx, conversationID, mallocSize)
			if err != nil {
				return 0, 0, err
			}
			if lastSeq == seq {
				s.setSeqRetry(ctx, key, states[3], currSeq+size, seq+mallocSize, mill)
				return currSeq, states[4], nil
			} else {
				log.ZWarn(ctx, "malloc seq not equal cache last seq", nil, "conversationID", conversationID, "currSeq", currSeq, "lastSeq", lastSeq, "mallocSeq", seq)
				s.setSeqRetry(ctx, key, states[3], seq+size, seq+mallocSize, mill)
				return seq, mill, nil
			}
		default:
			log.ZError(ctx, "malloc seq unknown state", nil, "state", states[0], "conversationID", conversationID, "size", size)
			return 0, 0, errs.New(fmt.Sprintf("unknown state: %d", states[0]))
		}
	}
	log.ZError(ctx, "malloc seq retrying still failed", nil, "conversationID", conversationID, "size", size)
	return 0, 0, errs.New("malloc seq waiting for lock timeout", "conversationID", conversationID, "size", size)
}

func (s *seqConversationCacheRedis) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	return s.Malloc(ctx, conversationID, 0)
}

func (s *seqConversationCacheRedis) GetMaxSeqWithTime(ctx context.Context, conversationID string) (database.SeqTime, error) {
	seq, mill, err := s.mallocTime(ctx, conversationID, 0)
	if err != nil {
		return database.SeqTime{}, err
	}
	return database.SeqTime{Seq: seq, Time: mill}, nil
}

func (s *seqConversationCacheRedis) SetMinSeqs(ctx context.Context, seqs map[string]int64) error {
	keys := make([]string, 0, len(seqs))
	for conversationID, seq := range seqs {
		keys = append(keys, s.getMinSeqKey(conversationID))
		if err := s.mgo.SetMinSeq(ctx, conversationID, seq); err != nil {
			return err
		}
	}
	return DeleteCacheBySlot(ctx, s.rocks, keys)
}

// GetCacheMaxSeqWithTime only get the existing cache, if there is no cache, no cache will be generated
func (s *seqConversationCacheRedis) GetCacheMaxSeqWithTime(ctx context.Context, conversationIDs []string) (map[string]database.SeqTime, error) {
	if len(conversationIDs) == 0 {
		return map[string]database.SeqTime{}, nil
	}
	key2conversationID := make(map[string]string)
	keys := make([]string, 0, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		key := s.getSeqMallocKey(conversationID)
		if _, ok := key2conversationID[key]; ok {
			continue
		}
		key2conversationID[key] = conversationID
		keys = append(keys, key)
	}
	slotKeys, err := groupKeysBySlot(ctx, s.rdb, keys)
	if err != nil {
		return nil, err
	}
	res := make(map[string]database.SeqTime)
	for _, keys := range slotKeys {
		if len(keys) == 0 {
			continue
		}
		pipe := s.rdb.Pipeline()
		cmds := make([]*redis.SliceCmd, 0, len(keys))
		for _, key := range keys {
			cmds = append(cmds, pipe.HMGet(ctx, key, "CURR", "TIME"))
		}
		if _, err := pipe.Exec(ctx); err != nil {
			return nil, errs.Wrap(err)
		}
		for i, cmd := range cmds {
			val, err := cmd.Result()
			if err != nil {
				return nil, err
			}
			if len(val) != 2 {
				return nil, errs.WrapMsg(err, "GetCacheMaxSeqWithTime invalid result", "key", keys[i], "res", val)
			}
			if val[0] == nil {
				continue
			}
			seq, err := s.parseInt64(val[0])
			if err != nil {
				return nil, err
			}
			mill, err := s.parseInt64(val[1])
			if err != nil {
				return nil, err
			}
			conversationID := key2conversationID[keys[i]]
			res[conversationID] = database.SeqTime{Seq: seq, Time: mill}
		}
	}
	return res, nil
}

func (s *seqConversationCacheRedis) parseInt64(val any) (int64, error) {
	switch v := val.(type) {
	case nil:
		return 0, nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errs.WrapMsg(err, "invalid string not int64", "value", v)
		}
		return res, nil
	default:
		return 0, errs.New("invalid result not int64", "resType", fmt.Sprintf("%T", v), "value", v)
	}
}
