package startrpc

import (
	"OpenIM/internal/common/network"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/mw"
	"OpenIM/pkg/common/prome"
	"OpenIM/pkg/discoveryregistry"
	"fmt"
	"github.com/OpenIMSDK/openKeeper"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"net"
)

func start(rpcPort int, rpcRegisterName string, prometheusPort int, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error, options []grpc.ServerOption) error {
	fmt.Println("start group rpc server, port: ", rpcPort, ", OpenIM version: ", config.Version)
	log.NewPrivateLog(constant.LogFileName)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Config.ListenIP, rpcPort))
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema)
	zkClient, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema, 10, "", "")
	if err != nil {
		return err
	}
	defer zkClient.Close()
	registerIP, err := network.GetRpcRegisterIP(config.Config.RpcRegisterIP)
	if err != nil {
		return err
	}
	options = append(options, mw.GrpcServer()) // ctx 中间件
	if config.Config.Prometheus.Enable {
		prome.NewGrpcRequestCounter()
		prome.NewGrpcRequestFailedCounter()
		prome.NewGrpcRequestSuccessCounter()
		options = append(options, []grpc.ServerOption{
			//grpc.UnaryInterceptor(prome.UnaryServerInterceptorPrometheus),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(options...)
	defer srv.GracefulStop()
	err = zkClient.Register(rpcRegisterName, registerIP, rpcPort)
	if err != nil {
		return err
	}
	if config.Config.Prometheus.Enable && prometheusPort != 0 {
		err := prome.StartPrometheusSrv(prometheusPort)
		if err != nil {
			return err
		}
	}
	return rpcFn(zkClient, srv)
}

func Start(rpcPort int, rpcRegisterName string, prometheusPort int, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error, options ...grpc.ServerOption) error {
	return start(rpcPort, rpcRegisterName, prometheusPort, rpcFn, options)
}
