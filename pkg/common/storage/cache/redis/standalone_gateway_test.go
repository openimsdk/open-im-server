package redis

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandaloneGatewayRedisRegisterGateway(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	cache := NewStandaloneGatewayRedis(rdb, time.Second*10)

	mock.Regexp().ExpectHSet(standaloneGatewayHashKey, "127.0.0.1:10001", `^[0-9]+$`).SetVal(1)
	mock.ExpectExpire(standaloneGatewayHashKey, time.Second*20).SetVal(true)

	err := cache.RegisterGateway(context.Background(), "127.0.0.1:10001")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStandaloneGatewayRedisUnregisterGateway(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	cache := NewStandaloneGatewayRedis(rdb, time.Second)

	mock.ExpectHDel(standaloneGatewayHashKey, "127.0.0.1:10001").SetVal(1)

	err := cache.UnregisterGateway(context.Background(), "127.0.0.1:10001")
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStandaloneGatewayRedisGetGatewayAddrs(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	cache := NewStandaloneGatewayRedis(rdb, time.Second*10)

	now := time.Now()
	mock.ExpectHGetAll(standaloneGatewayHashKey).SetVal(map[string]string{
		"127.0.0.1:10001": strconv.FormatInt(now.Add(-time.Second).UnixMilli(), 10),
		"127.0.0.1:10002": strconv.FormatInt(now.Add(-time.Second*20).UnixMilli(), 10),
		"127.0.0.1:10003": strconv.FormatInt(now.Add(time.Second).UnixMilli(), 10),
	})

	addrs, err := cache.GetGatewayAddrs(context.Background())
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"127.0.0.1:10001", "127.0.0.1:10003"}, addrs)
	assert.NoError(t, mock.ExpectationsWereMet())
}
