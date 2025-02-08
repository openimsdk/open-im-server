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
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/network"
	"google.golang.org/grpc/status"

	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discovery"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Start[T any](ctx context.Context, discovery *conf.Discovery, prometheusConfig *conf.Prometheus, listenIP,
	registerIP string, autoSetPorts bool, rpcPorts []int, index int, rpcRegisterName string, notification *conf.Notification, config T,
	watchConfigNames []string, watchServiceNames []string,
	rpcFn func(ctx context.Context, config T, client discovery.Conn, server grpc.ServiceRegistrar) error,
	options ...grpc.ServerOption) error {

	if notification != nil {
		conf.InitNotification(notification)
	}

	if discovery.Enable == conf.Standalone {
		return nil
	}

	options = append(options, mw.GrpcServer())

	registerIP, err := network.GetRpcRegisterIP(registerIP)
	if err != nil {
		return err
	}
	var (
		rpcListenAddr        string
		prometheusListenAddr string
	)
	if autoSetPorts {
		rpcListenAddr = net.JoinHostPort(listenIP, "0")
		prometheusListenAddr = net.JoinHostPort("", "0")
	} else {
		rpcPort, err := datautil.GetElemByIndex(rpcPorts, index)
		if err != nil {
			return err
		}
		prometheusPort, err := datautil.GetElemByIndex(prometheusConfig.Ports, index)
		if err != nil {
			return err
		}
		rpcListenAddr = net.JoinHostPort(listenIP, strconv.Itoa(rpcPort))
		prometheusListenAddr = net.JoinHostPort("", strconv.Itoa(prometheusPort))
	}

	log.CInfo(ctx, "RPC server is initializing", "rpcRegisterName", rpcRegisterName, "rpcAddr", rpcListenAddr, "prometheusAddr", prometheusListenAddr)

	watchConfigNames = append(watchConfigNames, conf.LogConfigFileName)

	client, err := kdisc.NewDiscoveryRegister(discovery, watchServiceNames)
	if err != nil {
		return err
	}

	defer client.Close()
	client.AddOption(
		mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")),
	)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	var gsrv grpcServiceRegistrar

	err = rpcFn(ctx, config, client, &gsrv)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancelCause(ctx)

	if prometheusListenAddr != "" {
		options = append(
			options,
			prommetricsUnaryInterceptor(rpcRegisterName),
			prommetricsStreamInterceptor(rpcRegisterName),
		)
		prometheusListener, prometheusPort, err := listenTCP(prometheusListenAddr)
		if err != nil {
			return err
		}
		if err := client.Register(ctx, "prometheus_"+rpcRegisterName, registerIP, prometheusPort); err != nil {
			return err
		}

		cs := prommetrics.GetGrpcCusMetrics(rpcRegisterName, discovery)
		go func() {
			err := prommetrics.RpcInit(cs, prometheusListener)
			if err == nil {
				err = fmt.Errorf("serve end")
			}
			cancel(fmt.Errorf("prommetrics %s %w", rpcRegisterName, err))
		}()
	}

	var rpcGracefulStop chan struct{}

	if len(gsrv.services) > 0 {
		rpcListener, rpcPort, err := listenTCP(rpcListenAddr)
		if err != nil {
			return err
		}
		srv := grpc.NewServer(options...)

		for _, service := range gsrv.services {
			srv.RegisterService(service.desc, service.impl)
		}
		grpcOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
		if err := client.Register(ctx, rpcRegisterName, registerIP, rpcPort, grpcOpt); err != nil {
			return err
		}

		rpcGracefulStop = make(chan struct{})

		go func() {
			err := srv.Serve(rpcListener)
			if err == nil {
				err = fmt.Errorf("serve end")
			}
			cancel(fmt.Errorf("rpc %s %w", rpcRegisterName, err))
		}()

		go func() {
			<-ctx.Done()
			srv.GracefulStop()
			close(rpcGracefulStop)
		}()
	}

	select {
	case val := <-sigs:
		log.ZDebug(ctx, "recv exit", "signal", val.String())
		cancel(fmt.Errorf("signal %s", val.String()))
	case <-ctx.Done():
	}
	if rpcGracefulStop != nil {
		timeout := time.NewTimer(time.Second * 15)
		defer timeout.Stop()
		select {
		case <-timeout.C:
			log.ZWarn(ctx, "rcp graceful stop timeout", nil)
		case <-rpcGracefulStop:
			log.ZDebug(ctx, "rcp graceful stop done")
		}
	}
	return context.Cause(ctx)
}

func listenTCP(addr string) (net.Listener, int, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, 0, errs.WrapMsg(err, "listen err", "addr", addr)
	}
	return listener, listener.Addr().(*net.TCPAddr).Port, nil
}

func prommetricsUnaryInterceptor(rpcRegisterName string) grpc.ServerOption {
	getCode := func(err error) int {
		if err == nil {
			return 0
		}
		rpcErr, ok := err.(interface{ GRPCStatus() *status.Status })
		if !ok {
			return -1
		}
		return int(rpcErr.GRPCStatus().Code())
	}
	return grpc.ChainUnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		prommetrics.RPCCall(rpcRegisterName, info.FullMethod, getCode(err))
		return resp, err
	})
}

func prommetricsStreamInterceptor(rpcRegisterName string) grpc.ServerOption {
	return grpc.ChainStreamInterceptor()
}

type grpcService struct {
	desc *grpc.ServiceDesc
	impl any
}

type grpcServiceRegistrar struct {
	services []*grpcService
}

func (x *grpcServiceRegistrar) RegisterService(desc *grpc.ServiceDesc, impl any) {
	x.services = append(x.services, &grpcService{
		desc: desc,
		impl: impl,
	})
}
