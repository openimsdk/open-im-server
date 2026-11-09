package cron

import (
	"context"
	"time"
)

const (
	lockLeaseTTL = time.Second * 300
)

type Locker interface {
	ExecuteWithLock(ctx context.Context, taskName string, task func())
}
