package utils

import (
	"Open_IM/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PlatformIDToName(t *testing.T) {
	assert.Equal(t, utils.PlatformIDToName(1), "IOS")
	assert.Equal(t, utils.PlatformIDToName(2), "Android")
	assert.Equal(t, utils.PlatformIDToName(3), "Windows")
	assert.Equal(t, utils.PlatformIDToName(4), "OSX")
	assert.Equal(t, utils.PlatformIDToName(5), "Web")
	assert.Equal(t, utils.PlatformIDToName(6), "MiniWeb")
	assert.Equal(t, utils.PlatformIDToName(7), "Linux")

	assert.Equal(t, utils.PlatformIDToName(0), "")
}

func Test_PlatformNameToID(t *testing.T) {
	assert.Equal(t, utils.PlatformNameToID("IOS"), int32(1))
	assert.Equal(t, utils.PlatformNameToID("Android"), int32(2))
	assert.Equal(t, utils.PlatformNameToID("Windows"), int32(3))
	assert.Equal(t, utils.PlatformNameToID("OSX"), int32(4))
	assert.Equal(t, utils.PlatformNameToID("Web"), int32(5))
	assert.Equal(t, utils.PlatformNameToID("MiniWeb"), int32(6))
	assert.Equal(t, utils.PlatformNameToID("Linux"), int32(7))

	assert.Equal(t, utils.PlatformNameToID("UnknownDevice"), int32(0))
	assert.Equal(t, utils.PlatformNameToID(""), int32(0))
}

func Test_PlatformNameToClass(t *testing.T) {
	assert.Equal(t, utils.PlatformNameToClass("IOS"), "Mobile")
	assert.Equal(t, utils.PlatformNameToClass("Android"), "Mobile")
	assert.Equal(t, utils.PlatformNameToClass("OSX"), "PC")
	assert.Equal(t, utils.PlatformNameToClass("Windows"), "PC")
	assert.Equal(t, utils.PlatformNameToClass("Web"), "PC")
	assert.Equal(t, utils.PlatformNameToClass("MiniWeb"), "Mobile")
	assert.Equal(t, utils.PlatformNameToClass("Linux"), "PC")

	assert.Equal(t, utils.PlatformNameToClass("UnknownDevice"), "")
	assert.Equal(t, utils.PlatformNameToClass(""), "")
}
