package token_verify

import (
	"Open_IM/pkg/common/constant"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseToken(t *testing.T) {
	token, _, err := CreateToken("", constant.IOSPlatformID)
	require.NoError(t, err)

	_, err = GetClaimFromToken(token)
	assert.NoError(t, err)
}
