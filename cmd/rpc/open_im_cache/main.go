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
	rpcCache "Open_IM/internal/rpc/cache"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"

	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImCachePort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.CachePrometheusPort[0], "cachePrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start cache rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := rpcCache.NewCacheServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
