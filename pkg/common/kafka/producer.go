package kafka

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	log "OpenIM/pkg/common/log"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils"
	"context"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"

	prome "OpenIM/pkg/common/prome"
)

type Producer struct {
	topic    string
	addr     []string
	config   *sarama.Config
	producer sarama.SyncProducer
}

func NewKafkaProducer(addr []string, topic string) *Producer {
	p := Producer{}
	p.config = sarama.NewConfig()             //Instantiate a sarama Config
	p.config.Producer.Return.Successes = true //Whether to enable the successes channel to be notified after the message is sent successfully
	p.config.Producer.Return.Errors = true
	p.config.Producer.RequiredAcks = sarama.WaitForAll        //Set producer Message Reply level 0 1 all
	p.config.Producer.Partitioner = sarama.NewHashPartitioner //Set the hash-key automatic hash partition. When sending a message, you must specify the key value of the message. If there is no key, the partition will be selected randomly
	if config.Config.Kafka.SASLUserName != "" && config.Config.Kafka.SASLPassword != "" {
		p.config.Net.SASL.Enable = true
		p.config.Net.SASL.User = config.Config.Kafka.SASLUserName
		p.config.Net.SASL.Password = config.Config.Kafka.SASLPassword
	}
	p.addr = addr
	p.topic = topic

	producer, err := sarama.NewSyncProducer(p.addr, p.config) //Initialize the client
	if err != nil {
		panic(err.Error())
		return nil
	}
	p.producer = producer
	return &p
}

func (p *Producer) SendMessage(ctx context.Context, m proto.Message, key string) (int32, int64, error) {
	operationID := tracelog.GetOperationID(ctx)
	log.Info(operationID, "SendMessage", "key ", key, m.String(), p.producer)
	kMsg := &sarama.ProducerMessage{}
	kMsg.Topic = p.topic
	kMsg.Key = sarama.StringEncoder(key)
	bMsg, err := proto.Marshal(m)
	if err != nil {
		log.Error(operationID, "", "proto marshal err = %s", err.Error())
		return -1, -1, err
	}
	if len(bMsg) == 0 {
		log.Error(operationID, "len(bMsg) == 0 ")
		return 0, 0, errors.New("len(bMsg) == 0 ")
	}
	kMsg.Value = sarama.ByteEncoder(bMsg)
	log.Info(operationID, "ByteEncoder SendMessage begin", "key ", kMsg, p.producer, "len: ", kMsg.Key.Length(), kMsg.Value.Length())
	if kMsg.Key.Length() == 0 || kMsg.Value.Length() == 0 {
		log.Error(operationID, "kMsg.Key.Length() == 0 || kMsg.Value.Length() == 0 ", kMsg)
		return -1, -1, errors.New("key or value == 0")
	}
	kMsg.Metadata = ctx
	kMsg.Headers = []sarama.RecordHeader{{Key: []byte(constant.OperationID), Value: []byte(operationID)}}
	partition, offset, err := p.producer.SendMessage(kMsg)
	log.Info(operationID, "ByteEncoder SendMessage end", "key ", kMsg.Key.Length(), kMsg.Value.Length(), p.producer)
	if err == nil {
		prome.Inc(prome.SendMsgCounter)
	}
	return partition, offset, utils.Wrap(err, "")
}
