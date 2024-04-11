package main

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"os"
)

func main() {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get virtual memory info: %v\n", err)
		os.Exit(1)
	}

	freeMemoryGB := float64(vMem.Free) / float64(1024*1024*1024)

	if freeMemoryGB < 1.0 {
		fmt.Fprintf(os.Stderr, "System free memory is less than 4GB: %.2fGB\n", freeMemoryGB)
		os.Exit(1)
	} else {
		fmt.Printf("System free memory is sufficient: %.2fGB\n", freeMemoryGB)
	}
}
