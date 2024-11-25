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
	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/encrypt"
	"github.com/redis/go-redis/v9"
	"math/big"
	"strings"
	"time"
)

const (
	conversationExpireTime = time.Second * 60 * 60 * 12
)

func NewConversationRedis(rdb redis.UniversalClient, localCache *config.LocalCache, opts *rockscache.Options, db database.Conversation) cache.ConversationCache {
	batchHandler := NewBatchDeleterRedis(rdb, opts, []string{localCache.Conversation.Topic})
	c := localCache.Conversation
	log.ZDebug(context.Background(), "conversation local cache init", "Topic", c.Topic, "SlotNum", c.SlotNum, "SlotSize", c.SlotSize, "enable", c.Enable())
	return &ConversationRedisCache{
		BatchDeleter:   batchHandler,
		rcClient:       rockscache.NewClient(rdb, *opts),
		conversationDB: db,
		expireTime:     conversationExpireTime,
	}
}

type ConversationRedisCache struct {
	cache.BatchDeleter
	rcClient       *rockscache.Client
	conversationDB database.Conversation
	expireTime     time.Duration
}

func (c *ConversationRedisCache) CloneConversationCache() cache.ConversationCache {
	return &ConversationRedisCache{
		BatchDeleter:   c.BatchDeleter.Clone(),
		rcClient:       c.rcClient,
		conversationDB: c.conversationDB,
		expireTime:     c.expireTime,
	}
}

func (c *ConversationRedisCache) getConversationKey(ownerUserID, conversationID string) string {
	return cachekey.GetConversationKey(ownerUserID, conversationID)
}

func (c *ConversationRedisCache) getConversationIDsKey(ownerUserID string) string {
	return cachekey.GetConversationIDsKey(ownerUserID)
}

func (c *ConversationRedisCache) getNotNotifyConversationIDsKey(ownerUserID string) string {
	return cachekey.GetNotNotifyConversationIDsKey(ownerUserID)
}

func (c *ConversationRedisCache) getPinnedConversationIDsKey(ownerUserID string) string {
	return cachekey.GetPinnedConversationIDs(ownerUserID)
}

func (c *ConversationRedisCache) getSuperGroupRecvNotNotifyUserIDsKey(groupID string) string {
	return cachekey.GetSuperGroupRecvNotNotifyUserIDsKey(groupID)
}

func (c *ConversationRedisCache) getRecvMsgOptKey(ownerUserID, conversationID string) string {
	return cachekey.GetRecvMsgOptKey(ownerUserID, conversationID)
}

func (c *ConversationRedisCache) getSuperGroupRecvNotNotifyUserIDsHashKey(groupID string) string {
	return cachekey.GetSuperGroupRecvNotNotifyUserIDsHashKey(groupID)
}

func (c *ConversationRedisCache) getConversationHasReadSeqKey(ownerUserID, conversationID string) string {
	return cachekey.GetConversationHasReadSeqKey(ownerUserID, conversationID)
}

func (c *ConversationRedisCache) getConversationNotReceiveMessageUserIDsKey(conversationID string) string {
	return cachekey.GetConversationNotReceiveMessageUserIDsKey(conversationID)
}

func (c *ConversationRedisCache) getUserConversationIDsHashKey(ownerUserID string) string {
	return cachekey.GetUserConversationIDsHashKey(ownerUserID)
}

func (c *ConversationRedisCache) getConversationUserMaxVersionKey(ownerUserID string) string {
	return cachekey.GetConversationUserMaxVersionKey(ownerUserID)
}

func (c *ConversationRedisCache) GetUserConversationIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	return getCache(ctx, c.rcClient, c.getConversationIDsKey(ownerUserID), c.expireTime, func(ctx context.Context) ([]string, error) {
		return c.conversationDB.FindUserIDAllConversationID(ctx, ownerUserID)
	})
}

