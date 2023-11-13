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
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"

	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/mw"
	"github.com/OpenIMSDK/tools/network"
	"github.com/OpenIMSDK/tools/utils"
)

// Start rpc server.
func Start(
	rpcPort int,
	rpcRegisterName string,
	prometheusPort int,
	rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error,
	options ...grpc.ServerOption,
) error {
	fmt.Printf("start %s server, port: %d, prometheusPort: %d, OpenIM version: %s\n",
		rpcRegisterName, rpcPort, prometheusPort, config.Version)
	listener, err := net.Listen(
		"tcp",
		net.JoinHostPort(network.GetListenIP(config.Config.Rpc.ListenIP), strconv.Itoa(rpcPort)),
	)
	if err != nil {
		return err
	}
	defer listener.Close()
	client, err := kdisc.NewDiscoveryRegister(config.Config.Envs.Discovery)
	if err != nil {
		return utils.Wrap1(err)
	}
	defer client.Close()
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	registerIP, err := network.GetRpcRegisterIP(config.Config.Rpc.RegisterIP)
	if err != nil {
		return err
	}
	var reg *prometheus.Registry
	var metric *grpcprometheus.ServerMetrics
	// ctx 中间件
	if config.Config.Prometheus.Enable {
		//////////////////////////
		cusMetrics := prommetrics.GetGrpcCusMetrics(rpcRegisterName)
		reg, metric, err = prommetrics.NewGrpcPromObj(cusMetrics)
		options = append(options, mw.GrpcServer(), grpc.StreamInterceptor(metric.StreamServerInterceptor()),
			grpc.UnaryInterceptor(metric.UnaryServerInterceptor()))
	} else {
		options = append(options, mw.GrpcServer())
	}
	srv := grpc.NewServer(options...)
	defer srv.GracefulStop()
	err = rpcFn(client, srv)
	if err != nil {
		return utils.Wrap1(err)
	}
	err = client.Register(
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
			metric.InitializeMetrics(srv)
			// Create a HTTP server for prometheus.
			httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", prometheusPort)}
			if err := httpServer.ListenAndServe(); err != nil {
				log.Fatal("Unable to start a http server.")
			}
		}
	}()

	return utils.Wrap1(srv.Serve(listener))
}
