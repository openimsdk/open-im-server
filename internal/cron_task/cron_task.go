package cronTask

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

const cronTaskOperationID = "cronTaskOperationID-"

func StartCronTask(userID, workingGroupID string) {
	log.NewPrivateLog("cron")
	log.NewInfo(utils.OperationIDGenerator(), "start cron task", "cron config", config.Config.Mongo.ChatRecordsClearTime)
	fmt.Println("cron config", config.Config.Mongo.ChatRecordsClearTime)
	if userID != "" {
		operationID := getCronTaskOperationID()
		StartClearMsg(operationID, []string{userID})
	}
	if workingGroupID != "" {
		operationID := getCronTaskOperationID()
		StartClearWorkingGroupMsg(operationID, []string{workingGroupID})
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

func getCronTaskOperationID() string {
	return cronTaskOperationID + utils.OperationIDGenerator()
}

func ClearAll() {
	operationID := getCronTaskOperationID()
	log.NewInfo(operationID, "====================== start del cron task ======================")
	var err error
	userIDList, err := im_mysql_model.SelectAllUserID()
	if err == nil {
		StartClearMsg(operationID, userIDList)
	} else {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}

	// working group msg clear
	workingGroupIDList, err := im_mysql_model.GetGroupIDListByGroupType(constant.WorkingGroup)
	if err == nil {
		StartClearWorkingGroupMsg(operationID, workingGroupIDList)
	} else {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}

	log.NewInfo(operationID, "====================== start del cron finished ======================")
}

func StartClearMsg(operationID string, userIDList []string) {
	log.NewDebug(operationID, utils.GetSelfFuncName(), "userIDList: ", userIDList)
	for _, userID := range userIDList {
		if err := DeleteMongoMsgAndResetRedisSeq(operationID, userID); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), userID)
		}
		if err := checkMaxSeqWithMongo(operationID, userID, constant.WriteDiffusion); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), userID, err)
		}
	}
}

func StartClearWorkingGroupMsg(operationID string, workingGroupIDList []string) {
	log.NewDebug(operationID, utils.GetSelfFuncName(), "workingGroupIDList: ", workingGroupIDList)
	for _, groupID := range workingGroupIDList {
		userIDList, err := rocksCache.GetGroupMemberIDListFromCache(groupID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID)
			continue
		}
		log.NewDebug(operationID, utils.GetSelfFuncName(), "groupID:", groupID, "workingGroupIDList:", userIDList)
		if err := ResetUserGroupMinSeq(operationID, groupID, userIDList); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID, userIDList)
		}
		if err := checkMaxSeqWithMongo(operationID, groupID, constant.ReadDiffusion); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), groupID, err)
		}
	}
}
