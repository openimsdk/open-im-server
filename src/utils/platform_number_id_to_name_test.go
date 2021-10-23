package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PlatformIDToName(t *testing.T) {
	assert.Equal(t, PlatformIDToName(1), "IOS")
	assert.Equal(t, PlatformIDToName(2), "Android")
	assert.Equal(t, PlatformIDToName(3), "Windows")
	assert.Equal(t, PlatformIDToName(4), "OSX")
	assert.Equal(t, PlatformIDToName(5), "Web")
	assert.Equal(t, PlatformIDToName(6), "MiniWeb")
	assert.Equal(t, PlatformIDToName(7), "Linux")

	assert.Equal(t, PlatformIDToName(0), "")
}

func Test_PlatformNameToID(t *testing.T) {
	assert.Equal(t, PlatformNameToID("IOS"), int32(1))
	assert.Equal(t, PlatformNameToID("Android"), int32(2))
	assert.Equal(t, PlatformNameToID("Windows"), int32(3))
	assert.Equal(t, PlatformNameToID("OSX"), int32(4))
	assert.Equal(t, PlatformNameToID("Web"), int32(5))
	assert.Equal(t, PlatformNameToID("MiniWeb"), int32(6))
	assert.Equal(t, PlatformNameToID("Linux"), int32(7))

	assert.Equal(t, PlatformNameToID("UnknownDevice"), int32(0))
	assert.Equal(t, PlatformNameToID(""), int32(0))
}

func Test_PlatformNameToClass(t *testing.T) {
	assert.Equal(t, PlatformNameToClass("IOS"), "Mobile")
	assert.Equal(t, PlatformNameToClass("Android"), "Mobile")
	assert.Equal(t, PlatformNameToClass("OSX"), "PC")
	assert.Equal(t, PlatformNameToClass("Windows"), "PC")
	assert.Equal(t, PlatformNameToClass("Web"), "PC")
	assert.Equal(t, PlatformNameToClass("MiniWeb"), "Mobile")
	assert.Equal(t, PlatformNameToClass("Linux"), "PC")

	assert.Equal(t, PlatformNameToClass("UnknownDevice"), "")
	assert.Equal(t, PlatformNameToClass(""), "")
}
