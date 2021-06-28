/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package logic

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/db"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/kafka"
	"Open_IM/src/common/log"
	"Open_IM/src/utils"
	"time"
)

var (
	rpcServer    RPCServer
	pushCh       PushConsumerHandler
	pushTerminal []int32
	producer     *kafka.Producer
)

func Init(rpcPort int) {
	log.NewPrivateLog(config.Config.ModuleName.PushName)
	rpcServer.Init(rpcPort)
	pushCh.Init()
	pushTerminal = []int32{utils.IOSPlatformID}
}
func init() {
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
}

func Run() {
	go rpcServer.run()
	go scheduleDelete()
	go pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&pushCh)
}

func scheduleDelete() {
	//uid, _ := im_mysql_model.SelectAllUID()
	//db.DB.DelHistoryChat(0, uid)
	//log.Info("", "", "sssssssssss")
	//if err != nil {
	//	db.DB.DelHistoryChat(0, uid)
	//}

	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C

		uid, err := im_mysql_model.SelectAllUID()
		if err != nil {
			db.DB.DelHistoryChat(int64(config.Config.Mongo.DBRetainChatRecords), uid)
		}

		//以下为定时执行的操作
		scheduleDelete()
	}
}
