package main

import (
	"OpenIM/pkg/common/cmd"
	"fmt"
	"os"
)

func main() {
	msgTransferCmd := cmd.NewMsgTransferCmd()
	msgTransferCmd.AddPrometheusPortFlag()
	if err := msgTransferCmd.Exec(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
