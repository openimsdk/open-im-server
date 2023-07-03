package tokenverify

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
)

func Test_ParseToken(t *testing.T) {
	config.Config.Secret = "OpenIM_server"
	claims1 := BuildClaims("123456", constant.AndroidPadPlatformID, 10)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims1)
	tokenString, err := token.SignedString([]byte(config.Config.Secret))
	if err != nil {
		t.Fatal(err)
	}
	claim2, err := GetClaimFromToken(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claim2)
}
