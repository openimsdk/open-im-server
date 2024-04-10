package mageutil

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	OpenIMRoot               string
	OpenIMOutputConfig       string
	OpenIMOutput             string
	OpenIMOutputTools        string
	OpenIMOutputTmp          string
	OpenIMOutputLogs         string
	OpenIMOutputBin          string
	OpenIMOutputBinPath      string
	OpenIMOutputBinToolPath  string
	OpenIMInitErrLogFile     string
	OpenIMInitLogFile        string
	OpenIMOutputHostBin      string
	OpenIMOutputHostBinTools string
)

func init() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic("Error getting current directory: " + err.Error())
	}

	OpenIMRoot = currentDir

	OpenIMOutputConfig = filepath.Join(OpenIMRoot, "config") + string(filepath.Separator)
	OpenIMOutput = filepath.Join(OpenIMRoot, "_output") + string(filepath.Separator)

	OpenIMOutputTools = filepath.Join(OpenIMOutput, "tools") + string(filepath.Separator)
	OpenIMOutputTmp = filepath.Join(OpenIMOutput, "tmp") + string(filepath.Separator)
	OpenIMOutputLogs = filepath.Join(OpenIMOutput, "logs") + string(filepath.Separator)
	OpenIMOutputBin = filepath.Join(OpenIMOutput, "bin") + string(filepath.Separator)

	OpenIMOutputBinPath = filepath.Join(OpenIMOutputBin, "platforms") + string(filepath.Separator)
	OpenIMOutputBinToolPath = filepath.Join(OpenIMOutputBin, "tools") + string(filepath.Separator)

	OpenIMInitErrLogFile = filepath.Join(OpenIMOutputLogs, "openim-init-err.log")
	OpenIMInitLogFile = filepath.Join(OpenIMOutputLogs, "openim-init.log")

	OpenIMOutputHostBin = filepath.Join(OpenIMOutputBinPath, OsArch()) + string(filepath.Separator)
	OpenIMOutputHostBinTools = filepath.Join(OpenIMOutputBinToolPath, OsArch()) + string(filepath.Separator)

	dirs := []string{
		OpenIMOutputConfig,
		OpenIMOutput,
		OpenIMOutputTools,
		OpenIMOutputTmp,
		OpenIMOutputLogs,
		OpenIMOutputBin,
		OpenIMOutputBinPath,
		OpenIMOutputBinToolPath,
		OpenIMOutputHostBin,
		OpenIMOutputHostBinTools,
	}

	for _, dir := range dirs {
		createDirIfNotExist(dir)
	}
}

func createDirIfNotExist(dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", dir, err)
		os.Exit(1)
	}
}

// GetBinFullPath constructs and returns the full path for the given binary name.
func GetBinFullPath(binName string) string {
	binFullPath := filepath.Join(OpenIMOutputHostBin, binName)
	return binFullPath
}

// GetToolFullPath constructs and returns the full path for the given tool name.
func GetToolFullPath(toolName string) string {
	toolFullPath := filepath.Join(OpenIMOutputHostBinTools, toolName)
	return toolFullPath
}
