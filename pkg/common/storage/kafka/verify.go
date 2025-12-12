// Copyright © 2024 OpenIM open source community. All rights reserved.
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
	"fmt"

	"github.com/IBM/sarama"
	"github.com/openimsdk/tools/errs"
)

func CheckTopics(ctx context.Context, conf *Config, topics []string) error {
	kfk, err := BuildConsumerGroupConfig(conf, sarama.OffsetNewest, false)
	if err != nil {
		return err
	}
	cli, err := sarama.NewClient(conf.Addr, kfk)
	if err != nil {
		return errs.WrapMsg(err, "NewClient failed", "config: ", fmt.Sprintf("%+v", conf))
	}
	defer cli.Close()

	existingTopics, err := cli.Topics()
	if err != nil {
		return errs.WrapMsg(err, "Failed to list topics")
	}

	existingTopicsMap := make(map[string]bool)
	for _, t := range existingTopics {
		existingTopicsMap[t] = true
	}

	for _, topic := range topics {
		if !existingTopicsMap[topic] {
			return errs.New("topic not exist", "topic", topic).Wrap()
		}
	}
	return nil
}

func CheckHealth(ctx context.Context, conf *Config) error {
	kfk, err := BuildConsumerGroupConfig(conf, sarama.OffsetNewest, false)
	if err != nil {
		return err
	}
	cli, err := sarama.NewClient(conf.Addr, kfk)
	if err != nil {
		return errs.WrapMsg(err, "NewClient failed", "config: ", fmt.Sprintf("%+v", conf))
	}
	defer cli.Close()

	// Get broker list
	brokers := cli.Brokers()
	if len(brokers) == 0 {
		return errs.New("no brokers found").Wrap()
	}

	// Check if all brokers are reachable
	for _, broker := range brokers {
		if err := broker.Open(kfk); err != nil {
			return errs.WrapMsg(err, "failed to open broker", "broker", broker.Addr())
		}
	}

	return nil
}
