package utils

import (
	"Open_IM/pkg/utils"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../..")
)

func Test_GenSmallImage(t *testing.T) {
	output := Root + "/out-test"
	defer func() {
		os.RemoveAll(output)
	}()

	require.NoError(t, os.Mkdir(output, 0755))

	for _, tt := range []struct {
		in  string
		out string
	}{
		{in: Root + "/docs/open-im-logo.png", out: output + "/open-im-logo-test.png"},
		{in: Root + "/docs/Architecture.jpg", out: output + "/Architecture-test.png"},
	} {
		t.Run(tt.in, func(t *testing.T) {
			err := utils.GenSmallImage(tt.in, tt.out)
			assert.NoError(t, err)
		})
	}
}
