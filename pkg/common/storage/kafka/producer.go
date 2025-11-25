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
	"github.com/IBM/sarama"
	"github.com/openimsdk/tools/errs"
	"google.golang.org/protobuf/proto"
)

// Producer represents a Kafka producer.
type Producer struct {
	addr     []string
	topic    string
	config   *sarama.Config
	producer sarama.SyncProducer
}

func NewKafkaProducer(config *sarama.Config, addr []string, topic string) (*Producer, error) {
	producer, err := NewProducer(config, addr)
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

// SendMessage sends a message to the Kafka topic configured in the Producer.
func (p *Producer) SendMessage(ctx context.Context, key string, msg proto.Message) (int32, int64, error) {
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
		return 0, 0, errs.WrapMsg(err, "p.producer.SendMessage error")
	}

	return partition, offset, nil
}
