package main

import (
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"os"
)

func main() {
	// 使用gopsutil获取虚拟内存信息
	vMem, err := mem.VirtualMemory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get virtual memory info: %v\n", err)
		os.Exit(1)
	}

	// 将可用内存从字节转换为GB
	freeMemoryGB := float64(vMem.Free) / float64(1024*1024*1024)

	if freeMemoryGB < 4.0 {
		fmt.Fprintf(os.Stderr, "System free memory is less than 4GB: %.2fGB\n", freeMemoryGB)
		os.Exit(1)
	} else {
		fmt.Printf("System free memory is sufficient: %.2fGB\n", freeMemoryGB)
	}
}
