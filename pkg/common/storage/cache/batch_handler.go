package cache

import (
	"context"
)

// BatchDeleter interface defines a set of methods for batch deleting cache and publishing deletion information.
type BatchDeleter interface {
	//ChainExecDel method is used for chain calls and must call Clone to prevent memory pollution.
	ChainExecDel(ctx context.Context) error
	//ExecDelWithKeys method directly takes keys for deletion.
	ExecDelWithKeys(ctx context.Context, keys []string) error
	//Clone method creates a copy of the BatchDeleter to avoid modifying the original object.
	Clone() BatchDeleter
	//AddKeys method adds keys to be deleted.
	AddKeys(keys ...string)
}
