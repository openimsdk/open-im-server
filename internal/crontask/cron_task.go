package cronTask

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

const cronTaskOperationID = "cronTaskOperationID-"
const moduleName = "cron"

func StartCronTask(userID, workingGroupID string) {
	log.NewPrivateLog(moduleName)
	log.NewInfo(utils.OperationIDGenerator(), "start cron task", "cron config", config.Config.Mongo.ChatRecordsClearTime)
	fmt.Println("cron task start, config", config.Config.Mongo.ChatRecordsClearTime)
	if userID != "" {
		operationID := getCronTaskOperationID()
		ClearUsersMsg(operationID, []string{userID})
	}
	if workingGroupID != "" {
		operationID := getCronTaskOperationID()
		ClearSuperGroupMsg(operationID, []string{workingGroupID})
	}
	if userID != "" || workingGroupID != "" {
		fmt.Println("clear msg finished")
		return
	}
	c := cron.New()
	_, err := c.AddFunc(config.Config.Mongo.ChatRecordsClearTime, ClearAll)
	if err != nil {
		fmt.Println("start cron failed", err.Error(), config.Config.Mongo.ChatRecordsClearTime)
		panic(err)
	}
	c.Start()
	fmt.Println("start cron task success")
	for {
		time.Sleep(10 * time.Second)
	}
}
