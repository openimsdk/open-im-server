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
	"os"
	"os/signal"
	"syscall"

	disetcd "github.com/openimsdk/open-im-server/v3/pkg/common/discovery/etcd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/mqbuild"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/mq"
	"github.com/openimsdk/tools/utils/runtimeenv"

	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/system/program"
	"google.golang.org/grpc"
)

type MsgTransfer struct {
	historyConsumer      mq.Consumer
	historyMongoConsumer mq.Consumer
	// This consumer aggregated messages, subscribed to the topic:toRedis,
	//  the message is stored in redis, Incr Redis, and then the message is sent to toPush topic for push,
	// and the message is sent to toMongo topic for persistence
	historyHandler *OnlineHistoryRedisConsumerHandler
	//This consumer handle message to mongo
	historyMongoHandler *OnlineHistoryMongoConsumerHandler
	ctx                 context.Context
	cancel              context.CancelFunc
}

type Config struct {
	MsgTransfer    conf.MsgTransfer
	RedisConfig    conf.Redis
	MongodbConfig  conf.Mongo
	KafkaConfig    conf.Kafka
	Share          conf.Share
	WebhooksConfig conf.Webhooks
	Discovery      conf.Discovery
	Index          conf.Index
}

func Start(ctx context.Context, config *Config, client discovery.Conn, server grpc.ServiceRegistrar) error {
	builder := mqbuild.NewBuilder(&config.KafkaConfig)

	log.CInfo(ctx, "MSG-TRANSFER server is initializing", "runTimeEnv", runtimeenv.RuntimeEnvironment(), "prometheusPorts",
		config.MsgTransfer.Prometheus.Ports, "index", config.Index)
	dbb := dbbuild.NewBuilder(&config.MongodbConfig, &config.RedisConfig)
	mgocli, err := dbb.Mongo(ctx)
	if err != nil {
		return err
	}
	rdb, err := dbb.Redis(ctx)
	if err != nil {
		return err
	}

	if config.Discovery.Enable == conf.ETCD {
		cm := disetcd.NewConfigManager(client.(*etcd.SvcDiscoveryRegistryImpl).GetClient(), []string{
			config.MsgTransfer.GetConfigFileName(),
			config.RedisConfig.GetConfigFileName(),
			config.MongodbConfig.GetConfigFileName(),
			config.KafkaConfig.GetConfigFileName(),
			config.Share.GetConfigFileName(),
			config.WebhooksConfig.GetConfigFileName(),
			config.Discovery.GetConfigFileName(),
			conf.LogConfigFileName,
		})
		cm.Watch(ctx)
	}
	mongoProducer, err := builder.GetTopicProducer(ctx, config.KafkaConfig.ToMongoTopic)
	if err != nil {
		return err
	}
	pushProducer, err := builder.GetTopicProducer(ctx, config.KafkaConfig.ToPushTopic)
	if err != nil {
		return err
	}
	msgDocModel, err := mgo.NewMsgMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	msgModel := redis.NewMsgCache(rdb, msgDocModel)
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
	msgTransferDatabase, err := controller.NewMsgTransferDatabase(msgDocModel, msgModel, seqUserCache, seqConversationCache, mongoProducer, pushProducer)
	if err != nil {
		return err
	}
	historyConsumer, err := builder.GetTopicConsumer(ctx, config.KafkaConfig.ToRedisTopic)
	if err != nil {
		return err
	}
	historyMongoConsumer, err := builder.GetTopicConsumer(ctx, config.KafkaConfig.ToMongoTopic)
	if err != nil {
		return err
	}
	historyHandler, err := NewOnlineHistoryRedisConsumerHandler(ctx, client, config, msgTransferDatabase)
	if err != nil {
		return err
	}
	historyMongoHandler := NewOnlineHistoryMongoConsumerHandler(msgTransferDatabase)

	msgTransfer := &MsgTransfer{
		historyConsumer:      historyConsumer,
		historyMongoConsumer: historyMongoConsumer,
		historyHandler:       historyHandler,
		historyMongoHandler:  historyMongoHandler,
	}

	return msgTransfer.Start(ctx, int(config.Index), config, client)
}

func (m *MsgTransfer) Start(ctx context.Context, index int, config *Config, client discovery.Conn) error {
	m.ctx, m.cancel = context.WithCancel(context.Background())
	var (
		netDone = make(chan struct{}, 1)
		netErr  error
	)

	go func() {
		for {
			if err := m.historyConsumer.Subscribe(m.ctx, m.historyHandler.HandlerRedisMessage); err != nil {
				log.ZError(m.ctx, "historyConsumer err", err)
				return
			}
		}
	}()

	go func() {
		fn := func(ctx context.Context, key string, value []byte) error {
			m.historyMongoHandler.HandleChatWs2Mongo(ctx, key, value)
			return nil
		}
		for {
			if err := m.historyMongoConsumer.Subscribe(m.ctx, fn); err != nil {
				log.ZError(m.ctx, "historyMongoConsumer err", err)
				return
			}
		}
	}()

	go m.historyHandler.HandleUserHasReadSeqMessages(m.ctx)
	err := m.historyHandler.redisMessageBatches.Start()
	if err != nil {
		return err
	}

	// todo
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	select {
	case <-sigs:
		program.SIGTERMExit()
		// graceful close kafka client.
		_ = m.historyConsumer.Close()
		_ = m.historyMongoConsumer.Close()
		m.cancel()
		m.historyHandler.redisMessageBatches.Close()
		m.historyHandler.Close()
		return nil
	case <-netDone:
		_ = m.historyConsumer.Close()
		_ = m.historyMongoConsumer.Close()
		m.cancel()
		m.historyHandler.redisMessageBatches.Close()
		m.historyHandler.Close()
		close(netDone)
		return netErr
	}
}
