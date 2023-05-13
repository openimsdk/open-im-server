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

package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"

	"Open_IM/pkg/statistics"
	"fmt"
	"sync"

	promePkg "Open_IM/pkg/common/prometheus"

	"github.com/go-playground/validator/v10"
)

var (
	rwLock              *sync.RWMutex
	validate            *validator.Validate
	ws                  WServer
	rpcSvr              RPCServer
	sendMsgAllCount     uint64
	sendMsgFailedCount  uint64
	sendMsgSuccessCount uint64
	userCount           uint64

	sendMsgAllCountLock sync.RWMutex
)

func Init(rpcPort, wsPort int) {
	rwLock = new(sync.RWMutex)
	validate = validator.New()
	statistics.NewStatistics(&sendMsgAllCount, config.Config.ModuleName.LongConnSvrName, fmt.Sprintf("%d second recv to msg_gateway sendMsgCount", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	statistics.NewStatistics(&userCount, config.Config.ModuleName.LongConnSvrName, fmt.Sprintf("%d second add user conn", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	ws.onInit(wsPort)
	rpcSvr.onInit(rpcPort)
	initPrometheus()
}

func Run(promethuesPort int) {
	go ws.run()
	go rpcSvr.run()
	go func() {
		err := promePkg.StartPromeSrv(promethuesPort)
		if err != nil {
			panic(err)
		}
	}()
}
