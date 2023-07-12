package tools

import (
	"context"
	"fmt"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/robfig/cron/v3"
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
	log.ZInfo(context.Background(), "start chatRecordsClearTime cron task", "cron config", config.Config.ChatRecordsClearTime)
	_, err = c.AddFunc(config.Config.ChatRecordsClearTime, msgTool.AllConversationClearMsgAndFixSeq)
	if err != nil {
		fmt.Println("start allConversationClearMsgAndFixSeq cron failed", err.Error(), config.Config.ChatRecordsClearTime)
		panic(err)
	}
	log.ZInfo(context.Background(), "start msgDestruct cron task", "cron config", config.Config.ChatRecordsClearTime)
	_, err = c.AddFunc(config.Config.MsgDestructTime, msgTool.ConversationsDestructMsgs)
	if err != nil {
		fmt.Println("start conversationsDestructMsgs cron failed", err.Error(), config.Config.ChatRecordsClearTime)
		panic(err)
	}
	c.Start()
	wg.Wait()
	return nil
}
