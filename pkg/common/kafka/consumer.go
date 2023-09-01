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
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	addr          []string
	WG            sync.WaitGroup
	Topic         string
	PartitionList []int32
	Consumer      sarama.Consumer
}

func NewKafkaConsumer(addr []string, topic string) *Consumer {
	p := Consumer{}
	p.Topic = topic
	p.addr = addr
	consumerConfig := sarama.NewConfig()
	if config.Config.Kafka.Username != "" && config.Config.Kafka.Password != "" {
		consumerConfig.Net.SASL.Enable = true
		consumerConfig.Net.SASL.User = config.Config.Kafka.Username
		consumerConfig.Net.SASL.Password = config.Config.Kafka.Password
	}
	SetupTLSConfig(consumerConfig)
	consumer, err := sarama.NewConsumer(p.addr, consumerConfig)
	if err != nil {
		panic(err.Error())
	}
	p.Consumer = consumer

	partitionList, err := consumer.Partitions(p.Topic)
	if err != nil {
		panic(err.Error())
	}
	p.PartitionList = partitionList

	return &p
}
