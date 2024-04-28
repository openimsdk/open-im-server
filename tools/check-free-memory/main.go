package main

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/mem"
)

func main() {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get virtual memory info: %v\n", err)
		os.Exit(1)
	}

	// Use the Available field to get the available memory
	availableMemoryGB := float64(vMem.Available) / float64(1024*1024*1024)

	if availableMemoryGB < 1.0 {
		fmt.Fprintf(os.Stderr, "System available memory is less than 1GB: %.2fGB\n", availableMemoryGB)
		os.Exit(1)
	} else {
		fmt.Printf("System available memory is sufficient: %.2fGB\n", availableMemoryGB)
	}
}
