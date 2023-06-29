package utils

import (
	"Open_IM/pkg/utils"
	"net"
	"testing"
)

func TestServerIP(t *testing.T) {
	if net.ParseIP(utils.ServerIP) == nil {
		t.Fail()
	}
}
