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
	"github.com/openimsdk/open-im-server/v3/internal/rpc/friend"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type FriendRpcCmd struct {
	*RootCmd
	ctx          context.Context
	configMap    map[string]any
	friendConfig *friend.Config
}

func NewFriendRpcCmd() *FriendRpcCmd {
	var friendConfig friend.Config
	ret := &FriendRpcCmd{friendConfig: &friendConfig}
	ret.configMap = map[string]any{
		OpenIMRPCFriendCfgFileName: &friendConfig.RpcConfig,
		RedisConfigFileName:        &friendConfig.RedisConfig,
		MongodbConfigFileName:      &friendConfig.MongodbConfig,
		ShareFileName:              &friendConfig.Share,
		NotificationFileName:       &friendConfig.NotificationConfig,
		WebhooksConfigFileName:     &friendConfig.WebhooksConfig,
		LocalCacheConfigFileName:   &friendConfig.LocalCacheConfig,
		DiscoveryConfigFilename:    &friendConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *FriendRpcCmd) Exec() error {
	return a.Execute()
}

func (a *FriendRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.friendConfig.Discovery, &a.friendConfig.RpcConfig.Prometheus, a.friendConfig.RpcConfig.RPC.ListenIP,
		a.friendConfig.RpcConfig.RPC.RegisterIP, a.friendConfig.RpcConfig.RPC.Ports,
		a.Index(), a.friendConfig.Share.RpcRegisterName.Friend, &a.friendConfig.Share, a.friendConfig, friend.Start)
}
