//go:build linux || darwin

package startrpc

import (
	"context"
	"github.com/openimsdk/tools/log"
	"net"
	"syscall"
)

func createListener() net.ListenConfig {
	lc := net.ListenConfig{
		Control: func(network, address string, conn syscall.RawConn) error {
			return conn.Control(func(fd uintptr) {
				err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				if err != nil {
					log.ZError(context.Background(), "Failed to set socket flag to SO_REUSEADDR", err)
					return
				}
			})
		},
	}

	return lc
}
