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

	"github.com/spf13/cobra"

	"github.com/openimsdk/open-im-server/v3/internal/msggateway"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/errs"

	v3config "github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

type MsgGatewayCmd struct {
	*RootCmd
}

func NewMsgGatewayCmd() *MsgGatewayCmd {
	ret := &MsgGatewayCmd{NewRootCmd("msgGateway")}
	ret.SetRootCmdPt(ret)
	return ret
}

func (m *MsgGatewayCmd) AddWsPortFlag() {
	m.Command.Flags().IntP(constant.FlagWsPort, "w", 0, "ws server listen port")
}

func (m *MsgGatewayCmd) getWsPortFlag(cmd *cobra.Command) (int, error) {
	port, err := cmd.Flags().GetInt(constant.FlagWsPort)
	if err != nil {
		return 0, errs.Wrap(err, "error getting ws port flag")
	}
	if port == 0 {
		port, _ = m.PortFromConfig(constant.FlagWsPort)
	}
	return port, nil
}

func (m *MsgGatewayCmd) addRunE() {
	m.Command.RunE = func(cmd *cobra.Command, args []string) error {
		wsPort, err := m.getWsPortFlag(cmd)
		if err != nil {
			return errs.Wrap(err, "failed to get WS port flag")
		}
		port, err := m.getPortFlag(cmd)
		if err != nil {
			return err
		}
		prometheusPort, err := m.getPrometheusPortFlag(cmd)
		if err != nil {
			return err
		}
		return msggateway.RunWsAndServer(port, wsPort, prometheusPort)
	}
}

func (m *MsgGatewayCmd) Exec() error {
	m.addRunE()
	return m.Execute()
}

func (m *MsgGatewayCmd) GetPortFromConfig(portType string) (int, error) {
	var port int
	var exists bool

	switch portType {
	case constant.FlagWsPort:
		if len(v3config.Config.LongConnSvr.OpenImWsPort) > 0 {
			port = v3config.Config.LongConnSvr.OpenImWsPort[0]
			exists = true
		}

	case constant.FlagPort:
		if len(v3config.Config.LongConnSvr.OpenImMessageGatewayPort) > 0 {
			port = v3config.Config.LongConnSvr.OpenImMessageGatewayPort[0]
			exists = true
		}

	case constant.FlagPrometheusPort:
		if len(v3config.Config.Prometheus.MessageGatewayPrometheusPort) > 0 {
			port = v3config.Config.Prometheus.MessageGatewayPrometheusPort[0]
			exists = true
		}
	}

	if !exists {
		return 0, errs.Wrap(errors.New("port type '%s' not found in configuration"), portType)
	}

	return port, nil
}
