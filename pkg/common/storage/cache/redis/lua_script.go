package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

var (
	setBatchWithCommonExpireScript = redis.NewScript(`
local expire = tonumber(ARGV[1])
for i, key in ipairs(KEYS) do
    redis.call('SET', key, ARGV[i + 1])
    redis.call('EXPIRE', key, expire)
end
return #KEYS
`)

	setBatchWithIndividualExpireScript = redis.NewScript(`
local n = #KEYS
for i = 1, n do
    redis.call('SET', KEYS[i], ARGV[i])
    redis.call('EXPIRE', KEYS[i], ARGV[i + n])
end
return n
`)

	deleteBatchScript = redis.NewScript(`
for i, key in ipairs(KEYS) do
    redis.call('DEL', key)
end
return #KEYS
`)

	getBatchScript = redis.NewScript(`
local values = {}
for i, key in ipairs(KEYS) do
    local value = redis.call('GET', key)
    table.insert(values, value)
end
return values
`)
)

func callLua(ctx context.Context, rdb redis.Scripter, script *redis.Script, keys []string, args []any) (any, error) {
	log.ZDebug(ctx, "callLua args", "scriptHash", script.Hash(), "keys", keys, "args", args)
	r := script.EvalSha(ctx, rdb, keys, args)
	if redis.HasErrorPrefix(r.Err(), "NOSCRIPT") {
		if err := script.Load(ctx, rdb).Err(); err != nil {
			r = script.Eval(ctx, rdb, keys, args)
		} else {
			r = script.EvalSha(ctx, rdb, keys, args)
		}
	}
	v, err := r.Result()
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	return v, errs.WrapMsg(err, "call lua err", "scriptHash", script.Hash(), "keys", keys, "args", args)
}

func LuaSetBatchWithCommonExpire(ctx context.Context, rdb redis.Scripter, keys []string, values []string, expire int) error {
	// Check if the lengths of keys and values match
	if len(keys) != len(values) {
		return errs.New("keys and values length mismatch").Wrap()
	}

	// Ensure allocation size does not overflow
	maxAllowedLen := (1 << 31) - 1 // 2GB limit (maximum address space for 32-bit systems)

	if len(values) > maxAllowedLen-1 {
		return fmt.Errorf("values length is too large, causing overflow")
	}
	var vals = make([]any, 0, 1+len(values))
	vals = append(vals, expire)
	for _, v := range values {
		vals = append(vals, v)
	}
	_, err := callLua(ctx, rdb, setBatchWithCommonExpireScript, keys, vals)
	return err
}

func LuaSetBatchWithIndividualExpire(ctx context.Context, rdb redis.Scripter, keys []string, values []string, expires []int) error {
	// Check if the lengths of keys, values, and expires match
	if len(keys) != len(values) || len(keys) != len(expires) {
		return errs.New("keys and values length mismatch").Wrap()
	}

	// Ensure the allocation size does not overflow
	maxAllowedLen := (1 << 31) - 1 // 2GB limit (maximum address space for 32-bit systems)

	if len(values) > maxAllowedLen-1 {
		return errs.New(fmt.Sprintf("values length %d exceeds the maximum allowed length %d", len(values), maxAllowedLen-1)).Wrap()
	}
	var vals = make([]any, 0, len(values)+len(expires))
	for _, v := range values {
		vals = append(vals, v)
	}
	for _, ex := range expires {
		vals = append(vals, ex)
	}
	_, err := callLua(ctx, rdb, setBatchWithIndividualExpireScript, keys, vals)
	return err
}

func LuaDeleteBatch(ctx context.Context, rdb redis.Scripter, keys []string) error {
	_, err := callLua(ctx, rdb, deleteBatchScript, keys, nil)
	return err
}

func LuaGetBatch(ctx context.Context, rdb redis.Scripter, keys []string) ([]any, error) {
	v, err := callLua(ctx, rdb, getBatchScript, keys, nil)
	if err != nil {
		return nil, err
	}
	values, ok := v.([]any)
	if !ok {
		return nil, servererrs.ErrArgs.WrapMsg("invalid lua get batch result")
	}
	return values, nil

}
