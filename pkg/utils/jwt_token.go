package utils

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	TokenExpired     = errors.New("token is timed out, please log in again")
	TokenInvalid     = errors.New("token has been invalidated")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenUnknown     = errors.New("couldn't handle this token")
)

type Claims struct {
	UID      string
	Platform string //login platform
	jwt.StandardClaims
}

func BuildClaims(uid, accountAddr, platform string, ttl int64) Claims {
	now := time.Now().Unix()
	//if ttl=-1 Permanent token
	if ttl == -1 {
		return Claims{
			UID:      uid,
			Platform: platform,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: -1,
				IssuedAt:  now,
				NotBefore: now,
			}}
	}
	return Claims{
		UID:      uid,
		Platform: platform,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now + ttl, //Expiration time
			IssuedAt:  now,       //Issuing time
			NotBefore: now,       //Begin Effective time
		}}
}

func CreateToken(userID, accountAddr string, platform int32) (string, int64, error) {
	claims := BuildClaims(userID, accountAddr, PlatformIDToName(platform), config.Config.TokenPolicy.AccessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.TokenPolicy.AccessSecret))

	return tokenString, claims.ExpiresAt, err
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenPolicy.AccessSecret), nil
	}
}

func ParseToken(tokensString string) (claims *Claims, err error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenUnknown
			}
		}
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		//	1.check userid and platform class   0 not exists and  1 exists
		existsInterface, err := db.DB.ExistsUserIDAndPlatform(claims.UID, Platform2class[claims.Platform])
		if err != nil {
			return nil, err
		}
		exists := existsInterface.(int64)
		//get config multi login policy
		if config.Config.MultiLoginPolicy.OnlyOneTerminalAccess {
			//OnlyOneTerminalAccess policy need to check all terminal
			//When only one end is allowed to log in, there is a situation that needs to be paid attention to. After PC login,
			//mobile login should check two platform times. One of them is less than the redis storage time, which is the invalid token.
			if Platform2class[claims.Platform] == "PC" {
				existsInterface, err = db.DB.ExistsUserIDAndPlatform(claims.UID, "Mobile")
				if err != nil {
					return nil, err
				}
				exists = existsInterface.(int64)
				if exists == 1 {
					res, err := MakeTheTokenInvalid(*claims, "Mobile")
					if err != nil {
						return nil, err
					}
					if res {
						return nil, TokenInvalid
					}
				}
			} else {
				existsInterface, err = db.DB.ExistsUserIDAndPlatform(claims.UID, "PC")
				if err != nil {
					return nil, err
				}
				exists = existsInterface.(int64)
				if exists == 1 {
					res, err := MakeTheTokenInvalid(*claims, "PC")
					if err != nil {
						return nil, err
					}
					if res {
						return nil, TokenInvalid
					}
				}
			}

			if exists == 1 {
				res, err := MakeTheTokenInvalid(*claims, Platform2class[claims.Platform])
				if err != nil {
					return nil, err
				}
				if res {
					return nil, TokenInvalid
				}
			}

		} else if config.Config.MultiLoginPolicy.MobileAndPCTerminalAccessButOtherTerminalKickEachOther {
			if exists == 1 {
				res, err := MakeTheTokenInvalid(*claims, Platform2class[claims.Platform])
				if err != nil {
					return nil, err
				}
				if res {
					return nil, TokenInvalid
				}
			}
		}
		return claims, nil
	}
	return nil, TokenUnknown
}

func MakeTheTokenInvalid(currentClaims Claims, platformClass string) (bool, error) {
	storedRedisTokenInterface, err := db.DB.GetPlatformToken(currentClaims.UID, platformClass)
	if err != nil {
		return false, err
	}
	storedRedisPlatformClaims, err := ParseRedisInterfaceToken(storedRedisTokenInterface)
	if err != nil {
		return false, err
	}
	//if issue time less than redis token then make this token invalid
	if currentClaims.IssuedAt < storedRedisPlatformClaims.IssuedAt {
		return true, TokenInvalid
	}
	return false, nil
}
func ParseRedisInterfaceToken(redisToken interface{}) (*Claims, error) {
	token, err := jwt.ParseWithClaims(string(redisToken.([]uint8)), &Claims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

//Validation token, false means failure, true means successful verification
func VerifyToken(token, uid string) bool {
	claims, err := ParseToken(token)
	if err != nil {
		return false
	} else if claims.UID != uid {
		return false
	} else {
		return true
	}
}
