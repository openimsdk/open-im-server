package kafka

import (
	"context"
	"errors"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	log "github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"

	prome "github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
)

var errEmptyMsg = errors.New("binary msg is empty")

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
	if config.Config.Kafka.Username != "" && config.Config.Kafka.Password != "" {
		p.config.Net.SASL.Enable = true
		p.config.Net.SASL.User = config.Config.Kafka.Username
		p.config.Net.SASL.Password = config.Config.Kafka.Password
	}
	p.addr = addr
	p.topic = topic
	producer, err := sarama.NewSyncProducer(p.addr, p.config) //Initialize the client
	if err != nil {
		panic(err.Error())
	}
	p.producer = producer
	return &p
}

func GetMQHeaderWithContext(ctx context.Context) ([]sarama.RecordHeader, error) {
	operationID, opUserID, platform, connID, err := mcontext.GetCtxInfos(ctx)
	if err != nil {
		return nil, err
	}
	return []sarama.RecordHeader{
		{Key: []byte(constant.OperationID), Value: []byte(operationID)},
		{Key: []byte(constant.OpUserID), Value: []byte(opUserID)},
		{Key: []byte(constant.OpUserPlatform), Value: []byte(platform)},
		{Key: []byte(constant.ConnID), Value: []byte(connID)}}, err
}

func GetContextWithMQHeader(header []*sarama.RecordHeader) context.Context {
	var values []string
	for _, recordHeader := range header {
		values = append(values, string(recordHeader.Value))
	}
	return mcontext.WithMustInfoCtx(values) // TODO
}

func (p *Producer) SendMessage(ctx context.Context, key string, msg proto.Message) (int32, int64, error) {
	log.ZDebug(ctx, "SendMessage", "msg", msg, "topic", p.topic, "key", key)
	kMsg := &sarama.ProducerMessage{}
	kMsg.Topic = p.topic
	kMsg.Key = sarama.StringEncoder(key)
	bMsg, err := proto.Marshal(msg)
	if err != nil {
		return 0, 0, utils.Wrap(err, "kafka proto Marshal err")
	}
	if len(bMsg) == 0 {
		return 0, 0, utils.Wrap(errEmptyMsg, "")
	}
	kMsg.Value = sarama.ByteEncoder(bMsg)
	if kMsg.Key.Length() == 0 || kMsg.Value.Length() == 0 {
		return 0, 0, utils.Wrap(errEmptyMsg, "")
	}
	kMsg.Metadata = ctx
	header, err := GetMQHeaderWithContext(ctx)
	if err != nil {
		return 0, 0, utils.Wrap(err, "")
	}
	kMsg.Headers = header
	partition, offset, err := p.producer.SendMessage(kMsg)
	log.ZDebug(ctx, "ByteEncoder SendMessage end", "key ", kMsg.Key, "key length", kMsg.Value.Length())
	if err == nil {
		prome.Inc(prome.SendMsgCounter)
	}
	return partition, offset, utils.Wrap(err, "")
}
