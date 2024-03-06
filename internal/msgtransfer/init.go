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

package msgtransfer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mw"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	util "github.com/openimsdk/open-im-server/v3/pkg/util/genutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MsgTransfer struct {
	// This consumer aggregated messages, subscribed to the topic:ws2ms_chat,
	// the modification notification is sent to msg_to_modify topic, the message is stored in redis, Incr Redis,
	// and then the message is sent to ms2pschat topic for push, and the message is sent to msg_to_mongo topic for persistence
	historyCH      *OnlineHistoryRedisConsumerHandler
	historyMongoCH *OnlineHistoryMongoConsumerHandler
	// mongoDB batch insert, delete messages in redis after success,
	// and handle the deletion notification message deleted subscriptions topic: msg_to_mongo
	ctx    context.Context
	cancel context.CancelFunc
	config *config.GlobalConfig
}

func StartTransfer(config *config.GlobalConfig, prometheusPort int) error {
	rdb, err := cache.NewRedis(config)
	if err != nil {
		return err
	}

	mongo, err := unrelation.NewMongo(config)
	if err != nil {
		return err
	}

	if err = mongo.CreateMsgIndex(); err != nil {
		return err
	}
	client, err := kdisc.NewDiscoveryRegister(config)
	if err != nil {
		return err
	}

	if err := client.CreateRpcRootNodes(config.GetServiceNames()); err != nil {
		return err
	}

	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	msgModel := cache.NewMsgCacheModel(rdb, config)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase(config.Mongo.Database))
	msgDatabase, err := controller.NewCommonMsgDatabase(msgDocModel, msgModel, config)
	if err != nil {
		return err
	}
	conversationRpcClient := rpcclient.NewConversationRpcClient(client, config)
	groupRpcClient := rpcclient.NewGroupRpcClient(client, config)
	msgTransfer, err := NewMsgTransfer(config, msgDatabase, &conversationRpcClient, &groupRpcClient)
	if err != nil {
		return err
	}
	return msgTransfer.Start(prometheusPort, config)
}

func NewMsgTransfer(
	config *config.GlobalConfig,
	msgDatabase controller.CommonMsgDatabase,
	conversationRpcClient *rpcclient.ConversationRpcClient,
	groupRpcClient *rpcclient.GroupRpcClient,
) (*MsgTransfer, error) {
	historyCH, err := NewOnlineHistoryRedisConsumerHandler(config, msgDatabase, conversationRpcClient, groupRpcClient)
	if err != nil {
		return nil, err
	}
	historyMongoCH, err := NewOnlineHistoryMongoConsumerHandler(config, msgDatabase)
	if err != nil {
		return nil, err
	}

	return &MsgTransfer{
		historyCH:      historyCH,
		historyMongoCH: historyMongoCH,
		config:         config,
	}, nil
}

func (m *MsgTransfer) Start(prometheusPort int, config *config.GlobalConfig) error {
	fmt.Println("start msg transfer", "prometheusPort:", prometheusPort)
	if prometheusPort <= 0 {
		return errs.Wrap(errors.New("prometheusPort not correct"))
	}
	m.ctx, m.cancel = context.WithCancel(context.Background())

	var (
		netDone = make(chan struct{}, 1)
		netErr  error
	)

	go m.historyCH.historyConsumerGroup.RegisterHandleAndConsumer(m.ctx, m.historyCH)
	go m.historyMongoCH.historyConsumerGroup.RegisterHandleAndConsumer(m.ctx, m.historyMongoCH)

	if config.Prometheus.Enable {
		go func() {
			proreg := prometheus.NewRegistry()
			proreg.MustRegister(
				collectors.NewGoCollector(),
			)
			proreg.MustRegister(prommetrics.GetGrpcCusMetrics("Transfer", config)...)
			http.Handle("/metrics", promhttp.HandlerFor(proreg, promhttp.HandlerOpts{Registry: proreg}))
			err := http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil)
			if err != nil && err != http.ErrServerClosed {
				netErr = errs.Wrap(err, fmt.Sprintf("prometheus start err: %d", prometheusPort))
				netDone <- struct{}{}
			}
		}()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		util.SIGTERMExit()
		// graceful close kafka client.
		m.cancel()
		m.historyCH.historyConsumerGroup.Close()
		m.historyMongoCH.historyConsumerGroup.Close()
		return nil
	case <-netDone:
		m.cancel()
		m.historyCH.historyConsumerGroup.Close()
		m.historyMongoCH.historyConsumerGroup.Close()
		close(netDone)
		return netErr
	}
}
