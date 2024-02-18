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
	util "github.com/openimsdk/open-im-server/v3/pkg/util/genutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/OpenIMSDK/tools/errs"

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
	apiCmd.AddPrometheusPortFlag()
	apiCmd.AddApi(run)
	if err := apiCmd.Execute(); err != nil {
		util.ExitWithError(err)
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
	var (
		netDone = make(chan struct{}, 1)
		netErr  error
	)
	router := api.NewGinRouter(client, rdb)
	if config.Config.Prometheus.Enable {
		go func() {
			p := ginprom.NewPrometheus("app", prommetrics.GetGinCusMetrics("Api"))
			p.SetListenAddress(fmt.Sprintf(":%d", proPort))
			if err = p.Use(router); err != nil && err != http.ErrServerClosed {
				netErr = errs.Wrap(err, fmt.Sprintf("prometheus start err: %d", proPort))
				netDone <- struct{}{}
			}
		}()

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
			netErr = errs.Wrap(err, fmt.Sprintf("api start err: %s", server.Addr))
			netDone <- struct{}{}

		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	select {
	case <-sigs:
		util.SIGUSR1Exit()
		err := server.Shutdown(ctx)
		if err != nil {
			return errs.Wrap(err, "shutdown err")
		}
	case <-netDone:
		close(netDone)
		return netErr
	}
	return nil
}
