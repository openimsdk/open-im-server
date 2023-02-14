package controller

import (
	"Open_IM/pkg/common/db/cache"
	"context"
	"github.com/go-redis/redis/v8"
)

type AuthInterface interface {
	//结果为空 不返回错误
	GetTokensWithoutError(ctx context.Context, userID, platform string) (map[string]int, error)

	//创建token
	CreateToken(ctx context.Context, userID string, platform string) (string, error)
}

type AuthController struct {
	database *cache.TokenRedis
}

func NewAuthController(rdb redis.UniversalClient, accessSecret string, accessExpire int64) *AuthController {
	return &AuthController{database: cache.NewTokenRedis(cache.NewRedisClient(rdb), accessSecret, accessExpire)}
}

// 结果为空 不返回错误
func (a *AuthController) GetTokensWithoutError(ctx context.Context, userID, platform string) (map[string]int, error) {
	return a.database.GetTokensWithoutError(ctx, userID, platform)
}

// 创建token
func (a *AuthController) CreateToken(ctx context.Context, userID string, platform string) (string, error) {
	return a.database.CreateToken(ctx, userID, platform)
}
