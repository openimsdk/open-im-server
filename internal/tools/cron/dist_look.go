package cron

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/openimsdk/tools/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const (
	lockLeaseTTL = 300
)

type EtcdLocker struct {
	client     *clientv3.Client
	instanceID string
}

// NewEtcdLocker creates a new etcd distributed lock
func NewEtcdLocker(client *clientv3.Client) (*EtcdLocker, error) {
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	instanceID := fmt.Sprintf("%s-pid-%d-%d", hostname, pid, time.Now().UnixNano())

	locker := &EtcdLocker{
		client:     client,
		instanceID: instanceID,
	}

	return locker, nil
}

func (e *EtcdLocker) ExecuteWithLock(ctx context.Context, taskName string, task func()) {
	session, err := concurrency.NewSession(e.client, concurrency.WithTTL(lockLeaseTTL))
	if err != nil {
		log.ZWarn(ctx, "Failed to create etcd session", err,
			"taskName", taskName,
			"instanceID", e.instanceID)
		return
	}
	defer session.Close()

	lockKey := fmt.Sprintf("openim/crontask/%s", taskName)
	mutex := concurrency.NewMutex(session, lockKey)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err = mutex.TryLock(ctxWithTimeout)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.ZDebug(ctx, "Task is being executed by another instance, skipping",
				"taskName", taskName,
				"instanceID", e.instanceID)
		} else {
			log.ZWarn(ctx, "Failed to acquire task lock", err,
				"taskName", taskName,
				"instanceID", e.instanceID)
		}
		return
	}

	defer func() {
		if err := mutex.Unlock(ctx); err != nil {
			log.ZWarn(ctx, "Failed to release task lock", err,
				"taskName", taskName,
				"instanceID", e.instanceID)
		} else {
			log.ZInfo(ctx, "Successfully released task lock",
				"taskName", taskName,
				"instanceID", e.instanceID)
		}
	}()

	log.ZInfo(ctx, "Successfully acquired task lock, starting execution",
		"taskName", taskName,
		"instanceID", e.instanceID,
		"sessionID", session.Lease())

	task()

	log.ZInfo(ctx, "Task execution completed",
		"taskName", taskName,
		"instanceID", e.instanceID)
}
