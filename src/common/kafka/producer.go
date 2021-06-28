package kafka

import (
	log2 "Open_IM/src/common/log"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type Producer struct {
	topic    string
	addr     []string
	config   *sarama.Config
	producer sarama.SyncProducer
}

func NewKafkaProducer(addr []string, topic string) *Producer {
	p := Producer{}
	p.config = sarama.NewConfig()                             //实例化个sarama的Config
	p.config.Producer.Return.Successes = true                 //是否开启消息发送成功后通知 successes channel
	p.config.Producer.RequiredAcks = sarama.WaitForAll        //设置生产者 消息 回复等级 0 1 all
	p.config.Producer.Partitioner = sarama.NewHashPartitioner //过设置 hash-key 自动 hash 分区,在发送消息的时候必须指定消息的key值,如果没有key，则随机选取分区

	p.addr = addr
	p.topic = topic

	producer, err := sarama.NewSyncProducer(p.addr, p.config) //初始化客户端
	if err != nil {
		panic(err)
		return nil
	}
	p.producer = producer
	return &p
}

func (p *Producer) SendMessage(m proto.Message, key ...string) (int32, int64, error) {
	kMsg := &sarama.ProducerMessage{}
	kMsg.Topic = p.topic
	if len(key) == 1 {
		kMsg.Key = sarama.StringEncoder(key[0])
	}
	bMsg, err := proto.Marshal(m)
	if err != nil {
		log2.Error("", "", "proto marshal err = %s", err.Error())
		return -1, -1, err
	}
	kMsg.Value = sarama.ByteEncoder(bMsg)

	return p.producer.SendMessage(kMsg)
}
