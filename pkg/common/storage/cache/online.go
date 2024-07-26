package cache

import "context"

type OnlineCache interface {
	GetOnline(ctx context.Context, userID string) ([]int32, error)
	SetUserOnline(ctx context.Context, userID string, online, offline []int32) error
}
