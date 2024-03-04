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
	"errors"
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

// AddApi configures the API command to run with specified ports for the API and Prometheus monitoring.
// It ensures error handling for port retrieval and only proceeds if both port numbers are successfully obtained.
func (a *ApiCmd) AddApi(f func(port int, promPort int) error) {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		port, err := a.getPortFlag(cmd)
		if err != nil {
			return err
		}

		promPort, err := a.getPrometheusPortFlag(cmd)
		if err != nil {
			return err
		}

		return f(port, promPort)
	}
}

func (a *ApiCmd) GetPortFromConfig(portType string) (int, error) {
	if portType == constant.FlagPort {
		if len(config2.Config.Api.OpenImApiPort) > 0 {
			return config2.Config.Api.OpenImApiPort[0], nil
		}
		return 0, errors.New("API port configuration is empty or missing")
	} else if portType == constant.FlagPrometheusPort {
		if len(config2.Config.Prometheus.ApiPrometheusPort) > 0 {
			return config2.Config.Prometheus.ApiPrometheusPort[0], nil
		}
		return 0, errors.New("Prometheus port configuration is empty or missing")
	}
	return 0, fmt.Errorf("unknown port type: %s", portType)
}
