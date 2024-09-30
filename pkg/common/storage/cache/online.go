package cache

import "context"

type OnlineCache interface {
	GetOnline(ctx context.Context, userID string) ([]int32, error)
	SetUserOnline(ctx context.Context, userID string, online, offline []int32) error
	GetAllOnlineUsers(ctx context.Context, cursor uint64) (map[string][]int32, uint64, error)
}
