package utils

type DistributedLock interface {
	Lock()
	UnLock()
}
