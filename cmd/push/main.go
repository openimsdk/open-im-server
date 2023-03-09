package main

import (
	"OpenIM/internal/push"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	pushCmd := cmd.NewPushCmd()
	pushCmd.AddPortFlag()
	pushCmd.AddPrometheusPortFlag()
	if err := pushCmd.Exec(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := pushCmd.StartSvr(config.Config.RpcRegisterName.OpenImPushName, push.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
