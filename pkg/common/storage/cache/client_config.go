package cache

import "context"

type ClientConfigCache interface {
	DeleteUserCache(ctx context.Context, userIDs []string) error
	GetUserConfig(ctx context.Context, userID string) (map[string]string, error)
}
