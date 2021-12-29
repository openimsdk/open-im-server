package utils

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_BuildClaims(t *testing.T) {
	uid := "1"
	platform := "PC"
	ttl := int64(-1)
	claim := token_verify.BuildClaims(uid, platform, ttl)
	now := time.Now().Unix()

	assert.Equal(t, claim.UID, uid, "uid should equal")
	assert.Equal(t, claim.Platform, platform, "platform should equal")
	assert.Equal(t, claim.RegisteredClaims.ExpiresAt, int64(-1), "StandardClaims.ExpiresAt should be equal")
	// time difference within 1s
	assert.Equal(t, claim.RegisteredClaims.IssuedAt, now, "StandardClaims.IssuedAt should be equal")
	assert.Equal(t, claim.RegisteredClaims.NotBefore, now, "StandardClaims.NotBefore should be equal")

	ttl = int64(60)
	now = time.Now().Unix()
	claim = token_verify.BuildClaims(uid, platform, ttl)
	// time difference within 1s
	assert.Equal(t, claim.RegisteredClaims.ExpiresAt, int64(60)+now, "StandardClaims.ExpiresAt should be equal")
	assert.Equal(t, claim.RegisteredClaims.IssuedAt, now, "StandardClaims.IssuedAt should be equal")
	assert.Equal(t, claim.RegisteredClaims.NotBefore, now, "StandardClaims.NotBefore should be equal")
}

func Test_CreateToken(t *testing.T) {
	uid := "1"
	platform := int32(1)
	now := time.Now().Unix()

	tokenString, expiresAt, err := token_verify.CreateToken(uid, platform)

	assert.NotEmpty(t, tokenString)
	assert.Equal(t, expiresAt, 604800+now)
	assert.Nil(t, err)
}

func Test_VerifyToken(t *testing.T) {
	uid := "1"
	platform := int32(1)
	tokenString, _, _ := token_verify.CreateToken(uid, platform)
	result, _ := token_verify.VerifyToken(tokenString, uid)
	assert.True(t, result)
	result, _ = token_verify.VerifyToken(tokenString, "2")
	assert.False(t, result)
}

func Test_ParseRedisInterfaceToken(t *testing.T) {
	uid := "1"
	platform := int32(1)
	tokenString, _, _ := token_verify.CreateToken(uid, platform)

	claims, err := token_verify.ParseRedisInterfaceToken([]uint8(tokenString))
	assert.Nil(t, err)
	assert.Equal(t, claims.UID, uid)

	// timeout
	config.Config.TokenPolicy.AccessExpire = -80
	tokenString, _, _ = token_verify.CreateToken(uid, platform)
	claims, err = token_verify.ParseRedisInterfaceToken([]uint8(tokenString))
	assert.Equal(t, err, constant.ExpiredToken)
	assert.Nil(t, claims)
}

func Test_ParseToken(t *testing.T) {
	uid := "1"
	platform := int32(1)
	tokenString, _, _ := token_verify.CreateToken(uid, platform)
	claims, err := token_verify.ParseToken(tokenString)
	if err == nil {
		assert.Equal(t, claims.UID, uid)
	}
}
func Test_GetClaimFromToken(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiJvcGVuSU0xMjM0NTYiLCJQbGF0Zm9ybSI6IiIsImV4cCI6MTYzODg0NjQ3NiwibmJmIjoxNjM4MjQxNjc2LCJpYXQiOjE2MzgyNDE2NzZ9.W8RZB7ec5ySFj-rGE2Aho2z32g3MprQMdCyPiQu_C2I"
	c, err := token_verify.GetClaimFromToken(token)
	assert.Nil(t, c)
	assert.Nil(t, err)
}
