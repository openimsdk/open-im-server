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
	"net"
	_ "net/http/pprof"
	"strconv"

	ginProm "github.com/openimsdk/open-im-server/v3/pkg/common/ginprometheus"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/internal/api"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
)

func main() {
	apiCmd := cmd.NewApiCmd()
	apiCmd.AddPortFlag()
	apiCmd.AddApi(run)
	if err := apiCmd.Execute(); err != nil {
		log.ZError(context.Background(), "API command execution failed", err)
		panic(err.Error())
	}
}

func run(port int, proPort int) error {
	log.ZInfo(context.Background(), "Openim api port:", "port", port, "proPort", proPort)

	if port == 0 || proPort == 0 {
		err := "port or proPort is empty:" + strconv.Itoa(port) + "," + strconv.Itoa(proPort)
		log.ZError(context.Background(), err, nil)

		return fmt.Errorf(err)
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		log.ZError(context.Background(), "Failed to initialize Redis", err)

		return err
	}
	log.ZInfo(context.Background(), "api start init discov client")

	var client discoveryregistry.SvcDiscoveryRegistry

	// Determine whether zk is passed according to whether it is a clustered deployment
	client, err = kdisc.NewDiscoveryRegister(config.Config.Envs.Discovery)
	if err != nil {
		log.ZError(context.Background(), "Failed to initialize discovery register", err)

		return err
	}
	if err = client.CreateRpcRootNodes(config.Config.GetServiceNames()); err != nil {
		log.ZError(context.Background(), "Failed to create RPC root nodes", err)

		return err
	}
	log.ZInfo(context.Background(), "api register public config to discov")
	if err = client.RegisterConf2Registry(constant.OpenIMCommonConfigKey, config.Config.EncodeConfig()); err != nil {
		log.ZError(context.Background(), "Failed to register public config to discov", err)

		return err
	}
	log.ZInfo(context.Background(), "api register public config to discov success")
	router := api.NewGinRouter(client, rdb)
	//////////////////////////////
	if config.Config.Prometheus.Enable {
		p := ginProm.NewPrometheus("app", prommetrics.GetGinCusMetrics("Api"))
		p.SetListenAddress(fmt.Sprintf(":%d", proPort))
		p.Use(router)
	}
	/////////////////////////////////
	log.ZInfo(context.Background(), "api init router success")
	var address string
	if config.Config.Api.ListenIP != "" {
		address = net.JoinHostPort(config.Config.Api.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}
	log.ZInfo(context.Background(), "start api server", "address", address, "OpenIM version", config.Version)

	err = router.Run(address)
	if err != nil {
		log.ZError(context.Background(), "api run failed", err, "address", address)

		return err
	}

	return nil
}
