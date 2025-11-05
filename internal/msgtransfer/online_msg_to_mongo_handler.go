package msgtransfer

import (
	"github.com/openimsdk/tools/mq"

	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/log"
	"google.golang.org/protobuf/proto"
)

type OnlineHistoryMongoConsumerHandler struct {
	msgTransferDatabase controller.MsgTransferDatabase
	config              *Config
	webhookClient       *webhook.Client
}

func NewOnlineHistoryMongoConsumerHandler(database controller.MsgTransferDatabase, config *Config) *OnlineHistoryMongoConsumerHandler {
	return &OnlineHistoryMongoConsumerHandler{
		msgTransferDatabase: database,
		config:              config,
		webhookClient:       webhook.NewWebhookClient(config.WebhooksConfig.URL),
	}
}

func (mc *OnlineHistoryMongoConsumerHandler) HandleChatWs2Mongo(val mq.Message) {
	ctx := val.Context()
	key := val.Key()
	msg := val.Value()
	msgFromMQ := pbmsg.MsgDataToMongoByMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.ZError(ctx, "unmarshall failed", err, "key", key, "len", len(msg))
		return
	}
	if len(msgFromMQ.MsgData) == 0 {
		log.ZError(ctx, "msgFromMQ.MsgData is empty", nil, "key", key, "msg", msg)
		return
	}
	log.ZDebug(ctx, "mongo consumer recv msg", "msgs", msgFromMQ.String())
	err = mc.msgTransferDatabase.BatchInsertChat2DB(ctx, msgFromMQ.ConversationID, msgFromMQ.MsgData, msgFromMQ.LastSeq)
	if err != nil {
		log.ZError(ctx, "batch data insert to mongo err", err, "msg", msgFromMQ.MsgData, "conversationID", msgFromMQ.ConversationID)
		prommetrics.MsgInsertMongoFailedCounter.Inc()
	} else {
		prommetrics.MsgInsertMongoSuccessCounter.Inc()
		val.Mark()
	}

	for _, msgData := range msgFromMQ.MsgData {
		mc.webhookAfterMsgSaveDB(ctx, &mc.config.WebhooksConfig.AfterMsgSaveDB, msgData)
	}
}
