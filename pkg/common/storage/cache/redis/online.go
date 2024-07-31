package redis

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func NewUserOnline(rdb redis.UniversalClient) cache.OnlineCache {
	return &userOnline{
		rdb:         rdb,
		expire:      cachekey.OnlineExpire,
		channelName: cachekey.OnlineChannel,
	}
}

type userOnline struct {
	rdb         redis.UniversalClient
	expire      time.Duration
	channelName string
}

func (s *userOnline) getUserOnlineKey(userID string) string {
	return cachekey.GetOnlineKey(userID)
}

func (s *userOnline) GetOnline(ctx context.Context, userID string) ([]int32, error) {
	members, err := s.rdb.ZRangeByScore(ctx, s.getUserOnlineKey(userID), &redis.ZRangeBy{
		Min: strconv.FormatInt(time.Now().Unix(), 10),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	platformIDs := make([]int32, 0, len(members))
	for _, member := range members {
		val, err := strconv.Atoi(member)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		platformIDs = append(platformIDs, int32(val))
	}
	return platformIDs, nil
}

func (s *userOnline) SetUserOnline(ctx context.Context, userID string, online, offline []int32) error {
	script := `
	local key = KEYS[1]
	local score = ARGV[3]
	local num1 = redis.call("ZCARD", key)
	redis.call("ZREMRANGEBYSCORE", key, "-inf", ARGV[2])
	for i = 5, tonumber(ARGV[4])+4 do
		redis.call("ZREM", key, ARGV[i])
	end
	local num2 = redis.call("ZCARD", key)
	for i = 5+tonumber(ARGV[4]), #ARGV do
		redis.call("ZADD", key, score, ARGV[i])
	end
	redis.call("EXPIRE", key, ARGV[1])
	local num3 = redis.call("ZCARD", key)
	local change = (num1 ~= num2) or (num2 ~= num3)
	if change then
		local members = redis.call("ZRANGE", key, 0, -1)
		table.insert(members, KEYS[2])
		redis.call("PUBLISH", KEYS[3], table.concat(members, ":"))
		return 1
	else
		return 0
	end
`
	now := time.Now()
	argv := make([]any, 0, 2+len(online)+len(offline))
	argv = append(argv, int32(s.expire/time.Second), now.Unix(), now.Add(s.expire).Unix(), int32(len(offline)))
	for _, platformID := range offline {
		argv = append(argv, platformID)
	}
	for _, platformID := range online {
		argv = append(argv, platformID)
	}
	keys := []string{s.getUserOnlineKey(userID), userID, s.channelName}
	status, err := s.rdb.Eval(ctx, script, keys, argv).Result()
	if err != nil {
		log.ZError(ctx, "redis SetUserOnline", err, "userID", userID, "online", online, "offline", offline)
		return err
	}
	log.ZDebug(ctx, "redis SetUserOnline", "userID", userID, "online", online, "offline", offline, "status", status)
	return nil
}
