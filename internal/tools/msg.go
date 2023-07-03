package tools

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type MsgTool struct {
	msgDatabase          controller.CommonMsgDatabase
	conversationDatabase controller.ConversationDatabase
	userDatabase         controller.UserDatabase
	groupDatabase        controller.GroupDatabase
}

var errSeq = errors.New("cache max seq and mongo max seq is diff > 10")

func NewMsgTool(
	msgDatabase controller.CommonMsgDatabase,
	userDatabase controller.UserDatabase,
	groupDatabase controller.GroupDatabase,
	conversationDatabase controller.ConversationDatabase,
) *MsgTool {
	return &MsgTool{
		msgDatabase:          msgDatabase,
		userDatabase:         userDatabase,
		groupDatabase:        groupDatabase,
		conversationDatabase: conversationDatabase,
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
	userDatabase := controller.NewUserDatabase(
		userDB,
		cache.NewUserCacheRedis(rdb, relation.NewUserGorm(db), cache.GetDefaultOpt()),
		tx.NewGorm(db),
	)
	groupDatabase := controller.InitGroupDatabase(db, rdb, mongo.GetDatabase())
	conversationDatabase := controller.NewConversationDatabase(
		relation.NewConversationGorm(db),
		cache.NewConversationRedis(rdb, cache.GetDefaultOpt(), relation.NewConversationGorm(db)),
		tx.NewGorm(db),
	)
	msgTool := NewMsgTool(msgDatabase, userDatabase, groupDatabase, conversationDatabase)
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
	for _, conversationID := range conversationIDs {
		conversationIDs = append(conversationIDs, utils.GetNotificationConversationIDByConversationID(conversationID))
	}
	c.ClearConversationsMsg(ctx, conversationIDs)
	log.ZInfo(ctx, "============================ start del cron finished ============================")
}

func (c *MsgTool) ClearConversationsMsg(ctx context.Context, conversationIDs []string) {
	for _, conversationID := range conversationIDs {
		if err := c.msgDatabase.DeleteConversationMsgsAndSetMinSeq(ctx, conversationID, int64(config.Config.RetainChatRecords*24*60*60)); err != nil {
			log.ZError(
				ctx,
				"DeleteUserSuperGroupMsgsAndSetMinSeq failed",
				err,
				"conversationID",
				conversationID,
				"DBRetainChatRecords",
				config.Config.RetainChatRecords,
			)
		}
		if err := c.checkMaxSeq(ctx, conversationID); err != nil {
			log.ZError(ctx, "fixSeq failed", err, "conversationID", conversationID)
		}

	}
}

func (c *MsgTool) checkMaxSeqWithMongo(ctx context.Context, conversationID string, maxSeqCache int64) error {
	maxSeqMongo, _, err := c.msgDatabase.GetMongoMaxAndMinSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	if math.Abs(float64(maxSeqMongo-maxSeqCache)) > 10 {
		return errSeq
	}
	return nil
}

func (c *MsgTool) checkMaxSeq(ctx context.Context, conversationID string) error {
	maxSeq, err := c.msgDatabase.GetMaxSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	if err := c.checkMaxSeqWithMongo(ctx, conversationID, maxSeq); err != nil {
		return err
	}
	return nil
}

func (c *MsgTool) FixAllSeq(ctx context.Context) error {
	conversationIDs, err := c.conversationDatabase.GetAllConversationIDs(ctx)
	if err != nil {
		log.ZError(ctx, "GetAllConversationIDs failed", err)
		return err
	}
	for _, conversationID := range conversationIDs {
		conversationIDs = append(conversationIDs, utils.GetNotificationConversationIDByConversationID(conversationID))
	}
	for _, conversationID := range conversationIDs {
		if err := c.checkMaxSeq(ctx, conversationID); err != nil {
			log.ZWarn(ctx, "fixSeq failed", err, "conversationID", conversationID)
		}
	}
	fmt.Println("fix all seq finished")
	return nil
}
