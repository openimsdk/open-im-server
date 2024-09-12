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

	"github.com/openimsdk/open-im-server/v3/internal/rpc/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type FriendRpcCmd struct {
	*RootCmd
	ctx            context.Context
	configMap      map[string]any
	relationConfig *relation.Config
}

func NewFriendRpcCmd() *FriendRpcCmd {
	var relationConfig relation.Config
	ret := &FriendRpcCmd{relationConfig: &relationConfig}
	ret.configMap = map[string]any{
		OpenIMRPCFriendCfgFileName: &relationConfig.RpcConfig,
		RedisConfigFileName:        &relationConfig.RedisConfig,
		MongodbConfigFileName:      &relationConfig.MongodbConfig,
		ShareFileName:              &relationConfig.Share,
		NotificationFileName:       &relationConfig.NotificationConfig,
		WebhooksConfigFileName:     &relationConfig.WebhooksConfig,
		LocalCacheConfigFileName:   &relationConfig.LocalCacheConfig,
		DiscoveryConfigFilename:    &relationConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *FriendRpcCmd) Exec() error {
	return a.Execute()
}

func (a *FriendRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.relationConfig.Discovery, &a.relationConfig.RpcConfig.Prometheus, a.relationConfig.RpcConfig.RPC.ListenIP,
		a.relationConfig.RpcConfig.RPC.RegisterIP, a.relationConfig.RpcConfig.RPC.Ports,
		a.Index(), a.relationConfig.Share.RpcRegisterName.Friend, &a.relationConfig.Share, a.relationConfig, relation.Start)
}
