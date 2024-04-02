package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	var sysInfo syscall.Sysinfo_t

	err := syscall.Sysinfo(&sysInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get system info: %v\n", err)
		os.Exit(-1)
	}

	// Sysinfo returns memory in bytes, convert it to gigabytes
	freeMemoryGB := float64(sysInfo.Freeram) / float64(1024*1024*1024)

	if freeMemoryGB < 4.0 {
		fmt.Fprintf(os.Stderr, "System free memory is less than 4GB: %.2fGB\n", freeMemoryGB)
		os.Exit(-1)
	} else {
		fmt.Printf("System free memory is sufficient: %.2fGB\n", freeMemoryGB)
	}
}
