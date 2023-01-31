package controller

import "context"

type AuthInterface interface {
	GetTokens(ctx context.Context, userID, platform string) (map[string]int, error)
	DeleteToken(ctx context.Context, userID, platform string) error
}
