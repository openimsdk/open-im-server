// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

func executeCommand(cmdName string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error executing %s: %v", cmdName, err)
	}
	return out.String(), nil
}

func printTime() string {
	currentTime := time.Now()

	// 使用 Format 函数以优雅的方式格式化日期和时间
	// 2006-01-02 15:04:05 是 Go 中的标准时间格式
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("Current Date & Time:", formattedTime)
}

func getGoVersion() string {
	version := runtime.Version()
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	return fmt.Sprintf("Go Version: %s\nOS: %s\nArchitecture: %s", version, goos, goarch)
}

func getDockerVersion() string {
	version, err := executeCommand("docker", "--version")
	if err != nil {
		return "Docker is not installed. Please install it to get the version."
	}
	return version
}

func getDockerComposeVersion() string {
	version, err := executeCommand("docker-compose", "--version")
	if err != nil {
		return "Docker Compose is not installed. Please install it to get the version."
	}
	return version
}

func getKubernetesVersion() string {
	version, err := executeCommand("kubectl", "version", "--client", "--short")
	if err != nil {
		return "Kubernetes is not installed. Please install it to get the version."
	}
	return version
}

func getGitVersion() string {
	version, err := executeCommand("git", "branch", "--show-current")
	if err != nil {
		return "Git is not installed. Please install it to get the version."
	}
	return version
}

// NOTE: You'll need to provide appropriate commands for OpenIM versions.
func getOpenIMServerVersion() string {
	// Placeholder
	return "OpenIM Server: v3.2"
}

func getOpenIMClientVersion() string {
	// Placeholder
	return "OpenIM Client: v3.2"
}

func main() {
	fmt.Println(printTime())
	fmt.Println("# Diagnostic Tool Result\n")
	fmt.Println("## Go Version")
	fmt.Println(getGoVersion())
	fmt.Println("## Branch Type")
	fmt.Println(getGitVersion())
	fmt.Println("## Docker Version")
	fmt.Println(getDockerVersion())
	fmt.Println("## Docker Compose Version")
	fmt.Println(getDockerComposeVersion())
	fmt.Println("## Kubernetes Version")
	fmt.Println(getKubernetesVersion())
	fmt.Println("## OpenIM Versions")
	fmt.Println(getOpenIMServerVersion())
	fmt.Println(getOpenIMClientVersion())
}
