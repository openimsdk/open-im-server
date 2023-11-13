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
	"math/rand"
	"time"

	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

//func (c *MsgTool) ConversationsDestructMsgs() {
//	log.ZInfo(context.Background(), "start msg destruct cron task")
//	ctx := mcontext.NewCtx(utils.GetSelfFuncName())
//	conversations, err := c.conversationDatabase.GetConversationIDsNeedDestruct(ctx)
//	if err != nil {
//		log.ZError(ctx, "get conversation id need destruct failed", err)
//		return
//	}
//	log.ZDebug(context.Background(), "nums conversations need destruct", "nums", len(conversations))
//	for _, conversation := range conversations {
//		ctx = mcontext.NewCtx(utils.GetSelfFuncName() + "-" + utils.OperationIDGenerator() + "-" + conversation.ConversationID + "-" + conversation.OwnerUserID)
//		log.ZDebug(
//			ctx,
//			"UserMsgsDestruct",
//			"conversationID",
//			conversation.ConversationID,
//			"ownerUserID",
//			conversation.OwnerUserID,
//			"msgDestructTime",
//			conversation.MsgDestructTime,
//			"lastMsgDestructTime",
//			conversation.LatestMsgDestructTime,
//		)
//		now := time.Now()
//		seqs, err := c.msgDatabase.UserMsgsDestruct(ctx, conversation.OwnerUserID, conversation.ConversationID, conversation.MsgDestructTime, conversation.LatestMsgDestructTime)
//		if err != nil {
//			log.ZError(ctx, "user msg destruct failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
//			continue
//		}
//		if len(seqs) > 0 {
// 			if err := c.conversationDatabase.UpdateUsersConversationFiled(ctx, []string{conversation.OwnerUserID}, conversation.ConversationID, map[string]interface{}{"latest_msg_destruct_time": now}); err
// != nil {
//				log.ZError(ctx, "updateUsersConversationFiled failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
//				continue
//			}
//			if err := c.msgNotificationSender.UserDeleteMsgsNotification(ctx, conversation.OwnerUserID, conversation.ConversationID, seqs); err != nil {
//				log.ZError(ctx, "userDeleteMsgsNotification failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
//			}
//		}
//	}
//}

func (c *MsgTool) ConversationsDestructMsgs() {
	log.ZInfo(context.Background(), "start msg destruct cron task")
	ctx := mcontext.NewCtx(utils.GetSelfFuncName())
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
		log.ZError(ctx, "PageConversationIDs failed", err, "pageNumber", pageNumber, "conversationIDsNum", len(conversationIDs), "conversationIDs", conversationIDs)
		if len(conversationIDs) == 0 {
			continue
		}
		conversations, err := c.conversationDatabase.GetConversationsByConversationID(ctx, conversationIDs)
		if err != nil {
			log.ZError(ctx, "GetConversationsByConversationID failed", err, "conversationIDs", conversationIDs)
			continue
		}
		temp := make([]*relation.ConversationModel, 0, len(conversations))
		for i, conversation := range conversations {
			if conversation.IsMsgDestruct && conversation.MsgDestructTime != 0 && (time.Now().Unix() > (conversation.MsgDestructTime+conversation.LatestMsgDestructTime.Unix()+8*60*60)) ||
				conversation.LatestMsgDestructTime.IsZero() {
				temp = append(temp, conversations[i])
			}
		}
		for _, conversation := range temp {
			ctx = mcontext.NewCtx(utils.GetSelfFuncName() + "-" + utils.OperationIDGenerator() + "-" + conversation.ConversationID + "-" + conversation.OwnerUserID)
			log.ZDebug(
				ctx,
				"UserMsgsDestruct",
				"conversationID",
				conversation.ConversationID,
				"ownerUserID",
				conversation.OwnerUserID,
				"msgDestructTime",
				conversation.MsgDestructTime,
				"lastMsgDestructTime",
				conversation.LatestMsgDestructTime,
			)
			now := time.Now()
			seqs, err := c.msgDatabase.UserMsgsDestruct(ctx, conversation.OwnerUserID, conversation.ConversationID, conversation.MsgDestructTime, conversation.LatestMsgDestructTime)
			if err != nil {
				log.ZError(ctx, "user msg destruct failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
				continue
			}
			if len(seqs) > 0 {
				if err := c.conversationDatabase.UpdateUsersConversationFiled(ctx, []string{conversation.OwnerUserID}, conversation.ConversationID, map[string]interface{}{"latest_msg_destruct_time": now}); err != nil {
					log.ZError(ctx, "updateUsersConversationFiled failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
					continue
				}
				if err := c.msgNotificationSender.UserDeleteMsgsNotification(ctx, conversation.OwnerUserID, conversation.ConversationID, seqs); err != nil {
					log.ZError(ctx, "userDeleteMsgsNotification failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
				}
			}
		}
	}
}
