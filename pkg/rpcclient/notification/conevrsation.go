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

package notification

import (
	"context"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/sdkws"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type ConversationNotificationSender struct {
	*rpcclient.NotificationSender
}

func NewConversationNotificationSender(msgRpcClient *rpcclient.MessageRpcClient) *ConversationNotificationSender {
	return &ConversationNotificationSender{rpcclient.NewNotificationSender(rpcclient.WithRpcClient(msgRpcClient))}
}

// SetPrivate调用.
func (c *ConversationNotificationSender) ConversationSetPrivateNotification(ctx context.Context, sendID, recvID string,
	isPrivateChat bool, conversationID string,
) error {
	tips := &sdkws.ConversationSetPrivateTips{
		RecvID:         recvID,
		SendID:         sendID,
		IsPrivate:      isPrivateChat,
		ConversationID: conversationID,
	}
	return c.Notification(ctx, sendID, recvID, constant.ConversationPrivateChatNotification, tips)
}

// 会话改变.
func (c *ConversationNotificationSender) ConversationChangeNotification(ctx context.Context, userID string, conversationIDs []string) error {
	tips := &sdkws.ConversationUpdateTips{
		UserID:             userID,
		ConversationIDList: conversationIDs,
	}
	return c.Notification(ctx, userID, userID, constant.ConversationChangeNotification, tips)
}

// 会话未读数同步.
func (c *ConversationNotificationSender) ConversationUnreadChangeNotification(
	ctx context.Context,
	userID, conversationID string,
	unreadCountTime, hasReadSeq int64,
) error {
	tips := &sdkws.ConversationHasReadTips{
		UserID:          userID,
		ConversationID:  conversationID,
		HasReadSeq:      hasReadSeq,
		UnreadCountTime: unreadCountTime,
	}
	return c.Notification(ctx, userID, userID, constant.ConversationUnreadNotification, tips)
}
