package utils

import (
	"net"
	"testing"
)

func TestServerIP(t *testing.T) {
	if net.ParseIP(ServerIP) == nil {
		t.Fail()
	}
}
