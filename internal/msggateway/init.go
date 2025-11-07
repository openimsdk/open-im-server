package msggateway

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/rpccache"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/runtimeenv"
	"google.golang.org/grpc"

	"github.com/openimsdk/tools/log"
)

type Config struct {
	MsgGateway     config.MsgGateway
	Share          config.Share
	RedisConfig    config.Redis
	WebhooksConfig config.Webhooks
	Discovery      config.Discovery
	Index          config.Index
}

// Start run ws server.
func Start(ctx context.Context, conf *Config, client discovery.SvcDiscoveryRegistry, server grpc.ServiceRegistrar) error {
	log.CInfo(ctx, "MSG-GATEWAY server is initializing", "runtimeEnv", runtimeenv.RuntimeEnvironment(),
		"rpcPorts", conf.MsgGateway.RPC.Ports,
		"wsPort", conf.MsgGateway.LongConnSvr.Ports, "prometheusPorts", conf.MsgGateway.Prometheus.Ports)
	wsPort, err := datautil.GetElemByIndex(conf.MsgGateway.LongConnSvr.Ports, int(conf.Index))
	if err != nil {
		return err
	}

	dbb := dbbuild.NewBuilder(nil, &conf.RedisConfig)
	rdb, err := dbb.Redis(ctx)
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

	if err := hubServer.InitServer(ctx, conf, client, server); err != nil {
		return err
	}

	go longServer.ChangeOnlineStatus(4)

	return hubServer.LongConnServer.Run(ctx)
}
