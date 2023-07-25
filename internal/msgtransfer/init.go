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
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/tools/config"
	openKeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mw"
)

type MsgTransfer struct {
	persistentCH   *PersistentConsumerHandler         // 聊天记录持久化到mysql的消费者 订阅的topic: ws2ms_chat
	historyCH      *OnlineHistoryRedisConsumerHandler // 这个消费者聚合消息, 订阅的topic：ws2ms_chat, 修改通知发往msg_to_modify topic, 消息存入redis后Incr Redis, 再发消息到ms2pschat topic推送， 发消息到msg_to_mongo topic持久化
	historyMongoCH *OnlineHistoryMongoConsumerHandler // mongoDB批量插入, 成功后删除redis中消息，以及处理删除通知消息删除的 订阅的topic: msg_to_mongo
	// modifyCH       *ModifyMsgConsumerHandler          // 负责消费修改消息通知的consumer, 订阅的topic: msg_to_modify
}

func StartTransfer(prometheusPort int) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&relationTb.ChatLogModel{}); err != nil {
		fmt.Printf("gorm: AutoMigrate ChatLogModel err: %v\n", err)
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}
	if err := mongo.CreateMsgIndex(); err != nil {
		return err
	}
	client, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithRoundRobin(), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.Username,
			config.Config.Zookeeper.Password), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		return err
	}
	if err := client.CreateRpcRootNodes(config.Config.GetServiceNames()); err != nil {
		return err
	}
	client.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	msgModel := cache.NewMsgCacheModel(rdb)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase())
	msgMysModel := relation.NewChatLogGorm(db)
	chatLogDatabase := controller.NewChatLogDatabase(msgMysModel)
	msgDatabase := controller.NewCommonMsgDatabase(msgDocModel, msgModel)
	conversationRpcClient := rpcclient.NewConversationRpcClient(client)
	groupRpcClient := rpcclient.NewGroupRpcClient(client)
	msgTransfer := NewMsgTransfer(chatLogDatabase, msgDatabase, &conversationRpcClient, &groupRpcClient)
	msgTransfer.initPrometheus()
	return msgTransfer.Start(prometheusPort)
}

func NewMsgTransfer(chatLogDatabase controller.ChatLogDatabase,
	msgDatabase controller.CommonMsgDatabase,
	conversationRpcClient *rpcclient.ConversationRpcClient, groupRpcClient *rpcclient.GroupRpcClient,
) *MsgTransfer {
	return &MsgTransfer{
		persistentCH: NewPersistentConsumerHandler(chatLogDatabase), historyCH: NewOnlineHistoryRedisConsumerHandler(msgDatabase, conversationRpcClient, groupRpcClient),
		historyMongoCH: NewOnlineHistoryMongoConsumerHandler(msgDatabase),
	}
}

func (m *MsgTransfer) initPrometheus() {
	prome.NewSeqGetSuccessCounter()
	prome.NewSeqGetFailedCounter()
	prome.NewSeqSetSuccessCounter()
	prome.NewSeqSetFailedCounter()
	prome.NewMsgInsertRedisSuccessCounter()
	prome.NewMsgInsertRedisFailedCounter()
	prome.NewMsgInsertMongoSuccessCounter()
	prome.NewMsgInsertMongoFailedCounter()
}

func (m *MsgTransfer) Start(prometheusPort int) error {
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("start msg transfer", "prometheusPort:", prometheusPort)
	if config.Config.ChatPersistenceMysql {
		// go m.persistentCH.persistentConsumerGroup.RegisterHandleAndConsumer(m.persistentCH)
	} else {
		fmt.Println("msg transfer not start mysql consumer")
	}
	go m.historyCH.historyConsumerGroup.RegisterHandleAndConsumer(m.historyCH)
	go m.historyMongoCH.historyConsumerGroup.RegisterHandleAndConsumer(m.historyMongoCH)
	// go m.modifyCH.modifyMsgConsumerGroup.RegisterHandleAndConsumer(m.modifyCH)
	err := prome.StartPrometheusSrv(prometheusPort)
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}
