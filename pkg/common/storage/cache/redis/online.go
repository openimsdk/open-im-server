package redis

import (
	"context"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
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

func (s *userOnline) GetAllOnlineUsers(ctx context.Context, cursor uint64) (map[string][]int32, uint64, error) {
	result := make(map[string][]int32)

	keys, nextCursor, err := s.scanOnlineKeys(ctx, cursor)
	if err != nil {
		return nil, 0, err
	}

	for _, key := range keys {
		userID := cachekey.GetOnlineKeyUserID(key)
		strValues, err := s.rdb.ZRange(ctx, key, 0, -1).Result()
		if err != nil {
			return nil, 0, err
		}

		values := make([]int32, 0, len(strValues))
		for _, value := range strValues {
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, 0, errs.Wrap(err)
			}
			values = append(values, int32(intValue))
		}

		result[userID] = values
	}

	return result, nextCursor, nil
}

func (s *userOnline) scanOnlineKeys(ctx context.Context, cursor uint64) ([]string, uint64, error) {
	const match = cachekey.OnlineKey + "*"
	cluster, ok := s.rdb.(*redis.ClusterClient)
	if !ok {
		return s.rdb.Scan(ctx, cursor, match, constant.ParamMaxLength).Result()
	}

	var primaries []*redis.Client
	var mu sync.Mutex
	if err := cluster.ForEachMaster(ctx, func(_ context.Context, primary *redis.Client) error {
		mu.Lock()
		primaries = append(primaries, primary)
		mu.Unlock()
		return nil
	}); err != nil {
		return nil, 0, err
	}
	// With a stable cluster topology, every RPC replica sees the same order.
	slices.SortFunc(primaries, func(a, b *redis.Client) int {
		return strings.Compare(a.Options().Addr, b.Options().Addr)
	})

	// Keep the RPC cursor stateless. Its upper 14 bits hold the primary
	// index (the cluster has 2^14 slots); the lower bits hold its SCAN cursor.
	const cursorBits = 50
	const cursorMask = 1<<cursorBits - 1
	primaryIndex := int(cursor >> cursorBits)
	if primaryIndex >= len(primaries) {
		return nil, 0, errs.New("online scan cursor references an unavailable Redis primary", "index", primaryIndex).Wrap()
	}
	keys, nextCursor, err := primaries[primaryIndex].Scan(ctx, cursor&cursorMask, match, constant.ParamMaxLength).Result()
	if err != nil {
		return nil, 0, err
	}
	if nextCursor > cursorMask {
		return nil, 0, errs.New("Redis SCAN cursor exceeds the online scan cursor range", "cursor", nextCursor).Wrap()
	}
	if nextCursor == 0 {
		primaryIndex++
		if primaryIndex == len(primaries) {
			return keys, 0, nil
		}
	}
	return keys, uint64(primaryIndex)<<cursorBits | nextCursor, nil
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
		table.insert(members, "1")
		return members
	else
		return {"0"}
	end
`
	now := time.Now()
	argv := make([]any, 0, 4+len(online)+len(offline))
	argv = append(argv, int32(s.expire/time.Second), now.Unix(), now.Add(s.expire).Unix(), int32(len(offline)))
	for _, platformID := range offline {
		argv = append(argv, platformID)
	}
	for _, platformID := range online {
		argv = append(argv, platformID)
	}
	keys := []string{s.getUserOnlineKey(userID)}
	platformIDs, err := s.rdb.Eval(ctx, script, keys, argv).StringSlice()
	if err != nil {
		log.ZError(ctx, "redis SetUserOnline", err, "userID", userID, "online", online, "offline", offline)
		return err
	}
	if len(platformIDs) == 0 {
		return errs.ErrInternalServer.WrapMsg("SetUserOnline redis lua invalid return value")
	}
	if platformIDs[len(platformIDs)-1] != "0" {
		log.ZDebug(ctx, "redis SetUserOnline push", "userID", userID, "online", online, "offline", offline, "platformIDs", platformIDs[:len(platformIDs)-1])
		platformIDs[len(platformIDs)-1] = userID
		msg := strings.Join(platformIDs, ":")
		if err := s.rdb.Publish(ctx, s.channelName, msg).Err(); err != nil {
			return errs.Wrap(err)
		}
	} else {
		log.ZDebug(ctx, "redis SetUserOnline not push", "userID", userID, "online", online, "offline", offline)
	}
	return nil
}
