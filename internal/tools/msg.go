package tools

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/go-redis/redis/v8"
)

type MsgTool struct {
	msgDatabase          controller.CommonMsgDatabase
	conversationDatabase controller.ConversationDataBase
	userDatabase         controller.UserDatabase
	groupDatabase        controller.GroupDatabase
}

var errSeq = errors.New("cache max seq and mongo max seq is diff > 10")

func NewMsgTool(msgDatabase controller.CommonMsgDatabase, userDatabase controller.UserDatabase, groupDatabase controller.GroupDatabase) *MsgTool {
	return &MsgTool{
		msgDatabase:   msgDatabase,
		userDatabase:  userDatabase,
		groupDatabase: groupDatabase,
	}
}

func InitMsgTool() (*MsgTool, error) {
	rdb, err := cache.NewRedis()
	if err != nil {
		return nil, err
	}
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return nil, err
	}
	db, err := relation.NewGormDB()
	if err != nil {
		return nil, err
	}
	userDB := relation.NewUserGorm(db)
	msgDatabase := controller.InitCommonMsgDatabase(rdb, mongo.GetDatabase())
	userDatabase := controller.NewUserDatabase(userDB, cache.NewUserCacheRedis(rdb, relation.NewUserGorm(db), cache.GetDefaultOpt()), tx.NewGorm(db))
	groupDatabase := controller.InitGroupDatabase(db, rdb, mongo.GetDatabase())
	msgTool := NewMsgTool(msgDatabase, userDatabase, groupDatabase)
	return msgTool, nil
}

func (c *MsgTool) AllConversationClearMsgAndFixSeq() {
	ctx := mcontext.NewCtx(utils.GetSelfFuncName())
	log.ZInfo(ctx, "============================ start del cron task ============================")
	conversationIDs, err := c.conversationDatabase.GetAllConversationIDs(ctx)
	if err != nil {
		log.ZError(ctx, "GetAllConversationIDs failed", err)
		return
	}
	c.ClearSuperGroupMsg(ctx, conversationIDs)
	log.ZInfo(ctx, "============================ start del cron finished ============================")
}

func (c *MsgTool) ClearUsersMsg(ctx context.Context, userIDs []string) {
	for _, userID := range userIDs {
		if err := c.msgDatabase.DeleteConversationMsgsAndSetMinSeq(ctx, userID, int64(config.Config.Mongo.DBRetainChatRecords*24*60*60)); err != nil {
			log.ZError(ctx, "DeleteUserMsgsAndSetMinSeq failed", err, "userID", userID, "DBRetainChatRecords", config.Config.Mongo.DBRetainChatRecords)
		}
		maxSeqCache, maxSeqMongo, err := c.GetAndFixUserSeqs(ctx, userID)
		if err != nil {
			continue
		}
		c.CheckMaxSeqWithMongo(ctx, userID, maxSeqCache, maxSeqMongo)
	}
}

func (c *MsgTool) ClearSuperGroupMsg(ctx context.Context, superGroupIDs []string) {
	for _, groupID := range superGroupIDs {
		userIDs, err := c.groupDatabase.FindGroupMemberUserID(ctx, groupID)
		if err != nil {
			log.ZError(ctx, "ClearSuperGroupMsg failed", err, "groupID", groupID)
			continue
		}
		if err := c.msgDatabase.DeleteConversationMsgsAndSetMinSeq(ctx, groupID, int64(config.Config.Mongo.DBRetainChatRecords*24*60*60)); err != nil {
			log.ZError(ctx, "DeleteUserSuperGroupMsgsAndSetMinSeq failed", err, "groupID", groupID, "userID", userIDs, "DBRetainChatRecords", config.Config.Mongo.DBRetainChatRecords)
		}
		if err := c.fixGroupSeq(ctx, groupID, userIDs); err != nil {
			log.ZError(ctx, "fixGroupSeq failed", err, "groupID", groupID, "userID", userIDs)
		}
	}
}

func (c *MsgTool) FixGroupSeq(ctx context.Context, groupID string) error {
	userIDs, err := c.groupDatabase.FindGroupMemberUserID(ctx, groupID)
	if err != nil {
		return err
	}
	return c.fixGroupSeq(ctx, groupID, userIDs)
}

func (c *MsgTool) fixGroupSeq(ctx context.Context, groupID string, userIDs []string) error {
	_, maxSeqMongo, _, maxSeqCache, err := c.msgDatabase.GetConversationMinMaxSeqInMongoAndCache(ctx, groupID)
	if err != nil {
		if err == unrelation.ErrMsgNotFound {
			return nil
		}
		return err
	}
	for _, userID := range userIDs {
		if _, err := c.GetAndFixGroupUserSeq(ctx, userID, groupID, maxSeqCache); err != nil {
			continue
		}
	}
	if err := c.CheckMaxSeqWithMongo(ctx, groupID, maxSeqCache, maxSeqMongo); err != nil {
		log.ZWarn(ctx, "cache max seq and mongo max seq is diff > 10", err, "groupID", groupID, "maxSeqCache", maxSeqCache, "maxSeqMongo", maxSeqMongo, "constant.WriteDiffusion", constant.WriteDiffusion)
	}
	return nil
}

