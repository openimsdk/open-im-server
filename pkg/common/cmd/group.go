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

	"github.com/spf13/cobra"

	"github.com/openimsdk/tools/system/program"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/group"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/versionctx"
	"github.com/openimsdk/open-im-server/v3/version"
)

type GroupRpcCmd struct {
	*RootCmd
	ctx         context.Context
	configMap   map[string]any
	groupConfig *group.Config
}

func NewGroupRpcCmd() *GroupRpcCmd {
	var groupConfig group.Config
	ret := &GroupRpcCmd{groupConfig: &groupConfig}
	ret.configMap = map[string]any{
		OpenIMRPCGroupCfgFileName: &groupConfig.RpcConfig,
		RedisConfigFileName:       &groupConfig.RedisConfig,
		MongodbConfigFileName:     &groupConfig.MongodbConfig,
		ShareFileName:             &groupConfig.Share,
		NotificationFileName:      &groupConfig.NotificationConfig,
		WebhooksConfigFileName:    &groupConfig.WebhooksConfig,
		LocalCacheConfigFileName:  &groupConfig.LocalCacheConfig,
		DiscoveryConfigFilename:   &groupConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *GroupRpcCmd) Exec() error {
	return a.Execute()
}

func (a *GroupRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.groupConfig.Discovery, &a.groupConfig.RpcConfig.Prometheus, a.groupConfig.RpcConfig.RPC.ListenIP,
		a.Index(), a.groupConfig.Share.RpcRegisterName.Group, &a.groupConfig.Share, a.groupConfig, group.Start, versionctx.EnableVersionCtx())
}
