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
	"Open_IM/internal/push/logic"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImPushPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessageTransferPrometheusPort[0], "PushrometheusPort default listen port")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	log.NewPrivateLog(constant.LogFileName)
	fmt.Println("start push rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	logic.Init(*rpcPort)
	logic.Run(*prometheusPort)
	wg.Wait()
}
