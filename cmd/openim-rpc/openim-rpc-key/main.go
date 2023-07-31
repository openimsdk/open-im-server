package main

import (
	"github.com/OpenIMSDK/Open-IM-Server/internal/rpc/key"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/cmd"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("key")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr("key", key.Start); err != nil {
		panic(err.Error())
	}
}