func (c *ConversationRedisCache) GetUserNotNotifyConversationIDs(ctx context.Context, userID string) ([]string, error) {
	return getCache(ctx, c.rcClient, c.getNotNotifyConversationIDsKey(userID), c.expireTime, func(ctx context.Context) ([]string, error) {
		return c.conversationDB.FindUserIDAllNotNotifyConversationID(ctx, userID)
	})
}

func (c *ConversationRedisCache) GetPinnedConversationIDs(ctx context.Context, userID string) ([]string, error) {
	return getCache(ctx, c.rcClient, c.getPinnedConversationIDsKey(userID), c.expireTime, func(ctx context.Context) ([]string, error) {
		return c.conversationDB.FindUserIDAllPinnedConversationID(ctx, userID)
	})
}

func (c *ConversationRedisCache) DelConversationIDs(userIDs ...string) cache.ConversationCache {
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, c.getConversationIDsKey(userID))
	}
	cache := c.CloneConversationCache()
	cache.AddKeys(keys...)

	return cache
}

func (c *ConversationRedisCache) GetUserConversationIDsHash(ctx context.Context, ownerUserID string) (hash uint64, err error) {
	return getCache(
		ctx,
		c.rcClient,
		c.getUserConversationIDsHashKey(ownerUserID),
		c.expireTime,
		func(ctx context.Context) (uint64, error) {
			conversationIDs, err := c.GetUserConversationIDs(ctx, ownerUserID)
			if err != nil {
				return 0, err
			}
			datautil.Sort(conversationIDs, true)
			bi := big.NewInt(0)
			bi.SetString(encrypt.Md5(strings.Join(conversationIDs, ";"))[0:8], 16)
			return bi.Uint64(), nil
		},
	)
}

func (c *ConversationRedisCache) DelUserConversationIDsHash(ownerUserIDs ...string) cache.ConversationCache {
	keys := make([]string, 0, len(ownerUserIDs))
	for _, ownerUserID := range ownerUserIDs {
		keys = append(keys, c.getUserConversationIDsHashKey(ownerUserID))
	}
	cache := c.CloneConversationCache()
	cache.AddKeys(keys...)

	return cache
}

func (c *ConversationRedisCache) GetConversation(ctx context.Context, ownerUserID, conversationID string) (*model.Conversation, error) {
	return getCache(ctx, c.rcClient, c.getConversationKey(ownerUserID, conversationID), c.expireTime, func(ctx context.Context) (*model.Conversation, error) {
		return c.conversationDB.Take(ctx, ownerUserID, conversationID)
	})
}

func (c *ConversationRedisCache) DelConversations(ownerUserID string, conversationIDs ...string) cache.ConversationCache {
	keys := make([]string, 0, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		keys = append(keys, c.getConversationKey(ownerUserID, conversationID))
	}
	cache := c.CloneConversationCache()
	cache.AddKeys(keys...)

	return cache
}

func (c *ConversationRedisCache) GetConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*model.Conversation, error) {
	return batchGetCache2(ctx, c.rcClient, c.expireTime, conversationIDs, func(conversationID string) string {
		return c.getConversationKey(ownerUserID, conversationID)
	}, func(conversation *model.Conversation) string {
		return conversation.ConversationID
	}, func(ctx context.Context, conversationIDs []string) ([]*model.Conversation, error) {
		return c.conversationDB.Find(ctx, ownerUserID, conversationIDs)
	})
}

func (c *ConversationRedisCache) GetUserAllConversations(ctx context.Context, ownerUserID string) ([]*model.Conversation, error) {
	conversationIDs, err := c.GetUserConversationIDs(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}
	return c.GetConversations(ctx, ownerUserID, conversationIDs)
}

func (c *ConversationRedisCache) GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (opt int, err error) {
	return getCache(ctx, c.rcClient, c.getRecvMsgOptKey(ownerUserID, conversationID), c.expireTime, func(ctx context.Context) (opt int, err error) {
		return c.conversationDB.GetUserRecvMsgOpt(ctx, ownerUserID, conversationID)
	})
}

