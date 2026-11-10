package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
)

func NewSeqUserCacheRedis(rdb redis.UniversalClient, mgo database.SeqUser) cache.SeqUser {
	return &seqUserCacheRedis{
		mgo:               mgo,
		readSeqWriteRatio: 100,
		expireTime:        time.Hour * 24 * 7,
		readExpireTime:    time.Hour * 24 * 30,
		rocks:             newRocksCacheClient(rdb),
	}
}

type seqUserCacheRedis struct {
	mgo               database.SeqUser
	rocks             *rocksCacheClient
	expireTime        time.Duration
	readExpireTime    time.Duration
	readSeqWriteRatio int64
}

func (s *seqUserCacheRedis) getSeqUserMaxSeqKey(conversationID string, userID string) string {
	return cachekey.GetSeqUserMaxSeqKey(conversationID, userID)
}

func (s *seqUserCacheRedis) getSeqUserMinSeqKey(conversationID string, userID string) string {
	return cachekey.GetSeqUserMinSeqKey(conversationID, userID)
}

func (s *seqUserCacheRedis) getSeqUserReadSeqKey(conversationID string, userID string) string {
	return cachekey.GetSeqUserReadSeqKey(conversationID, userID)
}

func (s *seqUserCacheRedis) GetUserMaxSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	return getCache(ctx, s.rocks, s.getSeqUserMaxSeqKey(conversationID, userID), s.expireTime, func(ctx context.Context) (int64, error) {
		return s.mgo.GetUserMaxSeq(ctx, conversationID, userID)
	})
}

func (s *seqUserCacheRedis) SetUserMaxSeq(ctx context.Context, conversationID string, userID string, seq int64) error {
	if err := s.mgo.SetUserMaxSeq(ctx, conversationID, userID, seq); err != nil {
		return err
	}
	return s.rocks.GetClient().TagAsDeleted2(ctx, s.getSeqUserMaxSeqKey(conversationID, userID))
}

func (s *seqUserCacheRedis) GetUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	return getCache(ctx, s.rocks, s.getSeqUserMinSeqKey(conversationID, userID), s.expireTime, func(ctx context.Context) (int64, error) {
		return s.mgo.GetUserMinSeq(ctx, conversationID, userID)
	})
}

func (s *seqUserCacheRedis) SetUserMinSeq(ctx context.Context, conversationID string, userID string, seq int64) error {
	return s.SetUserMinSeqs(ctx, userID, map[string]int64{conversationID: seq})
}

// Read sequences use their own monotonic fill path. A generic RocksCache
// refill can otherwise overwrite a concurrent read receipt with an older DB value.
func (s *seqUserCacheRedis) GetUserReadSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	if s.rocks.Disable() {
		return s.mgo.GetUserReadSeq(ctx, conversationID, userID)
	}
	key := s.getSeqUserReadSeqKey(conversationID, userID)
	value, err := s.rocks.GetRedis().HGet(ctx, key, "value").Result()
	if err == nil && value != "" {
		return strconv.ParseInt(value, 10, 64)
	}
	if err != nil && err != redis.Nil {
		return 0, errs.Wrap(err)
	}
	seq, err := s.mgo.GetUserReadSeq(ctx, conversationID, userID)
	if err != nil {
		return 0, err
	}
	return s.mergeReadSeq(ctx, key, seq)
}

// SetUserReadSeq returns the effective sequence, including any concurrent advance.
func (s *seqUserCacheRedis) SetUserReadSeq(ctx context.Context, conversationID string, userID string, seq int64) (int64, error) {
	if seq < 0 {
		return 0, errs.ErrArgs.WrapMsg("read seq must not be negative")
	}
	if s.rocks.Disable() {
		if err := s.SetUserReadSeqToDB(ctx, conversationID, userID, seq); err != nil {
			return 0, err
		}
		return s.mgo.GetUserReadSeq(ctx, conversationID, userID)
	}
	// Preserve the MongoDB baseline on a cold cache, as well as the value read
	// before a key expires between this read and the atomic update.
	current, err := s.GetUserReadSeq(ctx, conversationID, userID)
	if err != nil {
		return 0, err
	}
	return s.mergeReadSeq(ctx, s.getSeqUserReadSeqKey(conversationID, userID), max(current, seq))
}

// Compare decimal strings rather than Lua doubles to preserve the full int64 range.
// Removing the refill owner also fences any RocksCache fill started by an older
// server instance before this write during a rolling upgrade.
var mergeReadSeqScript = redis.NewScript(`
local incoming = ARGV[1]
local current = redis.call('HGET', KEYS[1], 'value')
if current and current ~= '' then
    if not string.match(current, '^%d+$') then
        return redis.error_reply('invalid cached read sequence')
    end
    if #current > #incoming or (#current == #incoming and current > incoming) then
        incoming = current
    end
end
redis.call('HSET', KEYS[1], 'value', incoming)
redis.call('HDEL', KEYS[1], 'lockUntil', 'lockOwner')
redis.call('PEXPIRE', KEYS[1], ARGV[2])
return incoming
`)

