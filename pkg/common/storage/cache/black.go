package cache

import (
	"context"
)

type BlackCache interface {
	BatchDeleter
	CloneBlackCache() BlackCache
	GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error)
	// del user's blackIDs msgCache, exec when a user's black list changed
	DelBlackIDs(ctx context.Context, userID string) BlackCache
}
