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
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/db/redisutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

type CronTaskConfig struct {
	CronTask        config.CronTask
	RedisConfig     config.Redis
	MongodbConfig   config.Mongo
	ZookeeperConfig config.ZooKeeper
	Share           config.Share
	KafkaConfig     config.Kafka
}

func Start(ctx context.Context, config *CronTaskConfig) error {

	log.CInfo(ctx, "CRON-TASK server is initializing", "chatRecordsClearTime",
		config.CronTask.ChatRecordsClearTime, "msgDestructTime", config.CronTask.MsgDestructTime)

	msgTool, err := InitMsgTool(ctx, config)
	if err != nil {
		return err
	}

	msgTool.convertTools()

	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}

	// register cron tasks
	var crontab = cron.New()

	_, err = crontab.AddFunc(config.CronTask.ChatRecordsClearTime,
		cronWrapFunc(config, rdb, "cron_clear_msg_and_fix_seq", msgTool.AllConversationClearMsgAndFixSeq))
	if err != nil {
		return errs.Wrap(err)
	}

	_, err = crontab.AddFunc(config.CronTask.MsgDestructTime,
		cronWrapFunc(config, rdb, "cron_conversations_destruct_msgs", msgTool.ConversationsDestructMsgs))
	if err != nil {
		return errs.WrapMsg(err, "cron_conversations_destruct_msgs")
	}

	// start crontab
	crontab.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs

	// stop crontab, Wait for the running task to exit.
	cronCtx := crontab.Stop()

	select {
	case <-cronCtx.Done():
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

func cronWrapFunc(config *CronTaskConfig, rdb redis.UniversalClient, key string, fn func()) func() {
	enableCronLocker := config.CronTask.EnableCronLocker
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
