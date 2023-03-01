package controller

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/cache"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/utils"
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type AuthDatabase interface {
	//结果为空 不返回错误
	GetTokensWithoutError(ctx context.Context, userID, platform string) (map[string]int, error)
	//创建token
	CreateToken(ctx context.Context, userID string, platform string) (string, error)
}

type authDatabase struct {
	cache cache.Cache

	accessSecret string
	accessExpire int64
}

func NewAuthDatabase(cache cache.Cache, accessSecret string, accessExpire int64) AuthDatabase {
	return &authDatabase{cache: cache, accessSecret: accessSecret, accessExpire: accessExpire}
}

// 结果为空 不返回错误
func (a *authDatabase) GetTokensWithoutError(ctx context.Context, userID, platform string) (map[string]int, error) {
	return a.cache.GetTokensWithoutError(ctx, userID, platform)
}

// 创建token
func (a *authDatabase) CreateToken(ctx context.Context, userID string, platform string) (string, error) {
	tokens, err := a.cache.GetTokensWithoutError(ctx, userID, platform)
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
		err := a.cache.DeleteTokenByUidPid(ctx, userID, platform, deleteTokenKey)
		if err != nil {
			return "", err
		}
	}
	claims := tokenverify.BuildClaims(userID, platform, a.accessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.accessSecret))
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	return tokenString, a.cache.AddTokenFlag(ctx, userID, constant.PlatformNameToID(platform), tokenString, constant.NormalToken)
}
