// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"context"
	"errors"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/stringutil"
	"github.com/redis/go-redis/v9"
	"sync"
)

func NewSeqCache(rdb redis.UniversalClient) cache.SeqCache {
	return &seqCache{rdb: rdb}
}

type seqCache struct {
	rdb redis.UniversalClient
}

func (c *seqCache) getMaxSeqKey(conversationID string) string {
	return cachekey.GetMaxSeqKey(conversationID)
}

func (c *seqCache) getMinSeqKey(conversationID string) string {
	return cachekey.GetMinSeqKey(conversationID)
}

func (c *seqCache) getHasReadSeqKey(conversationID string, userID string) string {
	return cachekey.GetHasReadSeqKey(conversationID, userID)
}

func (c *seqCache) getConversationUserMinSeqKey(conversationID, userID string) string {
	return cachekey.GetConversationUserMinSeqKey(conversationID, userID)
}

func (c *seqCache) setSeq(ctx context.Context, conversationID string, seq int64, getkey func(conversationID string) string) error {
	return errs.Wrap(c.rdb.Set(ctx, getkey(conversationID), seq, 0).Err())
}

func (c *seqCache) getSeq(ctx context.Context, conversationID string, getkey func(conversationID string) string) (int64, error) {
	val, err := c.rdb.Get(ctx, getkey(conversationID)).Int64()
	if err != nil {
		return 0, errs.Wrap(err)
	}
	return val, nil
}

func (c *seqCache) getSeqs(ctx context.Context, items []string, getkey func(s string) string) (m map[string]int64, err error) {
	m = make(map[string]int64, len(items))
	var (
		reverseMap = make(map[string]string, len(items))
		keys       = make([]string, len(items))
		lock       sync.Mutex
	)

	for i, v := range items {
		keys[i] = getkey(v)
		reverseMap[getkey(v)] = v
	}

	manager := NewRedisShardManager(c.rdb)
	if err = manager.ProcessKeysBySlot(ctx, keys, func(ctx context.Context, _ int64, keys []string) error {
		res, err := c.rdb.MGet(ctx, keys...).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return errs.Wrap(err)
		}

		// len(res) <= len(items)
		for i := range res {
			strRes, ok := res[i].(string)
			if !ok {
				continue
			}
			val := stringutil.StringToInt64(strRes)
			if val != 0 {
				lock.Lock()
				m[reverseMap[keys[i]]] = val
				lock.Unlock()
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return m, nil
}

func (c *seqCache) SetMaxSeq(ctx context.Context, conversationID string, maxSeq int64) error {
	return c.setSeq(ctx, conversationID, maxSeq, c.getMaxSeqKey)
}

func (c *seqCache) GetMaxSeqs(ctx context.Context, conversationIDs []string) (m map[string]int64, err error) {
	return c.getSeqs(ctx, conversationIDs, c.getMaxSeqKey)
}

func (c *seqCache) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	return c.getSeq(ctx, conversationID, c.getMaxSeqKey)
}

func (c *seqCache) SetMinSeq(ctx context.Context, conversationID string, minSeq int64) error {
	return c.setSeq(ctx, conversationID, minSeq, c.getMinSeqKey)
}

func (c *seqCache) setSeqs(ctx context.Context, seqs map[string]int64, getkey func(key string) string) error {
	for conversationID, seq := range seqs {
		if err := c.rdb.Set(ctx, getkey(conversationID), seq, 0).Err(); err != nil {
			return errs.Wrap(err)
		}
	}
	return nil
}

func (c *seqCache) SetMinSeqs(ctx context.Context, seqs map[string]int64) error {
	return c.setSeqs(ctx, seqs, c.getMinSeqKey)
}

func (c *seqCache) GetMinSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	return c.getSeqs(ctx, conversationIDs, c.getMinSeqKey)
}

func (c *seqCache) GetMinSeq(ctx context.Context, conversationID string) (int64, error) {
	return c.getSeq(ctx, conversationID, c.getMinSeqKey)
}

func (c *seqCache) GetConversationUserMinSeq(ctx context.Context, conversationID string, userID string) (int64, error) {
	val, err := c.rdb.Get(ctx, c.getConversationUserMinSeqKey(conversationID, userID)).Int64()
	if err != nil {
		return 0, errs.Wrap(err)
	}
	return val, nil
}

func (c *seqCache) GetConversationUserMinSeqs(ctx context.Context, conversationID string, userIDs []string) (m map[string]int64, err error) {
	return c.getSeqs(ctx, userIDs, func(userID string) string {
		return c.getConversationUserMinSeqKey(conversationID, userID)
	})
}

func (c *seqCache) SetConversationUserMinSeq(ctx context.Context, conversationID string, userID string, minSeq int64) error {
	return errs.Wrap(c.rdb.Set(ctx, c.getConversationUserMinSeqKey(conversationID, userID), minSeq, 0).Err())
}

func (c *seqCache) SetConversationUserMinSeqs(ctx context.Context, conversationID string, seqs map[string]int64) (err error) {
	return c.setSeqs(ctx, seqs, func(userID string) string {
		return c.getConversationUserMinSeqKey(conversationID, userID)
	})
}

func (c *seqCache) SetUserConversationsMinSeqs(ctx context.Context, userID string, seqs map[string]int64) (err error) {
	return c.setSeqs(ctx, seqs, func(conversationID string) string {
		return c.getConversationUserMinSeqKey(conversationID, userID)
	})
}

func (c *seqCache) SetHasReadSeq(ctx context.Context, userID string, conversationID string, hasReadSeq int64) error {
	return errs.Wrap(c.rdb.Set(ctx, c.getHasReadSeqKey(conversationID, userID), hasReadSeq, 0).Err())
}

func (c *seqCache) SetHasReadSeqs(ctx context.Context, conversationID string, hasReadSeqs map[string]int64) error {
	return c.setSeqs(ctx, hasReadSeqs, func(userID string) string {
		return c.getHasReadSeqKey(conversationID, userID)
	})
}

func (c *seqCache) UserSetHasReadSeqs(ctx context.Context, userID string, hasReadSeqs map[string]int64) error {
	return c.setSeqs(ctx, hasReadSeqs, func(conversationID string) string {
		return c.getHasReadSeqKey(conversationID, userID)
	})
}

func (c *seqCache) GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	return c.getSeqs(ctx, conversationIDs, func(conversationID string) string {
		return c.getHasReadSeqKey(conversationID, userID)
	})
}

func (c *seqCache) GetHasReadSeq(ctx context.Context, userID string, conversationID string) (int64, error) {
	val, err := c.rdb.Get(ctx, c.getHasReadSeqKey(conversationID, userID)).Int64()
	if err != nil {
		return 0, err
	}
	return val, nil
}
