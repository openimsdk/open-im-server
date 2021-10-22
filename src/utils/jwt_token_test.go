package utils

import (
	"Open_IM/src/common/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_BuildClaims(t *testing.T) {
	uid := "1"
	accountAddr := "accountAddr"
	platform := "PC"
	ttl := int64(-1)
	claim := BuildClaims(uid, accountAddr, platform, ttl)
	now := time.Now().Unix()

	assert.Equal(t, claim.UID, uid, "uid should equal")
	assert.Equal(t, claim.Platform, platform, "platform should equal")
	assert.Equal(t, claim.StandardClaims.ExpiresAt, int64(-1), "StandardClaims.ExpiresAt should be equal")
	// time difference within 1s
	assert.Equal(t, claim.StandardClaims.IssuedAt, now, "StandardClaims.IssuedAt should be equal")
	assert.Equal(t, claim.StandardClaims.NotBefore, now, "StandardClaims.NotBefore should be equal")

	ttl = int64(60)
	now = time.Now().Unix()
	claim = BuildClaims(uid, accountAddr, platform, ttl)
	// time difference within 1s
	assert.Equal(t, claim.StandardClaims.ExpiresAt, int64(60)+now, "StandardClaims.ExpiresAt should be equal")
	assert.Equal(t, claim.StandardClaims.IssuedAt, now, "StandardClaims.IssuedAt should be equal")
	assert.Equal(t, claim.StandardClaims.NotBefore, now, "StandardClaims.NotBefore should be equal")
}

func Test_CreateToken(t *testing.T) {
	uid := "1"
	accountAddr := "accountAddr"
	platform := int32(1)
	now := time.Now().Unix()

	tokenString, expiresAt, err := CreateToken(uid, accountAddr, platform)

	assert.NotEmpty(t, tokenString)
	assert.Equal(t, expiresAt, 604800+now)
	assert.Nil(t, err)
}

func Test_VerifyToken(t *testing.T) {
	uid := "1"
	accountAddr := "accountAddr"
	platform := int32(1)
	tokenString, _, _ := CreateToken(uid, accountAddr, platform)
	result := VerifyToken(tokenString, uid)
	assert.True(t, result)
	result = VerifyToken(tokenString, "2")
	assert.False(t, result)
}

func Test_ParseRedisInterfaceToken(t *testing.T) {
	uid := "1"
	accountAddr := "accountAddr"
	platform := int32(1)
	tokenString, _, _ := CreateToken(uid, accountAddr, platform)

	claims, err := ParseRedisInterfaceToken([]uint8(tokenString))
	assert.Nil(t, err)
	assert.Equal(t, claims.UID, uid)

	// timeout
	config.Config.TokenPolicy.AccessExpire = -80
	tokenString, _, _ = CreateToken(uid, accountAddr, platform)
	claims, err = ParseRedisInterfaceToken([]uint8(tokenString))
	assert.Equal(t, err, TokenExpired)
	assert.Nil(t, claims)
}

func Test_ParseToken(t *testing.T) {
	uid := "1"
	accountAddr := "accountAddr"
	platform := int32(1)
	tokenString, _, _ := CreateToken(uid, accountAddr, platform)
	claims, err := ParseToken(tokenString)
	if err == nil {
		assert.Equal(t, claims.UID, uid)
	}
}
