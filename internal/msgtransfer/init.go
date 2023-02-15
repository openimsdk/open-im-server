package msgtransfer

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/prome"
	"fmt"
)

type MsgTransfer struct {
	persistentCH   PersistentConsumerHandler         // 聊天记录持久化到mysql的消费者 订阅的topic: ws2ms_chat
	historyCH      OnlineHistoryRedisConsumerHandler // 这个消费者聚合消息, 订阅的topic：ws2ms_chat, 修改通知发往msg_to_modify topic, 消息存入redis后Incr Redis, 再发消息到ms2pschat topic推送， 发消息到msg_to_mongo topic持久化
	historyMongoCH OnlineHistoryMongoConsumerHandler // mongoDB批量插入, 成功后删除redis中消息，以及处理删除通知消息删除的 订阅的topic: msg_to_mongo
	modifyCH       ModifyMsgConsumerHandler          // 负责消费修改消息通知的consumer, 订阅的topic: msg_to_modify
}

func NewMsgTransfer() *MsgTransfer {
	msgTransfer := &MsgTransfer{}
	msgTransfer.persistentCH.Init()
	msgTransfer.historyCH.Init()
	msgTransfer.historyMongoCH.Init()
	msgTransfer.modifyCH.Init()
	if config.Config.Prometheus.Enable {
		msgTransfer.initPrometheus()
	}
	return msgTransfer
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

func (m *MsgTransfer) Run(promePort int) {
	if config.Config.ChatPersistenceMysql {
		go m.persistentCH.persistentConsumerGroup.RegisterHandleAndConsumer(&m.persistentCH)
	} else {
		fmt.Println("msg transfer not start mysql consumer")
	}
	go m.historyCH.historyConsumerGroup.RegisterHandleAndConsumer(&m.historyCH)
	go m.historyMongoCH.historyConsumerGroup.RegisterHandleAndConsumer(&m.historyMongoCH)
	go m.modifyCH.modifyMsgConsumerGroup.RegisterHandleAndConsumer(&m.modifyCH)
	go func() {
		err := prome.StartPromeSrv(promePort)
		if err != nil {
			panic(err)
		}
	}()
}
