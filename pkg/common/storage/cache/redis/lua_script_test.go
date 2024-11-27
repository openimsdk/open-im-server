package redis

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLuaSetBatchWithCommonExpire(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	ctx := context.Background()

	keys := []string{"key1", "key2"}
	values := []string{"value1", "value2"}
	expire := 10

	mock.ExpectEvalSha(setBatchWithCommonExpireScript.Hash(), keys, []any{expire, "value1", "value2"}).SetVal(int64(len(keys)))

	err := LuaSetBatchWithCommonExpire(ctx, rdb, keys, values, expire)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLuaSetBatchWithIndividualExpire(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	ctx := context.Background()

	keys := []string{"key1", "key2"}
	values := []string{"value1", "value2"}
	expires := []int{10, 20}

	args := make([]any, 0, len(values)+len(expires))
	for _, v := range values {
		args = append(args, v)
	}
	for _, ex := range expires {
		args = append(args, ex)
	}

	mock.ExpectEvalSha(setBatchWithIndividualExpireScript.Hash(), keys, args).SetVal(int64(len(keys)))

	err := LuaSetBatchWithIndividualExpire(ctx, rdb, keys, values, expires)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLuaDeleteBatch(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	ctx := context.Background()

	keys := []string{"key1", "key2"}

	mock.ExpectEvalSha(deleteBatchScript.Hash(), keys, []any{}).SetVal(int64(len(keys)))

	err := LuaDeleteBatch(ctx, rdb, keys)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLuaGetBatch(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	ctx := context.Background()

	keys := []string{"key1", "key2"}
	expectedValues := []any{"value1", "value2"}

	mock.ExpectEvalSha(getBatchScript.Hash(), keys, []any{}).SetVal(expectedValues)

	values, err := LuaGetBatch(ctx, rdb, keys)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedValues, values)
}
