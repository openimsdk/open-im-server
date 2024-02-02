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

package main

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/openimsdk/open-im-server/v3/internal/api"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	ginprom "github.com/openimsdk/open-im-server/v3/pkg/common/ginprometheus"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
)

func main() {
	apiCmd := cmd.NewApiCmd()
	apiCmd.AddPortFlag()
	apiCmd.AddApi(run)
	if err := apiCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "\n\nexit -1: \n%+v\n\n", err)
		os.Exit(-1)
	}
}

func run(port int, proPort int) error {
	if port == 0 || proPort == 0 {
		err := "port or proPort is empty:" + strconv.Itoa(port) + "," + strconv.Itoa(proPort)
		return errs.Wrap(fmt.Errorf(err))
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}

	var client discoveryregistry.SvcDiscoveryRegistry

	// Determine whether zk is passed according to whether it is a clustered deployment
	client, err = kdisc.NewDiscoveryRegister(config.Config.Envs.Discovery)
	if err != nil {
		return errs.Wrap(err, "register discovery err")
	}

	if err = client.CreateRpcRootNodes(config.Config.GetServiceNames()); err != nil {
		return errs.Wrap(err, "create rpc root nodes error")
	}

	if err = client.RegisterConf2Registry(constant.OpenIMCommonConfigKey, config.Config.EncodeConfig()); err != nil {
		return err
	}

	router := api.NewGinRouter(client, rdb)
	if config.Config.Prometheus.Enable {
		p := ginprom.NewPrometheus("app", prommetrics.GetGinCusMetrics("Api"))
		p.SetListenAddress(fmt.Sprintf(":%d", proPort))
		p.Use(router)
	}

	var address string
	if config.Config.Api.ListenIP != "" {
		address = net.JoinHostPort(config.Config.Api.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}

	server := http.Server{Addr: address, Handler: router}
	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// graceful shutdown operation.
	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
