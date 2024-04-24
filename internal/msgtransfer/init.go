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
	"fmt"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/utils/datautil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/system/program"
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
}

type Config struct {
	MsgTransfer     config.MsgTransfer
	RedisConfig     config.Redis
	MongodbConfig   config.Mongo
	KafkaConfig     config.Kafka
	ZookeeperConfig config.ZooKeeper
	Share           config.Share
	WebhooksConfig  config.Webhooks
}

func Start(ctx context.Context, index int, config *Config) error {
	log.CInfo(ctx, "MSG-TRANSFER server is initializing", "prometheusPorts",
		config.MsgTransfer.Prometheus.Ports, "index", index)
	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}
	client, err := kdisc.NewDiscoveryRegister(&config.ZookeeperConfig, &config.Share)
	if err != nil {
		return err
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	//todo MsgCacheTimeout
	msgModel := cache.NewMsgCache(rdb, config.RedisConfig.EnablePipeline)
	seqModel := cache.NewSeqCache(rdb)
	msgDocModel, err := mgo.NewMsgMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	msgDatabase, err := controller.NewCommonMsgDatabase(msgDocModel, msgModel, seqModel, &config.KafkaConfig)
	if err != nil {
		return err
	}
	conversationRpcClient := rpcclient.NewConversationRpcClient(client, config.Share.RpcRegisterName.Conversation)
	groupRpcClient := rpcclient.NewGroupRpcClient(client, config.Share.RpcRegisterName.Group)
	msgTransfer, err := NewMsgTransfer(&config.KafkaConfig, msgDatabase, &conversationRpcClient, &groupRpcClient)
	if err != nil {
		return err
	}
	return msgTransfer.Start(index, config)
}

func NewMsgTransfer(kafkaConf *config.Kafka, msgDatabase controller.CommonMsgDatabase,
	conversationRpcClient *rpcclient.ConversationRpcClient, groupRpcClient *rpcclient.GroupRpcClient) (*MsgTransfer, error) {
	historyCH, err := NewOnlineHistoryRedisConsumerHandler(kafkaConf, msgDatabase, conversationRpcClient, groupRpcClient)
	if err != nil {
		return nil, err
	}
	historyMongoCH, err := NewOnlineHistoryMongoConsumerHandler(kafkaConf, msgDatabase)
	if err != nil {
		return nil, err
	}

	return &MsgTransfer{
		historyCH:      historyCH,
		historyMongoCH: historyMongoCH,
	}, nil
}

func (m *MsgTransfer) Start(index int, config *Config) error {
	prometheusPort, err := datautil.GetElemByIndex(config.MsgTransfer.Prometheus.Ports, index)
	if err != nil {
		return err
	}
	m.ctx, m.cancel = context.WithCancel(context.Background())

	var (
		netDone = make(chan struct{}, 1)
		netErr  error
	)

	go m.historyCH.historyConsumerGroup.RegisterHandleAndConsumer(m.ctx, m.historyCH)
	go m.historyMongoCH.historyConsumerGroup.RegisterHandleAndConsumer(m.ctx, m.historyMongoCH)

	if config.MsgTransfer.Prometheus.Enable {
		go func() {
			proreg := prometheus.NewRegistry()
			proreg.MustRegister(
				collectors.NewGoCollector(),
			)
			proreg.MustRegister(prommetrics.GetGrpcCusMetrics("Transfer", &config.Share)...)
			http.Handle("/metrics", promhttp.HandlerFor(proreg, promhttp.HandlerOpts{Registry: proreg}))
			err := http.ListenAndServe(fmt.Sprintf(":%d", prometheusPort), nil)
			if err != nil && err != http.ErrServerClosed {
				netErr = errs.WrapMsg(err, "prometheus start error", "prometheusPort", prometheusPort)
				netDone <- struct{}{}
			}
		}()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		program.SIGTERMExit()
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
