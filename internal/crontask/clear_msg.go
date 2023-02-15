package cronTask

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/mongo"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	sdkws "Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
)

type SeqCheckInterface interface {
	ClearAll() error
}

type ClearMsgCronTask struct {
	msgModel   controller.MsgInterface
	userModel  controller.UserInterface
	groupModel controller.GroupInterface
	cache      cache.Cache
}

func (c *ClearMsgCronTask) getCronTaskOperationID() string {
	return cronTaskOperationID + utils.OperationIDGenerator()
}

func (c *ClearMsgCronTask) ClearAll() {
	operationID := c.getCronTaskOperationID()
	ctx := context.Background()
	tracelog.SetOperationID(ctx, operationID)
	log.NewInfo(operationID, "========================= start del cron task =========================")
	var err error
	userIDList, err := c.userModel.GetAllUserID(ctx)
	if err == nil {
		c.StartClearMsg(operationID, userIDList)
	} else {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}
	// working group msg clear
	workingGroupIDList, err := im_mysql_model.GetGroupIDListByGroupType(constant.WorkingGroup)
	if err == nil {
		c.StartClearWorkingGroupMsg(operationID, workingGroupIDList)
	} else {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}

	log.NewInfo(operationID, "========================= start del cron finished =========================")
}

func (c *ClearMsgCronTask) StartClearMsg(operationID string, userIDList []string) {
	log.NewDebug(operationID, utils.GetSelfFuncName(), "userIDList: ", userIDList)
	for _, userID := range userIDList {
		if err := DeleteUserMsgsAndSetMinSeq(operationID, userID); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), userID)
		}
		if err := checkMaxSeqWithMongo(operationID, userID, constant.WriteDiffusion); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), userID, err)
		}
	}
}

func (c *ClearMsgCronTask) StartClearWorkingGroupMsg(operationID string, workingGroupIDList []string) {
	log.NewDebug(operationID, utils.GetSelfFuncName(), "workingGroupIDList: ", workingGroupIDList)
	for _, groupID := range workingGroupIDList {
		userIDList, err := rocksCache.GetGroupMemberIDListFromCache(groupID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID)
			continue
		}
		log.NewDebug(operationID, utils.GetSelfFuncName(), "groupID:", groupID, "workingGroupIDList:", userIDList)
		if err := DeleteUserSuperGroupMsgsAndSetMinSeq(operationID, groupID, userIDList); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID, userIDList)
		}
		if err := checkMaxSeqWithMongo(operationID, groupID, constant.ReadDiffusion); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), groupID, err)
		}
	}
}

func checkMaxSeqWithMongo(operationID, sourceID string, diffusionType int) error {
	var seqRedis uint64
	var err error
	if diffusionType == constant.WriteDiffusion {
		seqRedis, err = db.DB.GetUserMaxSeq(sourceID)
	} else {
		seqRedis, err = db.DB.GetGroupMaxSeq(sourceID)
	}
	if err != nil {
		if err == goRedis.Nil {
			return nil
		}
		return utils.Wrap(err, "GetUserMaxSeq failed")
	}
	msg, err := db.DB.GetNewestMsg(sourceID)
	if err != nil {
		return utils.Wrap(err, "GetNewestMsg failed")
	}
	if msg == nil {
		return nil
	}
	if math.Abs(float64(msg.Seq-uint32(seqRedis))) > 10 {
		log.NewWarn(operationID, utils.GetSelfFuncName(), "seqMongo, seqRedis", msg.Seq, seqRedis, sourceID, "redis maxSeq is different with msg.Seq > 10", "status: ", msg.Status, msg.SendTime)
	} else {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "seqMongo, seqRedis", msg.Seq, seqRedis, sourceID, "seq and msg OK", "status:", msg.Status, msg.SendTime)
	}
	return nil
}
