// Copyright © 2024 OpenIM. All rights reserved.
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

package redispubsub

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Subscriber struct {
	client  redis.UniversalClient
	channel string
}

func NewSubscriber(client redis.UniversalClient, channel string) *Subscriber {
	return &Subscriber{client: client, channel: channel}
}

func (s *Subscriber) OnMessage(ctx context.Context, callback func(string)) error {
	messageChannel := s.client.Subscribe(ctx, s.channel).Channel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-messageChannel:
				callback(msg.Payload)
			}
		}
	}()

	return nil
}
