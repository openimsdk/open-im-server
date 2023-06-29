package push

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	kfk "github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	pbChat "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	pbPush "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/push"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

type ConsumerHandler struct {
	pushConsumerGroup *kfk.MConsumerGroup
	pusher            *Pusher
}

func NewConsumerHandler(pusher *Pusher) *ConsumerHandler {
	var consumerHandler ConsumerHandler
	consumerHandler.pusher = pusher
	consumerHandler.pushConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.MsgToPush.Topic}, config.Config.Kafka.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)
	return &consumerHandler
}

func (c *ConsumerHandler) handleMs2PsChat(ctx context.Context, msg []byte) {
	msgFromMQ := pbChat.PushMsgDataToMQ{}
	if err := proto.Unmarshal(msg, &msgFromMQ); err != nil {
		log.ZError(ctx, "push Unmarshal msg err", err, "msg", string(msg))
		return
	}
	pbData := &pbPush.PushMsgReq{
		MsgData:        msgFromMQ.MsgData,
		ConversationID: msgFromMQ.ConversationID,
	}
	sec := msgFromMQ.MsgData.SendTime / 1000
	nowSec := utils.GetCurrentTimestampBySecond()
	if nowSec-sec > 10 {
		return
	}
	var err error
	switch msgFromMQ.MsgData.SessionType {
	case constant.SuperGroupChatType:
		err = c.pusher.Push2SuperGroup(ctx, pbData.MsgData.GroupID, pbData.MsgData)
	default:
		err = c.pusher.Push2User(ctx, []string{pbData.MsgData.SendID, pbData.MsgData.RecvID}, pbData.MsgData)
	}
	if err != nil {
		if err == errNoOfflinePusher {
			log.ZWarn(ctx, "offline push failed", err, "msg", pbData.String())
		} else {
			log.ZError(ctx, "push failed", err, "msg", pbData.String())
		}
	}
}
func (ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c *ConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := c.pushConsumerGroup.GetContextFromMsg(msg)
		c.handleMs2PsChat(ctx, msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}
