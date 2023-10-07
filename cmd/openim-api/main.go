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
	"strconv"
	"time"

	_ "net/http/pprof"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	openkeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/internal/api"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
)

func main() {
	apiCmd := cmd.NewApiCmd()
	apiCmd.AddPortFlag()
	apiCmd.AddApi(run)
	if err := apiCmd.Execute(); err != nil {
		panic(err.Error())
	}
}

func run(port int) error {
	if port == 0 {
		return fmt.Errorf("port is empty")
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	fmt.Println("api start init discov client")
	var client discoveryregistry.SvcDiscoveryRegistry
	client, err = openkeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openkeeper.WithFreq(time.Hour), openkeeper.WithUserNameAndPassword(
			config.Config.Zookeeper.Username,
			config.Config.Zookeeper.Password,
		), openkeeper.WithRoundRobin(), openkeeper.WithTimeout(10), openkeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		return err
	}
	if err = client.CreateRpcRootNodes(config.Config.GetServiceNames()); err != nil {
		return err
	}
	fmt.Println("api register public config to discov")
	if err = client.RegisterConf2Registry(constant.OpenIMCommonConfigKey, config.Config.EncodeConfig()); err != nil {
		return err
	}
	fmt.Println("api register public config to discov success")
	router := api.NewGinRouter(client, rdb)
	fmt.Println("api init router success")
	var address string
	if config.Config.Api.ListenIP != "" {
		address = net.JoinHostPort(config.Config.Api.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}
	fmt.Println("start api server, address: ", address, ", OpenIM version: ", config.Version)
	log.ZInfo(context.Background(), "start server success", "address", address, "version", config.Version)
	err = router.Run(address)
	if err != nil {
		log.ZError(context.Background(), "api run failed ", err, "address", address)

		return err
	}

	return nil
}
