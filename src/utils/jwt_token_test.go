package utils

import (
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
