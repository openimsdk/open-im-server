package cache

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/tokenverify"
	"Open_IM/pkg/utils"
	"context"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
)

const (
	uidPidToken = "UID_PID_TOKEN_STATUS:"
)

type Token interface {
	//结果为空 不返回错误
	GetTokensWithoutError(ctx context.Context, userID, platform string) (map[string]int, error)
	//创建token
	CreateToken(ctx context.Context, userID string, platformID int) (string, error)
}

type TokenRedis struct {
	RedisClient  *RedisClient
	AccessSecret string
	AccessExpire int64
}

func NewTokenRedis(redisClient *RedisClient, accessSecret string, accessExpire int64) *TokenRedis {
	return &TokenRedis{redisClient, accessSecret, accessExpire}
}

// 结果为空 不返回错误
func (t *TokenRedis) GetTokensWithoutError(ctx context.Context, userID, platform string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platform
	m, err := t.RedisClient.GetClient().HGetAll(context.Background(), key).Result()
	if err != nil && err == go_redis.Nil {
		return nil, nil
	}
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}
	return mm, utils.Wrap(err, "")
}

// 创建token
func (t *TokenRedis) CreateToken(ctx context.Context, userID string, platform string) (string, error) {
	tokens, err := t.GetTokensWithoutError(ctx, userID, platform)
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
		key := uidPidToken + userID + ":" + platform
		err := t.RedisClient.GetClient().HDel(context.Background(), key, deleteTokenKey...).Err()
		if err != nil {
			return "", err
		}
	}
	claims := tokenverify.BuildClaims(userID, platform, t.AccessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(t.AccessSecret))
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	key := uidPidToken + userID + ":" + platform
	return "", utils.Wrap(t.RedisClient.GetClient().HSet(context.Background(), key, tokenString, constant.NormalToken).Err(), "")
}
