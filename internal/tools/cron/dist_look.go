package cron

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/openimsdk/tools/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	lockLeaseTTL = 3000
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

func (e *EtcdLocker) tryAcquireTaskLock(ctx context.Context, taskName string) (clientv3.LeaseID, bool, error) {
	lockKey := fmt.Sprintf("openim/crontask/%s-lock", taskName)

	lease, err := e.client.Grant(ctx, lockLeaseTTL)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create lease: %w", err)
	}

	txnResp, err := e.client.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, e.instanceID, clientv3.WithLease(lease.ID))).
		Else(clientv3.OpGet(lockKey)).
		Commit()

	if err != nil {
		e.client.Revoke(ctx, lease.ID)
		return 0, false, fmt.Errorf("transaction failed: %w", err)
	}

	if !txnResp.Succeeded {
		rangeResp := txnResp.Responses[0].GetResponseRange()
		if len(rangeResp.Kvs) > 0 {
			currentOwner := string(rangeResp.Kvs[0].Value)
			log.ZInfo(ctx, "Task lock already owned, skipping execution",
				"taskName", taskName,
				"instanceID", e.instanceID,
				"currentOwner", currentOwner)
		}
		e.client.Revoke(ctx, lease.ID)
		return 0, false, nil
	}

	log.ZInfo(ctx, "Successfully acquired task lock",
		"taskName", taskName,
		"instanceID", e.instanceID,
		"leaseID", lease.ID)
	return lease.ID, true, nil
}

func (e *EtcdLocker) releaseTaskLock(ctx context.Context, taskName string, leaseID clientv3.LeaseID) {
	if leaseID == 0 {
		return
	}

	_, err := e.client.Revoke(ctx, leaseID)
	if err != nil {
		log.ZWarn(ctx, "Failed to revoke task lease", err,
			"taskName", taskName,
			"instanceID", e.instanceID,
			"leaseID", leaseID)
	} else {
		log.ZInfo(ctx, "Successfully released task lock",
			"taskName", taskName,
			"instanceID", e.instanceID)
	}
}

func (e *EtcdLocker) startLeaseKeepAlive(ctx context.Context, taskName string, leaseID clientv3.LeaseID) (context.CancelFunc, error) {
	keepAliveCtx, cancel := context.WithCancel(ctx)

	keepAliveCh, err := e.client.KeepAlive(keepAliveCtx, leaseID)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start keepalive: %w", err)
	}

	go func() {
		defer cancel()
		for {
			select {
			case _, ok := <-keepAliveCh:
				if !ok {
					log.ZWarn(keepAliveCtx, "KeepAlive channel closed, lease may have expired", nil,
						"taskName", taskName,
						"instanceID", e.instanceID,
						"leaseID", leaseID)
					return
				}

			case <-keepAliveCtx.Done():
				log.ZDebug(keepAliveCtx, "KeepAlive stopped",
					"taskName", taskName,
					"instanceID", e.instanceID)
				return
			}
		}
	}()

	return cancel, nil
}

func (e *EtcdLocker) ExecuteWithLock(ctx context.Context, taskName string, task func()) {
	leaseID, acquired, err := e.tryAcquireTaskLock(ctx, taskName)
	if err != nil {
		log.ZWarn(ctx, "Failed to acquire task lock", err,
			"taskName", taskName,
			"instanceID", e.instanceID)
		return
	}

	if !acquired {
		log.ZDebug(ctx, "Task is being executed by another instance, skipping",
			"taskName", taskName,
			"instanceID", e.instanceID)
		return
	}

	cancelKeepAlive, err := e.startLeaseKeepAlive(ctx, taskName, leaseID)
	if err != nil {
		log.ZWarn(ctx, "Failed to start lease keepalive", err,
			"taskName", taskName,
			"instanceID", e.instanceID)
		e.releaseTaskLock(ctx, taskName, leaseID)
		return
	}

	defer func() {
		cancelKeepAlive()
		e.releaseTaskLock(ctx, taskName, leaseID)
	}()

	log.ZInfo(ctx, "Starting task execution with lease keepalive",
		"taskName", taskName,
		"instanceID", e.instanceID,
		"leaseID", leaseID)

	task()

	log.ZInfo(ctx, "Task execution completed",
		"taskName", taskName,
		"instanceID", e.instanceID)
}
