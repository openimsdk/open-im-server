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

	"github.com/openimsdk/open-im-server/v3/internal/api"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type ApiCmd struct {
	*RootCmd
	ctx       context.Context
	configMap map[string]any
	apiConfig *api.Config
}

func NewApiCmd() *ApiCmd {
	var apiConfig api.Config
	ret := &ApiCmd{apiConfig: &apiConfig}
	ret.configMap = map[string]any{
		config.DiscoveryConfigFilename:          &apiConfig.Discovery,
		config.KafkaConfigFileName:              &apiConfig.Kafka,
		config.LocalCacheConfigFileName:         &apiConfig.LocalCache,
		config.LogConfigFileName:                &apiConfig.Log,
		config.MinioConfigFileName:              &apiConfig.Minio,
		config.MongodbConfigFileName:            &apiConfig.Mongo,
		config.NotificationFileName:             &apiConfig.Notification,
		config.OpenIMAPICfgFileName:             &apiConfig.API,
		config.OpenIMCronTaskCfgFileName:        &apiConfig.CronTask,
		config.OpenIMMsgGatewayCfgFileName:      &apiConfig.MsgGateway,
		config.OpenIMMsgTransferCfgFileName:     &apiConfig.MsgTransfer,
		config.OpenIMPushCfgFileName:            &apiConfig.Push,
		config.OpenIMRPCAuthCfgFileName:         &apiConfig.Auth,
		config.OpenIMRPCConversationCfgFileName: &apiConfig.Conversation,
		config.OpenIMRPCFriendCfgFileName:       &apiConfig.Friend,
		config.OpenIMRPCGroupCfgFileName:        &apiConfig.Group,
		config.OpenIMRPCMsgCfgFileName:          &apiConfig.Msg,
		config.OpenIMRPCThirdCfgFileName:        &apiConfig.Third,
		config.OpenIMRPCUserCfgFileName:         &apiConfig.User,
		config.RedisConfigFileName:              &apiConfig.Redis,
		config.ShareFileName:                    &apiConfig.Share,
		config.WebhooksConfigFileName:           &apiConfig.Webhooks,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		apiConfig.ConfigPath = config.Path(ret.configPath)
		return ret.runE()
	}
	return ret
}

func (a *ApiCmd) Exec() error {
	return a.Execute()
}

func (a *ApiCmd) runE() error {
	a.apiConfig.Index = config.Index(a.Index())
	var prometheus config.Prometheus
	return startrpc.Start(
		a.ctx, &a.apiConfig.Discovery,
		&prometheus,
		a.apiConfig.API.Api.ListenIP, "",
		false,
		nil, int(a.apiConfig.Index),
		a.apiConfig.Discovery.RpcService.MessageGateway,
		&a.apiConfig.Notification,
		a.apiConfig,
		[]string{},
		[]string{},
		api.Start,
	)
}
