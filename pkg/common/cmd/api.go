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

package cmd

import "github.com/spf13/cobra"

type ApiCmd struct {
	*RootCmd
}

func NewApiCmd() *ApiCmd {
	return &ApiCmd{NewRootCmd("api")}
}

func (a *ApiCmd) AddApi(f func(port int) error) {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f(a.getPortFlag(cmd))
	}
}
