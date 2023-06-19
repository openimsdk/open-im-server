package cache

import (
	"context"
	"time"

	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
)

const (
	friendExpireTime    = time.Second * 60 * 60 * 12
	friendIDsKey        = "FRIEND_IDS:"
	TwoWayFriendsIDsKey = "COMMON_FRIENDS_IDS:"
	friendKey           = "FRIEND_INFO:"
)

// args fn will exec when no data in msgCache
type FriendCache interface {
	metaCache
	NewCache() FriendCache
	GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error)
	// call when friendID List changed
	DelFriendIDs(ownerUserID ...string) FriendCache
	// get single friendInfo from msgCache
	GetFriend(ctx context.Context, ownerUserID, friendUserID string) (friend *relationTb.FriendModel, err error)
	// del friend when friend info changed
	DelFriend(ownerUserID, friendUserID string) FriendCache
}

type FriendCacheRedis struct {
	metaCache
	friendDB   relationTb.FriendModelInterface
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewFriendCacheRedis(rdb redis.UniversalClient, friendDB relationTb.FriendModelInterface, options rockscache.Options) FriendCache {
	rcClient := rockscache.NewClient(rdb, options)
	return &FriendCacheRedis{
		metaCache:  NewMetaCacheRedis(rcClient),
		friendDB:   friendDB,
		expireTime: friendExpireTime,
		rcClient:   rcClient,
	}
}

func (c *FriendCacheRedis) NewCache() FriendCache {
	return &FriendCacheRedis{rcClient: c.rcClient, metaCache: NewMetaCacheRedis(c.rcClient, c.metaCache.GetPreDelKeys()...), friendDB: c.friendDB, expireTime: c.expireTime}
}

func (f *FriendCacheRedis) getFriendIDsKey(ownerUserID string) string {
	return friendIDsKey + ownerUserID
}

func (f *FriendCacheRedis) getTwoWayFriendsIDsKey(ownerUserID string) string {
	return TwoWayFriendsIDsKey + ownerUserID
}

func (f *FriendCacheRedis) getFriendKey(ownerUserID, friendUserID string) string {
	return friendKey + ownerUserID + "-" + friendUserID
}

func (f *FriendCacheRedis) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
	return getCache(ctx, f.rcClient, f.getFriendIDsKey(ownerUserID), f.expireTime, func(ctx context.Context) ([]string, error) {
		return f.friendDB.FindFriendUserIDs(ctx, ownerUserID)
	})
}

func (f *FriendCacheRedis) DelFriendIDs(ownerUserID ...string) FriendCache {
	new := f.NewCache()
	var keys []string
	for _, userID := range ownerUserID {
		keys = append(keys, f.getFriendIDsKey(userID))
	}
	new.AddKeys(keys...)
	return new
}

// todo
func (f *FriendCacheRedis) GetTwoWayFriendIDs(ctx context.Context, ownerUserID string) (twoWayFriendIDs []string, err error) {
	friendIDs, err := f.GetFriendIDs(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}
	for _, friendID := range friendIDs {
		friendFriendID, err := f.GetFriendIDs(ctx, friendID)
		if err != nil {
			return nil, err
		}
		if utils.IsContain(ownerUserID, friendFriendID) {
			twoWayFriendIDs = append(twoWayFriendIDs, ownerUserID)
		}
	}
	return twoWayFriendIDs, nil
}

func (f *FriendCacheRedis) DelTwoWayFriendIDs(ctx context.Context, ownerUserID string) FriendCache {
	new := f.NewCache()
	new.AddKeys(f.getTwoWayFriendsIDsKey(ownerUserID))
	return new
}

func (f *FriendCacheRedis) GetFriend(ctx context.Context, ownerUserID, friendUserID string) (friend *relationTb.FriendModel, err error) {
	return getCache(ctx, f.rcClient, f.getFriendKey(ownerUserID, friendUserID), f.expireTime, func(ctx context.Context) (*relationTb.FriendModel, error) {
		return f.friendDB.Take(ctx, ownerUserID, friendUserID)
	})
}

func (f *FriendCacheRedis) DelFriend(ownerUserID, friendUserID string) FriendCache {
	new := f.NewCache()
	new.AddKeys(f.getFriendKey(ownerUserID, friendUserID))
	return new
}
