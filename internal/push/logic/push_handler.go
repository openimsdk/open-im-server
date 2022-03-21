/*
** description("").
** copyright('Open_IM,www.Open_IM.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/13 10:33).
 */
package logic

import (
	"Open_IM/pkg/common/config"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/chat"
	pbPush "Open_IM/pkg/proto/push"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type fcb func(msg []byte)

type PushConsumerHandler struct {
	msgHandle         map[string]fcb
	pushConsumerGroup *kfk.MConsumerGroup
}

func (ms *PushConsumerHandler) Init() {
	ms.msgHandle = make(map[string]fcb)
	ms.msgHandle[config.Config.Kafka.Ms2pschat.Topic] = ms.handleMs2PsChat
	ms.pushConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ms2pschat.Topic}, config.Config.Kafka.Ms2pschat.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)
}
func (ms *PushConsumerHandler) handleMs2PsChat(msg []byte) {
	log.InfoByKv("msg come from kafka  And push!!!", "", "msg", string(msg))
	msgFromMQ := pbChat.PushMsgDataToMQ{}
	if err := proto.Unmarshal(msg, &msgFromMQ); err != nil {
		log.ErrorByKv("push Unmarshal msg err", "", "msg", string(msg), "err", err.Error())
		return
	}
	//Call push module to send message to the user
	MsgToUser((*pbPush.PushMsgReq)(&msgFromMQ))
}
func (PushConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (PushConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (ms *PushConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.InfoByKv("kafka get info to mysql", "", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
		ms.msgHandle[msg.Topic](msg.Value)
	}
	return nil
}
