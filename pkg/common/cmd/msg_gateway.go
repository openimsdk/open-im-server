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

	"github.com/openimsdk/open-im-server/v3/internal/msggateway"
	"github.com/openimsdk/open-im-server/v3/version"
)

type MsgGatewayCmd struct {
	*RootCmd
	ctx              context.Context
	configMap        map[string]any
	msgGatewayConfig *msggateway.Config
}

func NewMsgGatewayCmd() *MsgGatewayCmd {
	var msgGatewayConfig msggateway.Config
	ret := &MsgGatewayCmd{msgGatewayConfig: &msgGatewayConfig}
	ret.configMap = map[string]any{
		OpenIMMsgGatewayCfgFileName: &msgGatewayConfig.MsgGateway,
		ShareFileName:               &msgGatewayConfig.Share,
		RedisConfigFileName:         &msgGatewayConfig.RedisConfig,
		WebhooksConfigFileName:      &msgGatewayConfig.WebhooksConfig,
		DiscoveryConfigFilename:     &msgGatewayConfig.Discovery,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (m *MsgGatewayCmd) Exec() error {
	return m.Execute()
}

func (m *MsgGatewayCmd) runE() error {
	return msggateway.Start(m.ctx, m.Index(), m.msgGatewayConfig)
}
