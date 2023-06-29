package tokenverify

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	UserID     string
	PlatformID int //login platform
	jwt.RegisteredClaims
}

func BuildClaims(uid string, platformID int, ttl int64) Claims {
	now := time.Now()
	before := now.Add(-time.Minute * 5)
	return Claims{
		UserID:     uid,
		PlatformID: platformID,
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
				return nil, utils.Wrap(errs.ErrTokenMalformed, "")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, utils.Wrap(errs.ErrTokenExpired, "")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, utils.Wrap(errs.ErrTokenNotValidYet, "")
			} else {
				return nil, utils.Wrap(errs.ErrTokenUnknown, "")
			}
		} else {
			return nil, utils.Wrap(errs.ErrTokenUnknown, "")
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, utils.Wrap(errs.ErrTokenUnknown, "")
	}
}

func CheckAccessV3(ctx context.Context, ownerUserID string) (err error) {
	opUserID := mcontext.GetOpUserID(ctx)
	if utils.IsContain(opUserID, config.Config.Manager.UserID) {
		return nil
	}
	if opUserID == ownerUserID {
		return nil
	}
	return errs.ErrIdentity.Wrap(utils.GetSelfFuncName())
}

func IsAppManagerUid(ctx context.Context) bool {
	return utils.IsContain(mcontext.GetOpUserID(ctx), config.Config.Manager.UserID)
}

func CheckAdmin(ctx context.Context) error {
	if utils.IsContain(mcontext.GetOpUserID(ctx), config.Config.Manager.UserID) {
		return nil
	}
	return errs.ErrIdentity.Wrap(fmt.Sprintf("user %s is not admin userID", mcontext.GetOpUserID(ctx)))
}

func ParseRedisInterfaceToken(redisToken interface{}) (*Claims, error) {
	return GetClaimFromToken(string(redisToken.([]uint8)))
}
func IsManagerUserID(opUserID string) bool {
	return utils.IsContain(opUserID, config.Config.Manager.UserID)
}
func WsVerifyToken(token, userID string, platformID int) error {
	claim, err := GetClaimFromToken(token)
	if err != nil {
		return err
	}
	if claim.UserID != userID {
		return errs.ErrTokenInvalid.Wrap(fmt.Sprintf("token uid %s != userID %s", claim.UserID, userID))
	}
	if claim.PlatformID != platformID {
		return errs.ErrTokenInvalid.Wrap(fmt.Sprintf("token platform %d != %d", claim.PlatformID, platformID))
	}
	return nil
}