func (s *seqUserCacheRedis) mergeReadSeq(ctx context.Context, key string, seq int64) (int64, error) {
	if seq < 0 {
		return 0, errs.ErrArgs.WrapMsg("read seq must not be negative")
	}
	value, err := mergeReadSeqScript.Run(ctx, s.rocks.GetRedis(), []string{key},
		strconv.FormatInt(seq, 10), s.readExpireTime.Milliseconds()).Text()
	if err != nil {
		return 0, errs.Wrap(err)
	}
	return strconv.ParseInt(value, 10, 64)
}

func (s *seqUserCacheRedis) mergeReadSeqs(ctx context.Context, userID string, seqs map[string]int64) (map[string]int64, error) {
	pipe := s.rocks.GetRedis().Pipeline()
	commands := make(map[string]*redis.Cmd, len(seqs))
	for conversationID, seq := range seqs {
		if seq < 0 {
			return nil, errs.ErrArgs.WrapMsg("read seq must not be negative")
		}
		// Each script touches a single key, so ClusterClient can route the
		// pipeline across slots without a cross-slot transaction.
		commands[conversationID] = mergeReadSeqScript.Eval(ctx, pipe,
			[]string{s.getSeqUserReadSeqKey(conversationID, userID)},
			strconv.FormatInt(seq, 10), s.readExpireTime.Milliseconds())
	}
	if len(commands) > 0 {
		if _, err := pipe.Exec(ctx); err != nil {
			return nil, errs.Wrap(err)
		}
	}
	result := make(map[string]int64, len(commands))
	for conversationID, cmd := range commands {
		value, err := cmd.Text()
		if err != nil {
			return nil, errs.Wrap(err)
		}
		seq, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		result[conversationID] = seq
	}
	return result, nil
}

func (s *seqUserCacheRedis) SetUserReadSeqToDB(ctx context.Context, conversationID string, userID string, seq int64) error {
	return s.mgo.SetUserReadSeq(ctx, conversationID, userID, seq)
}

func (s *seqUserCacheRedis) SetUserMinSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	keys := make([]string, 0, len(seqs))
	for conversationID, seq := range seqs {
		if err := s.mgo.SetUserMinSeq(ctx, conversationID, userID, seq); err != nil {
			return err
		}
		keys = append(keys, s.getSeqUserMinSeqKey(conversationID, userID))
	}
	return DeleteCacheBySlot(ctx, s.rocks, keys)
}

func (s *seqUserCacheRedis) SetUserReadSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	conversationIDs := make([]string, 0, len(seqs))
	for conversationID, seq := range seqs {
		if seq < 0 {
			return errs.ErrArgs.WrapMsg("read seq must not be negative")
		}
		conversationIDs = append(conversationIDs, conversationID)
	}
	if s.rocks.Disable() {
		for conversationID, seq := range seqs {
			if err := s.SetUserReadSeqToDB(ctx, conversationID, userID, seq); err != nil {
				return err
			}
		}
		return nil
	}
	current, err := s.GetUserReadSeqs(ctx, userID, conversationIDs)
	if err != nil {
		return err
	}
	updates := make(map[string]int64, len(seqs))
	for conversationID, seq := range seqs {
		updates[conversationID] = max(current[conversationID], seq)
	}
	_, err = s.mergeReadSeqs(ctx, userID, updates)
	return err
}

func (s *seqUserCacheRedis) GetUserReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	if s.rocks.Disable() {
		return s.mgo.GetUserReadSeqs(ctx, userID, conversationIDs)
	}
	pipe := s.rocks.GetRedis().Pipeline()
	commands := make(map[string]*redis.StringCmd, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		if _, ok := commands[conversationID]; !ok {
			commands[conversationID] = pipe.HGet(ctx, s.getSeqUserReadSeqKey(conversationID, userID), "value")
		}
	}
	if len(commands) > 0 {
		if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
			return nil, errs.Wrap(err)
		}
	}
	result := make(map[string]int64, len(commands))
	missing := make([]string, 0)
	for conversationID, cmd := range commands {
		value, err := cmd.Result()
		if err == redis.Nil || (err == nil && value == "") {
			missing = append(missing, conversationID)
			continue
		}
		if err != nil {
			return nil, errs.Wrap(err)
		}
		seq, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		result[conversationID] = seq
	}
	if len(missing) == 0 {
		return result, nil
	}
	stored, err := s.mgo.GetUserReadSeqs(ctx, userID, missing)
	if err != nil {
		return nil, err
	}
	fills := make(map[string]int64, len(missing))
	for _, conversationID := range missing {
		fills[conversationID] = stored[conversationID]
	}
	filled, err := s.mergeReadSeqs(ctx, userID, fills)
	if err != nil {
		return nil, err
	}
	for conversationID, seq := range filled {
		result[conversationID] = seq
	}
	return result, nil
}