func (c *ConversationRedisCache) DelUsersConversation(conversationID string, ownerUserIDs ...string) cache.ConversationCache {
	keys := make([]string, 0, len(ownerUserIDs))
	for _, ownerUserID := range ownerUserIDs {
		keys = append(keys, c.getConversationKey(ownerUserID, conversationID))
	}
	cache := c.CloneConversationCache()
	cache.AddKeys(keys...)

	return cache
}

func (c *ConversationRedisCache) DelUserRecvMsgOpt(ownerUserID, conversationID string) cache.ConversationCache {
	cache := c.CloneConversationCache()
	cache.AddKeys(c.getRecvMsgOptKey(ownerUserID, conversationID))

	return cache
}

func (c *ConversationRedisCache) DelSuperGroupRecvMsgNotNotifyUserIDs(groupID string) cache.ConversationCache {
	cache := c.CloneConversationCache()
	cache.AddKeys(c.getSuperGroupRecvNotNotifyUserIDsKey(groupID))

	return cache
}

func (c *ConversationRedisCache) DelSuperGroupRecvMsgNotNotifyUserIDsHash(groupID string) cache.ConversationCache {
	cache := c.CloneConversationCache()
	cache.AddKeys(c.getSuperGroupRecvNotNotifyUserIDsHashKey(groupID))

	return cache
}

func (c *ConversationRedisCache) DelUserAllHasReadSeqs(ownerUserID string, conversationIDs ...string) cache.ConversationCache {
	cache := c.CloneConversationCache()
	for _, conversationID := range conversationIDs {
		cache.AddKeys(c.getConversationHasReadSeqKey(ownerUserID, conversationID))
	}

	return cache
}

func (c *ConversationRedisCache) GetConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) ([]string, error) {
	return getCache(ctx, c.rcClient, c.getConversationNotReceiveMessageUserIDsKey(conversationID), c.expireTime, func(ctx context.Context) ([]string, error) {
		return c.conversationDB.GetConversationNotReceiveMessageUserIDs(ctx, conversationID)
	})
}

func (c *ConversationRedisCache) DelConversationNotReceiveMessageUserIDs(conversationIDs ...string) cache.ConversationCache {
	cache := c.CloneConversationCache()
	for _, conversationID := range conversationIDs {
		cache.AddKeys(c.getConversationNotReceiveMessageUserIDsKey(conversationID))
	}
	return cache
}

func (c *ConversationRedisCache) DelConversationNotNotifyMessageUserIDs(userIDs ...string) cache.ConversationCache {
	cache := c.CloneConversationCache()
	for _, userID := range userIDs {
		cache.AddKeys(c.getNotNotifyConversationIDsKey(userID))
	}
	return cache
}

func (c *ConversationRedisCache) DelConversationPinnedMessageUserIDs(userIDs ...string) cache.ConversationCache {
	cache := c.CloneConversationCache()
	for _, userID := range userIDs {
		cache.AddKeys(c.getPinnedConversationIDsKey(userID))
	}
	return cache
}

func (c *ConversationRedisCache) DelConversationVersionUserIDs(userIDs ...string) cache.ConversationCache {
	cache := c.CloneConversationCache()
	for _, userID := range userIDs {
		cache.AddKeys(c.getConversationUserMaxVersionKey(userID))
	}
	return cache
}

func (c *ConversationRedisCache) FindMaxConversationUserVersion(ctx context.Context, userID string) (*model.VersionLog, error) {
	return getCache(ctx, c.rcClient, c.getConversationUserMaxVersionKey(userID), c.expireTime, func(ctx context.Context) (*model.VersionLog, error) {
		return c.conversationDB.FindConversationUserVersion(ctx, userID, 0, 0)
	})
}
