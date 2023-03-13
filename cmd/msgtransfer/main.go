package main

import (
	"OpenIM/pkg/common/cmd"
)

func main() {
	msgTransferCmd := cmd.NewMsgTransferCmd()
	msgTransferCmd.AddPrometheusPortFlag()
	if err := msgTransferCmd.Exec(); err != nil {
		panic(err.Error())
	}
}
