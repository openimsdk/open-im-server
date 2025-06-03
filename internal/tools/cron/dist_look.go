package cron

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/openimsdk/tools/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	lockKey           = "openim/crontask/dist-lock"
	lockLeaseTTL      = 15 // Lease TTL in seconds
	acquireRetryDelay = 500 * time.Millisecond
)

type EtcdLocker struct {
	client       *clientv3.Client
	instanceID   string
	leaseID      clientv3.LeaseID
	isLockOwner  int32 // Using atomic for lock ownership check
	watchCh      clientv3.WatchChan
	watchCancel  context.CancelFunc
	leaseTTL     int64
	stopCh       chan struct{}
	stoppedCh    chan struct{}
	acquireDelay time.Duration
}

// NewEtcdLocker creates a new etcd distributed lock
func NewEtcdLocker(client *clientv3.Client) (*EtcdLocker, error) {
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	instanceID := fmt.Sprintf("%s-pid-%d-%d", hostname, pid, time.Now().UnixNano())

	locker := &EtcdLocker{
		client:       client,
		instanceID:   instanceID,
		leaseTTL:     lockLeaseTTL,
		stopCh:       make(chan struct{}),
		stoppedCh:    make(chan struct{}),
		acquireDelay: acquireRetryDelay,
	}

	return locker, nil
}

func (e *EtcdLocker) Start(ctx context.Context) error {
	log.ZInfo(ctx, "Starting etcd distributed lock", "instanceID", e.instanceID)
	go e.runLockLoop(ctx)
	return nil
}

func (e *EtcdLocker) runLockLoop(ctx context.Context) {
	defer close(e.stoppedCh)
	for {
		select {
		case <-e.stopCh:
			e.releaseLock(ctx)
			return
		case <-ctx.Done():
			e.releaseLock(ctx)
			return
		default:
			acquired, err := e.tryAcquireLock(ctx)
			if err != nil {
				log.ZWarn(ctx, "Failed to acquire lock", err, "instanceID", e.instanceID)
				time.Sleep(e.acquireDelay)
				continue
			}

			if acquired {
				e.runKeepAlive(ctx)
				time.Sleep(e.acquireDelay)
			} else {
				e.watchLock(ctx)
			}
		}
	}
}

func (e *EtcdLocker) tryAcquireLock(ctx context.Context) (bool, error) {
	lease, err := e.client.Grant(ctx, e.leaseTTL)
	if err != nil {
		return false, fmt.Errorf("failed to create lease: %w", err)
	}

	txnResp, err := e.client.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, e.instanceID, clientv3.WithLease(lease.ID))).
		Else(clientv3.OpGet(lockKey)).
		Commit()

	if err != nil {
		e.client.Revoke(ctx, lease.ID)
		return false, fmt.Errorf("transaction failed: %w", err)
	}

	if !txnResp.Succeeded {
		rangeResp := txnResp.Responses[0].GetResponseRange()
		if len(rangeResp.Kvs) > 0 {
			currentOwner := string(rangeResp.Kvs[0].Value)
			log.ZInfo(ctx, "Lock already owned", "instanceID", e.instanceID, "owner", currentOwner)
		}

		e.client.Revoke(ctx, lease.ID)
		return false, nil
	}

	e.leaseID = lease.ID
	atomic.StoreInt32(&e.isLockOwner, 1)
	log.ZInfo(ctx, "Successfully acquired lock", "instanceID", e.instanceID, "leaseID", lease.ID)
	return true, nil
}

func (e *EtcdLocker) runKeepAlive(ctx context.Context) {
	keepAliveCh, err := e.client.KeepAlive(ctx, e.leaseID)
	if err != nil {
		log.ZError(ctx, "Failed to start lease keepalive", err, "instanceID", e.instanceID)
		e.releaseLock(ctx)
		return
	}

	for {
		select {
		case _, ok := <-keepAliveCh:
			if !ok {
				log.ZWarn(ctx, "Keepalive channel closed, lock lost", nil, "instanceID", e.instanceID)
				atomic.StoreInt32(&e.isLockOwner, 0) // Set to false atomically
				return
			}
		case <-ctx.Done():
			log.ZInfo(ctx, "Context canceled, releasing lock", "instanceID", e.instanceID)
			e.releaseLock(ctx)
			return
		case <-e.stopCh:
			log.ZInfo(ctx, "Stop signal received, releasing lock", "instanceID", e.instanceID)
			e.releaseLock(ctx)
			return
		}
	}
}

