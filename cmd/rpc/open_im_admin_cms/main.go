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
	rpcMessageCMS "Open_IM/internal/rpc/admin_cms"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImAdminCmsPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.AdminCmsPrometheusPort[0], "adminCMSPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start cms rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := rpcMessageCMS.NewAdminCMSServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
