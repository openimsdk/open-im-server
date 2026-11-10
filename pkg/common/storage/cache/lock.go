package cache

import (
	"context"
	"time"
)

type Lock interface {
	Lock(ctx context.Context, key string, timeout time.Duration) (string, error)
	Unlock(ctx context.Context, key, value string)
}
