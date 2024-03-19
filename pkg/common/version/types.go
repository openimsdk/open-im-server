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

// Info contains versioning information.
// TODO: Add []string of api versions supported? It's still unclear
// how we'll want to distribute that information.
type Info struct {
	Major        string `json:"major,omitempty"`
	Minor        string `json:"minor,omitempty"`
	GitVersion   string `json:"gitVersion"`
	GitTreeState string `json:"gitTreeState,omitempty"`
	GitCommit    string `json:"gitCommit,omitempty"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

type Output struct {
	OpenIMServerVersion Info                 `json:"OpenIMServerVersion,omitempty" yaml:"OpenIMServerVersion,omitempty"`
	OpenIMClientVersion *OpenIMClientVersion `json:"OpenIMClientVersion,omitempty" yaml:"OpenIMClientVersion,omitempty"`
}

type OpenIMClientVersion struct {
	ClientVersion string `json:"clientVersion,omitempty" yaml:"clientVersion,omitempty"` // sdk core version
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return info.GitVersion
}
