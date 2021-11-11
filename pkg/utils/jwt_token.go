package utils

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"errors"
	"github.com/golang-jwt/jwt/v4"
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

func CreateToken(userID string, platform int32) (string, int64, error) {
	claims := BuildClaims(userID, PlatformIDToName(platform), config.Config.TokenPolicy.AccessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.TokenPolicy.AccessSecret))

	return tokenString, claims.ExpiresAt.Time.Unix(), err
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenPolicy.AccessSecret), nil
	}
}

func getClaimFromToken(tokensString string) (*Claims, error) {
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
		return claims, nil
	}
	return nil, err
}

func ParseToken(tokensString string) (claims *Claims, err error) {
	claims, err = getClaimFromToken(tokensString)

	if err != nil {
		return nil, err
	}

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
		platform := "PC"
		if Platform2class[claims.Platform] == "PC" {
			platform = "Mobile"
		}

		existsInterface, err = db.DB.ExistsUserIDAndPlatform(claims.UID, platform)
		if err != nil {
			return nil, err
		}

		exists = existsInterface.(int64)
		if exists == 1 {
			res, err := MakeTheTokenInvalid(claims, platform)
			if err != nil {
				return nil, err
			}
			if res {
				return nil, TokenInvalid
			}
		}
	}
	// config.Config.MultiLoginPolicy.MobileAndPCTerminalAccessButOtherTerminalKickEachOther == true
	// or  PC/Mobile validate success
	// final check
	if exists == 1 {
		res, err := MakeTheTokenInvalid(claims, Platform2class[claims.Platform])
		if err != nil {
			return nil, err
		}
		if res {
			return nil, TokenInvalid
		}
	}
	return claims, nil
}

func MakeTheTokenInvalid(currentClaims *Claims, platformClass string) (bool, error) {
	storedRedisTokenInterface, err := db.DB.GetPlatformToken(currentClaims.UID, platformClass)
	if err != nil {
		return false, err
	}
	storedRedisPlatformClaims, err := ParseRedisInterfaceToken(storedRedisTokenInterface)
	if err != nil {
		return false, err
	}
	//if issue time less than redis token then make this token invalid
	if currentClaims.IssuedAt.Time.Unix() < storedRedisPlatformClaims.IssuedAt.Time.Unix() {
		return true, TokenInvalid
	}
	return false, nil
}

func ParseRedisInterfaceToken(redisToken interface{}) (*Claims, error) {
	return getClaimFromToken(string(redisToken.([]uint8)))
}

//Validation token, false means failure, true means successful verification
func VerifyToken(token, uid string) bool {
	claims, err := ParseToken(token)
	if err != nil || claims.UID != uid {
		return false
	}
	return true
}
