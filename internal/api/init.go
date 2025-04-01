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
	"syscall"
	"time"

	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discovery"
	disetcd "github.com/openimsdk/open-im-server/v3/pkg/common/discovery/etcd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/system/program"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/jsonutil"
	"github.com/openimsdk/tools/utils/network"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	*conf.AllConfig

	RuntimeEnv string
	ConfigPath string
}

func Start(ctx context.Context, index int, config *Config) error {
	apiPort, err := datautil.GetElemByIndex(config.API.Api.Ports, index)
	if err != nil {
		return err
	}

	config.RuntimeEnv = runtimeenv.RuntimeEnvironment()

	client, err := kdisc.NewDiscoveryRegister(&config.Discovery, config.RuntimeEnv, []string{
		config.Discovery.RpcService.MessageGateway,
	})
	if err != nil {
		return errs.WrapMsg(err, "failed to register discovery service")
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))

	var (
		netDone        = make(chan struct{}, 1)
		netErr         error
		prometheusPort int
	)

	registerIP, err := network.GetRpcRegisterIP("")
	if err != nil {
		return err
	}

	getAutoPort := func() (net.Listener, int, error) {
		registerAddr := net.JoinHostPort(registerIP, "0")
		listener, err := net.Listen("tcp", registerAddr)
		if err != nil {
			return nil, 0, errs.WrapMsg(err, "listen err", "registerAddr", registerAddr)
		}
		_, portStr, _ := net.SplitHostPort(listener.Addr().String())
		port, _ := strconv.Atoi(portStr)
		return listener, port, nil
	}

	if config.API.Prometheus.AutoSetPorts && config.Discovery.Enable != conf.ETCD {
		return errs.New("only etcd support autoSetPorts", "RegisterName", "api").Wrap()
	}

	router, err := newGinRouter(ctx, client, config)
	if err != nil {
		return err
	}
	if config.API.Prometheus.Enable {
		var (
			listener net.Listener
		)

		if config.API.Prometheus.AutoSetPorts {
			listener, prometheusPort, err = getAutoPort()
			if err != nil {
				return err
			}

			etcdClient := client.(*etcd.SvcDiscoveryRegistryImpl).GetClient()

			_, err = etcdClient.Put(ctx, prommetrics.BuildDiscoveryKey(prommetrics.APIKeyName), jsonutil.StructToJsonString(prommetrics.BuildDefaultTarget(registerIP, prometheusPort)))
			if err != nil {
				return errs.WrapMsg(err, "etcd put err")
			}
		} else {
			prometheusPort, err = datautil.GetElemByIndex(config.API.Prometheus.Ports, index)
			if err != nil {
				return err
			}
			listener, err = net.Listen("tcp", fmt.Sprintf(":%d", prometheusPort))
			if err != nil {
				return errs.WrapMsg(err, "listen err", "addr", fmt.Sprintf(":%d", prometheusPort))
			}
		}

		go func() {
			if err := prommetrics.ApiInit(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
				netErr = errs.WrapMsg(err, fmt.Sprintf("api prometheus start err: %d", prometheusPort))
				netDone <- struct{}{}
			}
		}()

	}
	address := net.JoinHostPort(network.GetListenIP(config.API.Api.ListenIP), strconv.Itoa(apiPort))

	server := http.Server{Addr: address, Handler: router}
	log.CInfo(ctx, "API server is initializing", "runtimeEnv", config.RuntimeEnv, "address", address, "apiPort", apiPort, "prometheusPort", prometheusPort)
	go func() {
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			netErr = errs.WrapMsg(err, fmt.Sprintf("api start err: %s", server.Addr))
			netDone <- struct{}{}
		}
	}()

	//if config.Discovery.Enable == conf.ETCD {
	//	cm := disetcd.NewConfigManager(client.(*etcd.SvcDiscoveryRegistryImpl).GetClient(), config.GetConfigNames())
	//	cm.Watch(ctx)
	//}
	//sigs := make(chan os.Signal, 1)
	//signal.Notify(sigs, syscall.SIGTERM)
	//select {
	//case val := <-sigs:
	//	log.ZDebug(ctx, "recv exit", "signal", val.String())
	//	cancel(fmt.Errorf("signal %s", val.String()))
	//case <-ctx.Done():
	//}
	<-apiCtx.Done()
	exitCause := context.Cause(apiCtx)
	log.ZWarn(ctx, "api server exit", exitCause)
	timer := time.NewTimer(time.Second * 15)
	defer timer.Stop()
	select {
	case <-sigs:
		program.SIGTERMExit()
		if err := shutdown(); err != nil {
			return err
		}
	case <-netDone:
		close(netDone)
		return netErr
	}
	return nil
}
