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
	"github.com/openimsdk/open-im-server/v3/internal/rpc/group"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type GroupRpcCmd struct {
	*RootCmd
	ctx         context.Context
	configMap   map[string]StructEnvPrefix
	groupConfig GroupConfig
}
type GroupConfig struct {
	RpcConfig          config.Group
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	ZookeeperConfig    config.ZooKeeper
	NotificationConfig config.Notification
	Share              config.Share
	WebhooksConfig     config.Webhooks
}

func NewGroupRpcCmd() *GroupRpcCmd {
	var groupConfig GroupConfig
	ret := &GroupRpcCmd{groupConfig: groupConfig}
	ret.configMap = map[string]StructEnvPrefix{
		OpenIMRPCGroupCfgFileName: {EnvPrefix: groupEnvPrefix, ConfigStruct: &groupConfig.RpcConfig},
		RedisConfigFileName:       {EnvPrefix: redisEnvPrefix, ConfigStruct: &groupConfig.RedisConfig},
		ZookeeperConfigFileName:   {EnvPrefix: zoopkeeperEnvPrefix, ConfigStruct: &groupConfig.ZookeeperConfig},
		MongodbConfigFileName:     {EnvPrefix: mongodbEnvPrefix, ConfigStruct: &groupConfig.MongodbConfig},
		ShareFileName:             {EnvPrefix: shareEnvPrefix, ConfigStruct: &groupConfig.Share},
		NotificationFileName:      {EnvPrefix: notificationEnvPrefix, ConfigStruct: &groupConfig.NotificationConfig},
		WebhooksConfigFileName:    {EnvPrefix: webhooksEnvPrefix, ConfigStruct: &groupConfig.WebhooksConfig},
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return ret.preRunE()
	}
	return ret
}

func (a *GroupRpcCmd) Exec() error {
	return a.Execute()
}

func (a *GroupRpcCmd) preRunE() error {
	return startrpc.Start(a.ctx, &a.groupConfig.ZookeeperConfig, &a.groupConfig.RpcConfig.Prometheus, a.groupConfig.RpcConfig.RPC.ListenIP,
		a.groupConfig.RpcConfig.RPC.RegisterIP, a.groupConfig.RpcConfig.RPC.Ports,
		a.Index(), a.groupConfig.Share.RpcRegisterName.Auth, &a.groupConfig, group.Start)
}
