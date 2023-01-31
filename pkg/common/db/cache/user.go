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
	UserExpireTime = time.Second * 60 * 60 * 12
	userInfoKey    = "USER_INFO:"
)

type UserCache struct {
	userDB *relation.User

	expireTime  time.Duration
	redisClient *RedisClient
	rcClient    *rockscache.Client
}

func NewUserCache(rdb redis.UniversalClient, userDB *relation.User, options rockscache.Options) *UserCache {
	return &UserCache{
		userDB:      userDB,
		expireTime:  UserExpireTime,
		redisClient: NewRedisClient(rdb),
		rcClient:    rockscache.NewClient(rdb, options),
	}
}

func (u *UserCache) getUserInfoKey(userID string) string {
	return userInfoKey + userID
}

func (u *UserCache) GetUserInfo(ctx context.Context, userID string) (userInfo *relation.User, err error) {
	getUserInfo := func() (string, error) {
		userInfo, err := u.userDB.Take(ctx, userID)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(userInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "userInfo", *userInfo)
	}()
	userInfoStr, err := u.rcClient.Fetch(u.getUserInfoKey(userID), time.Second*30*60, getUserInfo)
	if err != nil {
		return nil, err
	}
	userInfo = &relation.User{}
	err = json.Unmarshal([]byte(userInfoStr), userInfo)
	return userInfo, utils.Wrap(err, "")
}

func (u *UserCache) GetUsersInfo(ctx context.Context, userIDs []string) ([]*relation.User, error) {
	var users []*relation.User
	for _, userID := range userIDs {
		user, err := GetUserInfoFromCache(ctx, userID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserCache) DelUserInfo(ctx context.Context, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID)
	}()
	return u.rcClient.TagAsDeleted(u.getUserInfoKey(userID) + userID)
}

func (u *UserCache) DelUsersInfo(ctx context.Context, userIDs []string) (err error) {
	for _, userID := range userIDs {
		if err := u.DelUserInfo(ctx, userID); err != nil {
			return err
		}
	}
	return nil
}
