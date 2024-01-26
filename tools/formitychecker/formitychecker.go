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
	"flag"
	"fmt"

	"github.com/openimsdk/open-im-server/tools/formitychecker/checker"
	"github.com/openimsdk/open-im-server/tools/formitychecker/config"
)

func main() {
	defaultTargetDirs := "."
	defaultIgnoreDirs := "components,.git"

	var targetDirs string
	var ignoreDirs string
	flag.StringVar(&targetDirs, "target", defaultTargetDirs, "Directories to check (default: current directory)")
	flag.StringVar(&ignoreDirs, "ignore", defaultIgnoreDirs, "Directories to ignore (default: A/, B/)")
	flag.Parse()

	conf := config.NewConfig(targetDirs, ignoreDirs)

	err := checker.CheckDirectory(conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
