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

package version

import (
	"fmt"
	"runtime"

	"gopkg.in/src-d/go-git.v4"
)

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in ./base.go
	return Info{
		Major:        gitMajor,
		Minor:        gitMinor,
		GitVersion:   gitVersion,
		GitTreeState: gitTreeState,
		GitCommit:    gitCommit,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetClientVersion returns the git version of the OpenIM client repository.
func GetClientVersion() (*OpenIMClientVersion, error) {
	clientVersion, err := getClientVersion()
	if err != nil {
		return nil, err
	}
	return &OpenIMClientVersion{
		ClientVersion: clientVersion,
	}, nil
}

func getClientVersion() (string, error) {
	repo, err := git.PlainClone("/tmp/openim-sdk-core", false, &git.CloneOptions{
		URL: "https://github.com/OpenIMSDK/openim-sdk-core",
	})
	if err != nil {
		return "", fmt.Errorf("error cloning repository: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("error getting head reference: %w", err)
	}

	return ref.Hash().String(), nil
}

// GetSingleVersion returns single version of sealer.
func GetSingleVersion() string {
	return gitVersion
}
