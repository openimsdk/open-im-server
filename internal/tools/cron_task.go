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
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"

	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
)

func StartTask() error {
	fmt.Println("cron task start, config", config.Config.ChatRecordsClearTime)
	msgTool, err := InitMsgTool()
	if err != nil {
		return err
	}

	msgTool.convertTools()

	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}

	// register cron tasks
	var crontab = cron.New()
	log.ZInfo(context.Background(), "start chatRecordsClearTime cron task", "cron config", config.Config.ChatRecordsClearTime)
	_, err = crontab.AddFunc(config.Config.ChatRecordsClearTime, cronWrapFunc(rdb, "cron_clear_msg_and_fix_seq", msgTool.AllConversationClearMsgAndFixSeq))
	if err != nil {
		log.ZError(context.Background(), "start allConversationClearMsgAndFixSeq cron failed", err)
		panic(err)
	}

	log.ZInfo(context.Background(), "start msgDestruct cron task", "cron config", config.Config.MsgDestructTime)
	_, err = crontab.AddFunc(config.Config.MsgDestructTime, cronWrapFunc(rdb, "cron_conversations_destruct_msgs", msgTool.ConversationsDestructMsgs))
	if err != nil {
		log.ZError(context.Background(), "start conversationsDestructMsgs cron failed", err)
		panic(err)
	}

	// start crontab
	crontab.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs

	// stop crontab, Wait for the running task to exit.
	ctx := crontab.Stop()

	select {
	case <-ctx.Done():
		// graceful exit

	case <-time.After(15 * time.Second):
		// forced exit on timeout
	}

	return nil
}

// netlock redis lock.
func netlock(rdb redis.UniversalClient, key string, ttl time.Duration) bool {
	value := "used"
	ok, err := rdb.SetNX(context.Background(), key, value, ttl).Result() // nolint
	if err != nil {
		// when err is about redis server, return true.
		return false
	}

	return ok
}

func cronWrapFunc(rdb redis.UniversalClient, key string, fn func()) func() {
	enableCronLocker := config.Config.EnableCronLocker
	return func() {
		// if don't enable cron-locker, call fn directly.
		if !enableCronLocker {
			fn()
			return
		}

		// when acquire redis lock, call fn().
		if netlock(rdb, key, 5*time.Second) {
			fn()
		}
	}
}
