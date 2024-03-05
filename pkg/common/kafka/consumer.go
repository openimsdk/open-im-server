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

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"

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

func NewKafkaConsumer(addr []string, topic string, config *config.GlobalConfig) (*Consumer,error) {
	p := Consumer{}
	p.Topic = topic
	p.addr = addr
	consumerConfig := sarama.NewConfig()
	if config.Kafka.Username != "" && config.Kafka.Password != "" {
		consumerConfig.Net.SASL.Enable = true
		consumerConfig.Net.SASL.User = config.Kafka.Username
		consumerConfig.Net.SASL.Password = config.Kafka.Password
	}
	var tlsConfig *TLSConfig
	if config.Kafka.TLS != nil {
		tlsConfig = &TLSConfig{
			CACrt:              config.Kafka.TLS.CACrt,
			ClientCrt:          config.Kafka.TLS.ClientCrt,
			ClientKey:          config.Kafka.TLS.ClientKey,
			ClientKeyPwd:       config.Kafka.TLS.ClientKeyPwd,
			InsecureSkipVerify: false,
		}
	}
	err:=SetupTLSConfig(consumerConfig, tlsConfig)
	if err!=nil{
		return nil,err
	}
	consumer, err := sarama.NewConsumer(p.addr, consumerConfig)
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
