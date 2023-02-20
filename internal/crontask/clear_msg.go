package cronTask

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"math"
)

type ClearMsgTool struct {
	msgInterface   controller.MsgInterface
	userInterface  controller.UserInterface
	groupInterface controller.GroupInterface
}

func (c *ClearMsgTool) getCronTaskOperationID() string {
	return cronTaskOperationID + utils.OperationIDGenerator()
}

func (c *ClearMsgTool) ClearAll() {
	operationID := c.getCronTaskOperationID()
	ctx := context.Background()
	tracelog.SetOperationID(ctx, operationID)
	log.NewInfo(operationID, "============================ start del cron task ============================")
	var err error
	userIDList, err := c.userInterface.GetAllUserID(ctx)
	if err == nil {
		c.ClearUsersMsg(ctx, userIDList)
	} else {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}
	// working group msg clear
	workingGroupIDList, err := c.groupInterface.GetGroupIDsByGroupType(ctx, constant.WorkingGroup)
	if err == nil {
		c.ClearSuperGroupMsg(ctx, workingGroupIDList)
	} else {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}
	log.NewInfo(operationID, "============================ start del cron finished ============================")
}

func (c *ClearMsgTool) ClearUsersMsg(ctx context.Context, userIDList []string) {
	for _, userID := range userIDList {
		if err := c.msgInterface.DeleteUserMsgsAndSetMinSeq(ctx, userID, int64(config.Config.Mongo.DBRetainChatRecords*24*60*60)); err != nil {
			log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), err.Error(), userID)
		}
		_, maxSeqMongo, minSeqCache, maxSeqCache, err := c.msgInterface.GetUserMinMaxSeqInMongoAndCache(ctx, userID)
		if err != nil {
			log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), err.Error(), "GetUserMinMaxSeqInMongoAndCache failed", userID)
			continue
		}
		c.FixUserSeq(ctx, userID, minSeqCache, maxSeqCache)
		c.CheckMaxSeqWithMongo(ctx, userID, maxSeqCache, maxSeqMongo, constant.WriteDiffusion)
	}
}

func (c *ClearMsgTool) ClearSuperGroupMsg(ctx context.Context, workingGroupIDList []string) {
	for _, groupID := range workingGroupIDList {
		userIDs, err := c.groupInterface.FindGroupMemberUserID(ctx, groupID)
		if err != nil {
			log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "FindGroupMemberUserID", err.Error(), groupID)
			continue
		}
		if err := c.msgInterface.DeleteUserSuperGroupMsgsAndSetMinSeq(ctx, groupID, userIDs, int64(config.Config.Mongo.DBRetainChatRecords*24*60*60)); err != nil {
			log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), err.Error(), "DeleteUserSuperGroupMsgsAndSetMinSeq failed", groupID, userIDs, config.Config.Mongo.DBRetainChatRecords)
		}
		_, maxSeqMongo, maxSeqCache, err := c.msgInterface.GetSuperGroupMinMaxSeqInMongoAndCache(ctx, groupID)
		if err != nil {
			log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), err.Error(), "GetUserMinMaxSeqInMongoAndCache failed", groupID)
			continue
		}
		c.FixGroupUserSeq(ctx, userIDs, groupID)
		c.CheckMaxSeqWithMongo(ctx, groupID, maxSeqCache, maxSeqMongo, constant.WriteDiffusion)
	}
}

func (c *ClearMsgTool) FixUserSeq(ctx context.Context, userID string, minSeqCache, maxSeqCache int64) {
	if minSeqCache > maxSeqCache {
		if err := c.msgInterface.SetUserMinSeq(ctx, userID, maxSeqCache); err != nil {
			log.NewError(tracelog.GetOperationID(ctx), "SetUserMinSeq failed", userID, minSeqCache, maxSeqCache)
		} else {
			log.NewWarn(tracelog.GetOperationID(ctx), "SetUserMinSeq success", userID, minSeqCache, maxSeqCache)
		}
	}
}

func (c *ClearMsgTool) FixGroupUserSeq(ctx context.Context, userID string, groupID string, minSeqCache, maxSeqCache int64) {
	if minSeqCache > maxSeqCache {
		if err := c.msgInterface.SetGroupUserMinSeq(ctx, groupID, userID, maxSeqCache); err != nil {
			log.NewError(tracelog.GetOperationID(ctx), "SetGroupUserMinSeq failed", userID, minSeqCache, maxSeqCache)
		} else {
			log.NewWarn(tracelog.GetOperationID(ctx), "SetGroupUserMinSeq success", userID, minSeqCache, maxSeqCache)
		}
	}
}

func (c *ClearMsgTool) CheckMaxSeqWithMongo(ctx context.Context, sourceID string, maxSeqCache, maxSeqMongo int64, diffusionType int) {
	if math.Abs(float64(maxSeqMongo-maxSeqCache)) > 10 {
		log.NewWarn(tracelog.GetOperationID(ctx), "cache max seq and mongo max seq is diff > 10", sourceID, maxSeqCache, maxSeqMongo, diffusionType)
	}
}
