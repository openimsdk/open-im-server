//go:build mage && !windows
// +build mage,!windows

package main

import (
	"github.com/openimsdk/gomake/mageutil"
	"syscall"
)

func setMaxOpenFiles() error {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err
	}
	rLimit.Max = uint64(mageutil.MaxFileDescriptors)
	rLimit.Cur = uint64(mageutil.MaxFileDescriptors)
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}
