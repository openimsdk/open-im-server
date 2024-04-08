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
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"github.com/openimsdk/open-im-server/v3/internal/msggateway"

	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type MsgGatewayCmd struct {
	*RootCmd
	ctx              context.Context
	configMap        map[string]StructEnvPrefix
	msgGatewayConfig MsgGatewayConfig
}
type MsgGatewayConfig struct {
	MsgGateway      config.MsgGateway
	RedisConfig     config.Redis
	ZookeeperConfig config.ZooKeeper
	Share           config.Share
	WebhooksConfig  config.Webhooks
}

func NewMsgGatewayCmd() *MsgGatewayCmd {
	var msgGatewayConfig MsgGatewayConfig
	ret := &MsgGatewayCmd{msgGatewayConfig: msgGatewayConfig}
	ret.configMap = map[string]StructEnvPrefix{
		OpenIMAPICfgFileName:    {EnvPrefix: apiEnvPrefix, ConfigStruct: &msgGatewayConfig.MsgGateway},
		RedisConfigFileName:     {EnvPrefix: redisEnvPrefix, ConfigStruct: &msgGatewayConfig.RedisConfig},
		ZookeeperConfigFileName: {EnvPrefix: zoopkeeperEnvPrefix, ConfigStruct: &msgGatewayConfig.ZookeeperConfig},
		ShareFileName:           {EnvPrefix: shareEnvPrefix, ConfigStruct: &msgGatewayConfig.Share},
		WebhooksConfigFileName:  {EnvPrefix: webhooksEnvPrefix, ConfigStruct: &msgGatewayConfig.WebhooksConfig},
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return ret.preRunE()
	}
	return ret
}

func (m *MsgGatewayCmd) Exec() error {
	return m.Execute()
}

func (m *MsgGatewayCmd) preRunE() error {
	return msggateway.Start(m.ctx, m.Index(), &m.msgGatewayConfig)
}
