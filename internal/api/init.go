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

package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/network"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"google.golang.org/grpc"
)

type Config struct {
	conf.AllConfig

	ConfigPath conf.Path
	Index      conf.Index
}

func Start(ctx context.Context, config *Config, client discovery.Conn, service grpc.ServiceRegistrar) error {
	apiPort, err := datautil.GetElemByIndex(config.API.Api.Ports, int(config.Index))
	if err != nil {
		return err
	}

	//client, err := kdisc.NewDiscoveryRegister(&config.Discovery, []string{
	//	config.Discovery.RpcService.MessageGateway,
	//})
	//if err != nil {
	//	return errs.WrapMsg(err, "failed to register discovery service")
	//}
	//client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))

	//registerIP, err := network.GetRpcRegisterIP("")
	//if err != nil {
	//	return err
	//}
	//
	// todo
	//getAutoPort := func() (net.Listener, int, error) {
	//	registerAddr := net.JoinHostPort(registerIP, "0")
	//	listener, err := net.Listen("tcp", registerAddr)
	//	if err != nil {
	//		return nil, 0, errs.WrapMsg(err, "listen err", "registerAddr", registerAddr)
	//	}
	//	_, portStr, _ := net.SplitHostPort(listener.Addr().String())
	//	port, _ := strconv.Atoi(portStr)
	//	return listener, port, nil
	//}
	//
	//if config.API.Prometheus.AutoSetPorts && config.Discovery.Enable != conf.ETCD {
	//	return errs.New("only etcd support autoSetPorts", "RegisterName", "api").Wrap()
	//}

	router, err := newGinRouter(ctx, client, config)
	if err != nil {
		return err
	}
	//if config.API.Prometheus.Enable {
	//	var (
	//		listener net.Listener
	//	)
	//
	//	if config.API.Prometheus.AutoSetPorts {
	//		listener, prometheusPort, err = getAutoPort()
	//		if err != nil {
	//			return err
	//		}
	//
	//		etcdClient := client.(*etcd.SvcDiscoveryRegistryImpl).GetClient()
	//
	//		_, err = etcdClient.Put(ctx, prommetrics.BuildDiscoveryKey(prommetrics.APIKeyName), jsonutil.StructToJsonString(prommetrics.BuildDefaultTarget(registerIP, prometheusPort)))
	//		if err != nil {
	//			return errs.WrapMsg(err, "etcd put err")
	//		}
	//	} else {
	//		prometheusPort, err = datautil.GetElemByIndex(config.API.Prometheus.Ports, index)
	//		if err != nil {
	//			return err
	//		}
	//		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", prometheusPort))
	//		if err != nil {
	//			return errs.WrapMsg(err, "listen err", "addr", fmt.Sprintf(":%d", prometheusPort))
	//		}
	//	}
	//
	//	go func() {
	//		if err := prommetrics.ApiInit(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//			netErr = errs.WrapMsg(err, fmt.Sprintf("api prometheus start err: %d", prometheusPort))
	//			netDone <- struct{}{}
	//		}
	//	}()
	//
	//}
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancelCause(ctx)
	go func() {
		wg.Add(1)
		httpServer := &http.Server{
			Handler: router,
			Addr:    net.JoinHostPort(network.GetListenIP(config.API.Api.ListenIP), strconv.Itoa(apiPort)),
		}
		log.CInfo(ctx, "api server is init", "runtimeEnv", runtimeenv.RuntimeEnvironment(), "address", httpServer.Addr, "apiPort", apiPort)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			if err := httpServer.Shutdown(context.Background()); err != nil {
				log.ZWarn(ctx, "api server shutdown err", err)
			}
		}()
		err := httpServer.ListenAndServe()
		if err == nil {
			err = errors.New("api done")
		}
		cancel(err)
	}()

	//if config.Discovery.Enable == conf.ETCD {
	//	cm := disetcd.NewConfigManager(client.(*etcd.SvcDiscoveryRegistryImpl).GetClient(), config.GetConfigNames())
	//	cm.Watch(ctx)
	//}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case val := <-sigs:
		log.ZDebug(ctx, "recv exit", "signal", val.String())
		cancel(fmt.Errorf("signal %s", val.String()))
	case <-ctx.Done():
	}
	exitCause := context.Cause(ctx)
	log.ZWarn(ctx, "api server exit", exitCause)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	timer := time.NewTimer(time.Second * 15)
	defer timer.Stop()
	select {
	case <-timer.C:
		log.ZWarn(ctx, "api server graceful stop timeout", nil)
	case <-done:
		log.ZDebug(ctx, "api server graceful stop done")
	}
	return exitCause
}
