// Copyright © 2024 OpenIM. All rights reserved.
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

package rtc

import (
	"context"
	"time"

	lksdk "github.com/livekit/server-sdk-go/v2"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/rtc"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/discovery"
	"google.golang.org/grpc"
)

// Config aggregates all configuration needed by the RTC service.
type Config struct {
	RpcConfig     config.Rtc
	MongodbConfig config.Mongo
	Share         config.Share
	Discovery     config.Discovery
}

type rtcServer struct {
	rtc.UnimplementedRtcServiceServer
	config         *Config
	db             controller.RtcDatabase
	roomClient     *lksdk.RoomServiceClient
	msgClient      *rpcli.MsgClient
	userClient     *rpcli.UserClient
	relationClient *rpcli.RelationClient
	tokenExpiry    time.Duration
}

// Start initialises the RTC gRPC service and registers it with the gRPC server.
func Start(ctx context.Context, cfg *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgocli, err := mongoutil.NewMongoDB(ctx, cfg.MongodbConfig.Build())
	if err != nil {
		return err
	}

	signalDB, err := mgo.NewSignalMongo(mgocli.GetDB())
	if err != nil {
		return err
	}

	msgConn, err := client.GetConn(ctx, cfg.Share.RpcRegisterName.Msg)
	if err != nil {
		return err
	}

	userConn, err := client.GetConn(ctx, cfg.Share.RpcRegisterName.User)
	if err != nil {
		return err
	}

	friendConn, err := client.GetConn(ctx, cfg.Share.RpcRegisterName.Friend)
	if err != nil {
		return err
	}

	lk := cfg.RpcConfig.LiveKit
	roomClient := lksdk.NewRoomServiceClient(lk.InternalAddress, lk.APIKey, lk.APISecret)

	tokenExpiry := time.Duration(lk.TokenExpiry) * time.Second
	if tokenExpiry <= 0 {
		tokenExpiry = time.Hour
	}

	s := &rtcServer{
		config:         cfg,
		db:             controller.NewRtcDatabase(signalDB),
		roomClient:     roomClient,
		msgClient:      rpcli.NewMsgClient(msgConn),
		userClient:     rpcli.NewUserClient(userConn),
		relationClient: rpcli.NewRelationClient(friendConn),
		tokenExpiry:    tokenExpiry,
	}

	rtc.RegisterRtcServiceServer(server, s)
	return nil
}
