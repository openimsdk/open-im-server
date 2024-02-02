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

package main

import (
	"fmt"
	"github.com/openimsdk/open-im-server/v3/internal/tools"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"os"
)

func main() {
	cronTaskCmd := cmd.NewCronTaskCmd()
	if err := cronTaskCmd.Exec(tools.StartTask); err != nil {
		fmt.Fprintf(os.Stderr, "\n\nexit -1: \n%+v\n\n", err)
		os.Exit(-1)
	}
}
