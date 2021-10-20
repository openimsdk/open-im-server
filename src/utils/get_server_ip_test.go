package utils

import (
	"testing"
)

func TestServerIP(t *testing.T) {
	if ServerIP == "" {
		t.Fail()
	}
}
