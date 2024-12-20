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

	"github.com/openimsdk/open-im-server/v3/internal/rpc/third"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type ThirdRpcCmd struct {
	*RootCmd
	ctx         context.Context
	configMap   map[string]any
	thirdConfig *third.Config
}

func NewThirdRpcCmd() *ThirdRpcCmd {
	var thirdConfig third.Config
	ret := &ThirdRpcCmd{thirdConfig: &thirdConfig}
	ret.configMap = map[string]any{
		config.OpenIMRPCThirdCfgFileName: &thirdConfig.RpcConfig,
		config.RedisConfigFileName:       &thirdConfig.RedisConfig,
		config.MongodbConfigFileName:     &thirdConfig.MongodbConfig,
		config.ShareFileName:             &thirdConfig.Share,
		config.NotificationFileName:      &thirdConfig.NotificationConfig,
		config.MinioConfigFileName:       &thirdConfig.MinioConfig,
		config.LocalCacheConfigFileName:  &thirdConfig.LocalCacheConfig,
		config.DiscoveryConfigFilename:   &thirdConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *ThirdRpcCmd) Exec() error {
	return a.Execute()
}

func (a *ThirdRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.thirdConfig.Discovery, &a.thirdConfig.RpcConfig.Prometheus, a.thirdConfig.RpcConfig.RPC.ListenIP,
		a.thirdConfig.RpcConfig.RPC.RegisterIP, a.thirdConfig.RpcConfig.RPC.AutoSetPorts, a.thirdConfig.RpcConfig.RPC.Ports,
		a.Index(), a.thirdConfig.Discovery.RpcService.Third, &a.thirdConfig.NotificationConfig, a.thirdConfig, third.Start)
}
