// Copyright © 2024 OpenIM. All rights reserved.
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

package rpccache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

const (
	conversationWorkerCount = 20
)

func NewConversationLocalCache(client rpcclient.ConversationRpcClient, localCache *config.LocalCache, cli redis.UniversalClient) *ConversationLocalCache {
	lc := localCache.Conversation
	log.ZDebug(context.Background(), "ConversationLocalCache", "topic", lc.Topic, "slotNum", lc.SlotNum, "slotSize", lc.SlotSize, "enable", lc.Enable())
	x := &ConversationLocalCache{
		client: client,
		local: localcache.New[any](
			localcache.WithLocalSlotNum(lc.SlotNum),
			localcache.WithLocalSlotSize(lc.SlotSize),
			localcache.WithLinkSlotNum(lc.SlotNum),
			localcache.WithLocalSuccessTTL(lc.Success()),
			localcache.WithLocalFailedTTL(lc.Failed()),
		),
	}
	if lc.Enable() {
		go subscriberRedisDeleteCache(context.Background(), cli, lc.Topic, x.local.DelLocal)
	}
	return x
}

type ConversationLocalCache struct {
	client rpcclient.ConversationRpcClient
	local  localcache.Cache[any]
}

func (c *ConversationLocalCache) GetConversationIDs(ctx context.Context, ownerUserID string) (val []string, err error) {
	log.ZDebug(ctx, "ConversationLocalCache GetConversationIDs req", "ownerUserID", ownerUserID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "ConversationLocalCache GetConversationIDs return", "value", val)
		} else {
			log.ZError(ctx, "ConversationLocalCache GetConversationIDs return", err)
		}
	}()
	return localcache.AnyValue[[]string](c.local.Get(ctx, cachekey.GetConversationIDsKey(ownerUserID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "ConversationLocalCache GetConversationIDs rpc", "ownerUserID", ownerUserID)
		return c.client.GetConversationIDs(ctx, ownerUserID)
	}))
}

func (c *ConversationLocalCache) GetConversation(ctx context.Context, userID, conversationID string) (val *pbconversation.Conversation, err error) {
	log.ZDebug(ctx, "ConversationLocalCache GetConversation req", "userID", userID, "conversationID", conversationID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "ConversationLocalCache GetConversation return", "value", val)
		} else {
			log.ZError(ctx, "ConversationLocalCache GetConversation return", err)
		}
	}()
	return localcache.AnyValue[*pbconversation.Conversation](c.local.Get(ctx, cachekey.GetConversationKey(userID, conversationID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "ConversationLocalCache GetConversation rpc", "userID", userID, "conversationID", conversationID)
		return c.client.GetConversation(ctx, userID, conversationID)
	}))
}

func (c *ConversationLocalCache) GetSingleConversationRecvMsgOpt(ctx context.Context, userID, conversationID string) (int32, error) {
	conv, err := c.GetConversation(ctx, userID, conversationID)
	if err != nil {
		return 0, err
	}
	return conv.RecvMsgOpt, nil
}

func (c *ConversationLocalCache) GetConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*pbconversation.Conversation, error) {
	var (
		conversations     = make([]*pbconversation.Conversation, 0, len(conversationIDs))
		conversationsChan = make(chan *pbconversation.Conversation, len(conversationIDs))
	)

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(conversationWorkerCount)

	for _, conversationID := range conversationIDs {
		conversationID := conversationID
		g.Go(func() error {
			conversation, err := c.GetConversation(ctx, ownerUserID, conversationID)
			if err != nil {
				if errs.ErrRecordNotFound.Is(err) {
					return nil
				}
				return err
			}
			conversationsChan <- conversation
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(conversationsChan)
	for conversation := range conversationsChan {
		conversations = append(conversations, conversation)
	}
	return conversations, nil
}

func (c *ConversationLocalCache) getConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) (*listMap[string], error) {
	return localcache.AnyValue[*listMap[string]](c.local.Get(ctx, cachekey.GetConversationNotReceiveMessageUserIDsKey(conversationID), func(ctx context.Context) (any, error) {
		return newListMap(c.client.GetConversationNotReceiveMessageUserIDs(ctx, conversationID))
	}))
}

func (c *ConversationLocalCache) GetConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) ([]string, error) {
	res, err := c.getConversationNotReceiveMessageUserIDs(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	return res.List, nil
}

func (c *ConversationLocalCache) GetConversationNotReceiveMessageUserIDMap(ctx context.Context, conversationID string) (map[string]struct{}, error) {
	res, err := c.getConversationNotReceiveMessageUserIDs(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	return res.Map, nil
}
