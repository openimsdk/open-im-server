package main

import (
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/internal/rpc/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/cmd"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
)

func main() {
	fmt.Println("#######################################")
	rpcCmd := cmd.NewRpcCmd("third")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	name := config.Config.RpcRegisterName.OpenImThirdName
	fmt.Println("StartThirdRpc", "name:", name)
	if err := rpcCmd.StartSvr(name, third.Start); err != nil {
		panic(err.Error())
	}
}
