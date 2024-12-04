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

	"github.com/fatih/color"

	"github.com/openimsdk/tools/utils/timeutil"
)

func ExecuteCommand(cmdName string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing %s: %v, stderr: %s", cmdName, err, stderr.String())
	}
	return out.String(), nil
}

func printTime() string {
	formattedTime := timeutil.GetCurrentTimeFormatted()
	return fmt.Sprintf("Current Date & Time: %s", formattedTime)
}

func getGoVersion() string {
	version := runtime.Version()
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	return fmt.Sprintf("Go Version: %s\nOS: %s\nArchitecture: %s", version, goos, goarch)
}

func getDockerVersion() string {
	version, err := ExecuteCommand("docker", "--version")
	if err != nil {
		return "Docker is not installed. Please install it to get the version."
	}
	return version
}

func getKubernetesVersion() string {
	version, err := ExecuteCommand("kubectl", "version", "--client", "--short")
	if err != nil {
		return "Kubernetes is not installed. Please install it to get the version."
	}
	return version
}

func getGitVersion() string {
	version, err := ExecuteCommand("git", "branch", "--show-current")
	if err != nil {
		return "Git is not installed. Please install it to get the version."
	}
	return version
}

// // NOTE: You'll need to provide appropriate commands for OpenIM versions.
// func getOpenIMServerVersion() string {
// 	// Placeholder
// 	openimVersion := version.GetSingleVersion()
// 	return "OpenIM Server: " + openimVersion + "\n"
// }

// func getOpenIMClientVersion() (string, error) {
// 	openIMClientVersion, err := version.GetClientVersion()
// 	if err != nil {
// 		return "", err
// 	}
// 	return "OpenIM Client: " + openIMClientVersion.ClientVersion + "\n", nil
// }

func main() {
	// red := color.New(color.FgRed).SprintFunc()
	//	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	//	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println(blue("## Go Version"))
	fmt.Println(getGoVersion())
	fmt.Println(blue("## Branch Type"))
	fmt.Println(getGitVersion())
	fmt.Println(blue("## Docker Version"))
	fmt.Println(getDockerVersion())
	fmt.Println(blue("## Kubernetes Version"))
	fmt.Println(getKubernetesVersion())
	// fmt.Println(blue("## OpenIM Versions"))
	// fmt.Println(getOpenIMServerVersion())
	// clientVersion, err := getOpenIMClientVersion()
	// if err != nil {
	// 	fmt.Println(red("Error getting OpenIM Client Version: "), err)
	// } else {
	// 	fmt.Println(clientVersion)
	// }
}
