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
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"

	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
)

func (c *MsgTool) convertTools() {
	ctx := mcontext.NewCtx("convert")
	conversationIDs, err := c.conversationDatabase.GetAllConversationIDs(ctx)
	if err != nil {
		log.ZError(ctx, "get all conversation ids failed", err)
		return
	}
	for _, conversationID := range conversationIDs {
		conversationIDs = append(conversationIDs, msgprocessor.GetNotificationConversationIDByConversationID(conversationID))
	}
	_, userIDs, err := c.userDatabase.GetAllUserID(ctx, nil)
	if err != nil {
		log.ZError(ctx, "get all user ids failed", err)
		return
	}
	log.ZDebug(ctx, "all userIDs", "len userIDs", len(userIDs))
	for _, userID := range userIDs {
		conversationIDs = append(conversationIDs, msgprocessor.GetConversationIDBySessionType(constant.SingleChatType, userID, userID))
		conversationIDs = append(conversationIDs, msgprocessor.GetNotificationConversationID(constant.SingleChatType, userID, userID))
	}
	log.ZDebug(ctx, "all conversationIDs", "len userIDs", len(conversationIDs))
	c.msgDatabase.ConvertMsgsDocLen(ctx, conversationIDs)
}
