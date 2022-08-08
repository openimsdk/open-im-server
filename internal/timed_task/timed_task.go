package timedTask

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"github.com/robfig/cron/v3"
)

func main() {
	log.NewInfo(utils.OperationIDGenerator(), "start cron task")
	c := cron.New()
	_, err := c.AddFunc("30 3-6,20-23 * * *", func() {
		operationID := utils.OperationIDGenerator()
		if err := DeleteMongoMsgAndResetRedisSeq(operationID, "", constant.ReadDiffusion); err != nil {
			log.NewError(operationID)
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()
}
