package main

import (
	"OpenIM/internal/rpc/auth"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	authCmd := cmd.NewAuthCmd()
	authCmd.AddPortFlag()
	authCmd.AddPrometheusPortFlag()
	if err := authCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := authCmd.StartSvr(config.Config.RpcRegisterName.OpenImAuthName, auth.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
