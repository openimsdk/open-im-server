package mageutil

import (
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func OsArch() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	if os == "windows" {
		return fmt.Sprintf("%s\\%s", os, arch)
	}
	return fmt.Sprintf("%s/%s", os, arch)
}

func CheckProcessNames(processPath string, expectedCount int) error {
	processes, err := process.Processes()
	if err != nil {
		return fmt.Errorf("failed to get processes: %v", err)
	}

	runningCount := 0
	for _, p := range processes {
		exePath, err := p.Exe()
		if err != nil {
			continue
		}

		if strings.EqualFold(exePath, processPath) {
			runningCount++

		}
	}

	if runningCount == expectedCount {
		return nil
	} else {
		return fmt.Errorf("%s Expected %d processes, but %d running", processPath, expectedCount, runningCount)
	}
}

// CheckProcessNamesExist checks if there are any processes running that match the specified path.
func CheckProcessNamesExist(processPath string) bool {
	processes, err := process.Processes()
	if err != nil {
		fmt.Printf("Failed to get processes: %v\n", err)
		return false
	}

	for _, p := range processes {
		exePath, err := p.Exe()
		if err != nil {
			continue
		}

		if exePath == processPath {
			return true // 找到至少一个匹配的进程
		}
	}

	return false
}

// PrintBinaryPorts prints the ports listened by the process of a specified binary file along with its command line arguments.
func PrintBinaryPorts(binaryPath string) {
	pids, err := FindPIDsByBinaryPath(binaryPath)
	if err != nil {
		fmt.Println("Error finding PIDs:", err)
		return
	}

	if len(pids) == 0 {
		fmt.Printf("No running processes found for binary: %s\n", binaryPath)
		return
	}

	for _, pid := range pids {

		proc, err := process.NewProcess(int32(pid))
		if err != nil {
			fmt.Printf("Failed to create process object for PID %d: %v\n", pid, err)
			continue
		}

		cmdline, err := proc.Cmdline()
		if err != nil {
			fmt.Printf("Failed to get command line for PID %d: %v\n", pid, err)
			continue
		}

		connections, err := net.ConnectionsPid("all", int32(pid))
		if err != nil {
			fmt.Printf("Error getting connections for PID %d: %v\n", pid, err)
			continue
		}

		portsMap := make(map[string]struct{})
		for _, conn := range connections {
			if conn.Status == "LISTEN" {
				port := fmt.Sprintf("%d", conn.Laddr.Port)
				portsMap[port] = struct{}{}
			}
		}

		if len(portsMap) == 0 {
			PrintGreen(fmt.Sprintf("Cmdline: %s, PID: %d is not listening on any ports.", cmdline, pid))
		} else {
			ports := make([]string, 0, len(portsMap))
			for port := range portsMap {
				ports = append(ports, port)
			}
			PrintGreen(fmt.Sprintf("Cmdline: %s, PID: %d is listening on ports: %s", cmdline, pid, strings.Join(ports, ", ")))
		}
	}
}

// FindPIDsByBinaryPath finds all matching process PIDs by binary path.
func FindPIDsByBinaryPath(binaryPath string) ([]int, error) {
	var pids []int
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	for _, proc := range processes {
		exePath, err := proc.Exe()
		if err != nil {
			continue
		}

		if strings.EqualFold(exePath, binaryPath) {
			pids = append(pids, int(proc.Pid))
		}
	}

	return pids, nil
}

// KillExistBinary kills all processes matching the given binary file path.
func KillExistBinary(binaryPath string) {
	processes, err := process.Processes()
	if err != nil {
		fmt.Printf("Failed to get processes: %v\n", err)
		return
	}

	for _, p := range processes {
		exePath, err := p.Exe()
		if err != nil {
			continue
		}

		if strings.Contains(strings.ToLower(exePath), strings.ToLower(binaryPath)) {

			//if strings.EqualFold(exePath, binaryPath) {
			cmdline, err := p.Cmdline()
			if err != nil {
				fmt.Printf("Failed to get command line for process %d: %v\n", p.Pid, err)
				continue
			}

			err = p.Terminate()
			if err != nil {

				err = p.Kill()
				if err != nil {
					fmt.Printf("Failed to kill process cmdline: %s, pid: %d, err: %v\n", cmdline, p.Pid, err)
				} else {
					fmt.Printf("Killed process cmdline: %s, pid: %d\n", cmdline, p.Pid)
				}
			} else {
				fmt.Printf("Terminated process cmdline: %s, pid: %d\n", cmdline, p.Pid)
			}
		}
	}
}

// DetectPlatform detects the operating system and architecture.
func DetectPlatform() string {
	targetOS, targetArch := runtime.GOOS, runtime.GOARCH
	switch targetArch {
	case "amd64", "arm64":
	default:
		fmt.Printf("Unsupported architecture: %s\n", targetArch)
		os.Exit(1)
	}
	return fmt.Sprintf("%s_%s", targetOS, targetArch)
}

// rootDir gets the absolute path of the current directory.
func rootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
		os.Exit(1)
	}
	return dir
}

var rootDirPath = rootDir()
var platformsOutputBase = filepath.Join(rootDirPath, "_output/bin/platforms")
var toolsOutputBase = filepath.Join(rootDirPath, "_output/bin/tools")
