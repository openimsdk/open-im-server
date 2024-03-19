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
	"os"
	"path/filepath"
	"strings"

	"github.com/openimsdk/open-im-server/tools/formitychecker/config"
)

type Issue struct {
	Type    string
	Path    string
	Message string
}

type Checker struct {
	Config  *config.Config
	Summary struct {
		CheckedDirectories int
		CheckedFiles       int
		Issues             []Issue
	}
}

func (c *Checker) Check() error {
	return filepath.Walk(c.Config.BaseConfig.SearchDirectory, c.checkPath)
}

func (c *Checker) checkPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	relativePath, err := filepath.Rel(c.Config.BaseConfig.SearchDirectory, path)
	if err != nil {
		return err
	}

	if relativePath == "." {
		return nil
	}

	if info.IsDir() {
		c.Summary.CheckedDirectories++
		if c.isIgnoredDirectory(relativePath) {
			c.Summary.Issues = append(c.Summary.Issues, Issue{
				Type:    "ignoredDirectory",
				Path:    path,
				Message: "此目录已被忽略",
			})
			return filepath.SkipDir
		}
		if !c.checkDirectoryName(relativePath) {
			c.Summary.Issues = append(c.Summary.Issues, Issue{
				Type:    "directoryNaming",
				Path:    path,
				Message: "目录名称不符合规范",
			})
		}
	} else {
		if c.isIgnoredFile(path) {
			return nil
		}
		c.Summary.CheckedFiles++
		if !c.checkFileName(relativePath) {
			c.Summary.Issues = append(c.Summary.Issues, Issue{
				Type:    "fileNaming",
				Path:    path,
				Message: "文件名称不符合规范",
			})
		}
	}

	return nil
}

func (c *Checker) isIgnoredDirectory(path string) bool {
	for _, ignoredDir := range c.Config.IgnoreDirectories {
		if strings.Contains(path, ignoredDir) {
			return true
		}
	}
	return false
}

func (c *Checker) isIgnoredFile(path string) bool {
	ext := filepath.Ext(path)
	for _, format := range c.Config.IgnoreFormats {
		if ext == format {
			return true
		}
	}
	return false
}

func (c *Checker) checkDirectoryName(path string) bool {
	dirName := filepath.Base(path)
	if c.Config.DirectoryNaming.MustBeLowercase && (dirName != strings.ToLower(dirName)) {
		return false
	}
	if !c.Config.DirectoryNaming.AllowHyphens && strings.Contains(dirName, "-") {
		return false
	}
	if !c.Config.DirectoryNaming.AllowUnderscores && strings.Contains(dirName, "_") {
		return false
	}
	return true
}

func (c *Checker) checkFileName(path string) bool {
	fileName := filepath.Base(path)
	ext := filepath.Ext(fileName)

	allowHyphens := c.Config.FileNaming.AllowHyphens
	allowUnderscores := c.Config.FileNaming.AllowUnderscores
	mustBeLowercase := c.Config.FileNaming.MustBeLowercase

	if specificNaming, ok := c.Config.FileTypeSpecificNaming[ext]; ok {
		allowHyphens = specificNaming.AllowHyphens
		allowUnderscores = specificNaming.AllowUnderscores
		mustBeLowercase = specificNaming.MustBeLowercase
	}

	if mustBeLowercase && (fileName != strings.ToLower(fileName)) {
		return false
	}
	if !allowHyphens && strings.Contains(fileName, "-") {
		return false
	}
	if !allowUnderscores && strings.Contains(fileName, "_") {
		return false
	}

	return true
}
