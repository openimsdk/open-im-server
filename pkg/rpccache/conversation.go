package rpccache

import (
	"context"
	pbconversation "github.com/OpenIMSDK/protocol/conversation"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/redis/go-redis/v9"
)

func NewConversationLocalCache(client rpcclient.ConversationRpcClient, cli redis.UniversalClient) *ConversationLocalCache {
	lc := config.Config.LocalCache.Conversation
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
	conversations := make([]*pbconversation.Conversation, 0, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		conversation, err := c.GetConversation(ctx, ownerUserID, conversationID)
		if err != nil {
			if errs.ErrRecordNotFound.Is(err) {
				continue
			}
			return nil, err
		}
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
