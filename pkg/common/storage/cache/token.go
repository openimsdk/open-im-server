package cache

import (
	"context"
)

type TokenModel interface {
	SetTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error
	// SetTokenFlagEx set token and flag with expire time
	SetTokenFlagEx(ctx context.Context, userID string, platformID int, token string, flag int) error
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)
	SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error
	DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error
}
