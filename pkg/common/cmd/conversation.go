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
	"github.com/openimsdk/open-im-server/v3/internal/rpc/conversation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type ConversationRpcCmd struct {
	*RootCmd
	ctx                context.Context
	configMap          map[string]StructEnvPrefix
	conversationConfig ConversationConfig
}
type ConversationConfig struct {
	RpcConfig          config.Conversation
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	ZookeeperConfig    config.ZooKeeper
	NotificationConfig config.Notification
	Share              config.Share
}

func NewConversationRpcCmd() *ConversationRpcCmd {
	var conversationConfig ConversationConfig
	ret := &ConversationRpcCmd{conversationConfig: conversationConfig}
	ret.configMap = map[string]StructEnvPrefix{
		OpenIMRPCConversationCfgFileName: {EnvPrefix: conversationEnvPrefix, ConfigStruct: &conversationConfig.RpcConfig},
		RedisConfigFileName:              {EnvPrefix: redisEnvPrefix, ConfigStruct: &conversationConfig.RedisConfig},
		ZookeeperConfigFileName:          {EnvPrefix: zoopkeeperEnvPrefix, ConfigStruct: &conversationConfig.ZookeeperConfig},
		MongodbConfigFileName:            {EnvPrefix: mongodbEnvPrefix, ConfigStruct: &conversationConfig.MongodbConfig},
		ShareFileName:                    {EnvPrefix: shareEnvPrefix, ConfigStruct: &conversationConfig.Share},
		NotificationFileName:             {EnvPrefix: notificationEnvPrefix, ConfigStruct: &conversationConfig.NotificationConfig},
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return ret.preRunE()
	}
	return ret
}

func (a *ConversationRpcCmd) Exec() error {
	return a.Execute()
}

func (a *ConversationRpcCmd) preRunE() error {
	return startrpc.Start(a.ctx, &a.conversationConfig.ZookeeperConfig, &a.conversationConfig.RpcConfig.Prometheus, a.conversationConfig.RpcConfig.RPC.ListenIP,
		a.conversationConfig.RpcConfig.RPC.RegisterIP, a.conversationConfig.RpcConfig.RPC.Ports,
		a.Index(), a.conversationConfig.Share.RpcRegisterName.Auth, &a.conversationConfig, conversation.Start)
}
