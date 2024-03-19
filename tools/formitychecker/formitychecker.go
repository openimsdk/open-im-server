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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/openimsdk/open-im-server/tools/formitychecker/checker"
	"github.com/openimsdk/open-im-server/tools/formitychecker/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	if configPath == "" {
		configPath = "config.yaml"
		if _, err := os.Stat(".github/formitychecker.yaml"); err == nil {
			configPath = ".github/formitychecker.yaml"
		}
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	c := &checker.Checker{Config: cfg}
	err = c.Check()
	if err != nil {
		fmt.Println("Error during check:", err)
		os.Exit(1)
	}

	// if len(c.Errors) > 0 {
	// 	fmt.Println("Found errors:")
	// 	for _, errMsg := range c.Errors {
	// 		fmt.Println("-", errMsg)
	// 	}
	// 	os.Exit(1)
	// }

	summaryJSON, err := json.MarshalIndent(c.Summary, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling summary:", err)
		return
	}

	fmt.Println(string(summaryJSON))
}
