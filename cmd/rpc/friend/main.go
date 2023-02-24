package main

import (
	"OpenIM/internal/rpc/friend"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImFriendPort[0], config.Config.RpcRegisterName.OpenImFriendName, config.Config.Prometheus.FriendPrometheusPort[0], friend.Start)
}
