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
	"github.com/openimsdk/open-im-server/v3/internal/rpc/msg"
	"math"
	"math/rand"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/util/conversationutil"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/stringutil"
	"github.com/redis/go-redis/v9"
)

type MsgTool struct {
	msgDatabase           controller.CommonMsgDatabase
	conversationDatabase  controller.ConversationDatabase
	userDatabase          controller.UserDatabase
	groupDatabase         controller.GroupDatabase
	msgNotificationSender *msg.MsgNotificationSender
	config                *CronTaskConfig
}

func NewMsgTool(msgDatabase controller.CommonMsgDatabase, userDatabase controller.UserDatabase,
	groupDatabase controller.GroupDatabase, conversationDatabase controller.ConversationDatabase,
	msgNotificationSender *msg.MsgNotificationSender, config *CronTaskConfig,
) *MsgTool {
	return &MsgTool{
		msgDatabase:           msgDatabase,
		userDatabase:          userDatabase,
		groupDatabase:         groupDatabase,
		conversationDatabase:  conversationDatabase,
		msgNotificationSender: msgNotificationSender,
		config:                config,
	}
}

func InitMsgTool(ctx context.Context, config *CronTaskConfig) (*MsgTool, error) {
	ch := make(chan int)
	<-ch
	//mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	//if err != nil {
	//	return nil, err
	//}
	//rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	//if err != nil {
	//	return nil, err
	//}
	//discov, err := kdisc.NewDiscoveryRegister(&config.ZookeeperConfig, &config.Share)
	//if err != nil {
	//	return nil, err
	//}
	//discov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	//userDB, err := mgo.NewUserMongo(mgocli.GetDB())
	//if err != nil {
	//	return nil, err
	//}
	////msgDatabase, err := controller.InitCommonMsgDatabase(rdb, mgocli.GetDB(), config)
	//if err != nil {
	//	return nil, err
	//}
	//userMongoDB := mgo.NewUserMongoDriver(mgocli.GetDB())
	//userDatabase := controller.NewUserDatabase(
	//	userDB,
	//	cache.NewUserCacheRedis(rdb, userDB, cache.GetDefaultOpt()),
	//	mgocli.GetTx(),
	//	userMongoDB,
	//)
	//groupDB, err := mgo.NewGroupMongo(mgocli.GetDB())
	//if err != nil {
	//	return nil, err
	//}
	//groupMemberDB, err := mgo.NewGroupMember(mgocli.GetDB())
	//if err != nil {
	//	return nil, err
	//}
	//groupRequestDB, err := mgo.NewGroupRequestMgo(mgocli.GetDB())
	//if err != nil {
	//	return nil, err
	//}
	//conversationDB, err := mgo.NewConversationMongo(mgocli.GetDB())
	//if err != nil {
	//	return nil, err
	//}
	//groupDatabase := controller.NewGroupDatabase(rdb, groupDB, groupMemberDB, groupRequestDB, mgocli.GetTx(), nil)
	//conversationDatabase := controller.NewConversationDatabase(
	//	conversationDB,
	//	cache.NewConversationRedis(rdb, cache.GetDefaultOpt(), conversationDB),
	//	mgocli.GetTx(),
	//)
	//msgRpcClient := rpcclient.NewMessageRpcClient(discov, config.Share.RpcRegisterName.Msg)
	//msgNotificationSender := notification.NewMsgNotificationSender(config, rpcclient.WithRpcClient(&msgRpcClient))
	//msgTool := NewMsgTool(msgDatabase, userDatabase, groupDatabase, conversationDatabase, msgNotificationSender, config)
	//return msgTool, nil
	return nil, nil
}

// func (c *MsgTool) AllConversationClearMsgAndFixSeq() {
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
	ctx := mcontext.NewCtx(stringutil.GetSelfFuncName())
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
		pagination := &sdkws.RequestPagination{
			PageNumber: int32(pageNumber),
			ShowNumber: batchNum,
		}
		conversationIDs, err := c.conversationDatabase.PageConversationIDs(ctx, pagination)
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
		if err := c.msgDatabase.DeleteConversationMsgsAndSetMinSeq(ctx, conversationID, int64(c.config.CronTask.RetainChatRecords*24*60*60)); err != nil {
			log.ZError(ctx, "DeleteUserSuperGroupMsgsAndSetMinSeq failed", err, "conversationID",
				conversationID, "DBRetainChatRecords", c.config.CronTask.RetainChatRecords)
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
		err = fmt.Errorf("cache max seq and mongo max seq is diff > 10,  maxSeqMongo:%d,minSeqMongo:%d,maxSeqCache:%d,conversationID:%s", maxSeqMongo, minSeqMongo, maxSeqCache, conversationID)
		return errs.Wrap(err)
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
		return err
	}
	for _, conversationID := range conversationIDs {
		conversationIDs = append(conversationIDs, conversationutil.GetNotificationConversationIDByConversationID(conversationID))
	}
	for _, conversationID := range conversationIDs {
		if err := c.checkMaxSeq(ctx, conversationID); err != nil {
			log.ZWarn(ctx, "fixSeq failed", err, "conversationID", conversationID)
		}
	}
	fmt.Println("fix all seq finished")
	return nil
}
