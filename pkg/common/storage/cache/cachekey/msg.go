// Copyright © 2024 OpenIM. All rights reserved.
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

package cachekey

import (
	"strconv"

	"github.com/openimsdk/protocol/constant"
)

const (
	messageCache         = "MESSAGE_CACHE:"
	messageDelUserList   = "MESSAGE_DEL_USER_LIST:"
	userDelMessagesList  = "USER_DEL_MESSAGES_LIST:"
	sendMsgFailedFlag    = "SEND_MSG_FAILED_FLAG:"
	exTypeKeyLocker      = "EX_LOCK:"
	reactionExSingle     = "EX_SINGLE_"
	reactionWriteGroup   = "EX_GROUP_"
	reactionReadGroup    = "EX_SUPER_GROUP_"
	reactionNotification = "EX_NOTIFICATION_"
)

func GetMessageCacheKey(conversationID string, seq int64) string {
	return messageCache + conversationID + "_" + strconv.Itoa(int(seq))
}

func GetMessageDelUserListKey(conversationID string, seq int64) string {
	return messageDelUserList + conversationID + ":" + strconv.Itoa(int(seq))
}

func GetUserDelListKey(conversationID, userID string) string {
	return userDelMessagesList + conversationID + ":" + userID
}

func GetMessageReactionExKey(clientMsgID string, sessionType int32) string {
	switch sessionType {
	case constant.SingleChatType:
		return reactionExSingle + clientMsgID
	case constant.WriteGroupChatType:
		return reactionWriteGroup + clientMsgID
	case constant.ReadGroupChatType:
		return reactionReadGroup + clientMsgID
	case constant.NotificationChatType:
		return reactionNotification + clientMsgID
	}

	return ""
}
func GetLockMessageTypeKey(clientMsgID string, TypeKey string) string {
	return exTypeKeyLocker + clientMsgID + "_" + TypeKey
}

func GetSendMsgKey(id string) string {
	return sendMsgFailedFlag + id
}
