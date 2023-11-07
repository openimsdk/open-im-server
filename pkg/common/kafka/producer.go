// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/OpenIMSDK/protocol/constant"
	log "github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

const (
	maxRetry = 10 // number of retries
)

var errEmptyMsg = errors.New("binary msg is empty")

type Producer struct {
	topic    string
	addr     []string
	config   *sarama.Config
	producer sarama.SyncProducer
}

// NewKafkaProducer Initialize kafka producer.
func NewKafkaProducer(addr []string, topic string) *Producer {
	p := Producer{}
	p.config = sarama.NewConfig()             // Instantiate a sarama Config
	p.config.Producer.Return.Successes = true // Whether to enable the successes channel to be notified after the message is sent successfully
	p.config.Producer.Return.Errors = true
	p.config.Producer.RequiredAcks = sarama.WaitForAll        // Set producer Message Reply level 0 1 all
	p.config.Producer.Partitioner = sarama.NewHashPartitioner // Set the hash-key automatic hash partition. When sending a message, you must specify the key value of the message. If there is no key, the partition will be selected randomly
	if config.Config.Kafka.Username != "" && config.Config.Kafka.Password != "" {
		p.config.Net.SASL.Enable = true
		p.config.Net.SASL.User = config.Config.Kafka.Username
		p.config.Net.SASL.Password = config.Config.Kafka.Password
	}
	p.addr = addr
	p.topic = topic
	SetupTLSConfig(p.config)
	var producer sarama.SyncProducer
	var err error
	for i := 0; i <= maxRetry; i++ {
		producer, err = sarama.NewSyncProducer(p.addr, p.config) // Initialize the client
		if err == nil {
			p.producer = producer
			return &p
		}
		//TODO If the password is wrong, exit directly
		//if packetErr, ok := err.(*sarama.PacketEncodingError); ok {
		//if _, ok := packetErr.Err.(sarama.AuthenticationError); ok {
		//	fmt.Println("Kafka password is wrong.")
		//}
		//} else {
		//	fmt.Printf("Failed to create Kafka producer: %v\n", err)
		//}
		time.Sleep(time.Duration(1) * time.Second)
	}
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
		{Key: []byte(constant.ConnID), Value: []byte(connID)},
	}, err
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
	if err != nil {
		log.ZWarn(ctx, "p.producer.SendMessage error", err)
	}
	return partition, offset, utils.Wrap(err, "")
}
