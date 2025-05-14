package cache

import (
	"context"
)

type TokenModel interface {
	SetTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error
	// SetTokenFlagEx set token and flag with expire time
	SetTokenFlagEx(ctx context.Context, userID string, platformID int, token string, flag int) error
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)
	HasTemporaryToken(ctx context.Context, userID string, platformID int, token string) error
	GetAllTokensWithoutError(ctx context.Context, userID string) (map[int]map[string]int, error)
	SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error
	BatchSetTokenMapByUidPid(ctx context.Context, tokens map[string]map[string]any) error
	DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error
	DeleteTokenByTokenMap(ctx context.Context, userID string, tokens map[int][]string) error
	DeleteAndSetTemporary(ctx context.Context, userID string, platformID int, fields []string) error
}
