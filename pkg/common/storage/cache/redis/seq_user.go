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
	if err := s.mgo.SetMinSeq(ctx, conversationID, userID, seq); err != nil {
		return err
	}
	return s.rocks.TagAsDeleted2(ctx, s.getSeqUserMinSeqKey(conversationID, userID))
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
