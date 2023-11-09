// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tools

import (
	"context"
	"fmt"
	"math"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"

	"math/rand"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/mw"
	"github.com/OpenIMSDK/tools/tx"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient/notification"
)

type MsgTool struct {
	msgDatabase           controller.CommonMsgDatabase
	conversationDatabase  controller.ConversationDatabase
	userDatabase          controller.UserDatabase
	groupDatabase         controller.GroupDatabase
	msgNotificationSender *notification.MsgNotificationSender
}

func NewMsgTool(msgDatabase controller.CommonMsgDatabase, userDatabase controller.UserDatabase,
	groupDatabase controller.GroupDatabase, conversationDatabase controller.ConversationDatabase, msgNotificationSender *notification.MsgNotificationSender,
) *MsgTool {
	return &MsgTool{
		msgDatabase:           msgDatabase,
		userDatabase:          userDatabase,
		groupDatabase:         groupDatabase,
		conversationDatabase:  conversationDatabase,
		msgNotificationSender: msgNotificationSender,
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
	discov, err := kdisc.NewDiscoveryRegister(config.Config.Envs.Discovery)
	/*
		discov, err := zookeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
			zookeeper.WithFreq(time.Hour), zookeeper.WithRoundRobin(), zookeeper.WithUserNameAndPassword(config.Config.Zookeeper.Username,
				config.Config.Zookeeper.Password), zookeeper.WithTimeout(10), zookeeper.WithLogger(log.NewZkLogger()))*/
	if err != nil {
		return nil, err
	}
	discov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	userDB := relation.NewUserGorm(db)
	msgDatabase := controller.InitCommonMsgDatabase(rdb, mongo.GetDatabase())
	userMongoDB := unrelation.NewUserMongoDriver(mongo.GetDatabase())
	userDatabase := controller.NewUserDatabase(
		userDB,
		cache.NewUserCacheRedis(rdb, relation.NewUserGorm(db), cache.GetDefaultOpt()),
		tx.NewGorm(db),
		userMongoDB,
	)
	groupDatabase := controller.InitGroupDatabase(db, rdb, mongo.GetDatabase(), nil)
	conversationDatabase := controller.NewConversationDatabase(
		relation.NewConversationGorm(db),
		cache.NewConversationRedis(rdb, cache.GetDefaultOpt(), relation.NewConversationGorm(db)),
		tx.NewGorm(db),
	)
	msgRpcClient := rpcclient.NewMessageRpcClient(discov)
	msgNotificationSender := notification.NewMsgNotificationSender(rpcclient.WithRpcClient(&msgRpcClient))
	msgTool := NewMsgTool(msgDatabase, userDatabase, groupDatabase, conversationDatabase, msgNotificationSender)
	return msgTool, nil
}

//func (c *MsgTool) AllConversationClearMsgAndFixSeq() {
//	ctx := mcontext.NewCtx(utils.GetSelfFuncName())
//	log.ZInfo(ctx, "============================ start del cron task ============================")
//	conversationIDs, err := c.conversationDatabase.GetAllConversationIDs(ctx)
//	if err != nil {
//		log.ZError(ctx, "GetAllConversationIDs failed", err)
//		return
//	}
//	for _, conversationID := range conversationIDs {
//		conversationIDs = append(conversationIDs, utils.GetNotificationConversationIDByConversationID(conversationID))
//	}
//	c.ClearConversationsMsg(ctx, conversationIDs)
//	log.ZInfo(ctx, "============================ start del cron finished ============================")
//}

func (c *MsgTool) AllConversationClearMsgAndFixSeq() {
	ctx := mcontext.NewCtx(utils.GetSelfFuncName())
	log.ZInfo(ctx, "============================ start del cron task ============================")
	num, err := c.conversationDatabase.GetAllConversationIDsNumber(ctx)
	if err != nil {
		log.ZError(ctx, "GetAllConversationIDsNumber failed", err)
		return
	}
	const batchNum = 50
	log.ZDebug(ctx, "GetAllConversationIDsNumber", "num", num)
	if num == 0 {
		return
	}
	count := int(num/batchNum + num/batchNum/2)
	if count < 1 {
		count = 1
	}
	maxPage := 1 + num/batchNum
	if num%batchNum != 0 {
		maxPage++
	}
	for i := 0; i < count; i++ {
		pageNumber := rand.Int63() % maxPage
		conversationIDs, err := c.conversationDatabase.PageConversationIDs(ctx, int32(pageNumber), batchNum)
		if err != nil {
			log.ZError(ctx, "PageConversationIDs failed", err, "pageNumber", pageNumber)
			continue
		}
		log.ZDebug(ctx, "PageConversationIDs failed", "pageNumber", pageNumber, "conversationIDsNum", len(conversationIDs), "conversationIDs", conversationIDs)
		if len(conversationIDs) == 0 {
			continue
		}
		c.ClearConversationsMsg(ctx, conversationIDs)
	}
	log.ZInfo(ctx, "============================ start del cron finished ============================")
}

func (c *MsgTool) ClearConversationsMsg(ctx context.Context, conversationIDs []string) {
	for _, conversationID := range conversationIDs {
		if err := c.msgDatabase.DeleteConversationMsgsAndSetMinSeq(ctx, conversationID, int64(config.Config.RetainChatRecords*24*60*60)); err != nil {
			log.ZError(ctx, "DeleteUserSuperGroupMsgsAndSetMinSeq failed", err, "conversationID", conversationID, "DBRetainChatRecords", config.Config.RetainChatRecords)
		}
		if err := c.checkMaxSeq(ctx, conversationID); err != nil {
			log.ZError(ctx, "fixSeq failed", err, "conversationID", conversationID)
		}
	}
}

func (c *MsgTool) checkMaxSeqWithMongo(ctx context.Context, conversationID string, maxSeqCache int64) error {
	minSeqMongo, maxSeqMongo, err := c.msgDatabase.GetMongoMaxAndMinSeq(ctx, conversationID)
	if err != nil {
		return err
	}
	if math.Abs(float64(maxSeqMongo-maxSeqCache)) > 10 {
		log.ZError(ctx, "cache max seq and mongo max seq is diff > 10", nil, "maxSeqMongo", maxSeqMongo, "minSeqMongo", minSeqMongo, "maxSeqCache", maxSeqCache, "conversationID", conversationID)
	}
	return nil
}

func (c *MsgTool) checkMaxSeq(ctx context.Context, conversationID string) error {
	maxSeq, err := c.msgDatabase.GetMaxSeq(ctx, conversationID)
	if err != nil {
		if errs.Unwrap(err) == redis.Nil {
			return nil
		}
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
