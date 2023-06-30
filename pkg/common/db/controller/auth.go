package controller

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

type AuthDatabase interface {
	//结果为空 不返回错误
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)
	//创建token
	CreateToken(ctx context.Context, userID string, platformID int) (string, error)
}

type authDatabase struct {
	cache cache.MsgModel

	accessSecret string
	accessExpire int64
}

func NewAuthDatabase(cache cache.MsgModel, accessSecret string, accessExpire int64) AuthDatabase {
	return &authDatabase{cache: cache, accessSecret: accessSecret, accessExpire: accessExpire}
}

// 结果为空 不返回错误
func (a *authDatabase) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	return a.cache.GetTokensWithoutError(ctx, userID, platformID)
}

// 创建token
func (a *authDatabase) CreateToken(ctx context.Context, userID string, platformID int) (string, error) {
	tokens, err := a.cache.GetTokensWithoutError(ctx, userID, platformID)
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
		err := a.cache.DeleteTokenByUidPid(ctx, userID, platformID, deleteTokenKey)
		if err != nil {
			return "", err
		}
	}
	claims := tokenverify.BuildClaims(userID, platformID, a.accessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.accessSecret))
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	return tokenString, a.cache.AddTokenFlag(ctx, userID, platformID, tokenString, constant.NormalToken)
}
