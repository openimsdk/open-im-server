package utils

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../..")
)

func Test_GenSmallImage(t *testing.T) {
	println(Root)
	err := GenSmallImage(Root+"/docs/open-im-logo.png", Root+"/out-test/open-im-logo-test.png")
	assert.Nil(t, err)

	err = GenSmallImage(Root+"/docs/open-im-logo.png", "out-test/open-im-logo-test.png")
	assert.Nil(t, err)

	err = GenSmallImage(Root+"/docs/Architecture.jpg", "out-test/Architecture-test.jpg")
	assert.Nil(t, err)
}
