package token_verify

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	UID      string
	Platform string //login platform
	jwt.RegisteredClaims
}

func BuildClaims(uid, platform string, ttl int64) Claims {
	now := time.Now()
	before := now.Add(-time.Minute * 5)
	return Claims{
		UID:      uid,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttl*24) * time.Hour)), //Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                        //Issuing time
			NotBefore: jwt.NewNumericDate(before),                                     //Begin Effective time
		}}
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenPolicy.AccessSecret), nil
	}
}

func GetClaimFromToken(tokensString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, utils.Wrap(constant.ErrTokenMalformed, "")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, utils.Wrap(constant.ErrTokenExpired, "")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, utils.Wrap(constant.ErrTokenNotValidYet, "")
			} else {
				return nil, utils.Wrap(constant.ErrTokenUnknown, "")
			}
		} else {
			return nil, utils.Wrap(constant.ErrTokenUnknown, "")
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, utils.Wrap(constant.ErrTokenUnknown, "")
	}
}

func CheckAccessV3(ctx context.Context, ownerUserID string) (err error) {
	opUserID := tracelog.GetOpUserID(ctx)
	defer func() {
		tracelog.SetCtxInfo(ctx, utils.GetFuncName(1), err, "OpUserID", opUserID, "ownerUserID", ownerUserID)
	}()
	if utils.IsContain(opUserID, config.Config.Manager.AppManagerUid) {
		return nil
	}
	if opUserID == ownerUserID {
		return nil
	}
	return constant.ErrIdentity.Wrap(utils.GetSelfFuncName())
}

func CheckAdmin(ctx context.Context) error {
	if utils.IsContain(tracelog.GetOpUserID(ctx), config.Config.Manager.AppManagerUid) {
		return nil
	}
	return constant.ErrIdentity.Wrap()
}

func ParseRedisInterfaceToken(redisToken interface{}) (*Claims, error) {
	return GetClaimFromToken(string(redisToken.([]uint8)))
}
