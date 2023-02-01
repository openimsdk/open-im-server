package cache

import (
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	friendExpireTime = time.Second * 60 * 60 * 12
	friendIDsKey     = "FRIEND_IDS:"
)

type FriendCache struct {
	friendDB   *relation.Friend
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewFriendCache(rdb redis.UniversalClient, friendDB *relation.Friend, options rockscache.Options) *FriendCache {
	return &FriendCache{
		friendDB:   friendDB,
		expireTime: friendExpireTime,
		rcClient:   rockscache.NewClient(rdb, options),
	}
}

func (f *FriendCache) getFriendRelationKey(ownerUserID string) string {
	return friendIDsKey + ownerUserID
}

func (f *FriendCache) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
	getFriendIDList := func() (string, error) {
		friendIDList, err := f.friendDB.GetFriendIDs(ctx, ownerUserID)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(friendIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "friendIDs", friendIDs)
	}()
	friendIDListStr, err := f.rcClient.Fetch(f.getFriendRelationKey(ownerUserID), f.expireTime, getFriendIDList)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(friendIDListStr), &friendIDs)
	return friendIDs, utils.Wrap(err, "")
}

func (f *FriendCache) DelFriendIDs(ctx context.Context, ownerUserID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID)
	}()
	return f.rcClient.TagAsDeleted(f.getFriendRelationKey(ownerUserID))
}
