package utils

import (
	"Open_IM/pkg/common/constant"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PlatformIDToName(t *testing.T) {
	assert.Equal(t, constant.PlatformIDToName(1), "IOS")
	assert.Equal(t, constant.PlatformIDToName(2), "Android")
	assert.Equal(t, constant.PlatformIDToName(3), "Windows")
	assert.Equal(t, constant.PlatformIDToName(4), "OSX")
	assert.Equal(t, constant.PlatformIDToName(5), "Web")
	assert.Equal(t, constant.PlatformIDToName(6), "MiniWeb")
	assert.Equal(t, constant.PlatformIDToName(7), "Linux")

	assert.Equal(t, constant.PlatformIDToName(0), "")
}

func Test_PlatformNameToID(t *testing.T) {
	assert.Equal(t, constant.PlatformNameToID("IOS"), 1)
	assert.Equal(t, constant.PlatformNameToID("Android"), 2)
	assert.Equal(t, constant.PlatformNameToID("Windows"), 3)
	assert.Equal(t, constant.PlatformNameToID("OSX"), 4)
	assert.Equal(t, constant.PlatformNameToID("Web"), 5)
	assert.Equal(t, constant.PlatformNameToID("MiniWeb"), 6)
	assert.Equal(t, constant.PlatformNameToID("Linux"), 7)

	assert.Equal(t, constant.PlatformNameToID("UnknownDevice"), 0)
	assert.Equal(t, constant.PlatformNameToID(""), 0)
}

func Test_PlatformNameToClass(t *testing.T) {
	assert.Equal(t, constant.PlatformNameToClass("IOS"), "Mobile")
	assert.Equal(t, constant.PlatformNameToClass("Android"), "Mobile")
	assert.Equal(t, constant.PlatformNameToClass("OSX"), "PC")
	assert.Equal(t, constant.PlatformNameToClass("Windows"), "PC")
	assert.Equal(t, constant.PlatformNameToClass("Web"), "Web")
	assert.Equal(t, constant.PlatformNameToClass("MiniWeb"), "Web")
	assert.Equal(t, constant.PlatformNameToClass("Linux"), "PC")

	assert.Equal(t, constant.PlatformNameToClass("UnknownDevice"), "")
	assert.Equal(t, constant.PlatformNameToClass(""), "")
}
