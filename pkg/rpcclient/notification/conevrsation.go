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

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

// ConversationNotificationSender
type ConversationNotificationSender struct {
	*rpcclient.NotificationSender
}

// NewConversationNotificationSender
func NewConversationNotificationSender(msgRpcClient *rpcclient.MessageRpcClient) *ConversationNotificationSender {
	return &ConversationNotificationSender{rpcclient.NewNotificationSender(rpcclient.WithRpcClient(msgRpcClient))}
}

// ConversationSetPrivateNotification
func (c *ConversationNotificationSender) ConversationSetPrivateNotification(
	ctx context.Context,
	sendID, recvID string,
	isPrivateChat bool,
) error {
	tips := &sdkws.ConversationSetPrivateTips{
		RecvID:    recvID,
		SendID:    sendID,
		IsPrivate: isPrivateChat,
	}

	return c.Notification(ctx, sendID, recvID, constant.ConversationPrivateChatNotification, tips)
}

// ConversationChangeNotification
func (c *ConversationNotificationSender) ConversationChangeNotification(ctx context.Context, userID string) error {
	tips := &sdkws.ConversationUpdateTips{
		UserID: userID,
	}

	return c.Notification(ctx, userID, userID, constant.ConversationChangeNotification, tips)
}

// ConversationUnreadChangeNotification
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
