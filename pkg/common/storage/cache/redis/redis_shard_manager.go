package redis

import (
	"context"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

const (
	defaultBatchSize       = 50
	defaultConcurrentLimit = 3
)

// RedisShardManager is a class for sharding and processing keys
type RedisShardManager struct {
	redisClient redis.UniversalClient
	config      *Config
}
type Config struct {
	batchSize       int
	continueOnError bool
	concurrentLimit int
}

// Option is a function type for configuring Config
type Option func(c *Config)

//// NewRedisShardManager creates a new RedisShardManager instance
//func NewRedisShardManager(redisClient redis.UniversalClient, opts ...Option) *RedisShardManager {
//	config := &Config{
//		batchSize:       defaultBatchSize, // Default batch size is 50 keys
//		continueOnError: false,
//		concurrentLimit: defaultConcurrentLimit, // Default concurrent limit is 3
//	}
//	for _, opt := range opts {
//		opt(config)
//	}
//	rsm := &RedisShardManager{
//		redisClient: redisClient,
//		config:      config,
//	}
//	return rsm
//}
//
//// WithBatchSize sets the number of keys to process per batch
//func WithBatchSize(size int) Option {
//	return func(c *Config) {
//		c.batchSize = size
//	}
//}
//
//// WithContinueOnError sets whether to continue processing on error
//func WithContinueOnError(continueOnError bool) Option {
//	return func(c *Config) {
//		c.continueOnError = continueOnError
//	}
//}
//
//// WithConcurrentLimit sets the concurrency limit
//func WithConcurrentLimit(limit int) Option {
//	return func(c *Config) {
//		c.concurrentLimit = limit
//	}
//}
//
//// ProcessKeysBySlot groups keys by their Redis cluster hash slots and processes them using the provided function.
//func (rsm *RedisShardManager) ProcessKeysBySlot(
//	ctx context.Context,
//	keys []string,
//	processFunc func(ctx context.Context, slot int64, keys []string) error,
//) error {
//
//	// Group keys by slot
//	slots, err := groupKeysBySlot(ctx, rsm.redisClient, keys)
//	if err != nil {
//		return err
//	}
//
//	g, ctx := errgroup.WithContext(ctx)
//	g.SetLimit(rsm.config.concurrentLimit)
//
//	// Process keys in each slot using the provided function
//	for slot, singleSlotKeys := range slots {
//		batches := splitIntoBatches(singleSlotKeys, rsm.config.batchSize)
//		for _, batch := range batches {
//			slot, batch := slot, batch // Avoid closure capture issue
//			g.Go(func() error {
//				err := processFunc(ctx, slot, batch)
//				if err != nil {
//					log.ZWarn(ctx, "Batch processFunc failed", err, "slot", slot, "keys", batch)
//					if !rsm.config.continueOnError {
//						return err
//					}
//				}
//				return nil
//			})
//		}
//	}
//
//	if err := g.Wait(); err != nil {
//		return err
//	}
//	return nil
//}

// groupKeysBySlot groups keys by their Redis cluster hash slots.
func groupKeysBySlot(ctx context.Context, redisClient redis.UniversalClient, keys []string) (map[int64][]string, error) {
	slots := make(map[int64][]string)
	clusterClient, isCluster := redisClient.(*redis.ClusterClient)
	if isCluster && len(keys) > 1 {
		pipe := clusterClient.Pipeline()
		cmds := make([]*redis.IntCmd, len(keys))
		for i, key := range keys {
			cmds[i] = pipe.ClusterKeySlot(ctx, key)
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			return nil, errs.WrapMsg(err, "get slot err")
		}

		for i, cmd := range cmds {
			slot, err := cmd.Result()
			if err != nil {
				log.ZWarn(ctx, "some key get slot err", err, "key", keys[i])
				return nil, errs.WrapMsg(err, "get slot err", "key", keys[i])
			}
			slots[slot] = append(slots[slot], keys[i])
		}
	} else {
		// If not a cluster client, put all keys in the same slot (0)
		slots[0] = keys
	}

	return slots, nil
}

// splitIntoBatches splits keys into batches of the specified size
func splitIntoBatches(keys []string, batchSize int) [][]string {
	var batches [][]string
	for batchSize < len(keys) {
		keys, batches = keys[batchSize:], append(batches, keys[0:batchSize:batchSize])
	}
	return append(batches, keys)
}

// ProcessKeysBySlot groups keys by their Redis cluster hash slots and processes them using the provided function.
func ProcessKeysBySlot(
	ctx context.Context,
	redisClient redis.UniversalClient,
	keys []string,
	processFunc func(ctx context.Context, slot int64, keys []string) error,
	opts ...Option,
) error {

	config := &Config{
		batchSize:       defaultBatchSize,
		continueOnError: false,
		concurrentLimit: defaultConcurrentLimit,
	}
	for _, opt := range opts {
		opt(config)
	}

	// Group keys by slot
	slots, err := groupKeysBySlot(ctx, redisClient, keys)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(config.concurrentLimit)

	// Process keys in each slot using the provided function
	for slot, singleSlotKeys := range slots {
		batches := splitIntoBatches(singleSlotKeys, config.batchSize)
		for _, batch := range batches {
			slot, batch := slot, batch // Avoid closure capture issue
			g.Go(func() error {
				err := processFunc(ctx, slot, batch)
				if err != nil {
					log.ZWarn(ctx, "Batch processFunc failed", err, "slot", slot, "keys", batch)
					if !config.continueOnError {
						return err
					}
				}
				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func DeleteCacheBySlot(ctx context.Context, rcClient *rocksCacheClient, keys []string) error {
	switch len(keys) {
	case 0:
		return nil
	case 1:
		return rcClient.GetClient().TagAsDeletedBatch2(ctx, keys)
	default:
		return ProcessKeysBySlot(ctx, rcClient.GetRedis(), keys, func(ctx context.Context, slot int64, keys []string) error {
			return rcClient.GetClient().TagAsDeletedBatch2(ctx, keys)
		})
	}
}
