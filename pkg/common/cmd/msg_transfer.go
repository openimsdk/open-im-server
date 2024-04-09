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
	"github.com/openimsdk/open-im-server/v3/internal/msgtransfer"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type MsgTransferCmd struct {
	*RootCmd
	ctx               context.Context
	configMap         map[string]StructEnvPrefix
	msgTransferConfig msgtransfer.Config
}

func NewMsgTransferCmd() *MsgTransferCmd {
	var msgTransferConfig msgtransfer.Config
	ret := &MsgTransferCmd{msgTransferConfig: msgTransferConfig}
	ret.configMap = map[string]StructEnvPrefix{
		OpenIMMsgTransferCfgFileName: {EnvPrefix: msgTransferEnvPrefix, ConfigStruct: &msgTransferConfig.MsgTransfer},
		RedisConfigFileName:          {EnvPrefix: redisEnvPrefix, ConfigStruct: &msgTransferConfig.RedisConfig},
		MongodbConfigFileName:        {EnvPrefix: mongodbEnvPrefix, ConfigStruct: &msgTransferConfig.MongodbConfig},
		KafkaConfigFileName:          {EnvPrefix: kafkaEnvPrefix, ConfigStruct: &msgTransferConfig.KafkaConfig},
		ZookeeperConfigFileName:      {EnvPrefix: zoopkeeperEnvPrefix, ConfigStruct: &msgTransferConfig.ZookeeperConfig},
		ShareFileName:                {EnvPrefix: shareEnvPrefix, ConfigStruct: &msgTransferConfig.Share},
		WebhooksConfigFileName:       {EnvPrefix: webhooksEnvPrefix, ConfigStruct: &msgTransferConfig.WebhooksConfig},
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", config.Version)
	ret.Command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return ret.preRunE()
	}
	return ret
}

func (m *MsgTransferCmd) Exec() error {
	return m.Execute()
}

func (m *MsgTransferCmd) preRunE() error {
	return msgtransfer.Start(m.ctx, m.Index(), &m.msgTransferConfig)
}
