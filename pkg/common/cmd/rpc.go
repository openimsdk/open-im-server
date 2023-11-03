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

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

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
	a.Command.Run = func(cmd *cobra.Command, args []string) {
		a.port = a.getPortFlag(cmd)
		a.prometheusPort = a.getPrometheusPortFlag(cmd)
	}
	return a.Execute()
}

func (a *RpcCmd) StartSvr(
	name string,
	rpcFn func(discov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error,
) error {
	if a.GetPortFlag() == 0 {
		return errors.New("port is required")
	}
	return startrpc.Start(a.GetPortFlag(), name, a.GetPrometheusPortFlag(), rpcFn)
}
func (a *RpcCmd) GetPortFromConfig(portType string) int {
	switch a.Name {
	case RpcPushServer:
		if portType == constant.FlagPort {
			return config2.Config.RpcPort.OpenImPushPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return config2.Config.Prometheus.PushPrometheusPort[0]
		}
	case RpcAuthServer:
		if portType == constant.FlagPort {
			return config2.Config.RpcPort.OpenImAuthPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return config2.Config.Prometheus.AuthPrometheusPort[0]
		}
	case RpcConversationServer:
		if portType == constant.FlagPort {
			return config2.Config.RpcPort.OpenImConversationPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return config2.Config.Prometheus.ConversationPrometheusPort[0]
		}
	case RpcFriendServer:
		if portType == constant.FlagPort {
			return config2.Config.RpcPort.OpenImFriendPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return config2.Config.Prometheus.FriendPrometheusPort[0]
		}
	case RpcGroupServer:
		if portType == constant.FlagPort {
			return config2.Config.RpcPort.OpenImGroupPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return config2.Config.Prometheus.GroupPrometheusPort[0]
		}
	case RpcMsgServer:
		if portType == constant.FlagPort {
			return config2.Config.RpcPort.OpenImMessagePort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return config2.Config.Prometheus.MessagePrometheusPort[0]
		}
	case RpcThirdServer:
		if portType == constant.FlagPort {
			return config2.Config.RpcPort.OpenImThirdPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return config2.Config.Prometheus.ThirdPrometheusPort[0]
		}
	case RpcUserServer:
		if portType == constant.FlagPort {
			return config2.Config.RpcPort.OpenImUserPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return config2.Config.Prometheus.UserPrometheusPort[0]
		}
	}
	return 0
}
