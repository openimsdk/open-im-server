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

package tools

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

func TestDisLock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	assert.Equal(t, true, netlock(rdb, "cron-1", 1*time.Second))

	// if exists, get false
	assert.Equal(t, false, netlock(rdb, "cron-1", 1*time.Second))

	time.Sleep(2 * time.Second)

	// wait for key on timeout, get true
	assert.Equal(t, true, netlock(rdb, "cron-1", 2*time.Second))

	// set different key
	assert.Equal(t, true, netlock(rdb, "cron-2", 2*time.Second))
}

func TestCronWrapFunc(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	once := sync.Once{}
	done := make(chan struct{}, 1)
	cb := func() {
		once.Do(func() {
			close(done)
		})
	}

	start := time.Now()
	key := fmt.Sprintf("cron-%v", rand.Int31())
	crontab := cron.New(cron.WithSeconds())
	crontab.AddFunc("*/1 * * * * *", cronWrapFunc(rdb, key, cb))
	crontab.Start()
	<-done

	dur := time.Since(start)
	assert.LessOrEqual(t, dur.Seconds(), float64(2*time.Second))
	crontab.Stop()
}

func TestCronWrapFuncWithNetlock(t *testing.T) {
	config.Config.EnableCronLocker = true
	rdb := redis.NewClient(&redis.Options{})
	defer rdb.Close()

	done := make(chan string, 10)

	crontab := cron.New(cron.WithSeconds())

	key := fmt.Sprintf("cron-%v", rand.Int31())
	crontab.AddFunc("*/1 * * * * *", cronWrapFunc(rdb, key, func() {
		done <- "host1"
	}))
	crontab.AddFunc("*/1 * * * * *", cronWrapFunc(rdb, key, func() {
		done <- "host2"
	}))
	crontab.Start()

	time.Sleep(12 * time.Second)
	// the ttl of netlock is 5s, so expected value is 2.
	assert.Equal(t, len(done), 2)

	crontab.Stop()
}
