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

package rpccache

import (
	"context"
	"encoding/json"

	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

func subscriberRedisDeleteCache(ctx context.Context, client redis.UniversalClient, channel string, del func(ctx context.Context, key ...string)) {
	defer func() {
		if r := recover(); r != nil {
			log.ZPanic(ctx, "subscriberRedisDeleteCache Panic", r)
		}
	}()
	for message := range client.Subscribe(ctx, channel).Channel() {
		log.ZDebug(ctx, "subscriberRedisDeleteCache", "channel", channel, "payload", message.Payload)
		var keys []string
		if err := json.Unmarshal([]byte(message.Payload), &keys); err != nil {
			log.ZError(ctx, "subscriberRedisDeleteCache json.Unmarshal error", err)
			continue
		}
		if len(keys) == 0 {
			continue
		}
		del(ctx, keys...)
	}
}