func (c *MsgTool) GetAndFixUserSeqs(ctx context.Context, userID string) (maxSeqCache, maxSeqMongo int64, err error) {
	minSeqMongo, maxSeqMongo, minSeqCache, maxSeqCache, err := c.msgDatabase.GetConversationMinMaxSeqInMongoAndCache(ctx, userID)
	if err != nil {
		if err != unrelation.ErrMsgNotFound {
			log.ZError(ctx, "GetUserMinMaxSeqInMongoAndCache failed", err, "userID", userID)
		}
		return 0, 0, err
	}
	log.ZDebug(ctx, "userID", userID, "minSeqMongo", minSeqMongo, "maxSeqMongo", maxSeqMongo, "minSeqCache", minSeqCache, "maxSeqCache", maxSeqCache)
	if minSeqCache > maxSeqCache {
		if err := c.msgDatabase.SetMinSeq(ctx, userID, maxSeqCache); err != nil {
			log.ZError(ctx, "SetUserMinSeq failed", err, "userID", userID, "minSeqCache", minSeqCache, "maxSeqCache", maxSeqCache)
		} else {
			log.ZInfo(ctx, "SetUserMinSeq success", "userID", userID, "minSeqCache", minSeqCache, "maxSeqCache", maxSeqCache)
		}
	}
	return maxSeqCache, maxSeqMongo, nil
}

func (c *MsgTool) GetAndFixGroupUserSeq(ctx context.Context, userID string, groupID string, maxSeqCache int64) (minSeqCache int64, err error) {
	minSeqCache, err = c.msgDatabase.GetMinSeq(ctx, groupID)
	if err != nil {
		log.ZError(ctx, "GetGroupUserMinSeq failed", err, "groupID", groupID, "userID", userID)
		return 0, err
	}
	if minSeqCache > maxSeqCache {
		if err := c.msgDatabase.SetConversationUserMinSeq(ctx, groupID, userID, maxSeqCache); err != nil {
			log.ZError(ctx, "SetGroupUserMinSeq failed", err, "groupID", groupID, "userID", userID, "minSeqCache", minSeqCache, "maxSeqCache", maxSeqCache)
		} else {
			log.ZInfo(ctx, "SetGroupUserMinSeq success", "groupID", groupID, "userID", userID, "minSeqCache", minSeqCache, "maxSeqCache", maxSeqCache)
		}
	}
	return minSeqCache, nil
}

func (c *MsgTool) CheckMaxSeqWithMongo(ctx context.Context, conversationID string, maxSeqCache, maxSeqMongo int64) error {
	if math.Abs(float64(maxSeqMongo-maxSeqCache)) > 10 {
		return errSeq
	}
	return nil
}

func (c *MsgTool) ShowUserSeqs(ctx context.Context, userID string) {

}

func (c *MsgTool) ShowSuperGroupSeqs(ctx context.Context, groupID string) {

}

func (c *MsgTool) ShowSuperGroupUserSeqs(ctx context.Context, groupID, userID string) {

}

func (c *MsgTool) FixAllSeq(ctx context.Context) error {
	userIDs, err := c.userDatabase.GetAllUserID(ctx)
	if err != nil {
		return err
	}
	for _, userID := range userIDs {
		userCurrentMinSeq, err := c.msgDatabase.GetMinSeq(ctx, userID)
		if err != nil && err != redis.Nil {
			continue
		}
		userCurrentMaxSeq, err := c.msgDatabase.GetMaxSeq(ctx, userID)
		if err != nil && err != redis.Nil {
			continue
		}
		if userCurrentMinSeq > userCurrentMaxSeq {
			if err = c.msgDatabase.SetMinSeq(ctx, userID, userCurrentMaxSeq); err != nil {
				fmt.Println("SetUserMinSeq failed", userID, userCurrentMaxSeq)
			}
			fmt.Println("fix", userID, userCurrentMaxSeq)
		}
	}
	fmt.Println("fix users seq success")
	groupIDs, err := c.groupDatabase.GetGroupIDsByGroupType(ctx, constant.WorkingGroup)
	if err != nil {
		return err
	}
	for _, groupID := range groupIDs {
		maxSeq, err := c.msgDatabase.GetMaxSeq(ctx, groupID)
		if err != nil {
			fmt.Println("GetGroupMaxSeq failed", groupID)
			continue
		}
		userIDs, err := c.groupDatabase.FindGroupMemberUserID(ctx, groupID)
		if err != nil {
			fmt.Println("get groupID", groupID, "failed, try again later")
			continue
		}
		for _, userID := range userIDs {
			userMinSeq, err := c.msgDatabase.GetMinSeq(ctx, groupID)
			if err != nil && err != redis.Nil {
				fmt.Println("GetGroupUserMinSeq failed", groupID, userID)
				continue
			}
			if userMinSeq > maxSeq {
				if err = c.msgDatabase.SetMinSeq(ctx, groupID, maxSeq); err != nil {
					fmt.Println("SetGroupUserMinSeq failed", err.Error(), groupID, userID, maxSeq)
				}
				fmt.Println("fix", groupID, userID, maxSeq, userMinSeq)
			}
		}
	}
	fmt.Println("fix all seq finished")
	return nil
}
