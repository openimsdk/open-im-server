package main

import (
	"OpenIM/internal/rpc/friend"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	rpcCmd := cmd.NewRpcCmd(config.Config.RpcRegisterName.OpenImFriendName)
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(friend.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
