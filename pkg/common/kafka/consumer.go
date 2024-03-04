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

	"github.com/IBM/sarama"
	"github.com/OpenIMSDK/tools/errs"
)

type Consumer struct {
	addr          []string
	WG            sync.WaitGroup
	Topic         string
	PartitionList []int32
	Consumer      sarama.Consumer
}

func NewKafkaConsumer(addr []string, topic string, kafkaConfig *sarama.Config) (*Consumer, error) {
	p := Consumer{
		Topic: topic,
		addr:  addr,
	}

	if kafkaConfig.Net.SASL.User != "" && kafkaConfig.Net.SASL.Password != "" {
		kafkaConfig.Net.SASL.Enable = true
	}

	err := SetupTLSConfig(kafkaConfig)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(p.addr, kafkaConfig)
	if err != nil {
		return nil, errs.Wrap(err, "NewKafkaConsumer: creating consumer failed")
	}
	p.Consumer = consumer

	partitionList, err := consumer.Partitions(p.Topic)
	if err != nil {
		return nil, errs.Wrap(err, "NewKafkaConsumer: getting partitions failed")
	}
	p.PartitionList = partitionList

	return &p, nil

}
