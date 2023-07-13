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
	"sync"

	"github.com/robfig/cron/v3"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
)

func StartCronTask() error {
	fmt.Println("cron task start, config", config.Config.ChatRecordsClearTime)
	msgTool, err := InitMsgTool()
	if err != nil {
		return err
	}
	c := cron.New()
	var wg sync.WaitGroup
	wg.Add(1)
	log.ZInfo(
		context.Background(),
		"start chatRecordsClearTime cron task",
		"cron config",
		config.Config.ChatRecordsClearTime,
	)
	_, err = c.AddFunc(config.Config.ChatRecordsClearTime, msgTool.AllConversationClearMsgAndFixSeq)
	if err != nil {
		fmt.Println(
			"start allConversationClearMsgAndFixSeq cron failed",
			err.Error(),
			config.Config.ChatRecordsClearTime,
		)
		panic(err)
	}
	log.ZInfo(context.Background(), "start msgDestruct cron task", "cron config", config.Config.MsgDestructTime)
	_, err = c.AddFunc(config.Config.MsgDestructTime, msgTool.ConversationsDestructMsgs)
	if err != nil {
		fmt.Println("start conversationsDestructMsgs cron failed", err.Error(), config.Config.ChatRecordsClearTime)
		panic(err)
	}
	c.Start()
	wg.Wait()
	return nil
}
