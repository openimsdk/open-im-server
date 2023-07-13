// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package startrpc

import (
	"fmt"
	"net"
	"strconv"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/network"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	openKeeper "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func Start(
	rpcPort int,
	rpcRegisterName string,
	prometheusPort int,
	rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error,
	options ...grpc.ServerOption,
) error {
	rpcPort = 10030
	fmt.Println(
		"start",
		rpcRegisterName,
		"server, port: ",
		rpcPort,
		"prometheusPort:",
		prometheusPort,
		", OpenIM version: ",
		config.Version,
	)
	listener, err := net.Listen(
		"tcp",
		net.JoinHostPort(network.GetListenIP(config.Config.Rpc.ListenIP), strconv.Itoa(rpcPort)),
	)
	if err != nil {
		return err
	}
	defer listener.Close()
	zkClient, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(
			config.Config.Zookeeper.Username,
			config.Config.Zookeeper.Password,
		), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		return utils.Wrap1(err)
	}
	defer zkClient.CloseZK()
	zkClient.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	registerIP, err := network.GetRpcRegisterIP(config.Config.Rpc.RegisterIP)
	if err != nil {
		return err
	}
	// ctx 中间件
	if config.Config.Prometheus.Enable {
		prome.NewGrpcRequestCounter()
		prome.NewGrpcRequestFailedCounter()
		prome.NewGrpcRequestSuccessCounter()
		unaryInterceptor := mw.InterceptChain(grpcPrometheus.UnaryServerInterceptor, mw.RpcServerInterceptor)
		options = append(options, []grpc.ServerOption{
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(unaryInterceptor),
		}...)
	} else {
		options = append(options, mw.GrpcServer())
	}
	srv := grpc.NewServer(options...)
	defer srv.GracefulStop()
	err = rpcFn(zkClient, srv)
	if err != nil {
		return utils.Wrap1(err)
	}
	err = zkClient.Register(
		rpcRegisterName,
		registerIP,
		rpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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
	return utils.Wrap1(srv.Serve(listener))
}
