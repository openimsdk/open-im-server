package rpccache

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
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

func (a *AuthLocalCache) GetExistingToken(ctx context.Context, userID string, platformID int) (val map[string]int, err error) {
	resp, err := a.getExistingToken(ctx, userID, platformID)
	if err != nil {
		return nil, err
	}

	res := convert.TokenMapPb2DB(resp.TokenStates)

	return res, nil
}

func (a *AuthLocalCache) getExistingToken(ctx context.Context, userID string, platformID int) (val *auth.GetExistingTokenResp, err error) {
	log.ZDebug(ctx, "AuthLocalCache GetExistingToken req", "userID", userID, "platformID", platformID)
	defer func() {
		if err != nil {
			log.ZError(ctx, "AuthLocalCache GetExistingToken error", err, "userID", userID, "platformID", platformID)
		} else {
			log.ZDebug(ctx, "AuthLocalCache GetExistingToken resp", "userID", userID, "platformID", platformID, "val", val)
		}
	}()

	var cache cacheProto[auth.GetExistingTokenResp]

	return cache.Unmarshal(a.local.Get(ctx, cachekey.GetTokenKey(userID, platformID), func(ctx context.Context) ([]byte, error) {
		log.ZDebug(ctx, "AuthLocalCache GetExistingToken call rpc", "userID", userID, "platformID", platformID)
		return cache.Marshal(a.client.AuthClient.GetExistingToken(ctx, &auth.GetExistingTokenReq{UserID: userID, PlatformID: int32(platformID)}))
	}))
}

func (a *AuthLocalCache) RemoveLocalTokenCache(ctx context.Context, userID string, platformID int) {
	key := cachekey.GetTokenKey(userID, platformID)
	a.local.DelLocal(ctx, key)
	log.ZDebug(ctx, "AuthLocalCache RemoveLocalTokenCache", "userID", userID, "platformID", platformID, "key", key)
}
