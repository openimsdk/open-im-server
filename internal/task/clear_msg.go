package task

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/controller"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math"
)

type msgTool struct {
	msgInterface   controller.MsgDatabase
	userInterface  controller.UserDatabase
	groupInterface controller.GroupDatabase
}

func (c *msgTool) getCronTaskOperationID() string {
	return cronTaskOperationID + utils.OperationIDGenerator()
}

func (c *msgTool) ClearAll() {
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
	superGroupIDList, err := c.groupInterface.GetGroupIDsByGroupType(ctx, constant.WorkingGroup)
	if err == nil {
		c.ClearSuperGroupMsg(ctx, superGroupIDList)
	} else {
		log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
	}
	log.NewInfo(operationID, "============================ start del cron finished ============================")
}

func (c *msgTool) ClearUsersMsg(ctx context.Context, userIDList []string) {
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

func (c *msgTool) ClearSuperGroupMsg(ctx context.Context, superGroupIDList []string) {
	for _, groupID := range superGroupIDList {
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
		for _, userID := range userIDs {
			minSeqCache, err := c.msgInterface.GetGroupUserMinSeq(ctx, groupID, userID)
			if err != nil {
				log.NewError(tracelog.GetOperationID(ctx), "GetGroupUserMinSeq failed", groupID, userID)
				continue
			}
			c.FixGroupUserSeq(ctx, userID, groupID, minSeqCache, maxSeqCache)

		}
		c.CheckMaxSeqWithMongo(ctx, groupID, maxSeqCache, maxSeqMongo, constant.WriteDiffusion)
	}
}

func (c *msgTool) FixUserSeq(ctx context.Context, userID string, minSeqCache, maxSeqCache int64) {
	if minSeqCache > maxSeqCache {
		if err := c.msgInterface.SetUserMinSeq(ctx, userID, maxSeqCache); err != nil {
			log.NewError(tracelog.GetOperationID(ctx), "SetUserMinSeq failed", userID, minSeqCache, maxSeqCache)
		} else {
			log.NewWarn(tracelog.GetOperationID(ctx), "SetUserMinSeq success", userID, minSeqCache, maxSeqCache)
		}
	}
}

func (c *msgTool) FixGroupUserSeq(ctx context.Context, userID string, groupID string, minSeqCache, maxSeqCache int64) {
	if minSeqCache > maxSeqCache {
		if err := c.msgInterface.SetGroupUserMinSeq(ctx, groupID, userID, maxSeqCache); err != nil {
			log.NewError(tracelog.GetOperationID(ctx), "SetGroupUserMinSeq failed", userID, minSeqCache, maxSeqCache)
		} else {
			log.NewWarn(tracelog.GetOperationID(ctx), "SetGroupUserMinSeq success", userID, minSeqCache, maxSeqCache)
		}
	}
}

func (c *msgTool) CheckMaxSeqWithMongo(ctx context.Context, sourceID string, maxSeqCache, maxSeqMongo int64, diffusionType int) {
	if math.Abs(float64(maxSeqMongo-maxSeqCache)) > 10 {
		log.NewWarn(tracelog.GetOperationID(ctx), "cache max seq and mongo max seq is diff > 10", sourceID, maxSeqCache, maxSeqMongo, diffusionType)
	}
}

func (c *msgTool) FixAllSeq(ctx context.Context) {
	userIDs, err := c.userInterface.GetAllUserID(ctx)
	if err != nil {
		panic(err.Error())
	}
	for _, userID := range userIDs {
		userCurrentMinSeq, err := c.msgInterface.GetUserMinSeq(ctx, userID)
		if err != nil && err != redis.Nil {
			continue
		}
		userCurrentMaxSeq, err := c.msgInterface.GetUserMaxSeq(ctx, userID)
		if err != nil && err != redis.Nil {
			continue
		}
		if userCurrentMinSeq > userCurrentMaxSeq {
			if err = c.msgInterface.SetUserMinSeq(ctx, userID, userCurrentMaxSeq); err != nil {
				fmt.Println("SetUserMinSeq failed", userID, userCurrentMaxSeq)
			}
			fmt.Println("fix", userID, userCurrentMaxSeq)
		}
	}
	fmt.Println("fix users seq success")

	groupIDs, err := c.groupInterface.GetGroupIDsByGroupType(ctx, constant.WorkingGroup)
	if err != nil {
		panic(err.Error())
	}
	for _, groupID := range groupIDs {
		maxSeq, err := c.msgInterface.GetGroupMaxSeq(ctx, groupID)
		if err != nil {
			fmt.Println("GetGroupMaxSeq failed", groupID)
			continue
		}
		userIDs, err := c.groupInterface.FindGroupMemberUserID(ctx, groupID)
		if err != nil {
			fmt.Println("get groupID", groupID, "failed, try again later")
			continue
		}
		for _, userID := range userIDs {
			userMinSeq, err := c.msgInterface.GetGroupUserMinSeq(ctx, groupID, userID)
			if err != nil && err != redis.Nil {
				fmt.Println("GetGroupUserMinSeq failed", groupID, userID)
				continue
			}
			if userMinSeq > maxSeq {
				if err = c.msgInterface.SetGroupUserMinSeq(ctx, groupID, userID, maxSeq); err != nil {
					fmt.Println("SetGroupUserMinSeq failed", err.Error(), groupID, userID, maxSeq)
				}
				fmt.Println("fix", groupID, userID, maxSeq, userMinSeq)
			}
		}
	}
	fmt.Println("fix all seq finished")
}
