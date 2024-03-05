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
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Directory string   `yaml:"directory"`
	FileTypes []string `yaml:"file_types"`
	Languages []string `yaml:"languages"`
}

func ParseConfig() (Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "./", "Path to config file")
	flag.Parse()

	var config Config
	if configPath != "" {
		configFile, err := os.ReadFile(configPath)
		if err != nil {
			return Config{}, err
		}
		if err := yaml.Unmarshal(configFile, &config); err != nil {
			return Config{}, err
		}
	} else {
		log.Fatal("Config file must be provided")
	}
	return config, nil
}
