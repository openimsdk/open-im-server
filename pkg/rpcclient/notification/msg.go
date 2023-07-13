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

type MsgNotificationSender struct {
	*rpcclient.NotificationSender
}

func NewMsgNotificationSender(opts ...rpcclient.NotificationSenderOptions) *MsgNotificationSender {
	return &MsgNotificationSender{rpcclient.NewNotificationSender(opts...)}
}

func (m *MsgNotificationSender) UserDeleteMsgsNotification(ctx context.Context, userID, conversationID string, seqs []int64) error {
	tips := sdkws.DeleteMsgsTips{
		UserID:         userID,
		ConversationID: conversationID,
		Seqs:           seqs,
	}
	return m.Notification(ctx, userID, userID, constant.MsgDeleteNotification, &tips)
}

func (m *MsgNotificationSender) MarkAsReadNotification(ctx context.Context, conversationID string, sesstionType int32, sendID, recvID string, seqs []int64, hasReadSeq int64) error {
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
		HasReadSeq:       hasReadSeq,
	}
	return m.NotificationWithSesstionType(ctx, sendID, recvID, constant.HasReadReceipt, sesstionType, tips)
}
