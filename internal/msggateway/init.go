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
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpccache"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/runtimeenv"

	"github.com/openimsdk/tools/log"
)

type Config struct {
	MsgGateway     config.MsgGateway
	Share          config.Share
	RedisConfig    config.Redis
	WebhooksConfig config.Webhooks
	Discovery      config.Discovery
}

// Start run ws server.
func Start(ctx context.Context, index int, conf *Config) error {
	log.CInfo(ctx, "MSG-GATEWAY server is initializing", "runtimeEnv", runtimeenv.RuntimeEnvironment(),
		"rpcPorts", conf.MsgGateway.RPC.Ports,
		"wsPort", conf.MsgGateway.LongConnSvr.Ports, "prometheusPorts", conf.MsgGateway.Prometheus.Ports)
	wsPort, err := datautil.GetElemByIndex(conf.MsgGateway.LongConnSvr.Ports, index)
	if err != nil {
		return err
	}

	rdb, err := redisutil.NewRedisClient(ctx, conf.RedisConfig.Build())
	if err != nil {
		return err
	}
	longServer := NewWsServer(
		conf,
		WithPort(wsPort),
		WithMaxConnNum(int64(conf.MsgGateway.LongConnSvr.WebsocketMaxConnNum)),
		WithHandshakeTimeout(time.Duration(conf.MsgGateway.LongConnSvr.WebsocketTimeout)*time.Second),
		WithMessageMaxMsgLength(conf.MsgGateway.LongConnSvr.WebsocketMaxMsgLen),
	)

	hubServer := NewServer(longServer, conf, func(srv *Server) error {
		var err error
		longServer.online, err = rpccache.NewOnlineCache(srv.userClient, nil, rdb, false, longServer.subscriberUserOnlineStatusChanges)
		return err
	})

	go longServer.ChangeOnlineStatus(4)

	netDone := make(chan error)
	go func() {
		err = hubServer.Start(ctx, index, conf)
		netDone <- err
	}()
	return hubServer.LongConnServer.Run(netDone)
}
