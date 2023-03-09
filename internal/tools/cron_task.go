package tools

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
)

const cronTaskOperationID = "cronTaskOperationID-"
const moduleName = "cron"

func StartCronTask() error {
	log.NewPrivateLog(moduleName)
	log.NewInfo("StartCronTask", "start cron task", "cron config", config.Config.Mongo.ChatRecordsClearTime)
	fmt.Println("cron task start, config", config.Config.Mongo.ChatRecordsClearTime)
	msgTool, err := InitMsgTool()
	if err != nil {
		return err
	}
	c := cron.New()
	var wg sync.WaitGroup
	wg.Add(1)
	_, err = c.AddFunc(config.Config.Mongo.ChatRecordsClearTime, msgTool.AllUserClearMsgAndFixSeq)
	if err != nil {
		fmt.Println("start cron failed", err.Error(), config.Config.Mongo.ChatRecordsClearTime)
		return err
	}
	c.Start()
	wg.Wait()
	return nil
}
