package startrpc

import (
	"Open_IM/internal/common/network"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/middleware"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
	"github.com/OpenIMSDK/openKeeper"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"net"
)

func StartRpc(rpcPort int, rpcRegisterName string, prometheusPort int, fn func(server *grpc.Server), options ...grpc.ServerOption) {
	flagRpcPort := flag.Int("port", rpcPort, "get RpcGroupPort from cmd,default 16000 as port")
	flagPrometheusPort := flag.Int("prometheus_port", prometheusPort, "groupPrometheusPort default listen port")
	flag.Parse()
	rpcPort = *flagRpcPort
	prometheusPort = *flagPrometheusPort
	fmt.Println("start group rpc server, port: ", rpcPort, ", OpenIM version: ", constant.CurrentVersion)
	log.NewPrivateLog(constant.LogFileName)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Config.ListenIP, rpcPort))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	zkClient, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema, 10, "", "")
	if err != nil {
		panic(err.Error())
	}
	registerIP, err := network.GetRpcRegisterIP(config.Config.RpcRegisterIP)
	if err != nil {
		panic(err)
	}
	options = append(options, grpc.UnaryInterceptor(middleware.RpcServerInterceptor)) // ctx 中间件
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		options = append(options, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(options...)
	defer srv.GracefulStop()
	fn(srv)
	err = zkClient.Register(rpcRegisterName, registerIP, rpcPort)
	if err != nil {
		panic(err.Error())
	}
	if config.Config.Prometheus.Enable {
		err := promePkg.StartPromeSrv(prometheusPort)
		if err != nil {
			panic(err)
		}
	}
	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
