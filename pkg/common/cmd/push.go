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
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/push"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type PushRpcCmd struct {
	*RootCmd
	ctx        context.Context
	configMap  map[string]any
	pushConfig *push.Config
}

func NewPushRpcCmd() *PushRpcCmd {
	var pushConfig push.Config
	ret := &PushRpcCmd{pushConfig: &pushConfig}
	ret.configMap = map[string]any{
		config.OpenIMPushCfgFileName:    &pushConfig.RpcConfig,
		config.RedisConfigFileName:      &pushConfig.RedisConfig,
		config.KafkaConfigFileName:      &pushConfig.KafkaConfig,
		config.ShareFileName:            &pushConfig.Share,
		config.NotificationFileName:     &pushConfig.NotificationConfig,
		config.WebhooksConfigFileName:   &pushConfig.WebhooksConfig,
		config.LocalCacheConfigFileName: &pushConfig.LocalCacheConfig,
		config.DiscoveryConfigFilename:  &pushConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		ret.pushConfig.FcmConfigPath = ret.ConfigPath()
		return ret.runE()
	}
	return ret
}

func (a *PushRpcCmd) Exec() error {
	return a.Execute()
}

func (a *PushRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.pushConfig.Discovery, &a.pushConfig.RpcConfig.Prometheus, a.pushConfig.RpcConfig.RPC.ListenIP,
		a.pushConfig.RpcConfig.RPC.RegisterIP, a.pushConfig.RpcConfig.RPC.AutoSetPorts, a.pushConfig.RpcConfig.RPC.Ports,
		a.Index(), a.pushConfig.Discovery.RpcService.Push, &a.pushConfig.NotificationConfig, a.pushConfig,
		[]string{
			a.pushConfig.RpcConfig.GetConfigFileName(),
			a.pushConfig.RedisConfig.GetConfigFileName(),
			a.pushConfig.KafkaConfig.GetConfigFileName(),
			a.pushConfig.NotificationConfig.GetConfigFileName(),
			a.pushConfig.Share.GetConfigFileName(),
			a.pushConfig.WebhooksConfig.GetConfigFileName(),
			a.pushConfig.LocalCacheConfig.GetConfigFileName(),
			a.pushConfig.Discovery.GetConfigFileName(),
		},
		[]string{
			a.pushConfig.Discovery.RpcService.MessageGateway,
		},
		push.Start)
}
