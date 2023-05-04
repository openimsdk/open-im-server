package msgtransfer

import (
	"fmt"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/openKeeper"
)

type MsgTransfer struct {
	persistentCH   *PersistentConsumerHandler         // 聊天记录持久化到mysql的消费者 订阅的topic: ws2ms_chat
	historyCH      *OnlineHistoryRedisConsumerHandler // 这个消费者聚合消息, 订阅的topic：ws2ms_chat, 修改通知发往msg_to_modify topic, 消息存入redis后Incr Redis, 再发消息到ms2pschat topic推送， 发消息到msg_to_mongo topic持久化
	historyMongoCH *OnlineHistoryMongoConsumerHandler // mongoDB批量插入, 成功后删除redis中消息，以及处理删除通知消息删除的 订阅的topic: msg_to_mongo
	modifyCH       *ModifyMsgConsumerHandler          // 负责消费修改消息通知的consumer, 订阅的topic: msg_to_modify
}

func StartTransfer(prometheusPort int) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&relationTb.ChatLogModel{}); err != nil {
		return err
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}
	client, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema, 10, config.Config.Zookeeper.UserName, config.Config.Zookeeper.Password)
	if err != nil {
		return err
	}
	cacheModel := cache.NewCacheModel(rdb)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase())
	extendMsgModel := unrelation.NewExtendMsgSetMongoDriver(mongo.GetDatabase())
	extendMsgCache := cache.NewExtendMsgSetCacheRedis(rdb, extendMsgModel, cache.GetDefaultOpt())
	chatLogDatabase := controller.NewChatLogDatabase(relation.NewChatLogGorm(db))
	extendMsgDatabase := controller.NewExtendMsgDatabase(extendMsgModel, extendMsgCache, tx.NewMongo(mongo.GetClient()))
	msgDatabase := controller.NewMsgDatabase(msgDocModel, cacheModel)
	notificationDatabase := controller.NewNotificationDatabase(msgDocModel, cacheModel)
	conversationRpcClient := rpcclient.NewConversationClient(client)

	msgTransfer := NewMsgTransfer(chatLogDatabase, extendMsgDatabase, msgDatabase, notificationDatabase, conversationRpcClient)
	msgTransfer.initPrometheus()
	return msgTransfer.Start(prometheusPort)
}

func NewMsgTransfer(chatLogDatabase controller.ChatLogDatabase,
	extendMsgDatabase controller.ExtendMsgDatabase, msgDatabase controller.MsgDatabase, notificationDatabase controller.MsgDatabase,
	conversationRpcClient *rpcclient.ConversationClient) *MsgTransfer {
	return &MsgTransfer{persistentCH: NewPersistentConsumerHandler(chatLogDatabase), historyCH: NewOnlineHistoryRedisConsumerHandler(msgDatabase, conversationRpcClient),
		historyMongoCH: NewOnlineHistoryMongoConsumerHandler(msgDatabase, notificationDatabase), modifyCH: NewModifyMsgConsumerHandler(extendMsgDatabase)}
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
	wg.Add(4)
	fmt.Println("start msg transfer", "prometheusPort:", prometheusPort)
	if config.Config.ChatPersistenceMysql {
		go m.persistentCH.persistentConsumerGroup.RegisterHandleAndConsumer(m.persistentCH)
	} else {
		fmt.Println("msg transfer not start mysql consumer")
	}
	go m.historyCH.historyConsumerGroup.RegisterHandleAndConsumer(m.historyCH)
	go m.historyMongoCH.historyConsumerGroup.RegisterHandleAndConsumer(m.historyMongoCH)
	go m.modifyCH.modifyMsgConsumerGroup.RegisterHandleAndConsumer(m.modifyCH)
	err := prome.StartPrometheusSrv(prometheusPort)
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}
