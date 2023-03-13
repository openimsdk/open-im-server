package main

import (
	"OpenIM/internal/rpc/auth"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	authCmd := cmd.NewRpcCmd("auth")
	authCmd.AddPortFlag()
	authCmd.AddPrometheusPortFlag()
	if err := authCmd.Exec(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := authCmd.StartSvr(config.Config.RpcRegisterName.OpenImAuthName, auth.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
