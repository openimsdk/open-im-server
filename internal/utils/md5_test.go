package utils

import (
	"Open_IM/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Md5(t *testing.T) {
	result := utils.Md5("go")
	assert.Equal(t, result, "34d1f91fb2e514b8576fab1a75a89a6b")

	result2 := utils.Md5("go")
	assert.Equal(t, result, result2)
}
