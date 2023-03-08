package main

import (
	"OpenIM/pkg/common/cmd"
	"fmt"
	"os"
)

func main() {
	msgGatewayCmd := cmd.NewMsgGatewayCmd()
	msgGatewayCmd.AddPortFlag()
	msgGatewayCmd.AddPrometheusPortFlag()
	msgGatewayCmd.AddWsPortFlag()
	if err := msgGatewayCmd.Exec(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
