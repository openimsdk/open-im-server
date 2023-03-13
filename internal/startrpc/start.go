package startrpc

import (
	"OpenIM/internal/common/network"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/mw"
	"OpenIM/pkg/common/prome"
	"OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/utils"
	"fmt"
	"github.com/OpenIMSDK/openKeeper"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

func Start(rpcPort int, rpcRegisterName string, prometheusPort int, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error, options ...grpc.ServerOption) error {
	fmt.Println("start", rpcRegisterName, "rpc server, port: ", rpcPort, "prometheusPort:", prometheusPort, ", OpenIM version: ", config.Version)
	log.NewPrivateLog(constant.LogFileName)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Config.ListenIP, rpcPort))
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema, rpcRegisterName)
	zkClient, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema, 10, "", "")
	if err != nil {
		return utils.Wrap1(err)
	}
	defer zkClient.Close()
	zkClient.AddOption(grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	err = rpcFn(zkClient, srv)
	if err != nil {
		return utils.Wrap1(err)
	}
	err = zkClient.Register(rpcRegisterName, registerIP, rpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return utils.Wrap1(err)
	}
	go func() {
		if config.Config.Prometheus.Enable && prometheusPort != 0 {
			if err := prome.StartPrometheusSrv(prometheusPort); err != nil {
				panic(err.Error())
			}
		}
	}()
	err = srv.Serve(listener)
	if err != nil {
		return utils.Wrap1(err)
	}
	return nil
}

//func Start(rpcPort int, rpcRegisterName string, prometheusPort int, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error, options ...grpc.ServerOption) error {
//	return start(rpcPort, rpcRegisterName, prometheusPort, rpcFn, options)
//}
