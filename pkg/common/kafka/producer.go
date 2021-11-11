package kafka

import (
	log2 "Open_IM/pkg/common/log"
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
	p.config = sarama.NewConfig()                             //Instantiate a sarama Config
	p.config.Producer.Return.Successes = true                 //Whether to enable the successes channel to be notified after the message is sent successfully
	p.config.Producer.RequiredAcks = sarama.WaitForAll        //Set producer Message Reply level 0 1 all
	p.config.Producer.Partitioner = sarama.NewHashPartitioner //Set the hash-key automatic hash partition. When sending a message, you must specify the key value of the message. If there is no key, the partition will be selected randomly

	p.addr = addr
	p.topic = topic

	producer, err := sarama.NewSyncProducer(p.addr, p.config) //Initialize the client
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
