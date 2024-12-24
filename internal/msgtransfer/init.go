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
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	disetcd "github.com/openimsdk/open-im-server/v3/pkg/common/discovery/etcd"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/utils/jsonutil"
	"github.com/openimsdk/tools/utils/network"

	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/runtimeenv"

	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	discRegister "github.com/openimsdk/open-im-server/v3/pkg/common/discovery"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mw"
	"github.com/openimsdk/tools/system/program"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MsgTransfer struct {
	// This consumer aggregated messages, subscribed to the topic:toRedis,
	//  the message is stored in redis, Incr Redis, and then the message is sent to toPush topic for push,
	// and the message is sent to toMongo topic for persistence
	historyCH *OnlineHistoryRedisConsumerHandler
	//This consumer handle message to mongo
	historyMongoCH *OnlineHistoryMongoConsumerHandler
	ctx            context.Context
	cancel         context.CancelFunc

	runTimeEnv string
}

type Config struct {
	MsgTransfer    conf.MsgTransfer
	RedisConfig    conf.Redis
	MongodbConfig  conf.Mongo
	KafkaConfig    conf.Kafka
	Share          conf.Share
	WebhooksConfig conf.Webhooks
	Discovery      conf.Discovery
}

func Start(ctx context.Context, index int, config *Config) error {
	runTimeEnv := runtimeenv.PrintRuntimeEnvironment()

	log.CInfo(ctx, "MSG-TRANSFER server is initializing", "runTimeEnv", runTimeEnv, "prometheusPorts",
		config.MsgTransfer.Prometheus.Ports, "index", index)

	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}
	client, err := discRegister.NewDiscoveryRegister(&config.Discovery, runTimeEnv)
	if err != nil {
		return err
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))

	if config.Discovery.Enable == conf.ETCD {
		cm := disetcd.NewConfigManager(client.(*etcd.SvcDiscoveryRegistryImpl).GetClient(), []string{
			config.MsgTransfer.GetConfigFileName(),
			config.RedisConfig.GetConfigFileName(),
			config.MongodbConfig.GetConfigFileName(),
			config.KafkaConfig.GetConfigFileName(),
			config.Share.GetConfigFileName(),
			config.WebhooksConfig.GetConfigFileName(),
			config.Discovery.GetConfigFileName(),
		})
		cm.Watch(ctx)
	}

	msgModel := redis.NewMsgCache(rdb)
	msgDocModel, err := mgo.NewMsgMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	seqConversation, err := mgo.NewSeqConversationMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	seqConversationCache := redis.NewSeqConversationCacheRedis(rdb, seqConversation)
	seqUser, err := mgo.NewSeqUserMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	seqUserCache := redis.NewSeqUserCacheRedis(rdb, seqUser)
	msgTransferDatabase, err := controller.NewMsgTransferDatabase(msgDocModel, msgModel, seqUserCache, seqConversationCache, &config.KafkaConfig)
	if err != nil {
		return err
	}
	historyCH, err := NewOnlineHistoryRedisConsumerHandler(ctx, client, config, msgTransferDatabase)
	if err != nil {
		return err
	}
	historyMongoCH, err := NewOnlineHistoryMongoConsumerHandler(&config.KafkaConfig, msgTransferDatabase)
	if err != nil {
		return err
	}

	msgTransfer := &MsgTransfer{
		historyCH:      historyCH,
		historyMongoCH: historyMongoCH,
		runTimeEnv:     runTimeEnv,
	}

	return msgTransfer.Start(index, config, client)
}

func (m *MsgTransfer) Start(index int, config *Config, client discovery.SvcDiscoveryRegistry) error {
	m.ctx, m.cancel = context.WithCancel(context.Background())
	var (
		netDone = make(chan struct{}, 1)
		netErr  error
	)

	go m.historyCH.historyConsumerGroup.RegisterHandleAndConsumer(m.ctx, m.historyCH)
	go m.historyMongoCH.historyConsumerGroup.RegisterHandleAndConsumer(m.ctx, m.historyMongoCH)
	go m.historyCH.HandleUserHasReadSeqMessages(m.ctx)
	err := m.historyCH.redisMessageBatches.Start()
	if err != nil {
		return err
	}

	registerIP, err := network.GetRpcRegisterIP("")
	if err != nil {
		return err
	}

	getAutoPort := func() (net.Listener, int, error) {
		registerAddr := net.JoinHostPort(registerIP, "0")
		listener, err := net.Listen("tcp", registerAddr)
		if err != nil {
			return nil, 0, errs.WrapMsg(err, "listen err", "registerAddr", registerAddr)
		}
		_, portStr, _ := net.SplitHostPort(listener.Addr().String())
		port, _ := strconv.Atoi(portStr)
		return listener, port, nil
	}

	if config.MsgTransfer.Prometheus.AutoSetPorts && config.Discovery.Enable != conf.ETCD {
		return errs.New("only etcd support autoSetPorts", "RegisterName", "api").Wrap()
	}

	if config.MsgTransfer.Prometheus.Enable {
		var (
			listener       net.Listener
			prometheusPort int
		)

		if config.MsgTransfer.Prometheus.AutoSetPorts {
			listener, prometheusPort, err = getAutoPort()
			if err != nil {
				return err
			}

			etcdClient := client.(*etcd.SvcDiscoveryRegistryImpl).GetClient()

			_, err = etcdClient.Put(context.TODO(), prommetrics.BuildDiscoveryKey(prommetrics.MessageTransferKeyName), jsonutil.StructToJsonString(prommetrics.BuildDefaultTarget(registerIP, prometheusPort)))
			if err != nil {
				return errs.WrapMsg(err, "etcd put err")
			}
		} else {
			prometheusPort, err = datautil.GetElemByIndex(config.MsgTransfer.Prometheus.Ports, index)
			if err != nil {
				return err
			}
			listener, err = net.Listen("tcp", fmt.Sprintf(":%d", prometheusPort))
			if err != nil {
				return errs.WrapMsg(err, "listen err", "addr", fmt.Sprintf(":%d", prometheusPort))
			}
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.ZPanic(m.ctx, "MsgTransfer Start Panic", r)
				}
			}()
			if err := prommetrics.TransferInit(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
		m.historyCH.redisMessageBatches.Close()
		m.historyCH.Close()
		m.historyCH.historyConsumerGroup.Close()
		m.historyMongoCH.historyConsumerGroup.Close()
		return nil
	case <-netDone:
		m.cancel()
		m.historyCH.redisMessageBatches.Close()
		m.historyCH.Close()
		m.historyCH.historyConsumerGroup.Close()
		m.historyMongoCH.historyConsumerGroup.Close()
		close(netDone)
		return netErr
	}
}
