package check

import (
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	"context"
)

type BlackChecker struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func NewBlackChecker(zk discoveryRegistry.SvcDiscoveryRegistry) *BlackChecker {
	return &BlackChecker{
		zk: zk,
	}
}

// possibleBlackUserID是否被userID拉黑，也就是是否在userID的黑名单中
func (b *BlackChecker) IsBlocked(ctx context.Context, possibleBlackUserID, userID string) (bool, error) {

}
