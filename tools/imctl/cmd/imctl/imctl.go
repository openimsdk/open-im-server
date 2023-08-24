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

// iamctl is the command line tool for iam platform.
package main

import (
	"os"

	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd"
)

func main() {
	command := cmd.NewDefaultIMCtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
