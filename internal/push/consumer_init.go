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

import "context"

type Consumer struct {
	pushCh       ConsumerHandler
	successCount uint64
}

func NewConsumer(pusher *Pusher) (*Consumer, error) {
	c, err := NewConsumerHandler(pusher)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		pushCh: *c,
	}, nil
}

func (c *Consumer) Start() {

	go c.pushCh.pushConsumerGroup.RegisterHandleAndConsumer(context.Background(), &c.pushCh)
}
