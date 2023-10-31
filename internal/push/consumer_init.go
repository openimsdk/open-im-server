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

package push

type Consumer struct {
	pushCh       ConsumerHandler
	successCount uint64
}

func NewConsumer(pusher *Pusher) *Consumer {
	return &Consumer{
		pushCh: *NewConsumerHandler(pusher),
	}
}

func (c *Consumer) Start() {
	// statistics.NewStatistics(&c.successCount, config.Config.ModuleName.PushName, fmt.Sprintf("%d second push to
	// msg_gateway count", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	go c.pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&c.pushCh)
}
