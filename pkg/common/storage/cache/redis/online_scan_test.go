package redis

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/protocol/constant"
)

func TestGetAllOnlineUsersScansEachPrimary(t *testing.T) {
	const match = cachekey.OnlineKey + "*"
	const lowerAddr = "redis-a:6379"
	const upperAddr = "redis-b:6379"
	// {bar} hashes to slot 5061 and {foo} to slot 12182.
	lowerKey := cachekey.GetOnlineKey("{bar}-user")
	upperKey := cachekey.GetOnlineKey("{foo}-user")
	lower, lowerMock := redismock.NewClientMock()
	upper, upperMock := redismock.NewClientMock()
	lower.Options().Addr = lowerAddr
	upper.Options().Addr = upperAddr

	lowerMock.ExpectScan(0, match, constant.ParamMaxLength).SetVal(nil, 17)
	lowerMock.ExpectScan(17, match, constant.ParamMaxLength).SetVal([]string{lowerKey}, 0)
	lowerMock.ExpectZRange(lowerKey, 0, -1).SetVal([]string{"1"})
	upperMock.ExpectScan(0, match, constant.ParamMaxLength).SetVal([]string{upperKey}, 29)
	upperMock.ExpectZRange(upperKey, 0, -1).SetVal([]string{"2"})
	upperMock.ExpectScan(29, match, constant.ParamMaxLength).SetVal(nil, 0)

	cluster := redis.NewClusterClient(&redis.ClusterOptions{
		ClusterSlots: func(context.Context) ([]redis.ClusterSlot, error) {
			return []redis.ClusterSlot{
				{Start: 8192, End: 16383, Nodes: []redis.ClusterNode{{Addr: upperAddr}}},
				{Start: 0, End: 8191, Nodes: []redis.ClusterNode{{Addr: lowerAddr}}},
			}, nil
		},
		NewClient: func(opts *redis.Options) *redis.Client {
			return map[string]*redis.Client{lowerAddr: lower, upperAddr: upper}[opts.Addr]
		},
	})
	t.Cleanup(func() { _ = cluster.Close() })
	online := NewUserOnline(cluster)
	allUsers := make(map[string][]int32)
	var cursor uint64
	for page := range 4 {
		users, nextCursor, err := online.GetAllOnlineUsers(context.Background(), cursor)
		require.NoError(t, err)
		for id, platforms := range users {
			require.NotContains(t, allUsers, id)
			allUsers[id] = platforms
		}
		if page < 3 {
			require.NotZero(t, nextCursor)
		}
		cursor = nextCursor
	}
	require.Zero(t, cursor)
	require.Equal(t, map[string][]int32{"{bar}-user": {1}, "{foo}-user": {2}}, allUsers)
	require.NoError(t, lowerMock.ExpectationsWereMet())
	require.NoError(t, upperMock.ExpectationsWereMet())
}

func TestGetAllOnlineUsersStandaloneCursor(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	t.Cleanup(func() { _ = rdb.Close() })
	const match = cachekey.OnlineKey + "*"
	const cursor = 31
	mock.ExpectScan(0, match, constant.ParamMaxLength).SetVal(nil, cursor)
	mock.ExpectScan(cursor, match, constant.ParamMaxLength).SetVal(nil, 0)

	online := NewUserOnline(rdb)
	_, nextCursor, err := online.GetAllOnlineUsers(context.Background(), 0)
	require.NoError(t, err)
	require.EqualValues(t, cursor, nextCursor)
	_, nextCursor, err = online.GetAllOnlineUsers(context.Background(), nextCursor)
	require.NoError(t, err)
	require.Zero(t, nextCursor)
	require.NoError(t, mock.ExpectationsWereMet())
}
