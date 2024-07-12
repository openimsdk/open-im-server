package redis

import (
	"context"
	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func NewSeqUserCacheRedis(rdb redis.UniversalClient, mgo database.SeqUser) cache.SeqUser {
	return &seqUserCacheRedis{
		rdb:               rdb,
		mgo:               mgo,
		readSeqWriteRatio: 100,
		expireTime:        time.Hour * 24 * 7,
		readExpireTime:    time.Hour * 24 * 30,
		rocks:             rockscache.NewClient(rdb, *GetRocksCacheOptions()),
	}
}

type seqUserCacheRedis struct {
	rdb               redis.UniversalClient
	mgo               database.SeqUser
	rocks             *rockscache.Client
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

func (s *seqUserCacheRedis) GetMaxSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	return getCache(ctx, s.rocks, s.getSeqUserMaxSeqKey(conversationID, userID), s.expireTime, func(ctx context.Context) (int64, error) {
		return s.mgo.GetMaxSeq(ctx, conversationID, userID)
	})
}

func (s *seqUserCacheRedis) SetMaxSeq(ctx context.Context, conversationID string, userID string, seq int64) error {
	if err := s.mgo.SetMaxSeq(ctx, conversationID, userID, seq); err != nil {
		return err
	}
	return s.rocks.TagAsDeleted2(ctx, s.getSeqUserMaxSeqKey(conversationID, userID))
}

func (s *seqUserCacheRedis) GetMinSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	return getCache(ctx, s.rocks, s.getSeqUserMinSeqKey(conversationID, userID), s.expireTime, func(ctx context.Context) (int64, error) {
		return s.mgo.GetMaxSeq(ctx, conversationID, userID)
	})
}

func (s *seqUserCacheRedis) SetMinSeq(ctx context.Context, conversationID string, userID string, seq int64) error {
	return s.SetMinSeqs(ctx, userID, map[string]int64{conversationID: seq})
}

func (s *seqUserCacheRedis) GetReadSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	return getCache(ctx, s.rocks, s.getSeqUserReadSeqKey(conversationID, userID), s.readExpireTime, func(ctx context.Context) (int64, error) {
		return s.mgo.GetMaxSeq(ctx, conversationID, userID)
	})
}

func (s *seqUserCacheRedis) SetReadSeq(ctx context.Context, conversationID string, userID string, seq int64) error {
	if seq%s.readSeqWriteRatio == 0 {
		if err := s.mgo.SetReadSeq(ctx, conversationID, userID, seq); err != nil {
			return err
		}
	}
	if err := s.rocks.RawSet(ctx, s.getSeqUserReadSeqKey(conversationID, userID), strconv.Itoa(int(seq)), s.readExpireTime); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (s *seqUserCacheRedis) SetMinSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	keys := make([]string, 0, len(seqs))
	for conversationID, seq := range seqs {
		if err := s.mgo.SetMinSeq(ctx, conversationID, userID, seq); err != nil {
			return err
		}
		keys = append(keys, s.getSeqUserMinSeqKey(conversationID, userID))
	}
	return DeleteCacheBySlot(ctx, s.rocks, keys)
}

func (s *seqUserCacheRedis) setRedisReadSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	keys := make([]string, 0, len(seqs))
	keySeq := make(map[string]int64)
	for conversationID, seq := range seqs {
		key := s.getSeqUserReadSeqKey(conversationID, userID)
		keys = append(keys, key)
		keySeq[key] = seq
	}
	slotKeys, err := groupKeysBySlot(ctx, s.rdb, keys)
	if err != nil {
		return err
	}
	for _, keys := range slotKeys {
		pipe := s.rdb.Pipeline()
		for _, key := range keys {
			pipe.HSet(ctx, key, "value", strconv.FormatInt(keySeq[key], 10))
			pipe.Expire(ctx, key, s.readExpireTime)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *seqUserCacheRedis) SetReadSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	if len(seqs) == 0 {
		return nil
	}
	if err := s.setRedisReadSeqs(ctx, userID, seqs); err != nil {
		return err
	}
	for conversationID, seq := range seqs {
		if seq%s.readSeqWriteRatio == 0 {
			if err := s.mgo.SetReadSeq(ctx, conversationID, userID, seq); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *seqUserCacheRedis) GetReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	res, err := batchGetCache2(ctx, s.rocks, s.readExpireTime, conversationIDs, func(conversationID string) string {
		return s.getSeqUserReadSeqKey(conversationID, userID)
	}, func(v *readSeqModel) string {
		return v.ConversationID
	}, func(ctx context.Context, conversationIDs []string) ([]*readSeqModel, error) {
		seqs, err := s.mgo.GetReadSeqs(ctx, userID, conversationIDs)
		if err != nil {
			return nil, err
		}
		res := make([]*readSeqModel, 0, len(seqs))
		for conversationID, seq := range seqs {
			res = append(res, &readSeqModel{ConversationID: conversationID, Seq: seq})
		}
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	data := make(map[string]int64)
	for _, v := range res {
		data[v.ConversationID] = v.Seq
	}
	return data, nil
}

var _ BatchCacheCallback[string] = (*readSeqModel)(nil)

type readSeqModel struct {
	ConversationID string
	Seq            int64
}

func (r *readSeqModel) BatchCache(conversationID string) {
	r.ConversationID = conversationID
}

func (r *readSeqModel) UnmarshalJSON(bytes []byte) (err error) {
	r.Seq, err = strconv.ParseInt(string(bytes), 10, 64)
	return
}

func (r *readSeqModel) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(r.Seq, 10)), nil
}
