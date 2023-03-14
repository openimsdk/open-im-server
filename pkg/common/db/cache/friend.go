package cache

import (
	"OpenIM/pkg/common/db/relation"
	relationTb "OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	friendExpireTime    = time.Second * 60 * 60 * 12
	friendIDsKey        = "FRIEND_IDS:"
	TwoWayFriendsIDsKey = "COMMON_FRIENDS_IDS:"
	friendKey           = "FRIEND_INFO:"
)

// args fn will exec when no data in cache
type FriendCache interface {
	GetFriendIDs(ctx context.Context, ownerUserID string, fn func(ctx context.Context, ownerUserID string) (friendIDs []string, err error)) (friendIDs []string, err error)
	// call when friendID List changed
	DelFriendIDs(ctx context.Context, ownerUserID string) (err error)
	// get single friendInfo from cache
	GetFriend(ctx context.Context, ownerUserID, friendUserID string, fn func(ctx context.Context, ownerUserID, friendUserID string) (friend *relationTb.FriendModel, err error)) (friend *relationTb.FriendModel, err error)
	// del friend when friend info changed
	DelFriend(ctx context.Context, ownerUserID, friendUserID string) (err error)
}

type FriendCacheRedis struct {
	friendDB   *relation.FriendGorm
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewFriendCacheRedis(rdb redis.UniversalClient, friendDB *relation.FriendGorm, options rockscache.Options) *FriendCacheRedis {
	return &FriendCacheRedis{
		friendDB:   friendDB,
		expireTime: friendExpireTime,
		rcClient:   rockscache.NewClient(rdb, options),
	}
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
	return GetCache(ctx, f.rcClient, f.getFriendIDsKey(ownerUserID), f.expireTime, func(ctx context.Context) ([]string, error) {
		return f.friendDB.FindFriendUserIDs(ctx, ownerUserID)
	})
}

func (f *FriendCacheRedis) DelFriendIDs(ctx context.Context, ownerUserID string) (err error) {
	return f.rcClient.TagAsDeleted(f.getFriendIDsKey(ownerUserID))
}

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

func (f *FriendCacheRedis) DelTwoWayFriendIDs(ctx context.Context, ownerUserID string) (err error) {
	return f.rcClient.TagAsDeleted(f.getTwoWayFriendsIDsKey(ownerUserID))
}

func (f *FriendCacheRedis) GetFriend(ctx context.Context, ownerUserID, friendUserID string, fn func(ctx context.Context, ownerUserID, friendUserID string) (friend *relationTb.FriendModel, err error)) (friend *relationTb.FriendModel, err error) {
	getFriend := func() (string, error) {
		friend, err = f.friendDB.Take(ctx, ownerUserID, friendUserID)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(friend)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	friendStr, err := f.rcClient.Fetch(f.getFriendKey(ownerUserID, friendUserID), f.expireTime, getFriend)
	if err != nil {
		return nil, err
	}
	friend = &relationTb.FriendModel{}
	err = json.Unmarshal([]byte(friendStr), friend)
	return friend, utils.Wrap(err, "")
}

func (f *FriendCacheRedis) DelFriend(ctx context.Context, ownerUserID, friendUserID string) (err error) {
	return f.rcClient.TagAsDeleted(f.getFriendKey(ownerUserID, friendUserID))
}
