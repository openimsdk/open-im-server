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
	"github.com/openimsdk/tools/log"
	kfk "github.com/openimsdk/tools/mq/kafka"
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

func NewMConsumerGroup(conf *kfk.Config, groupID string, topics []string) (*MConsumerGroup, error) {
	config, err := kfk.BuildConsumerGroupConfig(conf, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}
	group, err := kfk.NewConsumerGroup(config, conf.Addr, groupID)
	if err != nil {
		return nil, err
	}
	return &MConsumerGroup{
		ConsumerGroup: group,
		groupID:       groupID,
		topics:        topics,
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
