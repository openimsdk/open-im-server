// Copyright © 2024 OpenIM. All rights reserved.
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

package config

import (
	"strings"
)

// Config holds all the configuration parameters for the checker.
type Config struct {
	TargetDirs []string // Directories to check
	IgnoreDirs []string // Directories to ignore
}

// NewConfig creates and returns a new Config instance.
func NewConfig(targetDirs, ignoreDirs string) *Config {
	return &Config{
		TargetDirs: parseDirs(targetDirs),
		IgnoreDirs: parseDirs(ignoreDirs),
	}
}

// parseDirs splits a comma-separated string into a slice of directory names.
func parseDirs(dirs string) []string {
	if dirs == "" {
		return nil
	}
	return strings.Split(dirs, ",")
}
