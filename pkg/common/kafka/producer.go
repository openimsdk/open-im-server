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

	"github.com/IBM/sarama"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	kfk "github.com/openimsdk/tools/mq/kafka"
	"google.golang.org/protobuf/proto"
)

var errEmptyMsg = errors.New("kafka binary msg is empty")

// Producer represents a Kafka producer.
type Producer struct {
	addr     []string
	topic    string
	config   *sarama.Config
	producer sarama.SyncProducer
}

type ProducerConfig struct {
	ProducerAck  string
	CompressType string
	Username     string
	Password     string
}

func BuildProducerConfig(conf kfk.Config) (*sarama.Config, error) {
	return kfk.BuildProducerConfig(conf)
}

func NewKafkaProducer(config *sarama.Config, addr []string, topic string) (*Producer, error) {
	producer, err := kfk.NewProducer(config, addr)
	if err != nil {
		return nil, err
	}
	return &Producer{
		addr:     addr,
		topic:    topic,
		config:   config,
		producer: producer,
	}, nil
}

// GetMQHeaderWithContext extracts message queue headers from the context.
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
	}, nil
}

// GetContextWithMQHeader creates a context from message queue headers.
func GetContextWithMQHeader(header []*sarama.RecordHeader) context.Context {
	var values []string
	for _, recordHeader := range header {
		values = append(values, string(recordHeader.Value))
	}
	return mcontext.WithMustInfoCtx(values) // Attach extracted values to context
}

// SendMessage sends a message to the Kafka topic configured in the Producer.
func (p *Producer) SendMessage(ctx context.Context, key string, msg proto.Message) (int32, int64, error) {
	log.ZDebug(ctx, "SendMessage", "msg", msg, "topic", p.topic, "key", key)

	// Marshal the protobuf message
	bMsg, err := proto.Marshal(msg)
	if err != nil {
		return 0, 0, errs.WrapMsg(err, "kafka proto Marshal err")
	}
	if len(bMsg) == 0 {
		return 0, 0, errs.WrapMsg(errEmptyMsg, "kafka proto Marshal err")
	}

	// Prepare Kafka message
	kMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(bMsg),
	}

	// Validate message key and value
	if kMsg.Key.Length() == 0 || kMsg.Value.Length() == 0 {
		return 0, 0, errs.Wrap(errEmptyMsg)
	}

	// Attach context metadata as headers
	header, err := GetMQHeaderWithContext(ctx)
	if err != nil {
		return 0, 0, err
	}
	kMsg.Headers = header

	// Send the message
	partition, offset, err := p.producer.SendMessage(kMsg)
	if err != nil {
		log.ZWarn(ctx, "p.producer.SendMessage error", err)
		return 0, 0, errs.Wrap(err)
	}

	log.ZDebug(ctx, "ByteEncoder SendMessage end", "key", kMsg.Key, "key length", kMsg.Value.Length())
	return partition, offset, nil
}
