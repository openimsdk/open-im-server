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

import (
	"fmt"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/spf13/cobra"

	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

type ApiCmd struct {
	*RootCmd
}

func NewApiCmd() *ApiCmd {
	ret := &ApiCmd{NewRootCmd("api")}
	ret.SetRootCmdPt(ret)

	return ret
}

func (a *ApiCmd) AddApi(f func(port int, promPort int) error) {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f(a.getPortFlag(cmd), a.getPrometheusPortFlag(cmd))
	}
}

func (a *ApiCmd) GetPortFromConfig(portType string) int {
	fmt.Println("GetPortFromConfig:", portType)
	if portType == constant.FlagPort {
		return config2.Config.Api.OpenImApiPort[0]
	} else if portType == constant.FlagPrometheusPort {
		return config2.Config.Prometheus.ApiPrometheusPort[0]
	}
	return 0
}
