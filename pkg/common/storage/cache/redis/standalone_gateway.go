package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/tools/errs"
)

const standaloneGatewayHashKey = "STANDALONE_GATEWAY_REGISTRY"

type StandaloneGatewayRedis struct {
	rdb       redis.UniversalClient
	validTime time.Duration
}

func NewStandaloneGatewayRedis(rdb redis.UniversalClient, validTime time.Duration) *StandaloneGatewayRedis {
	return &StandaloneGatewayRedis{rdb: rdb, validTime: validTime}
}

func (s *StandaloneGatewayRedis) RegisterGateway(ctx context.Context, addr string) error {
	pipe := s.rdb.Pipeline()
	pipe.HSet(ctx, standaloneGatewayHashKey, addr, strconv.FormatInt(time.Now().UnixMilli(), 10))
	pipe.Expire(ctx, standaloneGatewayHashKey, s.validTime*2)
	_, err := pipe.Exec(ctx)
	return errs.Wrap(err)
}

func (s *StandaloneGatewayRedis) UnregisterGateway(ctx context.Context, addr string) error {
	return errs.Wrap(s.rdb.HDel(ctx, standaloneGatewayHashKey, addr).Err())
}

func (s *StandaloneGatewayRedis) GetGatewayAddrs(ctx context.Context) ([]string, error) {
	gateways, err := s.rdb.HGetAll(ctx, standaloneGatewayHashKey).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}

	now := time.Now()
	addrs := make([]string, 0, len(gateways))
	for addr, registeredAt := range gateways {
		registeredAtMs, err := strconv.ParseInt(registeredAt, 10, 64)
		if err != nil {
			return nil, errs.WrapMsg(err, "redis gateway register time is not int64", "addr", addr, "value", registeredAt)
		}
		if now.Sub(time.UnixMilli(registeredAtMs)) <= s.validTime {
			addrs = append(addrs, addr)
		}
	}
	return addrs, nil
}
