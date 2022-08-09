package cronTask

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"github.com/robfig/cron/v3"
)

const cronTaskOperationID = "cronTaskOperationID-"

func StartCronTask() {
	log.NewInfo(utils.OperationIDGenerator(), "start cron task")
	c := cron.New()
	_, err := c.AddFunc("30 3-6,20-23 * * *", func() {
		operationID := getCronTaskOperationID()
		userIDList, err := im_mysql_model.SelectAllUserID()
		if err == nil {
			log.NewDebug(operationID, utils.GetSelfFuncName(), "userIDList: ", userIDList)
			for _, userID := range userIDList {
				if err := DeleteMongoMsgAndResetRedisSeq(operationID, userID); err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), userID)
				}
			}
		} else {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		}

		workingGroupIDList, err := im_mysql_model.GetGroupIDListByGroupType(constant.WorkingGroup)
		if err == nil {
			for _, groupID := range workingGroupIDList {
				userIDList, err = rocksCache.GetGroupMemberIDListFromCache(groupID)
				if err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID)
					continue
				}
				log.NewDebug(operationID, utils.GetSelfFuncName(), "groupID:", groupID, "userIDList:", userIDList)
				for _, userID := range userIDList {
					if err := ResetUserGroupMinSeq(operationID, groupID, userID); err != nil {
						log.NewError(operationID, utils.GetSelfFuncName(), operationID, groupID, userID, err.Error())
					}
				}
			}
		} else {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
			return
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}

func getCronTaskOperationID() string {
	return cronTaskOperationID + utils.OperationIDGenerator()
}
