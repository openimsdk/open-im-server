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
	"strings"

	"github.com/IBM/sarama"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
)

type MConsumerGroup struct {
	sarama.ConsumerGroup
	groupID string
	topics  []string
}

type MConsumerGroupConfig struct {
	KafkaVersion   sarama.KafkaVersion
	OffsetsInitial int64
	IsReturnErr    bool
	UserName       string
	Password       string
}

func NewMConsumerGroup(consumerConfig *MConsumerGroupConfig, topics, addrs []string, groupID string, tlsConfig *TLSConfig) (*MConsumerGroup, error) {
	consumerGroupConfig := sarama.NewConfig()
	consumerGroupConfig.Version = consumerConfig.KafkaVersion
	consumerGroupConfig.Consumer.Offsets.Initial = consumerConfig.OffsetsInitial
	consumerGroupConfig.Consumer.Return.Errors = consumerConfig.IsReturnErr
	if consumerConfig.UserName != "" && consumerConfig.Password != "" {
		consumerGroupConfig.Net.SASL.Enable = true
		consumerGroupConfig.Net.SASL.User = consumerConfig.UserName
		consumerGroupConfig.Net.SASL.Password = consumerConfig.Password
	}

	SetupTLSConfig(consumerGroupConfig, tlsConfig)
	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, consumerGroupConfig)
	if err != nil {
		return nil, errs.Wrap(err, strings.Join(topics, ","), strings.Join(addrs, ","), groupID, consumerConfig.UserName, consumerConfig.Password)
	}

	return &MConsumerGroup{
		consumerGroup,
		groupID,
		topics,
	}, nil
}

func (mc *MConsumerGroup) GetContextFromMsg(cMsg *sarama.ConsumerMessage) context.Context {
	return GetContextWithMQHeader(cMsg.Headers)
}

func (mc *MConsumerGroup) RegisterHandleAndConsumer(ctx context.Context, handler sarama.ConsumerGroupHandler) {
	log.ZDebug(ctx, "register consumer group", "groupID", mc.groupID)
	for {
		err := mc.ConsumerGroup.Consume(ctx, mc.topics, handler)
		if errors.Is(err, sarama.ErrClosedConsumerGroup) {
			return
		}
		if errors.Is(err, context.Canceled) {
			return
		}
		if err != nil {
			log.ZWarn(ctx, "consume err", err, "topic", mc.topics, "groupID", mc.groupID)
		}
	}
}

func (mc *MConsumerGroup) Close() error {
	return mc.ConsumerGroup.Close()
}
