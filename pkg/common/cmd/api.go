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
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/openimsdk/open-im-server/v3/internal/api"
	"github.com/spf13/cobra"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

type ApiCmd struct {
	*RootCmd
	initFunc func(config *config.GlobalConfig, port int, promPort int) error
}

func NewApiCmd() *ApiCmd {
	ret := &ApiCmd{RootCmd: NewRootCmd("api"), initFunc: api.Start}
	ret.SetRootCmdPt(ret)
	ret.addPreRun()
	ret.addRunE()
	return ret
}

func (a *ApiCmd) addPreRun() {
	a.Command.PreRun = func(cmd *cobra.Command, args []string) {
		a.port = a.getPortFlag(cmd)
		a.prometheusPort = a.getPrometheusPortFlag(cmd)
	}
}

func (a *ApiCmd) addRunE() {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return a.initFunc(a.config, a.port, a.prometheusPort)
	}
}

func (a *ApiCmd) GetPortFromConfig(portType string) int {
	if portType == constant.FlagPort {
		return a.config.Api.OpenImApiPort[0]
	} else if portType == constant.FlagPrometheusPort {
		return a.config.Prometheus.ApiPrometheusPort[0]
	}
	return 0
}
