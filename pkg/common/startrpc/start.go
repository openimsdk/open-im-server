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
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"

	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/jsonutil"
	"github.com/openimsdk/tools/utils/network"
	"google.golang.org/grpc/status"

	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discovery"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	grpccli "github.com/openimsdk/tools/mw/grpc/client"
	grpcsrv "github.com/openimsdk/tools/mw/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	prommetrics.RegistryAll()
}

func Start[T any](ctx context.Context, disc *conf.Discovery, prometheusConfig *conf.Prometheus, listenIP,
	registerIP string, autoSetPorts bool, rpcPorts []int, index int, rpcRegisterName string, notification *conf.Notification, config T,
	watchConfigNames []string, watchServiceNames []string,
	rpcFn func(ctx context.Context, config T, client discovery.SvcDiscoveryRegistry, server grpc.ServiceRegistrar) error,
	options ...grpc.ServerOption) error {

	if notification != nil {
		conf.InitNotification(notification)
	}

	maxRequestBody := getConfigRpcMaxRequestBody(reflect.ValueOf(config))
	shareConfig := getConfigShare(reflect.ValueOf(config))

	log.ZDebug(ctx, "rpc start", "rpcMaxRequestBody", maxRequestBody, "rpcRegisterName", rpcRegisterName, "registerIP", registerIP, "listenIP", listenIP)

	options = append(options,
		grpcsrv.GrpcServerMetadataContext(),
		grpcsrv.GrpcServerErrorConvert(),
		grpcsrv.GrpcServerLogger(),
		grpcsrv.GrpcServerRequestValidate(),
		grpcsrv.GrpcServerPanicCapture(),
	)
	if shareConfig != nil && len(shareConfig.IMAdminUser.UserIDs) > 0 {
		options = append(options, grpcServerIMAdminUserID(shareConfig.IMAdminUser.UserIDs))
	}
	var clientOptions []grpc.DialOption
	if maxRequestBody != nil {
		if maxRequestBody.RequestMaxBodySize > 0 {
			options = append(options, grpc.MaxRecvMsgSize(maxRequestBody.RequestMaxBodySize))
			clientOptions = append(clientOptions, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxRequestBody.RequestMaxBodySize)))
		}
		if maxRequestBody.ResponseMaxBodySize > 0 {
			options = append(options, grpc.MaxSendMsgSize(maxRequestBody.ResponseMaxBodySize))
			clientOptions = append(clientOptions, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxRequestBody.ResponseMaxBodySize)))
		}
	}

	registerIP, err := network.GetRpcRegisterIP(registerIP)
	if err != nil {
		return err
	}
	var prometheusListenAddr string
	if autoSetPorts {
		prometheusListenAddr = net.JoinHostPort(listenIP, "0")
	} else {
		prometheusPort, err := datautil.GetElemByIndex(prometheusConfig.Ports, index)
		if err != nil {
			return err
		}
		prometheusListenAddr = net.JoinHostPort(listenIP, strconv.Itoa(prometheusPort))
	}

	watchConfigNames = append(watchConfigNames, conf.LogConfigFileName)

	client, err := kdisc.NewDiscoveryRegister(disc, watchServiceNames)
	if err != nil {
		return err
	}

	defer client.Close()
	client.AddOption(
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")),

		grpccli.GrpcClientLogger(),
		grpccli.GrpcClientContext(),
		grpccli.GrpcClientErrorConvert(),
	)
	if len(clientOptions) > 0 {
		client.AddOption(clientOptions...)
	}

	ctx, cancel := context.WithCancelCause(ctx)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		select {
		case <-ctx.Done():
			return
		case val := <-sigs:
			log.ZDebug(ctx, "recv signal", "signal", val.String())
			cancel(fmt.Errorf("signal %s", val.String()))
		}
	}()

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
		log.ZDebug(ctx, "prometheus start", "addr", prometheusListener.Addr(), "rpcRegisterName", rpcRegisterName)
		target, err := jsonutil.JsonMarshal(prommetrics.BuildDefaultTarget(registerIP, prometheusPort))
		if err != nil {
			return err
		}
		if autoSetPorts {
			if err = client.SetWithLease(ctx, prommetrics.BuildDiscoveryKey(rpcRegisterName, index), target, prommetrics.TTL); err != nil {
				if !errors.Is(err, discovery.ErrNotSupported) {
					return err
				}
			}
		}
		go func() {
			err := prommetrics.Start(prometheusListener)
			if err == nil {
				err = fmt.Errorf("listener done")
			}
			cancel(fmt.Errorf("prommetrics %s %w", rpcRegisterName, err))
		}()
	}

	var (
		rpcServer       *grpc.Server
		rpcGracefulStop chan struct{}
	)

	onGrpcServiceRegistrar := func(desc *grpc.ServiceDesc, impl any) {
		if rpcServer != nil {
			rpcServer.RegisterService(desc, impl)
			return
		}
		var rpcListenAddr string
		if autoSetPorts {
			rpcListenAddr = net.JoinHostPort(listenIP, "0")
		} else {
			rpcPort, err := datautil.GetElemByIndex(rpcPorts, index)
			if err != nil {
				cancel(fmt.Errorf("rpcPorts index out of range %s %w", rpcRegisterName, err))
				return
			}
			rpcListenAddr = net.JoinHostPort(listenIP, strconv.Itoa(rpcPort))
		}
		rpcListener, err := net.Listen("tcp", rpcListenAddr)
		if err != nil {
			cancel(fmt.Errorf("listen rpc %s %s %w", rpcRegisterName, rpcListenAddr, err))
			return
		}

		rpcServer = grpc.NewServer(options...)
		rpcServer.RegisterService(desc, impl)
		rpcGracefulStop = make(chan struct{})
		rpcPort := rpcListener.Addr().(*net.TCPAddr).Port
		log.ZDebug(ctx, "rpc start register", "rpcRegisterName", rpcRegisterName, "registerIP", registerIP, "rpcPort", rpcPort)
		grpcOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
		rpcGracefulStop = make(chan struct{})
		go func() {
			<-ctx.Done()
			rpcServer.GracefulStop()
			close(rpcGracefulStop)
		}()
		if err := client.Register(ctx, rpcRegisterName, registerIP, rpcListener.Addr().(*net.TCPAddr).Port, grpcOpt); err != nil {
			cancel(fmt.Errorf("rpc register %s %w", rpcRegisterName, err))
			return
		}

		go func() {
			err := rpcServer.Serve(rpcListener)
			if err == nil {
				err = fmt.Errorf("serve end")
			}
			cancel(fmt.Errorf("rpc %s %w", rpcRegisterName, err))
		}()
	}

	err = rpcFn(ctx, config, client, &grpcServiceRegistrar{onRegisterService: onGrpcServiceRegistrar})
	if err != nil {
		return err
	}
	<-ctx.Done()
	log.ZDebug(ctx, "cmd wait done", "err", context.Cause(ctx))
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

type grpcServiceRegistrar struct {
	onRegisterService func(desc *grpc.ServiceDesc, impl any)
}

func (x *grpcServiceRegistrar) RegisterService(desc *grpc.ServiceDesc, impl any) {
	x.onRegisterService(desc, impl)
}
