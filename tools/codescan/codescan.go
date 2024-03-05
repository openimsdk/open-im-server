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
	"log"

	"github.com/openimsdk/open-im-server/tools/codescan/checker"
	"github.com/openimsdk/open-im-server/tools/codescan/config"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	err = checker.WalkDirAndCheckComments(cfg)
	if err != nil {
		panic(err)
	}
}
