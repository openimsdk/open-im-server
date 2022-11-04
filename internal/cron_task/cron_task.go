package cronTask

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

const cronTaskOperationID = "cronTaskOperationID-"

func StartCronTask() {
	log.NewPrivateLog("cron")
	log.NewInfo(utils.OperationIDGenerator(), "start cron task", "cron config", config.Config.Mongo.ChatRecordsClearTime)
	c := cron.New()
	fmt.Println("cron config", config.Config.Mongo.ChatRecordsClearTime)
	_, err := c.AddFunc(config.Config.Mongo.ChatRecordsClearTime, func() {
		// user msg clear
		operationID := getCronTaskOperationID()
		log.NewInfo(operationID, "====================== start del cron task ======================")
		userIDList, err := im_mysql_model.SelectAllUserID()
		if err == nil {
			log.NewDebug(operationID, utils.GetSelfFuncName(), "userIDList: ", userIDList)
			userIDList = []string{"4158779020"}
			for _, userID := range userIDList {
				if err := DeleteMongoMsgAndResetRedisSeq(operationID, userID); err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), userID)
				}
				//if err := checkMaxSeqWithMongo(operationID, userID, constant.WriteDiffusion); err != nil {
				//	log.NewError(operationID, utils.GetSelfFuncName(), userID, err)
				//}
			}
		} else {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		}

		// working group msg clear
		workingGroupIDList, err := im_mysql_model.GetGroupIDListByGroupType(constant.WorkingGroup)
		if err == nil {
			log.NewDebug(operationID, utils.GetSelfFuncName(), "workingGroupIDList: ", workingGroupIDList)
			for _, groupID := range workingGroupIDList {
				userIDList, err = rocksCache.GetGroupMemberIDListFromCache(groupID)
				if err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID)
					continue
				}
				log.NewDebug(operationID, utils.GetSelfFuncName(), "groupID:", groupID, "userIDList:", userIDList)
				if err := ResetUserGroupMinSeq(operationID, groupID, userIDList); err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID, userIDList)
				}
				//if err := checkMaxSeqWithMongo(operationID, groupID, constant.ReadDiffusion); err != nil {
				//	log.NewError(operationID, utils.GetSelfFuncName(), groupID, err)
				//}
			}
		} else {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		}

		log.NewInfo(operationID, "====================== start del cron finished ======================")
	})
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
