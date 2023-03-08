package cmd

import (
	//"OpenIM/internal/msggateway"
	"OpenIM/pkg/common/constant"
	"github.com/spf13/cobra"
)

type MsgGatewayCmd struct {
	*RootCmd
}

func NewMsgGatewayCmd() MsgGatewayCmd {
	return MsgGatewayCmd{NewRootCmd()}
}

func (m *MsgGatewayCmd) AddWsPortFlag() {
	m.Command.Flags().IntP(constant.FlagWsPort, "w", 0, "ws server listen port")
}

func (m *MsgGatewayCmd) getWsPortFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.FlagWsPort)
	return port
}

func (m *MsgGatewayCmd) addRun() {
	m.Command.Run = func(cmd *cobra.Command, args []string) {
		//msggateway.Init(m.getPortFlag(cmd), m.getWsPortFlag(cmd))
		//msggateway.Run(m.getPrometheusPortFlag(cmd))
	}
}

func (m *MsgGatewayCmd) Exec() error {
	m.addRun()
	return m.Execute()
}
