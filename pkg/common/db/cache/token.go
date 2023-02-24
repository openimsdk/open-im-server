package cache

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/utils"
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type Token interface {
	//结果为空 不返回错误
	GetTokensWithoutError(ctx context.Context, userID, platform string) (map[string]int, error)
	//创建token
	CreateToken(ctx context.Context, userID string, platformID int) (string, error)
}

type TokenRedis struct {
	redisClient  *RedisClient
	accessSecret string
	accessExpire int64
}

func NewTokenRedis(redisClient *RedisClient, accessSecret string, accessExpire int64) *TokenRedis {
	return &TokenRedis{redisClient, accessSecret, accessExpire}
}

// 结果为空 不返回错误
func (t *TokenRedis) GetTokensWithoutError(ctx context.Context, userID, platform string) (map[string]int, error) {
	return t.redisClient.GetTokensWithoutError(ctx, userID, platform)
}

// 创建token
func (t *TokenRedis) CreateToken(ctx context.Context, userID string, platform string) (string, error) {
	tokens, err := t.redisClient.GetTokensWithoutError(ctx, userID, platform)
	if err != nil {
		return "", err
	}
	var deleteTokenKey []string
	for k, v := range tokens {
		_, err = tokenverify.GetClaimFromToken(k)
		if err != nil || v != constant.NormalToken {
			deleteTokenKey = append(deleteTokenKey, k)
		}
	}
	if len(deleteTokenKey) != 0 {
		err := t.redisClient.DeleteTokenByUidPid(ctx, userID, platform, deleteTokenKey)
		if err != nil {
			return "", err
		}
	}
	claims := tokenverify.BuildClaims(userID, platform, t.accessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.accessSecret))
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	return tokenString, t.redisClient.AddTokenFlag(ctx, userID, constant.PlatformNameToID(platform), tokenString, constant.NormalToken)
}
