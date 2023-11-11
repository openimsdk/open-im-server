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
	"fmt"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/spf13/cobra"

	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"github.com/openimsdk/open-im-server/v3/internal/msgtransfer"
)

type MsgTransferCmd struct {
	*RootCmd
}

func NewMsgTransferCmd() *MsgTransferCmd {
	ret := &MsgTransferCmd{NewRootCmd("msgTransfer")}
	ret.SetRootCmdPt(ret)
	return ret
}

func (m *MsgTransferCmd) addRunE() {
	m.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return msgtransfer.StartTransfer(m.getPrometheusPortFlag(cmd))
	}
}

func (m *MsgTransferCmd) Exec() error {
	m.addRunE()
	return m.Execute()
}

func (m *MsgTransferCmd) GetPortFromConfig(portType string) int {
	fmt.Println("GetPortFromConfig:", portType)
	if portType == constant.FlagPort {
		return 0
	} else if portType == constant.FlagPrometheusPort {
		n := m.getTransferProgressFlagValue()
		return config2.Config.Prometheus.MessageTransferPrometheusPort[n]
	}
	return 0
}
func (m *MsgTransferCmd) AddTransferProgressFlag() {
	m.Command.Flags().IntP(constant.FlagTransferProgressIndex, "n", 0, "transfer progress index")
}
func (m *MsgTransferCmd) getTransferProgressFlagValue() int {
	nindex, err := m.Command.Flags().GetInt(constant.FlagTransferProgressIndex)
	if err != nil {
		fmt.Println("get transfercmd error,make sure it is k8s env or not")
		return 0
	}
	return nindex
}
