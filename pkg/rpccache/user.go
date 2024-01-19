package rpccache

import (
	"context"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/redis/go-redis/v9"
)

func NewUserLocalCache(client rpcclient.UserRpcClient, cli redis.UniversalClient) *UserLocalCache {
	lc := config.Config.LocalCache.User
	log.ZDebug(context.Background(), "UserLocalCache", "topic", lc.Topic, "slotNum", lc.SlotNum, "slotSize", lc.SlotSize, "enable", lc.Enable())
	x := &UserLocalCache{
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

type UserLocalCache struct {
	client rpcclient.UserRpcClient
	local  localcache.Cache[any]
}

func (u *UserLocalCache) GetUserInfo(ctx context.Context, userID string) (val *sdkws.UserInfo, err error) {
	log.ZDebug(ctx, "UserLocalCache GetUserInfo req", "userID", userID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "UserLocalCache GetUserInfo return", "value", val)
		} else {
			log.ZError(ctx, "UserLocalCache GetUserInfo return", err)
		}
	}()
	return localcache.AnyValue[*sdkws.UserInfo](u.local.Get(ctx, cachekey.GetUserInfoKey(userID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "UserLocalCache GetUserInfo rpc", "userID", userID)
		return u.client.GetUserInfo(ctx, userID)
	}))
}

func (u *UserLocalCache) GetUserGlobalMsgRecvOpt(ctx context.Context, userID string) (val int32, err error) {
	log.ZDebug(ctx, "UserLocalCache GetUserGlobalMsgRecvOpt req", "userID", userID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "UserLocalCache GetUserGlobalMsgRecvOpt return", "value", val)
		} else {
			log.ZError(ctx, "UserLocalCache GetUserGlobalMsgRecvOpt return", err)
		}
	}()
	return localcache.AnyValue[int32](u.local.Get(ctx, cachekey.GetUserGlobalRecvMsgOptKey(userID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "UserLocalCache GetUserGlobalMsgRecvOpt rpc", "userID", userID)
		return u.client.GetUserGlobalMsgRecvOpt(ctx, userID)
	}))
}

func (u *UserLocalCache) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	users := make([]*sdkws.UserInfo, 0, len(userIDs))
	for _, userID := range userIDs {
		user, err := u.GetUserInfo(ctx, userID)
		if err != nil {
			if errs.ErrRecordNotFound.Is(err) {
				continue
			}
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserLocalCache) GetUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error) {
	users := make(map[string]*sdkws.UserInfo, len(userIDs))
	for _, userID := range userIDs {
		user, err := u.GetUserInfo(ctx, userID)
		if err != nil {
			if errs.ErrRecordNotFound.Is(err) {
				continue
			}
			return nil, err
		}
		users[userID] = user
	}
	return users, nil
}
