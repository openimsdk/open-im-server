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
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/errs"

	"github.com/IBM/sarama"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"google.golang.org/protobuf/proto"
)

const maxRetry = 10 // number of retries

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

// NewKafkaProducer initializes a new Kafka producer.
func NewKafkaProducer(addr []string, topic string, producerConfig *ProducerConfig, tlsConfig *TLSConfig) (*Producer, error) {
	p := Producer{
		addr:   addr,
		topic:  topic,
		config: sarama.NewConfig(),
	}

	// Set producer return flags
	p.config.Producer.Return.Successes = true
	p.config.Producer.Return.Errors = true

	// Set partitioner strategy
	p.config.Producer.Partitioner = sarama.NewHashPartitioner

	// Configure producer acknowledgement level
	configureProducerAck(&p, producerConfig.ProducerAck)

	// Configure message compression
	configureCompression(&p, producerConfig.CompressType)

	// Get Kafka configuration from environment variables or fallback to config file
	kafkaUsername := getEnvOrConfig("KAFKA_USERNAME", producerConfig.Username)
	kafkaPassword := getEnvOrConfig("KAFKA_PASSWORD", producerConfig.Password)
	kafkaAddr := getKafkaAddrFromEnv(addr) // Updated to use the new function

	// Configure SASL authentication if credentials are provided
	if kafkaUsername != "" && kafkaPassword != "" {
		p.config.Net.SASL.Enable = true
		p.config.Net.SASL.User = kafkaUsername
		p.config.Net.SASL.Password = kafkaPassword
	}

	// Set the Kafka address
	p.addr = kafkaAddr

	// Set up TLS configuration (if required)
	SetupTLSConfig(p.config, tlsConfig)

	// Create the producer with retries
	var err error
	for i := 0; i <= maxRetry; i++ {
		p.producer, err = sarama.NewSyncProducer(p.addr, p.config)
		if err == nil {
			return &p, errs.Wrap(err)
		}
		time.Sleep(1 * time.Second) // Wait before retrying
	}
	// Panic if unable to create producer after retries
	if err != nil {
		return nil, errs.Wrap(errors.New("failed to create Kafka producer: " + err.Error()))
	}

	return &p, nil
}

// configureProducerAck configures the producer's acknowledgement level.
func configureProducerAck(p *Producer, ackConfig string) {
	switch strings.ToLower(ackConfig) {
	case "no_response":
		p.config.Producer.RequiredAcks = sarama.NoResponse
	case "wait_for_local":
		p.config.Producer.RequiredAcks = sarama.WaitForLocal
	case "wait_for_all":
		p.config.Producer.RequiredAcks = sarama.WaitForAll
	default:
		p.config.Producer.RequiredAcks = sarama.WaitForAll
	}
}

// configureCompression configures the message compression type for the producer.
func configureCompression(p *Producer, compressType string) {
	var compress = sarama.CompressionNone
	err := compress.UnmarshalText(bytes.ToLower([]byte(compressType)))
	if err != nil {
		fmt.Printf("Failed to configure compression: %v\n", err)
		return
	}
	p.config.Producer.Compression = compress
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
		return 0, 0, errs.Wrap(err, "kafka proto Marshal err")
	}
	if len(bMsg) == 0 {
		return 0, 0, errs.Wrap(errEmptyMsg, "")
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
