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
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/http"
	"sync"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/tools/mw"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type MsgTransfer struct {
	historyCH      *OnlineHistoryRedisConsumerHandler // 这个消费者聚合消息, 订阅的topic：ws2ms_chat, 修改通知发往msg_to_modify topic, 消息存入redis后Incr Redis, 再发消息到ms2pschat topic推送， 发消息到msg_to_mongo topic持久化
	historyMongoCH *OnlineHistoryMongoConsumerHandler // mongoDB批量插入, 成功后删除redis中消息，以及处理删除通知消息删除的 订阅的topic: msg_to_mongo
	// modifyCH       *ModifyMsgConsumerHandler          // 负责消费修改消息通知的consumer, 订阅的topic: msg_to_modify
}

func StartTransfer(prometheusPort int) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}

	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}

	if err = mongo.CreateMsgIndex(); err != nil {
		return err
	}
	client, err := kdisc.NewDiscoveryRegister(config.Config.Envs.Discovery)
	/*
		client, err := openkeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
			openkeeper.WithFreq(time.Hour), openkeeper.WithRoundRobin(), openkeeper.WithUserNameAndPassword(config.Config.Zookeeper.Username,
				config.Config.Zookeeper.Password), openkeeper.WithTimeout(10), openkeeper.WithLogger(log.NewZkLogger()))*/
	if err != nil {
		return err
	}

	if err := client.CreateRpcRootNodes(config.Config.GetServiceNames()); err != nil {
		return err
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	msgModel := cache.NewMsgCacheModel(rdb)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase())
	msgDatabase, err := controller.NewCommonMsgDatabase(msgDocModel, msgModel)
	if err != nil {
		return err
	}
	conversationRpcClient := rpcclient.NewConversationRpcClient(client)
	groupRpcClient := rpcclient.NewGroupRpcClient(client)
	msgTransfer, err := NewMsgTransfer(msgDatabase, &conversationRpcClient, &groupRpcClient)
	if err != nil {
		return err
	}
	return msgTransfer.Start(prometheusPort)
}

func NewMsgTransfer(msgDatabase controller.CommonMsgDatabase, conversationRpcClient *rpcclient.ConversationRpcClient, groupRpcClient *rpcclient.GroupRpcClient) (*MsgTransfer, error) {
	historyCH, err := NewOnlineHistoryRedisConsumerHandler(msgDatabase, conversationRpcClient, groupRpcClient)
	if err != nil {
		return nil, err
	}
	historyMongoCH, err := NewOnlineHistoryMongoConsumerHandler(msgDatabase)
	if err != nil {
		return nil, err
	}

	return &MsgTransfer{
		historyCH:      historyCH,
		historyMongoCH: historyMongoCH,
	}, nil
}

func (m *MsgTransfer) Start(prometheusPort int) error {
	ctx := context.Background()
	fmt.Println("start msg transfer", "prometheusPort:", prometheusPort)
	if prometheusPort <= 0 {
		return errs.Wrap(errors.New("prometheusPort not correct"))
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		m.historyCH.historyConsumerGroup.RegisterHandleAndConsumer(ctx, m.historyCH)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		m.historyMongoCH.historyConsumerGroup.RegisterHandleAndConsumer(ctx, m.historyMongoCH)
	}()

	if config.Config.Prometheus.Enable {
		go func() {
			proreg := prometheus.NewRegistry()
			proreg.MustRegister(
				collectors.NewGoCollector(),
			)
			proreg.MustRegister(prommetrics.GetGrpcCusMetrics("Transfer")...)
			http.Handle("/metrics", promhttp.HandlerFor(proreg, promhttp.HandlerOpts{Registry: proreg}))
			err := http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil)
			if err != nil && err != http.ErrServerClosed {
				panic(err)
			}
		}()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs

	// graceful close kafka client.
	go m.historyCH.historyConsumerGroup.Close()
	go m.historyMongoCH.historyConsumerGroup.Close()

	done := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.ZInfo(context.Background(), "msgtrasfer exit successfully")
	case <-time.After(15 * time.Second):
		log.ZError(context.Background(), "msgtransfer force to exit, timeout 15s", nil)
	}

	return nil
}