// Watch lock status directly in etcd
func (e *EtcdLocker) watchLock(ctx context.Context) {
	log.ZInfo(ctx, "Starting to watch lock status", "instanceID", e.instanceID)
	watchCtx, cancel := context.WithCancel(ctx)
	e.watchCancel = cancel
	defer e.cancelWatch()

	// Watch for changes to the lock key
	e.watchCh = e.client.Watch(watchCtx, lockKey)
	for {
		select {
		case resp, ok := <-e.watchCh:
			if !ok {
				log.ZWarn(ctx, "Watch channel closed", nil, "instanceID", e.instanceID)
				return
			}
			for _, event := range resp.Events {
				if event.Type == clientv3.EventTypeDelete {
					log.ZInfo(ctx, "Lock released, attempting to acquire", "instanceID", e.instanceID)
					return
				}
			}
		case <-ctx.Done():
			return
		case <-e.stopCh:
			return
		}
	}
}

// Release the lock
func (e *EtcdLocker) releaseLock(ctx context.Context) {
	if atomic.LoadInt32(&e.isLockOwner) == 0 {
		return
	}

	leaseID := e.leaseID
	atomic.StoreInt32(&e.isLockOwner, 0)
	e.leaseID = 0
	if leaseID != 0 {
		_, err := e.client.Revoke(context.Background(), leaseID)
		if err != nil {
			log.ZWarn(ctx, "Failed to revoke lease", err, "instanceID", e.instanceID, "error", err)
		} else {
			log.ZInfo(ctx, "Successfully released lock", "instanceID", e.instanceID)
		}
	}
}

func (e *EtcdLocker) CheckLockOwnership(ctx context.Context) (bool, error) {
	if atomic.LoadInt32(&e.isLockOwner) == 0 {
		return false, nil
	}

	resp, err := e.client.Get(ctx, lockKey)
	if err != nil {
		return false, fmt.Errorf("failed to check lock status: %w", err)
	}
	if len(resp.Kvs) > 0 && string(resp.Kvs[0].Value) == e.instanceID {
		return true, nil
	}
	if atomic.LoadInt32(&e.isLockOwner) == 1 {
		log.ZWarn(ctx, "Lock ownership lost unexpectedly", nil, "instanceID", e.instanceID)
		atomic.StoreInt32(&e.isLockOwner, 0)
	}

	return false, nil
}

func (e *EtcdLocker) cancelWatch() {
	if e.watchCancel != nil {
		e.watchCancel()
		e.watchCancel = nil
	}
}

func (e *EtcdLocker) Stop() {
	e.cancelWatch()
	close(e.stopCh)
	<-e.stoppedCh
}

func (e *EtcdLocker) IsLockOwner() bool {
	return atomic.LoadInt32(&e.isLockOwner) == 1
}

func (e *EtcdLocker) ExecuteWithLock(ctx context.Context, task func()) {
	if atomic.LoadInt32(&e.isLockOwner) == 0 {
		log.ZDebug(ctx, "Instance does not own lock (local check), skipping task execution", "instanceID", e.instanceID)
		return
	}
	isOwner, err := e.CheckLockOwnership(ctx)
	if err != nil {
		log.ZWarn(ctx, "Failed to verify lock ownership", err, "instanceID", e.instanceID)
		return
	}
	if !isOwner {
		log.ZDebug(ctx, "Instance does not own lock (etcd verification), skipping task execution", "instanceID", e.instanceID)
		return
	}

	log.ZInfo(ctx, "Starting lock-protected task execution", "instanceID", e.instanceID)
	task()
	log.ZInfo(ctx, "Lock-protected task execution completed", "instanceID", e.instanceID)
}
