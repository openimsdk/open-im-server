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

package msggateway

import (
	"fmt"
	"time"

	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

// RunWsAndServer run ws server
func RunWsAndServer(rpcPort, wsPort, prometheusPort int) error {
	fmt.Println(
		"start rpc/msg_gateway server, port: ",
		rpcPort,
		wsPort,
		prometheusPort,
		", OpenIM version: ",
		config.Version,
	)
	longServer, err := NewWsServer(
		WithPort(wsPort),
		WithMaxConnNum(int64(config.Config.LongConnSvr.WebsocketMaxConnNum)),
		WithHandshakeTimeout(time.Duration(config.Config.LongConnSvr.WebsocketTimeout)*time.Second),
		WithMessageMaxMsgLength(config.Config.LongConnSvr.WebsocketMaxMsgLen))
	if err != nil {
		return err
	}
	hubServer := NewServer(rpcPort, prometheusPort, longServer)
	go func() {
		err := hubServer.Start()
		if err != nil {
			panic(utils.Wrap1(err))
		}
	}()
	return hubServer.LongConnServer.Run()
}
