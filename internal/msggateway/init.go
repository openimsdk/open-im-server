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
	"context"
	"github.com/openimsdk/tools/log"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

// Start run ws server.
func Start(ctx context.Context, conf *config.GlobalConfig, rpcPort, wsPort, prometheusPort int) error {
	log.CInfo(ctx, "MSG-GATEWAY server is initializing", "rpcPort", rpcPort, "wsPort", wsPort,
		"prometheusPort", prometheusPort)
	longServer, err := NewWsServer(
		conf,
		WithPort(wsPort),
		WithMaxConnNum(int64(conf.LongConnSvr.WebsocketMaxConnNum)),
		WithHandshakeTimeout(time.Duration(conf.LongConnSvr.WebsocketTimeout)*time.Second),
		WithMessageMaxMsgLength(conf.LongConnSvr.WebsocketMaxMsgLen),
		WithWriteBufferSize(conf.LongConnSvr.WebsocketWriteBufferSize),
	)
	if err != nil {
		return err
	}

	hubServer := NewServer(rpcPort, prometheusPort, longServer, conf)
	netDone := make(chan error)
	go func() {
		err = hubServer.Start(ctx, conf)
		netDone <- err
	}()
	return hubServer.LongConnServer.Run(netDone)
}
