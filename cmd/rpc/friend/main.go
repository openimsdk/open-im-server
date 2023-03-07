package main

import (
	"OpenIM/internal/rpc/friend"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	if err := config.InitConfig(""); err != nil {
		panic(err.Error())
	}
	if err := startrpc.Start(config.Config.RpcPort.OpenImFriendPort[0], config.Config.RpcRegisterName.OpenImFriendName, config.Config.Prometheus.FriendPrometheusPort[0], friend.Start); err != nil {
		panic(err.Error())
	}
}
