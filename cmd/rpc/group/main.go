package main

import (
	"Open_IM/internal/rpc/group"
	"Open_IM/internal/startrpc"
	"Open_IM/pkg/common/config"
)

func main() {
	//defaultPorts := config.Config.RpcPort.OpenImGroupPort
	//rpcPort := flag.Int("port", defaultPorts[0], "get RpcGroupPort from cmd,default 16000 as port")
	//prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.GroupPrometheusPort[0], "groupPrometheusPort default listen port")
	//flag.Parse()
	//fmt.Println("start group rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	//rpcServer := group.NewGroupServer(*rpcPort)
	//go func() {
	//	err := promePkg.StartPromeSrv(*prometheusPort)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()
	//rpcServer.Run()

	startrpc.StartRpc(config.Config.RpcPort.OpenImGroupPort[0], config.Config.RpcRegisterName.OpenImGroupName, config.Config.Prometheus.GroupPrometheusPort[0], group.Start)

}
