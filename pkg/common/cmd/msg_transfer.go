package cmd

import (
	"github.com/OpenIMSDK/Open-IM-Server/internal/msgtransfer"
	"github.com/spf13/cobra"
)

type MsgTransferCmd struct {
	*RootCmd
}

func NewMsgTransferCmd() MsgTransferCmd {
	return MsgTransferCmd{NewRootCmd("msgTransfer")}
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
