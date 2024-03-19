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
	"os"

	"github.com/openimsdk/open-im-server/tools/codescan/config"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BaseConfig struct {
		SearchDirectory string `yaml:"searchDirectory"`
		IgnoreCase      bool   `yaml:"ignoreCase"`
	} `yaml:"baseConfig"`
	DirectoryNaming struct {
		AllowHyphens     bool `yaml:"allowHyphens"`
		AllowUnderscores bool `yaml:"allowUnderscores"`
		MustBeLowercase  bool `yaml:"mustBeLowercase"`
	} `yaml:"directoryNaming"`
	FileNaming struct {
		AllowHyphens     bool `yaml:"allowHyphens"`
		AllowUnderscores bool `yaml:"allowUnderscores"`
		MustBeLowercase  bool `yaml:"mustBeLowercase"`
	} `yaml:"fileNaming"`
	IgnoreFormats          []string                          `yaml:"ignoreFormats"`
	IgnoreDirectories      []string                          `yaml:"ignoreDirectories"`
	FileTypeSpecificNaming map[string]FileTypeSpecificNaming `yaml:"fileTypeSpecificNaming"`
}

type FileTypeSpecificNaming struct {
	AllowHyphens     bool `yaml:"allowHyphens"`
	AllowUnderscores bool `yaml:"allowUnderscores"`
	MustBeLowercase  bool `yaml:"mustBeLowercase"`
}

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
	Errors []string
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
