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
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type PushRpcCmd struct {
	*RootCmd
	ctx        context.Context
	configMap  map[string]StructEnvPrefix
	pushConfig PushConfig
}
type PushConfig struct {
	RpcConfig          config.Push
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	KafkaConfig        config.Kafka
	ZookeeperConfig    config.ZooKeeper
	NotificationConfig config.Notification
	Share              config.Share
	WebhooksConfig     config.Webhooks
}

func NewPushRpcCmd() *PushRpcCmd {
	var pushConfig PushConfig
	ret := &PushRpcCmd{pushConfig: pushConfig}
	ret.configMap = map[string]StructEnvPrefix{
		OpenIMPushCfgFileName:   {EnvPrefix: pushEnvPrefix, ConfigStruct: &pushConfig.RpcConfig},
		RedisConfigFileName:     {EnvPrefix: redisEnvPrefix, ConfigStruct: &pushConfig.RedisConfig},
		ZookeeperConfigFileName: {EnvPrefix: zoopkeeperEnvPrefix, ConfigStruct: &pushConfig.ZookeeperConfig},
		MongodbConfigFileName:   {EnvPrefix: mongodbEnvPrefix, ConfigStruct: &pushConfig.MongodbConfig},
		KafkaConfigFileName:     {EnvPrefix: kafkaEnvPrefix, ConfigStruct: &pushConfig.KafkaConfig},
		ShareFileName:           {EnvPrefix: shareEnvPrefix, ConfigStruct: &pushConfig.Share},
		NotificationFileName:    {EnvPrefix: notificationEnvPrefix, ConfigStruct: &pushConfig.NotificationConfig},
		WebhooksConfigFileName:  {EnvPrefix: webhooksEnvPrefix, ConfigStruct: &pushConfig.WebhooksConfig},
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return ret.preRunE()
	}
	return ret
}

func (a *PushRpcCmd) Exec() error {
	return a.Execute()
}

func (a *PushRpcCmd) preRunE() error {
	return startrpc.Start(a.ctx, &a.pushConfig.ZookeeperConfig, &a.pushConfig.RpcConfig.Prometheus, a.pushConfig.RpcConfig.RPC.ListenIP,
		a.pushConfig.RpcConfig.RPC.RegisterIP, a.pushConfig.RpcConfig.RPC.Ports,
		a.Index(), a.pushConfig.Share.RpcRegisterName.Auth, &a.pushConfig, push.Start)
}
