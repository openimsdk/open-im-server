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
	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"github.com/OpenIMSDK/tools/discoveryregistry"

	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
)

type rpcInitFuc func(config *config2.GlobalConfig, disCov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error

type RpcCmd struct {
	*RootCmd
	RpcRegisterName string
	initFunc        rpcInitFuc
}

func NewRpcCmd(name string, initFunc rpcInitFuc) *RpcCmd {
	ret := &RpcCmd{RootCmd: NewRootCmd(name), initFunc: initFunc}
	ret.addPreRun()
	ret.addRunE()
	ret.SetRootCmdPt(ret)
	return ret
}

func (a *RpcCmd) addPreRun() {
	a.Command.PreRun = func(cmd *cobra.Command, args []string) {
		a.port = a.getPortFlag(cmd)
		a.prometheusPort = a.getPrometheusPortFlag(cmd)
	}
}

func (a *RpcCmd) addRunE() {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		rpcRegisterName, err := a.GetRpcRegisterNameFromConfig()
		if err != nil {
			return err
		} else {
			return a.StartSvr(rpcRegisterName, a.initFunc)
		}
	}
}

func (a *RpcCmd) Exec() error {
	return a.Execute()
}

func (a *RpcCmd) StartSvr(name string, rpcFn func(config *config2.GlobalConfig, disCov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) error {
	if a.GetPortFlag() == 0 {
		return errs.Wrap(errors.New("port is required"))
	}
	return startrpc.Start(a.GetPortFlag(), name, a.GetPrometheusPortFlag(), a.config, rpcFn)
}

func (a *RpcCmd) GetPortFromConfig(portType string) int {
	switch a.Name {
	case RpcPushServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImPushPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.PushPrometheusPort[0]
		}
	case RpcAuthServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImAuthPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.AuthPrometheusPort[0]
		}
	case RpcConversationServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImConversationPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.ConversationPrometheusPort[0]
		}
	case RpcFriendServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImFriendPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.FriendPrometheusPort[0]
		}
	case RpcGroupServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImGroupPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.GroupPrometheusPort[0]
		}
	case RpcMsgServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImMessagePort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.MessagePrometheusPort[0]
		}
	case RpcThirdServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImThirdPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.ThirdPrometheusPort[0]
		}
	case RpcUserServer:
		if portType == constant.FlagPort {
			return a.config.RpcPort.OpenImUserPort[0]
		}
		if portType == constant.FlagPrometheusPort {
			return a.config.Prometheus.UserPrometheusPort[0]
		}
	}
	return 0
}

func (a *RpcCmd) GetRpcRegisterNameFromConfig() (string, error) {
	switch a.Name {
	case RpcPushServer:
		return a.config.RpcRegisterName.OpenImPushName, nil
	case RpcAuthServer:
		return a.config.RpcRegisterName.OpenImAuthName, nil
	case RpcConversationServer:
		return a.config.RpcRegisterName.OpenImConversationName, nil
	case RpcFriendServer:
		return a.config.RpcRegisterName.OpenImFriendName, nil
	case RpcGroupServer:
		return a.config.RpcRegisterName.OpenImGroupName, nil
	case RpcMsgServer:
		return a.config.RpcRegisterName.OpenImMsgName, nil
	case RpcThirdServer:
		return a.config.RpcRegisterName.OpenImThirdName, nil
	case RpcUserServer:
		return a.config.RpcRegisterName.OpenImUserName, nil
	}
	return "", errs.Wrap(errors.New("can not get rpc register name"), a.Name)
}
