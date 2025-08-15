package rpccache

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

func NewAuthLocalCache(client *rpcli.AuthClient, localCache *config.LocalCache, cli redis.UniversalClient) *AuthLocalCache {
	lc := localCache.Auth
	log.ZDebug(context.Background(), "AuthLocalCache", "topic", lc.Topic, "slotNum", lc.SlotNum, "slotSize", lc.SlotSize, "enable", lc.Enable())
	x := &AuthLocalCache{
		client: client,
		local: localcache.New[[]byte](
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

type AuthLocalCache struct {
	client *rpcli.AuthClient
	local  localcache.Cache[[]byte]
}

// 感觉有点问题 是应该保存token map，还是根据OperationID来保存一个bool

// 也不应该只绑定token 是不是还得绑定其他属性 确认说是这个用户在操作的

func (a *AuthLocalCache) ParseToken(ctx context.Context, token string) (val *auth.ParseTokenResp, err error) {
	log.ZDebug(ctx, "AuthLocalCache ParseToken req", "token", token)
	defer func() {
		if err != nil {
			log.ZError(ctx, "AuthLocalCache ParseToken error", err, "token", token, "err", err)
		} else {
			log.ZDebug(ctx, "AuthLocalCache ParseToken resp", "token", token, "val", val)
		}
	}()

	var cache cacheProto[auth.ParseTokenResp]
	return cache.Unmarshal(a.local.Get(ctx, token, func(ctx context.Context) ([]byte, error) {
		log.ZDebug(ctx, "AuthLocalCache ParseToken call rpc", "token", token)
		return cache.Marshal(a.client.ParseToken(ctx, token))
	}))
}
