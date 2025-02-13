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

func (s *seqUserCacheRedis) GetUserReadSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	return getCache(ctx, s.rocks, s.getSeqUserReadSeqKey(conversationID, userID), s.readExpireTime, func(ctx context.Context) (int64, error) {
		return s.mgo.GetUserReadSeq(ctx, conversationID, userID)
	})
}

func (s *seqUserCacheRedis) SetUserReadSeq(ctx context.Context, conversationID string, userID string, seq int64) error {
	if s.rocks.GetRedis() == nil {
		return s.SetUserReadSeqToDB(ctx, conversationID, userID, seq)
	}
	dbSeq, err := s.GetUserReadSeq(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if dbSeq < seq {
		if err := s.rocks.GetClient().RawSet(ctx, s.getSeqUserReadSeqKey(conversationID, userID), strconv.Itoa(int(seq)), s.readExpireTime); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
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

func (s *seqUserCacheRedis) setUserRedisReadSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	keys := make([]string, 0, len(seqs))
	keySeq := make(map[string]int64)
	for conversationID, seq := range seqs {
		key := s.getSeqUserReadSeqKey(conversationID, userID)
		keys = append(keys, key)
		keySeq[key] = seq
	}
	slotKeys, err := groupKeysBySlot(ctx, s.rocks.GetRedis(), keys)
	if err != nil {
		return err
	}
	for _, keys := range slotKeys {
		pipe := s.rocks.GetRedis().Pipeline()
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

func (s *seqUserCacheRedis) SetUserReadSeqs(ctx context.Context, userID string, seqs map[string]int64) error {
	if len(seqs) == 0 {
		return nil
	}
	if err := s.setUserRedisReadSeqs(ctx, userID, seqs); err != nil {
		return err
	}
	return nil
}

func (s *seqUserCacheRedis) GetUserReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	res, err := batchGetCache2(ctx, s.rocks, s.readExpireTime, conversationIDs, func(conversationID string) string {
		return s.getSeqUserReadSeqKey(conversationID, userID)
	}, func(v *readSeqModel) string {
		return v.ConversationID
	}, func(ctx context.Context, conversationIDs []string) ([]*readSeqModel, error) {
		seqs, err := s.mgo.GetUserReadSeqs(ctx, userID, conversationIDs)
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
