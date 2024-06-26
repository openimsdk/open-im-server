package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type userOnline struct {
	rdb         redis.UniversalClient
	expire      time.Duration
	channelName string
}

func (s *userOnline) getUserOnlineKey(userID string) string {
	return "USER_ONLINE:" + userID
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
	if err := s.rdb.Eval(ctx, script, keys, argv).Err(); err != nil {
		return err
	}
	return nil
}
