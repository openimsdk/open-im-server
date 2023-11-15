package main

import (
	"github.com/openimsdk/open-im-server/v3/internal/rpc/aes_key"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd(cmd.RpcAesKeyServer)
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImAesKeyName, aes_key.Start); err != nil {
		panic(err.Error())
	}
}
