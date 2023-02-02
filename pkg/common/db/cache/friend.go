package cache

import (
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
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

type FriendCache struct {
	friendDB   *relation.FriendGorm
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewFriendCache(rdb redis.UniversalClient, friendDB *relation.FriendGorm, options rockscache.Options) *FriendCache {
	return &FriendCache{
		friendDB:   friendDB,
		expireTime: friendExpireTime,
		rcClient:   rockscache.NewClient(rdb, options),
	}
}

func (f *FriendCache) getFriendIDsKey(ownerUserID string) string {
	return friendIDsKey + ownerUserID
}

func (f *FriendCache) getTwoWayFriendsIDsKey(ownerUserID string) string {
	return TwoWayFriendsIDsKey + ownerUserID
}

func (f *FriendCache) getFriendKey(ownerUserID, friendUserID string) string {
	return friendKey + ownerUserID + "-" + friendUserID
}

func (f *FriendCache) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
	getFriendIDs := func() (string, error) {
		friendIDs, err := f.friendDB.GetFriendIDs(ctx, ownerUserID)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(friendIDs)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "friendIDs", friendIDs)
	}()
	friendIDsStr, err := f.rcClient.Fetch(f.getFriendIDsKey(ownerUserID), f.expireTime, getFriendIDs)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(friendIDsStr), &friendIDs)
	return friendIDs, utils.Wrap(err, "")
}

func (f *FriendCache) DelFriendIDs(ctx context.Context, ownerUserID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID)
	}()
	return f.rcClient.TagAsDeleted(f.getFriendIDsKey(ownerUserID))
}

func (f *FriendCache) GetTwoWayFriendIDs(ctx context.Context, ownerUserID string) (twoWayFriendIDs []string, err error) {
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

func (f *FriendCache) DelTwoWayFriendIDs(ctx context.Context, ownerUserID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID)
	}()
	return f.rcClient.TagAsDeleted(f.getTwoWayFriendsIDsKey(ownerUserID))
}

func (f *FriendCache) GetFriend(ctx context.Context, ownerUserID, friendUserID string) (friend *table.FriendModel, err error) {
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
	friendInfo := &relation.FriendUser{}
	err = json.Unmarshal([]byte(groupMemberInfoStr), groupMember)
	return groupMember, utils.Wrap(err, "")
}
