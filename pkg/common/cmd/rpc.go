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
	"google.golang.org/grpc"

	"github.com/OpenIMSDK/tools/errs"

	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"github.com/OpenIMSDK/tools/discoveryregistry"

	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
)

type RpcCmd struct {
	*RootCmd
}

func NewRpcCmd(name string) *RpcCmd {
	ret := &RpcCmd{NewRootCmd(name)}
	ret.SetRootCmdPt(ret)
	return ret
}

func (a *RpcCmd) Exec() error {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		portFlag, err := a.getPortFlag(cmd)
		if err != nil {
			return err
		}
		a.port = portFlag

		prometheusPort, err := a.getPrometheusPortFlag(cmd)
		if err != nil {
			return err
		}
		a.prometheusPort = prometheusPort

		return nil
	}
	return a.Execute()
}

func (a *RpcCmd) StartSvr(name string, rpcFn func(discov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) error {
	portFlag, err := a.GetPortFlag()
	if err != nil {
		return err
	} else {
		a.port = portFlag
	}

	return startrpc.Start(portFlag, name, a.GetPrometheusPortFlag(), rpcFn)
}

func (a *RpcCmd) GetPortFromConfig(portType string) (int, error) {
	portConfigMap := map[string]map[string]int{
		RpcPushServer: {
			constant.FlagPort:           config2.Config.RpcPort.OpenImPushPort[0],
			constant.FlagPrometheusPort: config2.Config.Prometheus.PushPrometheusPort[0],
		},
		RpcAuthServer: {
			constant.FlagPort:           config2.Config.RpcPort.OpenImAuthPort[0],
			constant.FlagPrometheusPort: config2.Config.Prometheus.AuthPrometheusPort[0],
		},
		RpcConversationServer: {
			constant.FlagPort:           config2.Config.RpcPort.OpenImConversationPort[0],
			constant.FlagPrometheusPort: config2.Config.Prometheus.ConversationPrometheusPort[0],
		},
		RpcFriendServer: {
			constant.FlagPort:           config2.Config.RpcPort.OpenImFriendPort[0],
			constant.FlagPrometheusPort: config2.Config.Prometheus.FriendPrometheusPort[0],
		},
		RpcGroupServer: {
			constant.FlagPort:           config2.Config.RpcPort.OpenImGroupPort[0],
			constant.FlagPrometheusPort: config2.Config.Prometheus.GroupPrometheusPort[0],
		},
		RpcMsgServer: {
			constant.FlagPort:           config2.Config.RpcPort.OpenImMessagePort[0],
			constant.FlagPrometheusPort: config2.Config.Prometheus.MessagePrometheusPort[0],
		},
		RpcThirdServer: {
			constant.FlagPort:           config2.Config.RpcPort.OpenImThirdPort[0],
			constant.FlagPrometheusPort: config2.Config.Prometheus.ThirdPrometheusPort[0],
		},
		RpcUserServer: {
			constant.FlagPort:           config2.Config.RpcPort.OpenImUserPort[0],
			constant.FlagPrometheusPort: config2.Config.Prometheus.UserPrometheusPort[0],
		},
	}

	if portMap, ok := portConfigMap[a.Name]; ok {
		if port, ok := portMap[portType]; ok {
			return port, nil
		} else {
			return 0, errs.Wrap(errors.New("port type not found"), fmt.Sprintf("Failed to get port for %s", a.Name))
		}
	}

	return 0, errs.Wrap(fmt.Errorf("server name '%s' not found", a.Name), "Failed to get port configuration")
}
