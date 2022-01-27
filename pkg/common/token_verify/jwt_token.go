package token_verify

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	commonDB "Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"github.com/garyburd/redigo/redis"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

//var (
//	TokenExpired     = errors.New("token is timed out, please log in again")
//	TokenInvalid     = errors.New("token has been invalidated")
//	TokenNotValidYet = errors.New("token not active yet")
//	TokenMalformed   = errors.New("that's not even a token")
//	TokenUnknown     = errors.New("couldn't handle this token")
//)

type Claims struct {
	UID      string
	Platform string //login platform
	jwt.RegisteredClaims
}

func BuildClaims(uid, platform string, ttl int64) Claims {
	now := time.Now()
	return Claims{
		UID:      uid,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttl*24) * time.Hour)), //Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                        //Issuing time
			NotBefore: jwt.NewNumericDate(now),                                        //Begin Effective time
		}}
}

func CreateToken(userID string, platformID int32) (string, int64, error) {
	claims := BuildClaims(userID, constant.PlatformIDToName(platformID), config.Config.TokenPolicy.AccessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.TokenPolicy.AccessSecret))
	if err != nil {
		return "", 0, err
	}
	//remove Invalid token
	m, err := commonDB.DB.GetTokenMapByUidPid(userID, constant.PlatformIDToName(platformID))
	if err != nil && err != redis.ErrNil {
		return "", 0, err
	}
	var deleteTokenKey []string
	for k, v := range m {
		_, err = GetClaimFromToken(k)
		if err != nil || v != constant.NormalToken {
			deleteTokenKey = append(deleteTokenKey, k)
		}
	}
	if len(deleteTokenKey) != 0 {
		err = commonDB.DB.DeleteTokenByUidPid(userID, platformID, deleteTokenKey)
		if err != nil {
			return "", 0, err
		}
	}
	err = commonDB.DB.AddTokenFlag(userID, platformID, tokenString, constant.NormalToken)
	if err != nil {
		return "", 0, err
	}
	return tokenString, claims.ExpiresAt.Time.Unix(), err
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
				return nil, &constant.ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, &constant.ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, &constant.ErrTokenNotValidYet
			} else {
				return nil, &constant.ErrTokenUnknown
			}
		} else {
			return nil, &constant.ErrTokenNotValidYet
		}
	} else {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			log.NewDebug("", claims.UID, claims.Platform)
			return claims, nil
		}
		return nil, &constant.ErrTokenNotValidYet
	}
}

func IsAppManagerAccess(token string, OpUserID string) bool {
	claims, err := ParseToken(token)
	if err != nil {
		return false
	}
	if utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) && claims.UID == OpUserID {
		return true
	}
	return false
}

func IsMangerUserID(OpUserID string) bool {
	if utils.IsContain(OpUserID, config.Config.Manager.AppManagerUid) {
		return true
	} else {
		return false
	}
}

func CheckAccess(OpUserID string, OwnerUserID string) bool {
	if utils.IsContain(OpUserID, config.Config.Manager.AppManagerUid) {
		return true
	}
	if OpUserID == OwnerUserID {
		return true
	}
	return false
}

func GetUserIDFromToken(token string) (bool, string) {
	claims, err := ParseToken(token)
	if err != nil {
		return false, ""
	}
	return true, claims.UID
}

func ParseToken(tokensString string) (claims *Claims, err error) {

	claims, err = GetClaimFromToken(tokensString)
	if err != nil {
		log.NewError("", "token validate err", err.Error())
		return nil, err
	}

	m, err := commonDB.DB.GetTokenMapByUidPid(claims.UID, claims.Platform)
	if err != nil {
		log.NewError("", "get token from redis err", err.Error())
		return nil, &constant.ErrTokenInvalid
	}
	if m == nil {
		log.NewError("", "get token from redis err", "m is nil")
		return nil, &constant.ErrTokenInvalid
	}
	if v, ok := m[tokensString]; ok {
		switch v {
		case constant.NormalToken:
			log.NewDebug("", "this is normal return", claims)
			return claims, nil
		case constant.InValidToken:
			return nil, &constant.ErrTokenInvalid
		case constant.KickedToken:
			return nil, &constant.ErrTokenInvalid
		case constant.ExpiredToken:
			return nil, &constant.ErrTokenExpired
		default:
			return nil, &constant.ErrTokenUnknown
		}
	}
	return nil, &constant.ErrTokenUnknown
}

//func MakeTheTokenInvalid(currentClaims *Claims, platformClass string) (bool, error) {
//	storedRedisTokenInterface, err := db.DB.GetPlatformToken(currentClaims.UID, platformClass)
//	if err != nil {
//		return false, err
//	}
//	storedRedisPlatformClaims, err := ParseRedisInterfaceToken(storedRedisTokenInterface)
//	if err != nil {
//		return false, err
//	}
//	//if issue time less than redis token then make this token invalid
//	if currentClaims.IssuedAt.Time.Unix() < storedRedisPlatformClaims.IssuedAt.Time.Unix() {
//		return true, constant.TokenInvalid
//	}
//	return false, nil
//}

func ParseRedisInterfaceToken(redisToken interface{}) (*Claims, error) {
	return GetClaimFromToken(string(redisToken.([]uint8)))
}

//Validation token, false means failure, true means successful verification
func VerifyToken(token, uid string) (bool, error) {
	claims, err := ParseToken(token)
	if err != nil {
		return false, err
	}
	if claims.UID != uid {
		return false, &constant.ErrTokenUnknown
	}
	log.NewDebug("", claims.UID, claims.Platform)
	return true, nil
}
