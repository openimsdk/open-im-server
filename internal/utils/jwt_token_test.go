package utils

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/token_verify"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BuildClaims(t *testing.T) {
	uid := "1"
	platform := "PC"
	ttl := int64(-1)
	claim := token_verify.BuildClaims(uid, platform, ttl)
	now := time.Now()

	assert.Equal(t, claim.UID, uid, "uid should equal")
	assert.Equal(t, claim.Platform, platform, "platform should equal")
	assert.Equal(t, claim.RegisteredClaims.ExpiresAt.Unix(), now.AddDate(0, 0, int(ttl)).Unix(), "StandardClaims.ExpiresAt should be equal")
	// time difference within 1s
	assert.Equal(t, claim.RegisteredClaims.IssuedAt.Unix(), now.Unix(), "StandardClaims.IssuedAt should be equal")
	assert.Equal(t, claim.RegisteredClaims.NotBefore.Unix(), now.Unix()-300, "StandardClaims.NotBefore should be equal")

	ttl = int64(1)
	now = time.Now()
	claim = token_verify.BuildClaims(uid, platform, ttl)
	// time difference within 1s
	assert.Equal(t, claim.RegisteredClaims.ExpiresAt.Unix(), now.AddDate(0, 0, 1).Unix(), "StandardClaims.ExpiresAt should be equal")
	assert.Equal(t, claim.RegisteredClaims.IssuedAt.Unix(), now.Unix(), "StandardClaims.IssuedAt should be equal")
	assert.Equal(t, claim.RegisteredClaims.NotBefore.Unix(), now.Unix()-300, "StandardClaims.NotBefore should be equal")
}

func Test_CreateToken(t *testing.T) {
	uid := "1"
	platform := int32(1)
	now := time.Now()

	tokenString, expiresAt, err := token_verify.CreateToken(uid, int(platform))

	assert.NotEmpty(t, tokenString)
	assert.Equal(t, expiresAt, now.AddDate(0, 0, int(config.Config.TokenPolicy.AccessExpire)).Unix())
	assert.Nil(t, err)
}

func Test_VerifyToken(t *testing.T) {
	uid := "1"
	platform := int32(1)
	tokenString, _, _ := token_verify.CreateToken(uid, int(platform))
	result, _ := token_verify.VerifyToken(tokenString, uid)
	assert.True(t, result)
	result, _ = token_verify.VerifyToken(tokenString, "2")
	assert.False(t, result)
}

func Test_ParseRedisInterfaceToken(t *testing.T) {
	uid := "1"
	platform := int32(1)
	tokenString, _, _ := token_verify.CreateToken(uid, int(platform))

	claims, err := token_verify.ParseRedisInterfaceToken([]uint8(tokenString))
	assert.Nil(t, err)
	assert.Equal(t, claims.UID, uid)

	// timeout
	ttl := config.Config.TokenPolicy.AccessExpire
	config.Config.TokenPolicy.AccessExpire = -80
	defer func() { config.Config.TokenPolicy.AccessExpire = ttl }()

	tokenString, exp, err := token_verify.CreateToken(uid, int(platform))
	require.NoError(t, err)
	require.Less(t, exp, time.Now().Unix())

	claims, err = token_verify.ParseRedisInterfaceToken([]uint8(tokenString))
	assert.ErrorIs(t, err, constant.ErrTokenExpired)
	assert.Nil(t, claims)
}

func Test_ParseToken(t *testing.T) {
	uid := "1"
	platform := int32(1)
	tokenString, _, _ := token_verify.CreateToken(uid, int(platform))
	claims, err := token_verify.ParseToken(tokenString, "")
	if err == nil {
		assert.Equal(t, claims.UID, uid)
	}
}

func Test_GetClaimFromToken(t *testing.T) {
	token, _, err := token_verify.CreateToken("", constant.IOSPlatformID)
	require.NoError(t, err)

	c, err := token_verify.GetClaimFromToken(token)
	assert.NotNil(t, c)
	assert.NoError(t, err)
}
