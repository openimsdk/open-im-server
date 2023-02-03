package cache

import (
	"Open_IM/pkg/common/db/relation"
	relation2 "Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const (
	userExpireTime          = time.Second * 60 * 60 * 12
	userInfoKey             = "USER_INFO:"
	userGlobalRecvMsgOptKey = "USER_GLOBAL_RECV_MSG_OPT_KEY:"
)

type UserCache struct {
	userDB *relation.UserGorm

	expireTime  time.Duration
	redisClient *RedisClient
	rcClient    *rockscache.Client
}

func NewUserCache(rdb redis.UniversalClient, userDB *relation.UserGorm, options rockscache.Options) *UserCache {
	return &UserCache{
		userDB:      userDB,
		expireTime:  userExpireTime,
		redisClient: NewRedisClient(rdb),
		rcClient:    rockscache.NewClient(rdb, options),
	}
}

func (u *UserCache) getUserInfoKey(userID string) string {
	return userInfoKey + userID
}

func (u *UserCache) getUserGlobalRecvMsgOptKey(userID string) string {
	return userGlobalRecvMsgOptKey + userID
}

func (u *UserCache) GetUserInfo(ctx context.Context, userID string) (userInfo *relation2.UserModel, err error) {
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
	userInfoStr, err := u.rcClient.Fetch(u.getUserInfoKey(userID), u.expireTime, getUserInfo)
	if err != nil {
		return nil, err
	}
	userInfo = &relation2.UserModel{}
	err = json.Unmarshal([]byte(userInfoStr), userInfo)
	return userInfo, utils.Wrap(err, "")
}

func (u *UserCache) GetUsersInfo(ctx context.Context, userIDs []string) ([]*relation2.UserModel, error) {
	var users []*relation2.UserModel
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
	return u.rcClient.TagAsDeleted(u.getUserInfoKey(userID))
}

func (u *UserCache) DelUsersInfo(ctx context.Context, userIDs []string) (err error) {
	for _, userID := range userIDs {
		if err := u.DelUserInfo(ctx, userID); err != nil {
			return err
		}
	}
	return nil
}

func (u *UserCache) GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error) {
	getUserGlobalRecvMsgOpt := func() (string, error) {
		userInfo, err := u.userDB.Take(ctx, userID)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(userInfo.GlobalRecvMsgOpt)), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "opt", opt)
	}()
	optStr, err := u.rcClient.Fetch(u.getUserInfoKey(userID), u.expireTime, getUserGlobalRecvMsgOpt)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(optStr)
}

func (u *UserCache) DelUserGlobalRecvMsgOpt(ctx context.Context, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID)
	}()
	return u.rcClient.TagAsDeleted(u.getUserGlobalRecvMsgOptKey(userID))
}
