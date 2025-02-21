package database

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key []string) (map[string]string, error)
	Prefix(ctx context.Context, prefix string) (map[string]string, error)
	Set(ctx context.Context, key string, value string, expireAt time.Duration) error
	Incr(ctx context.Context, key string, value int) (int, error)
	Del(ctx context.Context, key []string) error
	Lock(ctx context.Context, key string, duration time.Duration) (string, error)
	Unlock(ctx context.Context, key string, value string) error
}
