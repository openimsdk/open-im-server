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

	"github.com/openimsdk/open-im-server/v3/internal/msgtransfer"
	"github.com/openimsdk/open-im-server/v3/version"
)

type MsgTransferCmd struct {
	*RootCmd
	ctx               context.Context
	configMap         map[string]any
	msgTransferConfig *msgtransfer.Config
}

func NewMsgTransferCmd() *MsgTransferCmd {
	var msgTransferConfig msgtransfer.Config
	ret := &MsgTransferCmd{msgTransferConfig: &msgTransferConfig}
	ret.configMap = map[string]any{
		OpenIMMsgTransferCfgFileName: &msgTransferConfig.MsgTransfer,
		RedisConfigFileName:          &msgTransferConfig.RedisConfig,
		MongodbConfigFileName:        &msgTransferConfig.MongodbConfig,
		KafkaConfigFileName:          &msgTransferConfig.KafkaConfig,
		ShareFileName:                &msgTransferConfig.Share,
		WebhooksConfigFileName:       &msgTransferConfig.WebhooksConfig,
		DiscoveryConfigFilename:      &msgTransferConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (m *MsgTransferCmd) Exec() error {
	return m.Execute()
}

func (m *MsgTransferCmd) runE() error {
	return msgtransfer.Start(m.ctx, m.Index(), m.msgTransferConfig)
}
