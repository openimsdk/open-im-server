/*
** description("").
** copyright('Open_IM,www.Open_IM.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/13 10:33).
 */
package push

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	pbChat "Open_IM/pkg/proto/msg"
	pbPush "Open_IM/pkg/proto/push"
	"Open_IM/pkg/utils"
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type ConsumerHandler struct {
	pushConsumerGroup *kfk.MConsumerGroup
	pusher            Pusher
}

func (c *ConsumerHandler) Init() {
	c.pushConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ms2pschat.Topic}, config.Config.Kafka.Ms2pschat.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)
}

func (c *ConsumerHandler) handleMs2PsChat(msg []byte) {
	log.NewDebug("", "msg come from kafka  And push!!!", "msg", string(msg))
	msgFromMQ := pbChat.PushMsgDataToMQ{}
	if err := proto.Unmarshal(msg, &msgFromMQ); err != nil {
		log.Error("", "push Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	pbData := &pbPush.PushMsgReq{
		MsgData:  msgFromMQ.MsgData,
		SourceID: msgFromMQ.PushToUserID,
	}
	sec := msgFromMQ.MsgData.SendTime / 1000
	nowSec := utils.GetCurrentTimestampBySecond()
	if nowSec-sec > 10 {
		return
	}
	ctx := context.Background()
	tracelog.SetOperationID(ctx, "")
	var err error
	switch msgFromMQ.MsgData.SessionType {
	case constant.SuperGroupChatType:
		err = c.pusher.MsgToSuperGroupUser(ctx, pbData.SourceID, pbData.MsgData)
	default:
		err = c.pusher.MsgToUser(ctx, pbData.SourceID, pbData.MsgData)
	}
	if err != nil {
		log.NewError("", "push failed", *pbData)
	}
}
func (ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c *ConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
		c.handleMs2PsChat(msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}
