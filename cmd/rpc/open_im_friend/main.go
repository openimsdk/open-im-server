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
	"Open_IM/internal/rpc/friend"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImFriendPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcFriendPort from cmd,default 12000 as port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.FriendPrometheusPort[0], "friendPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start friend rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := friend.NewFriendServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
