package main

import (
	"OpenIM/pkg/common/cmd"
	"fmt"
	"os"
)

func main() {
	pushCmd := cmd.NewPushCmd()
	pushCmd.AddPortFlag()
	pushCmd.AddPrometheusPortFlag()
	pushCmd.AddPush()
	if err := pushCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
