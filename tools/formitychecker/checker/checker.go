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

package checker

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/openimsdk/open-im-server/tools/formitychecker/config"
)

var (
	underscoreRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+\.[a-zA-Z0-9]+$`)
	hyphenRegex     = regexp.MustCompile(`^[a-zA-Z0-9\-]+\.[a-zA-Z0-9]+$`)
)

// CheckDirectory initiates the checking process for the specified directories using configuration from config.Config.
func CheckDirectory(cfg *config.Config) error {
	ignoreMap := make(map[string]struct{})
	for _, dir := range cfg.IgnoreDirs {
		ignoreMap[dir] = struct{}{}
	}

	for _, targetDir := range cfg.TargetDirs {
		err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errs.Wrap(err, fmt.Sprintf("error walking directory '%s'", targetDir))
			}

			// Skip if the directory is in the ignore list
			dirName := filepath.Base(filepath.Dir(path))
			if _, ok := ignoreMap[dirName]; ok && info.IsDir() {
				return filepath.SkipDir
			}

			// Check the naming convention
			if err := checkNamingConvention(path, info); err != nil {
				fmt.Println(err)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("error checking directory '%s': %w", targetDir, err)
		}
	}

	return nil
}

// checkNamingConvention checks if the file or directory name conforms to the standard naming conventions.
func checkNamingConvention(path string, info os.FileInfo) error {
	fileName := info.Name()

	// Handle special cases for directories like .git
	if info.IsDir() && strings.HasPrefix(fileName, ".") {
		return nil // Skip special directories
	}

	// Extract the main part of the name (without extension for files)
	mainName := fileName
	if !info.IsDir() {
		mainName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	}

	// Determine the type of file and apply corresponding naming rule
	switch {
	case info.IsDir():
		if !isValidName(mainName, "_") { // Directory names must only contain underscores
			return fmt.Errorf("!!! invalid directory name: %s", path)
		}
	case strings.HasSuffix(fileName, ".go"):
		if !isValidName(mainName, "_") { // Go files must only contain underscores
			return fmt.Errorf("!!! invalid Go file name: %s", path)
		}
	case strings.HasSuffix(fileName, ".yml"), strings.HasSuffix(fileName, ".yaml"), strings.HasSuffix(fileName, ".md"):
		if !isValidName(mainName, "-") { // YML, YAML, and Markdown files must only contain hyphens
			return fmt.Errorf("!!! invalid file name: %s", path)
		}
	}

	return nil
}

// isValidName checks if the file name conforms to the specified rule (underscore or hyphen).
func isValidName(name, charType string) bool {
	switch charType {
	case "_":
		return underscoreRegex.MatchString(name)
	case "-":
		return hyphenRegex.MatchString(name)
	default:
		return false
	}
}
