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
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type ConversationRpcCmd struct {
	*RootCmd
	ctx                context.Context
	configMap          map[string]any
	conversationConfig *conversation.Config
}

func NewConversationRpcCmd() *ConversationRpcCmd {
	var conversationConfig conversation.Config
	ret := &ConversationRpcCmd{conversationConfig: &conversationConfig}
	ret.configMap = map[string]any{
		OpenIMRPCConversationCfgFileName: &conversationConfig.RpcConfig,
		RedisConfigFileName:              &conversationConfig.RedisConfig,
		MongodbConfigFileName:            &conversationConfig.MongodbConfig,
		ShareFileName:                    &conversationConfig.Share,
		NotificationFileName:             &conversationConfig.NotificationConfig,
		LocalCacheConfigFileName:         &conversationConfig.LocalCacheConfig,
		DiscoveryConfigFilename:          &conversationConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (a *ConversationRpcCmd) Exec() error {
	return a.Execute()
}

func (a *ConversationRpcCmd) runE() error {
	return startrpc.Start(a.ctx, &a.conversationConfig.Discovery, &a.conversationConfig.RpcConfig.Prometheus, a.conversationConfig.RpcConfig.RPC.ListenIP,
		a.conversationConfig.RpcConfig.RPC.RegisterIP, a.conversationConfig.RpcConfig.RPC.AutoSetPorts, a.conversationConfig.RpcConfig.RPC.Ports,
		a.Index(), a.conversationConfig.Discovery.RpcService.Conversation, &a.conversationConfig.NotificationConfig, a.conversationConfig, conversation.Start)
}
