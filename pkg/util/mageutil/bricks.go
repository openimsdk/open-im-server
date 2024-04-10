package mageutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// StopBinaries iterates over all binary files and terminates their corresponding processes.
func StopBinaries() {
	for binary := range serviceBinaries {
		fullPath := GetBinFullPath(binary)
		KillExistBinary(fullPath)
	}
}

// StartBinaries Start all binary services.
func StartBinaries() error {
	for binary, count := range serviceBinaries {
		binFullPath := filepath.Join(OpenIMOutputHostBin, binary)
		for i := 0; i < count; i++ {
			args := []string{"-i", strconv.Itoa(i), "-c", OpenIMOutputConfig}
			cmd := exec.Command(binFullPath, args...)
			fmt.Printf("Starting %s\n", cmd.String())
			cmd.Dir = OpenIMOutputHostBin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to start %s with args %v: %v\n", binFullPath, args, err)
				return err
			}
		}
	}
	return nil
}

// StartTools starts all tool binaries.
func StartTools() error {
	for _, tool := range toolBinaries {
		toolFullPath := GetToolFullPath(tool)
		cmd := exec.Command(toolFullPath, "-c", OpenIMOutputConfig)
		fmt.Printf("Starting %s\n", cmd.String())
		cmd.Dir = OpenIMOutputHostBinTools
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			fmt.Printf("Failed to start %s with error: %v\n", toolFullPath, err)
			return err
		}

		if err := cmd.Wait(); err != nil {
			fmt.Printf("Failed to execute %s with exit code: %v\n", toolFullPath, err)
			return err
		}
		fmt.Printf("Starting %s successfully \n", cmd.String())
	}

	return nil
}

// KillExistBinaries iterates over all binary files and kills their corresponding processes.
func KillExistBinaries() {
	for binary := range serviceBinaries {
		fullPath := GetBinFullPath(binary)
		KillExistBinary(fullPath)
	}
}

// CheckBinariesStop checks if all binary files have stopped and returns an error if there are any binaries still running.
func CheckBinariesStop() error {
	var runningBinaries []string

	for binary := range serviceBinaries {
		fullPath := GetBinFullPath(binary)
		if CheckProcessNamesExist(fullPath) {
			runningBinaries = append(runningBinaries, binary)
		}
	}

	if len(runningBinaries) > 0 {
		return fmt.Errorf("the following binaries are still running: %s", strings.Join(runningBinaries, ", "))
	}

	return nil
}

// CheckBinariesRunning checks if all binary files are running as expected and returns any errors encountered.
func CheckBinariesRunning() error {
	var errorMessages []string

	for binary, expectedCount := range serviceBinaries {
		fullPath := GetBinFullPath(binary)
		err := CheckProcessNames(fullPath, expectedCount)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("binary %s is not running as expected: %v", binary, err))
		}
	}

	if len(errorMessages) > 0 {
		return fmt.Errorf(strings.Join(errorMessages, "\n"))
	}

	return nil
}

// PrintListenedPortsByBinaries iterates over all binary files and prints the ports they are listening on.
func PrintListenedPortsByBinaries() {
	for binary, _ := range serviceBinaries {
		basePath := GetBinFullPath(binary)
		fullPath := basePath
		PrintBinaryPorts(fullPath)
	}
}
